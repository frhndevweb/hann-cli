package system

import (
	"time"

	"github.com/frhndevweb/hann-cli/internal/environment"
)

// Info aggregates real-time hardware and operating system status.
type Info struct {
	OS           string
	Kernel       string
	Architecture string
	Hostname     string
	User         string
	Shell        string
	CPU          *CPUInfo
	Memory       *MemoryInfo
	Uptime       time.Duration
	WorkingDir   string
	IsPRoot      bool
	IsTermux     bool
}

// CollectSystemInfo gathers real-time system metrics without hardcoded values or external APIs.
func CollectSystemInfo() (*Info, error) {
	env := environment.DetectEnvironment()
	cpu := GetCPUInfo()

	mem, err := GetMemoryInfo()
	if err != nil {
		// Create a fallback zeroed memory struct so callers don't crash
		mem = &MemoryInfo{}
	}

	upt, err := GetUptime()
	if err != nil {
		upt = 0
	}

	return &Info{
		OS:           env.OS,
		Kernel:       env.Kernel,
		Architecture: env.Architecture,
		Hostname:     env.Hostname,
		User:         env.User,
		Shell:        env.Shell,
		CPU:          cpu,
		Memory:       mem,
		Uptime:       upt,
		WorkingDir:   env.WorkingDir,
		IsPRoot:      env.IsPRoot,
		IsTermux:     env.IsTermux,
	}, nil
}
