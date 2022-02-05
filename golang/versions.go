package golang

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/coreos/go-semver/semver"
)

type versionList []versionInfo

type versionInfo struct {
	binPath    string
	semver     *semver.Version
	original   string
	betaSuffix string
}

var lastVersionPartRe = regexp.MustCompile(`([0-9]+)(.*)`)

// fixVersion fixes version without patch part.
// Example: 1.17 -> 1.17.0.
func fixVersion(version string) string {
	versionParts := strings.Split(version, ".")
	lastPart := versionParts[len(versionParts)-1]
	match := lastVersionPartRe.FindStringSubmatch(lastPart)

	// fmt.Println(" fixVersion/match:", match)
	if match[2] != "" {
		// fmt.Println("lastPart:", match[2])
		versionParts[len(versionParts)-1] = match[1]
	}
	// fmt.Println(" fixVersion:", versionParts)

	if len(versionParts) == 2 {
		versionParts = append(versionParts, "0")
	}
	return strings.Join(versionParts, ".")
}

func parseVersionInfo(version string) (*versionInfo, error) {
	versionParts := strings.Split(version, ".")
	lastPart := versionParts[len(versionParts)-1]
	match := lastVersionPartRe.FindStringSubmatch(lastPart)

	if match[2] != "" {
		versionParts[len(versionParts)-1] = match[1]
	}

	if len(versionParts) == 1 {
		versionParts = append(versionParts, "0")
	}
	if len(versionParts) == 2 {
		versionParts = append(versionParts, "0")
	}

	cleanedVersion := strings.Join(versionParts, ".")
	semVer, err := semver.NewVersion(cleanedVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %v: %w", cleanedVersion, err)
	}

	return &versionInfo{
		semver:     semVer,
		original:   version,
		betaSuffix: match[2],
	}, nil
}

func (v versionList) Len() int {
	return len(v)
}

func (v versionList) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v versionList) Less(i, j int) bool {
	return !v[i].semver.LessThan(*v[j].semver)
}

// Sort sorts the given slice of Version
func (v versionList) Sort() {
	sort.Sort(v)
}
