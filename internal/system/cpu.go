package system

import (
	"bufio"
	"os"
	"runtime"
	"strings"
)

// CPUInfo represents processor hardware details.
type CPUInfo struct {
	Cores int
	Model string
}

// GetCPUInfo gathers processor statistics from /proc/cpuinfo and runtime.
func GetCPUInfo() *CPUInfo {
	info := &CPUInfo{
		Cores: runtime.NumCPU(),
		Model: "Generic Processor",
	}

	content, err := os.ReadFile("/proc/cpuinfo")
	if err == nil {
		parsed := ParseCPUInfo(string(content))
		if parsed.Cores > 0 {
			info.Cores = parsed.Cores
		}
		if parsed.Model != "" && parsed.Model != "Generic Processor" {
			info.Model = parsed.Model
		}
	}

	return info
}

// ParseCPUInfo parses raw /proc/cpuinfo content.
func ParseCPUInfo(content string) CPUInfo {
	var cores int
	var model, hardware string

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.Contains(line, ":") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		switch {
		case key == "processor":
			cores++
		case key == "model name" && model == "":
			model = val
		case (key == "Hardware" || key == "vendor_id") && hardware == "":
			hardware = val
		case key == "cpu model" && model == "":
			model = val
		}
	}

	finalModel := model
	if finalModel == "" {
		finalModel = hardware
	}
	if finalModel == "" {
		finalModel = "Generic Processor"
	}

	return CPUInfo{
		Cores: cores,
		Model: finalModel,
	}
}
