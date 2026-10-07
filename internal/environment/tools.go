package environment

import (
	"os/exec"
)

// DefaultToolsList represents the standard developer toolset checked by hann doctor.
var DefaultToolsList = []string{
	"git",
	"go",
	"node",
	"npm",
	"pnpm",
	"python",
	"python3",
	"curl",
	"wget",
	"docker",
	"gcc",
	"make",
}

// ToolStatus represents the status of an executable on the system.
type ToolStatus struct {
	Name      string
	Available bool
	Path      string
}

// CheckTool tests whether a specific binary is located in PATH using exec.LookPath.
func CheckTool(name string) ToolStatus {
	path, err := exec.LookPath(name)
	return ToolStatus{
		Name:      name,
		Available: err == nil,
		Path:      path,
	}
}

// CheckTools verifies availability for a given slice of binary names.
func CheckTools(names []string) []ToolStatus {
	results := make([]ToolStatus, len(names))
	for i, name := range names {
		results[i] = CheckTool(name)
	}
	return results
}

// CheckDefaultTools verifies all tools from the default suite.
func CheckDefaultTools() []ToolStatus {
	return CheckTools(DefaultToolsList)
}
