package golang

import (
	"os"

	"golang.org/x/mod/modfile"
)

// GoModParse parses go.mod file.
func goModParse(file string) (*modfile.File, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return modfile.Parse(file, b, nil)
}
