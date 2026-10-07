package system

import (
	"testing"
)

func TestCollectSystemInfoSmoke(t *testing.T) {
	info, err := CollectSystemInfo()
	if err != nil {
		t.Fatalf("CollectSystemInfo failed: %v", err)
	}

	if info.OS == "" {
		t.Error("Expected OS to be populated")
	}
	if info.Architecture == "" {
		t.Error("Expected Architecture to be populated")
	}
	if info.CPU == nil || info.CPU.Cores <= 0 {
		t.Errorf("Expected valid CPU info, got %+v", info.CPU)
	}
	if info.WorkingDir == "" {
		t.Error("Expected WorkingDir to be populated")
	}
}
