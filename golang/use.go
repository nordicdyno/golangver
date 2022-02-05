package golang

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/coreos/go-semver/semver"

	"github.com/nordicdyno/golangver/uitools"
)

// UseVersion sets goBinPath symlink to requested Go version (must be installed).
func UseVersion(goBinPath string, version string) error {
	goRoot, err := golangBinGOROOT("go" + version)
	if err != nil {
		return err
	}

	newBin := goBin(goRoot)
	if err := UseBinary(goBinPath, newBin); err != nil {
		return err
	}

	if err := patchSDKForIDEA(goRoot); err != nil {
		return fmt.Errorf("patch IDEA config is failed: %w", err)
	}
	if err := patchGoMod(version); err != nil {
		return fmt.Errorf("patch go.mod is failed: %w", err)
	}
	return nil
}

// UseBinary sets goBinPath symlink to requested binary path.
func UseBinary(goBinPath string, newBin string) error {
	currentPath, err := goBinCheckSymlink(goBinPath)
	if err != nil {
		return err
	}

	if currentPath != "" {
		if err := os.Remove(goBinPath); err != nil {
			return err
		}
	}

	if err := os.Symlink(newBin, goBinPath); err != nil {
		return fmt.Errorf("symlink %s -> %s failed: %w", goBinPath, newBin, err)
	}
	fmt.Printf("set symlink %s -> %s\n", goBinPath, newBin)
	fmt.Printf("(previous value was: %s)\n", currentPath)

	// detect and patch stage:
	return nil
}

const goModFile = "go.mod"

func patchGoMod(version string) error {
	version = fixVersion(version)
	sVer := semver.New(version)

	version = fmt.Sprintf("%d.%d", sVer.Major, sVer.Minor)
	if _, err := os.Stat(goModFile); err != nil {
		return nil
	}

	modInfo, err := goModParse(goModFile)
	if err != nil {
		return fmt.Errorf("read of %s is failed: %w", goModFile, err)
	}
	if modInfo.Go.Version == version {
		return nil
	}

	fmt.Println("\ngo.mod is detected:")
	fmt.Printf("  current value: %s\n", modInfo.Go.Version)
	yes, err := uitools.InputYesNo(fmt.Sprintf(
		"Do you want to set Go version = %s", version), false)
	if err != nil {
		return err
	}
	if !yes {
		return nil
	}

	if err := runCmd("go", "mod", "edit", "-go="+version); err != nil {
		return fmt.Errorf("go mod edit is failed: %w", err)
	}
	var modTidyArgs = []string{"mod", "tidy"}
	if sVer.Minor >= 17 {
		modTidyArgs = append(modTidyArgs, "-compat="+version, "-go="+version)
	}
	if err := runCmd("go", modTidyArgs...); err != nil {
		return fmt.Errorf("go mod tidy is failed: %w", err)
	}
	return nil
}

const ideaDir = ".idea"

func patchSDKForIDEA(goRoot string) error {
	_, err := os.ReadDir(ideaDir)
	if err != nil {
		// fmt.Println(".idea not found:", err)
		return nil
	}

	workplaceFile := filepath.Join(ideaDir, "workspace.xml")
	b, err := os.ReadFile(workplaceFile)
	if err != nil {
		// fmt.Println("ReadFile is failed", err)
		return nil
	}

	var buf bytes.Buffer
	decoder := xml.NewDecoder(bytes.NewReader(b))
	encoder := xml.NewEncoder(&buf)

	// project/ component name="GOROOT
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "error getting token: %v\n", err)
			break
		}

		switch v := token.(type) {
		case xml.StartElement:
			if v.Name.Local == "component" {
				// fmt.Println("found component")
				// fmt.Printf("Attrs %#v\n", v.Attr)
				var isGoroot bool
				for _, attr := range v.Attr {
					if attr.Name.Local == "name" && attr.Value == "GOROOT" {
						isGoroot = true
						break
					}
				}
				if !isGoroot {
					break
				}
				var foundSDK string
				for i, attr := range v.Attr {
					if attr.Name.Local == "url" {
						foundSDK = attr.Value
						setGoRoot := "file://" + substituteUserHomeDirInPathLikeIDEADoes(goRoot)
						if foundSDK == setGoRoot {
							return nil
						}
						// fmt.Println("found:", foundSDK)
						// fmt.Println("candidate:", setGoRoot)
						fmt.Printf("\n%s is detected\n", workplaceFile)
						fmt.Printf("  current value: %s\n", foundSDK)
						yes, err := uitools.InputYesNo(fmt.Sprintf(
							"Do you want to set Go SDK = %s?", goRoot), false)
						if err != nil {
							return err
						}

						if !yes {
							return nil
						}
						// fmt.Println(".idea detected do you want to set Go SDK to %% [y/N]")
						// fmt.Println("Go SDK:", foundSDK)
						// fmt.Printf("change it to %v?\n", goRoot)
						v.Attr[i].Value = "file://" + substituteUserHomeDirInPathLikeIDEADoes(goRoot)
						break
					}
				}
			}
		}

		if err := encoder.EncodeToken(xml.CopyToken(token)); err != nil {
			return fmt.Errorf("EncodeToken failed: %w", err)
		}
	}

	// must call flush, otherwise some elements will be missing
	if err := encoder.Flush(); err != nil {
		return fmt.Errorf("xml encoder flush failed: %w", err)
	}

	fmt.Println("overwrite", workplaceFile)
	fmt.Println("INFO: project reopening in IDEA is required")
	return os.WriteFile(workplaceFile, buf.Bytes(), 0644)
}

func substituteUserHomeDirInPathLikeIDEADoes(path string) string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	if !strings.HasPrefix(path, homedir) {
		return path
	}
	return "$USER_HOME$" + path[len(homedir):]
}
