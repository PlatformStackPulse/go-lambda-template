// Package version provides version information.
package version

import "fmt"

// Version will be set during build
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
	GoVersion = "unknown"
)

// Info contains version information
type Info struct {
	Version   string
	Commit    string
	BuildTime string
	GoVersion string
}

// Get returns the version info
func Get() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		BuildTime: BuildTime,
		GoVersion: GoVersion,
	}
}

// String returns a formatted version string
func (i Info) String() string {
	return fmt.Sprintf("Version: %s\nCommit: %s\nBuild Time: %s\nGo Version: %s",
		i.Version, i.Commit, i.BuildTime, i.GoVersion)
}
