package system

import (
	"testing"
)

func TestParseCPUInfo(t *testing.T) {
	sampleX86 := `processor	: 0
vendor_id	: GenuineIntel
cpu family	: 6
model name	: Intel(R) Core(TM) i7-10700K CPU @ 3.80GHz

processor	: 1
vendor_id	: GenuineIntel
cpu family	: 6
model name	: Intel(R) Core(TM) i7-10700K CPU @ 3.80GHz
`
	parsed := ParseCPUInfo(sampleX86)
	if parsed.Cores != 2 {
		t.Errorf("Expected 2 cores, got %d", parsed.Cores)
	}
	if parsed.Model != "Intel(R) Core(TM) i7-10700K CPU @ 3.80GHz" {
		t.Errorf("Unexpected model: %s", parsed.Model)
	}

	sampleARM := `processor	: 0
BogoMIPS	: 52.00
Features	: fp asimd evtstrm aes pmull sha1 sha2 crc32 atomics
CPU implementer	: 0x41
CPU architecture: 8
CPU variant	: 0x1
CPU part	: 0xd05
CPU revision	: 0

Hardware	: MT6877
`
	parsedARM := ParseCPUInfo(sampleARM)
	if parsedARM.Cores != 1 {
		t.Errorf("Expected 1 core, got %d", parsedARM.Cores)
	}
	if parsedARM.Model != "MT6877" {
		t.Errorf("Expected MT6877 hardware model, got %s", parsedARM.Model)
	}
}

func TestGetCPUInfoSmoke(t *testing.T) {
	info := GetCPUInfo()
	if info == nil {
		t.Fatal("Expected non-nil CPUInfo")
	}
	if info.Cores <= 0 {
		t.Errorf("Expected cores > 0, got %d", info.Cores)
	}
}
