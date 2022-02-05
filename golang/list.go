package golang

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/coreos/go-semver/semver"
)

// ListOpts controls List behaviour.
type ListOpts struct {
	// ShowRemotes adds to list available remote versions.
	ShowRemotes bool
	// ShowRemotes adds to list all remote versions, not only latest minor version for every major version.
	ShowAllRemotes bool
	// ShowOutdated adds to list even old Go versions (older than 1.13).
	ShowOutdated bool
}

// List shows Go versions available locally and remotely.
func List(linkPath string, opts ListOpts) error {
	currentTarget, err := goBinCheckSymlink(linkPath)
	if err != nil {
		return fmt.Errorf("check symlink %s is failed: %w", linkPath, err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("user home dir resolving is failed: %w", err)
	}

	// parse downloaded by IDEA SDK directories
	pathIdeaSDK := filepath.Join(homeDir, "go")
	namesIDEA, err := filepath.Glob(filepath.Join(pathIdeaSDK, "go1.*"))
	if err != nil {
		return fmt.Errorf("list of %v is failed: %w", pathIdeaSDK, err)
	}

	var ideaVersions versionList
	for _, binPath := range namesIDEA {
		goBinPath := filepath.Join(binPath, "bin", "go")

		_, name := filepath.Split(binPath)
		name = name[2:]
		v, err := parseVersionInfo(name)
		if err != nil {
			return fmt.Errorf("parse version failed %v: %w", name, err)
		}
		v.binPath = goBinPath
		ideaVersions = append(ideaVersions, *v)
	}
	ideaVersions.Sort()

	// go* binaries downloaded by `go install go*`
	var dlVersions versionList
	// parse downloaded Go SDK directories by `go install golang.org/dl/go1.*`
	namesDl, err := filepath.Glob(filepath.Join(golangBinDir(), "go1.*"))
	if err != nil {
		return fmt.Errorf("list of %v is failed: %w", golangBinDir(), err)
	}
	for _, binPath := range namesDl {
		_, name := filepath.Split(binPath)
		goBinPath, err := parseGolangBin(binPath)
		if err != nil {
			if errors.Is(err, ErrNotDownloaded) {
				continue
			}
			return fmt.Errorf("go bin path detection failed: %w", err)
		}

		name = name[2:]
		v, err := parseVersionInfo(name)
		if err != nil {
			return fmt.Errorf("failed parse version %v: %w", name, err)
		}
		v.binPath = goBinPath
		dlVersions = append(dlVersions, *v)
	}

	dlVersions.Sort()
	var currentFound bool
	printVersions := func(vl versionList) {
		for _, v := range vl {
			mark := " "
			if !currentFound && v.binPath == currentTarget {
				currentFound = true
				mark = "*"
			}
			out := fmt.Sprintf("%s %-10s  %s", mark, v.original, v.binPath)
			fmt.Println(out)
		}
	}

	fmt.Println(" downloaded by `go install`:")
	printVersions(dlVersions)
	if len(ideaVersions) > 0 {
		fmt.Print("\n downloaded by IDEA:\n")
		printVersions(ideaVersions)
	}

	if !currentFound {
		fmt.Println("currentTarget:", currentTarget)
	}

	if opts.ShowRemotes {
		return showRemoteGoVersions(opts.ShowAllRemotes, opts.ShowOutdated)
	}
	return nil
}

var lastNonOutdatedVersion = "1.13.0"

func showRemoteGoVersions(showAll bool, showOutdated bool) error {
	var output bytes.Buffer
	gitArgs := []string{"ls-remote", "-t", "https://github.com/golang/go"}
	cmd := exec.Command("git", gitArgs...)
	cmd.Stdout = &output
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command git %s failed: %w", strings.Join(gitArgs, " "), err)
	}

	var prefix = "refs/tags/go"
	var minVersion = semver.New(lastNonOutdatedVersion) // September 2019

	// var versions []*semver.Version
	var foundMaxVersion *semver.Version
	var versions versionList
	scanner := bufio.NewScanner(&output)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "\t")
		if len(parts) != 2 {
			continue
		}
		name := parts[1]
		if !strings.HasPrefix(name, prefix) {
			continue
		}

		name = name[len(prefix):]
		v, err := parseVersionInfo(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parseVersionInfo err=%v\n", err)
			continue
		}
		if !showOutdated && v.semver.LessThan(*minVersion) {
			continue
		}
		versions = append(versions, *v)

		if foundMaxVersion == nil || foundMaxVersion.LessThan(*v.semver) {
			foundMaxVersion = v.semver
		}
	}

	versions.Sort()
	fmt.Print("\n# remote Go versions:\n")

	// "https://golang.org/doc/go1.17"
	var lastMinor int64
	for _, v := range versions {
		var extraInfo = ""
		if v.betaSuffix != "" && v.semver.Minor != foundMaxVersion.Minor {
			continue
		}

		if v.semver.Minor == lastMinor {
			if !showAll {
				continue
			}
		} else {
			suffix := ""
			lastMinor = v.semver.Minor
			if lastMinor != 0 {
				suffix = "." + strconv.Itoa(int(lastMinor))
			}
			if v.betaSuffix == "" {
				extraInfo = "\thttps://golang.org/doc/devel/release#go1" + suffix
			}
		}
		if v.betaSuffix != "" && v.semver.Minor == 18 {
			extraInfo = "\thttps://go.dev/blog/go" + v.original
		}

		fmt.Printf("  %s%s\n", v.original, extraInfo)
	}

	return nil
}
