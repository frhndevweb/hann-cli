package environment

import (
	"testing"
)

func TestCheckToolGo(t *testing.T) {
	// "go" executable should definitely exist during test execution
	res := CheckTool("go")
	if !res.Available {
		t.Error("Expected 'go' binary to be detected")
	}
	if res.Path == "" {
		t.Error("Expected non-empty path for 'go' binary")
	}
}

func TestCheckToolMissing(t *testing.T) {
	res := CheckTool("non_existent_tool_xyz_998877")
	if res.Available {
		t.Error("Expected dummy tool to not be available")
	}
	if res.Path != "" {
		t.Errorf("Expected empty path for non-existent tool, got %q", res.Path)
	}
}

func TestCheckDefaultTools(t *testing.T) {
	tools := CheckDefaultTools()
	if len(tools) != len(DefaultToolsList) {
		t.Fatalf("Expected %d tools checked, got %d", len(DefaultToolsList), len(tools))
	}
}
