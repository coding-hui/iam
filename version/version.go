package version

import (
	"fmt"
	"github.com/hashicorp/go-version"
	"os"
	"runtime"
)

var (
	// GitRevision is the commit of repo
	GitRevision = "UNKNOWN"
	// IAMVersion is the version of cli.
	IAMVersion = "UNKNOWN"
	// Built shows the built time of the binary.
	Built      = "Not provided."
	apiVersion = "v1alpha1"
)

// PrintVersionAndExit prints versions from the array returned by Info() and exit.
func PrintVersionAndExit() {
	for _, i := range Info(apiVersion) {
		fmt.Printf("%v\n", i)
	}
	os.Exit(0)
}

// Info returns an array of various service versions.
func Info(apiVersion string) []string {
	return []string{
		fmt.Sprintf("API Version: %s", apiVersion),
		fmt.Sprintf("Version: %s", IAMVersion),
		fmt.Sprintf("GitRevision: %s", GitRevision),
		fmt.Sprintf("Built At: %s", Built),
		fmt.Sprintf("Go Version: %s", runtime.Version()),
		fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// IsOfficialIAMVersion checks whether the provided version string follows a IAM version pattern
func IsOfficialIAMVersion(versionStr string) bool {
	_, err := version.NewSemver(versionStr)
	return err == nil
}

// GetOfficialIAMVersion extracts the IAM version from the provided string
// More precisely, this method returns the segments and prerelease info w/o metadata
func GetOfficialIAMVersion(versionStr string) (string, error) {
	s, err := version.NewSemver(versionStr)
	if err != nil {
		return "", err
	}
	v := s.String()
	metadata := s.Metadata()
	if metadata != "" {
		metadata = "+" + metadata
	}
	return v[:len(v)-len(metadata)], nil
}
