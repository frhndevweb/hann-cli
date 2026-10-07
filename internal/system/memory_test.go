package system

import (
	"testing"
)

func TestParseMeminfo(t *testing.T) {
	sample := `MemTotal:        8000000 kB
MemFree:         1000000 kB
MemAvailable:    4000000 kB
Buffers:          200000 kB
Cached:          1800000 kB
`
	mem, err := ParseMeminfo(sample)
	if err != nil {
		t.Fatalf("Unexpected error parsing meminfo: %v", err)
	}

	expectedTotal := uint64(8000000 * 1024)
	if mem.TotalBytes != expectedTotal {
		t.Errorf("Expected total %d, got %d", expectedTotal, mem.TotalBytes)
	}

	expectedAvail := uint64(4000000 * 1024)
	if mem.AvailableBytes != expectedAvail {
		t.Errorf("Expected available %d, got %d", expectedAvail, mem.AvailableBytes)
	}

	expectedUsed := expectedTotal - expectedAvail
	if mem.UsedBytes != expectedUsed {
		t.Errorf("Expected used %d, got %d", expectedUsed, mem.UsedBytes)
	}

	if mem.UsedPercent != 50.0 {
		t.Errorf("Expected used percent 50.0, got %f", mem.UsedPercent)
	}
}

func TestParseMeminfoFallback(t *testing.T) {
	// Sample without MemAvailable (older kernels)
	sample := `MemTotal:        4000000 kB
MemFree:         1000000 kB
Buffers:          500000 kB
Cached:           500000 kB
`
	mem, err := ParseMeminfo(sample)
	if err != nil {
		t.Fatalf("Unexpected error parsing meminfo: %v", err)
	}

	// Available should be Free + Buffers + Cached = 2000000 kB
	expectedAvail := uint64(2000000 * 1024)
	if mem.AvailableBytes != expectedAvail {
		t.Errorf("Expected fallback available %d, got %d", expectedAvail, mem.AvailableBytes)
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes uint64
		want  string
	}{
		{500, "500 B"},
		{1024, "1.00 KB"},
		{1024 * 1024 * 100, "100.0 MB"},
		{1024 * 1024 * 1024 * 4, "4.00 GB"},
	}

	for _, tc := range tests {
		got := FormatBytes(tc.bytes)
		if got != tc.want {
			t.Errorf("FormatBytes(%d) = %q, want %q", tc.bytes, got, tc.want)
		}
	}
}
