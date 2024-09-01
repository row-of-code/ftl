package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/enclave-ci/aerie/internal/ssh"
)

func FindSSHKey(keyPath string, isRoot bool) ([]byte, error) {
	if keyPath != "" {
		return os.ReadFile(keyPath)
	}

	// Try to find key in .ssh directory
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %v", err)
	}

	sshDir := filepath.Join(home, ".ssh")
	keyNames := []string{"id_rsa", "id_ecdsa", "id_ed25519"}

	for _, name := range keyNames {
		path := filepath.Join(sshDir, name)
		if _, err := os.Stat(path); err == nil {
			key, err := os.ReadFile(path)
			if err == nil {
				return key, nil
			}
		}
	}

	if isRoot {
		return nil, fmt.Errorf("no suitable SSH key found in .ssh directory")
	}
	return nil, nil
}

func FindKeyAndConnect(host, rootKeyPath string) (*ssh.Client, []byte, error) {
	rootKey, err := FindSSHKey(rootKeyPath, true)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find root SSH key: %v", err)
	}

	client, err := ssh.Connect(host, rootKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to establish SSH connection: %v", err)
	}

	return client, rootKey, nil
}
