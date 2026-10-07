package system

import (
	"testing"
	"time"
)

func TestParseUptime(t *testing.T) {
	sample := "9104.25 18208.50\n"
	d, err := ParseUptime(sample)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedSeconds := 9104
	if int(d.Seconds()) != expectedSeconds {
		t.Errorf("Expected %d seconds, got %f", expectedSeconds, d.Seconds())
	}
}

func TestFormatUptime(t *testing.T) {
	tests := []struct {
		duration time.Duration
		want     string
	}{
		{45 * time.Second, "45s"},
		{15 * time.Minute, "15m"},
		{15*time.Minute + 30*time.Second, "15m 30s"},
		{2*time.Hour + 31*time.Minute, "2h 31m"},
		{2 * time.Hour, "2h"},
		{25*time.Hour + 10*time.Minute, "1d 1h 10m"},
		{48 * time.Hour, "2d 0h"},
	}

	for _, tc := range tests {
		got := FormatUptime(tc.duration)
		if got != tc.want {
			t.Errorf("FormatUptime(%v) = %q, want %q", tc.duration, got, tc.want)
		}
	}
}
