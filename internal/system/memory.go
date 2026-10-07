package system

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// MemoryInfo holds system memory metrics.
type MemoryInfo struct {
	TotalBytes     uint64
	AvailableBytes uint64
	UsedBytes      uint64
	FreeBytes      uint64
	BuffersBytes   uint64
	CachedBytes    uint64
	UsedPercent    float64
}

// GetMemoryInfo reads and parses memory metrics from /proc/meminfo.
func GetMemoryInfo() (*MemoryInfo, error) {
	content, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return nil, fmt.Errorf("unable to read /proc/meminfo: %w", err)
	}

	return ParseMeminfo(string(content))
}

// ParseMeminfo parses raw /proc/meminfo lines into MemoryInfo.
func ParseMeminfo(content string) (*MemoryInfo, error) {
	var (
		memTotal     uint64
		memFree      uint64
		memAvailable uint64
		buffers      uint64
		cached       uint64
	)

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		fields := strings.Fields(parts[1])
		if len(fields) == 0 {
			continue
		}

		valKb, err := strconv.ParseUint(fields[0], 10, 64)
		if err != nil {
			continue
		}

		switch key {
		case "MemTotal":
			memTotal = valKb
		case "MemFree":
			memFree = valKb
		case "MemAvailable":
			memAvailable = valKb
		case "Buffers":
			buffers = valKb
		case "Cached":
			cached = valKb
		}
	}

	if memTotal == 0 {
		return nil, fmt.Errorf("invalid /proc/meminfo: MemTotal is 0 or missing")
	}

	totalBytes := memTotal * 1024
	var availableBytes uint64
	if memAvailable > 0 {
		availableBytes = memAvailable * 1024
	} else {
		availableBytes = (memFree + buffers + cached) * 1024
	}

	if availableBytes > totalBytes {
		availableBytes = totalBytes
	}
	usedBytes := totalBytes - availableBytes

	var usedPercent float64
	if totalBytes > 0 {
		usedPercent = (float64(usedBytes) / float64(totalBytes)) * 100.0
	}

	return &MemoryInfo{
		TotalBytes:     totalBytes,
		AvailableBytes: availableBytes,
		UsedBytes:      usedBytes,
		FreeBytes:      memFree * 1024,
		BuffersBytes:   buffers * 1024,
		CachedBytes:    cached * 1024,
		UsedPercent:    usedPercent,
	}, nil
}

// FormatBytes converts raw bytes into clean human-readable representation.
func FormatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	val := float64(b) / float64(div)
	suffix := []string{"KB", "MB", "GB", "TB", "PB"}[exp]
	if val >= 10 {
		return fmt.Sprintf("%.1f %s", val, suffix)
	}
	return fmt.Sprintf("%.2f %s", val, suffix)
}

// Summary returns a formatted used/total memory string with percentage.
func (m *MemoryInfo) Summary() string {
	if m == nil || m.TotalBytes == 0 {
		return "N/A"
	}
	return fmt.Sprintf("%s / %s (%.0f%%)", FormatBytes(m.UsedBytes), FormatBytes(m.TotalBytes), m.UsedPercent)
}
