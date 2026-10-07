package environment

import (
	"testing"
)

func TestParseOSRelease(t *testing.T) {
	ubuntuSample := `NAME="Ubuntu"
VERSION="24.04 LTS (Noble Numbat)"
ID=ubuntu
ID_LIKE=debian
PRETTY_NAME="Ubuntu 24.04 LTS"
VERSION_ID="24.04"
`
	got := ParseOSRelease(ubuntuSample)
	expected := "Ubuntu 24.04 LTS"
	if got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}

	fallbackSample := `NAME="Alpine Linux"
VERSION_ID="3.19.1"
`
	gotFallback := ParseOSRelease(fallbackSample)
	if gotFallback != "Alpine Linux" {
		t.Errorf("Expected %q, got %q", "Alpine Linux", gotFallback)
	}
}

func TestDetectEnvironmentSmoke(t *testing.T) {
	env := DetectEnvironment()
	if env == nil {
		t.Fatal("Expected non-nil EnvironmentInfo")
	}
	if env.OS == "" {
		t.Error("Expected OS to not be empty")
	}
	if env.Architecture == "" {
		t.Error("Expected Architecture to not be empty")
	}
	if env.WorkingDir == "" {
		t.Error("Expected WorkingDir to not be empty")
	}
}

func TestDetectTermuxAndPRootFlags(t *testing.T) {
	// Temporarily override env vars to test detection logic
	t.Setenv("TERMUX_VERSION", "0.118.0")
	if !DetectTermux() {
		t.Error("Expected DetectTermux to return true when TERMUX_VERSION is set")
	}

	t.Setenv("container", "proot-distro")
	if !DetectPRoot() {
		t.Error("Expected DetectPRoot to return true when container=proot-distro")
	}
}
