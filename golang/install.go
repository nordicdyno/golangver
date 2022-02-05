package golang

import (
	"fmt"
	"os"
	"path/filepath"
)

// Install installs requested Golang version.
// does:
// go install golang.org/dl/go1.10.7@latest
// go1.10.7 download
func Install(version string, force bool) error {
	fmt.Printf("Install helper tool for %v...\n", version)
	if force {
		if err := removeIfExists(version); err != nil {
			return err
		}
	}
	if err := runCmd("go", "install", fmt.Sprintf("golang.org/dl/go%s@latest", version)); err != nil {
		return err
	}
	fmt.Printf("Download Go version %v...\n", version)
	return runCmd(fmt.Sprintf("go%s", version), "download")
}

func removeIfExists(version string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("user home dir resolving is failed: %w", err)
	}

	// parse downloaded Go SDK directories by `go install golang.org/dl/go1.*`
	pathDlSDK := filepath.Join(homeDir, "sdk", "go"+version)
	if _, err := os.Stat(pathDlSDK); os.IsNotExist(err) {
		return nil
	}
	fmt.Printf("os.RemoveAll(%s)\n", pathDlSDK)
	return os.RemoveAll(pathDlSDK)
}
