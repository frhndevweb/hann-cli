package version

import (
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	info := Get()
	if info.Version == "" {
		t.Error("Expected non-empty version")
	}
	if info.GoVersion == "" {
		t.Error("Expected non-empty GoVersion")
	}
	if info.Platform == "" {
		t.Error("Expected non-empty Platform")
	}
}

func TestFormatted(t *testing.T) {
	formatted := Formatted()
	if !strings.Contains(formatted, "hann") {
		t.Errorf("Formatted string missing hann name: %s", formatted)
	}
}
