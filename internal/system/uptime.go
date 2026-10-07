package system

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// GetUptime reads system uptime from /proc/uptime.
func GetUptime() (time.Duration, error) {
	content, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0, fmt.Errorf("unable to read /proc/uptime: %w", err)
	}

	return ParseUptime(string(content))
}

// ParseUptime parses raw /proc/uptime content into a time.Duration.
func ParseUptime(content string) (time.Duration, error) {
	fields := strings.Fields(content)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty /proc/uptime content")
	}

	seconds, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid /proc/uptime format: %w", err)
	}

	return time.Duration(seconds * float64(time.Second)), nil
}

// FormatUptime renders a Duration into a human-friendly format (e.g., 2h 31m).
func FormatUptime(d time.Duration) string {
	totalSeconds := int64(d.Seconds())
	if totalSeconds < 0 {
		return "0s"
	}

	days := totalSeconds / 86400
	hours := (totalSeconds % 86400) / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	switch {
	case days > 0:
		if minutes > 0 {
			return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
		}
		return fmt.Sprintf("%dd %dh", days, hours)
	case hours > 0:
		if minutes > 0 {
			return fmt.Sprintf("%dh %dm", hours, minutes)
		}
		return fmt.Sprintf("%dh", hours)
	case minutes > 0:
		if seconds > 0 {
			return fmt.Sprintf("%dm %ds", minutes, seconds)
		}
		return fmt.Sprintf("%dm", minutes)
	default:
		return fmt.Sprintf("%ds", seconds)
	}
}
