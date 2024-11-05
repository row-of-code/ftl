package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yarlson/ftl/pkg/ssh"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/yarlson/ftl/pkg/config"
	"github.com/yarlson/ftl/pkg/console"
	"github.com/yarlson/ftl/pkg/deployment"
	"github.com/yarlson/ftl/pkg/runner/remote"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy your application to configured servers",
	Long: `Deploy your application to all servers defined in ftl.yaml.
This command handles the entire deployment process, ensuring
zero-downtime updates of your services.`,
	Run: runDeploy,
}

func init() {
	rootCmd.AddCommand(deployCmd)
}

func runDeploy(cmd *cobra.Command, args []string) {
	cfg, err := parseConfig("ftl.yaml")
	if err != nil {
		console.Error("Failed to parse config file:", err)
		return
	}

	if err := deployToServers(cfg); err != nil {
		console.Error("Deployment failed:", err)
		return
	}
}

func parseConfig(filename string) (*config.Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg, err := config.ParseConfig(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cfg, nil
}

func deployToServers(cfg *config.Config) error {
	for _, server := range cfg.Servers {
		if err := deployToServer(cfg.Project.Name, cfg, server); err != nil {
			return fmt.Errorf("failed to deploy to server %s: %w", server.Host, err)
		}
		console.Success(fmt.Sprintf("Successfully deployed to server %s", server.Host))
	}

	return nil
}

func deployToServer(project string, cfg *config.Config, server config.Server) error {
	console.Info(fmt.Sprintf("Deploying to server %s...", server.Host))

	runner, err := connectToServer(server)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	defer runner.Close()

	deploy := deployment.NewDeployment(runner)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	events := deploy.Deploy(ctx, project, cfg)

	var spinner *pterm.SpinnerPrinter

	for event := range events {
		switch event.Type {
		case deployment.EventTypeStart:
			if spinner != nil {
				spinner.Success()
			}
			spinner = console.NewSpinner(event.Message)
		case deployment.EventTypeProgress:
			if spinner != nil {
				spinner.UpdateText(event.Message)
			} else {
				console.Info(event.Message)
			}
		case deployment.EventTypeFinish:
			if spinner != nil {
				spinner.Success(event.Message)
				spinner = nil
			} else {
				console.Success(event.Message)
			}
		case deployment.EventTypeError:
			if spinner != nil {
				spinner.Fail(fmt.Sprintf("Deployment error: %s", event.Message))
			} else {
				console.Error(fmt.Sprintf("Deployment error: %s", event.Message))
			}
			return fmt.Errorf("deployment error: %s", event.Message)
		case deployment.EventTypeComplete:
			if spinner != nil {
				spinner.Success(event.Message)
			} else {
				console.Success(event.Message)
			}
		default:
			if spinner != nil {
				spinner.UpdateText(event.Message)
			} else {
				console.Info(event.Message)
			}
		}
	}

	if spinner != nil {
		_ = spinner.Stop()
	}

	return nil
}

func connectToServer(server config.Server) (*remote.Runner, error) {
	sshKeyPath := filepath.Join(os.Getenv("HOME"), ".ssh", filepath.Base(server.SSHKey))
	sshClient, _, err := ssh.FindKeyAndConnectWithUser(server.Host, server.Port, server.User, sshKeyPath)
	if err != nil {
		return nil, err
	}

	return remote.NewRunner(sshClient), nil
}
