package version

import (
	"fmt"
	"runtime"
)

// Default build-time values, overridable with -ldflags.
var (
	Version   = "0.1.0"
	Commit    = "dev"
	BuildDate = "unknown"
)

// Info holds complete version and runtime metadata.
type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"build_date"`
	GoVersion string `json:"go_version"`
	Compiler  string `json:"compiler"`
	Platform  string `json:"platform"`
}

// Get returns the populated Info struct.
func Get() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		BuildDate: BuildDate,
		GoVersion: runtime.Version(),
		Compiler:  runtime.Compiler,
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// Formatted returns a clean multi-line representation of the version.
func Formatted() string {
	info := Get()
	return fmt.Sprintf("hann %s (%s)\nBuilt: %s\nTarget: %s\nGo: %s",
		info.Version, info.Commit, info.BuildDate, info.Platform, info.GoVersion)
}
