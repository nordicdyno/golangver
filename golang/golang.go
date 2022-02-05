// Package golang provides tools for Go version management on localhost.
package golang

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/build"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func runCmd(bin string, args ...string) error {
	cmd := exec.Command(bin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s %s: %w", bin, strings.Join(args, " "), err)
	}
	return nil
}

// BinaryPath returns path to go binary for provided version.
func BinaryPath(version string) (string, error) {
	goRoot, err := golangBinGOROOT("go" + version)
	if err != nil {
		return "", err
	}
	return goBin(goRoot), nil
}

// goBinCheckSymlink cheks is provided path symlink if exists
func goBinCheckSymlink(fpath string) (string, error) {
	fileInfo, err := os.Lstat(fpath)
	if err != nil {
		if _, ok := err.(*fs.PathError); !ok {
			return "", err
		}
		return "", nil
	}

	if fileInfo.Mode()&os.ModeSymlink != os.ModeSymlink {
		return "", fmt.Errorf("%s is not symlink", fpath)
	}

	originFile, err := os.Readlink(fpath)
	if err != nil {
		return "", err
	}
	return originFile, nil
}

var ErrNotDownloaded = fmt.Errorf("not downloaded")

func golangBinGOROOT(bin string) (string, error) {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	cmd := exec.Command(bin, "env", "-json", "GOROOT")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	// cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if strings.Contains(stderr.String(), "not downloaded.") {
			return "", ErrNotDownloaded
		}
		return "", err
	}

	type goEnv struct {
		GOROOT string
	}
	var goe goEnv
	if err := json.NewDecoder(&stdout).Decode(&goe); err != nil {
		return "", err
	}
	return goe.GOROOT, nil
}

func parseGolangBin(bin string) (string, error) {
	goRoot, err := golangBinGOROOT(bin)
	if err != nil {
		return "", err
	}
	return goBin(goRoot), nil
}

func goBin(goRoot string) string {
	return filepath.Join(goRoot, "bin", "go")
}

func golangBinDir() string {
	return filepath.Join(build.Default.GOPATH, "bin")
}
