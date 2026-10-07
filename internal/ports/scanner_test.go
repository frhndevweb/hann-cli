package ports

import (
	"testing"
)

func TestDecodeHexAddressIPv4(t *testing.T) {
	// 127.0.0.1:80 -> 0100007F:0050
	ip, port, err := DecodeHexAddress("0100007F:0050")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if ip != "127.0.0.1" {
		t.Errorf("Expected 127.0.0.1, got %s", ip)
	}
	if port != 80 {
		t.Errorf("Expected port 80, got %d", port)
	}

	// 0.0.0.0:8080 -> 00000000:1F90
	ip2, port2, err := DecodeHexAddress("00000000:1F90")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if ip2 != "0.0.0.0" {
		t.Errorf("Expected 0.0.0.0, got %s", ip2)
	}
	if port2 != 8080 {
		t.Errorf("Expected port 8080, got %d", port2)
	}
}

func TestDecodeHexAddressIPv6(t *testing.T) {
	// All zeros IPv6 -> :: on port 22 (0x0016)
	ip, port, err := DecodeHexAddress("00000000000000000000000000000000:0016")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if ip != "::" {
		t.Errorf("Expected ::, got %s", ip)
	}
	if port != 22 {
		t.Errorf("Expected port 22, got %d", port)
	}
}

func TestParseProcNetContent(t *testing.T) {
	mockTcp := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 00000000:0BB8 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 54321 1 0000000000000000 100 0 0 10 0
   1: 0100007F:1435 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 54322 1 0000000000000000 100 0 0 10 0
   2: 0100007F:0050 0100007F:9999 01 00000000:00000000 00:00000000 00000000     0        0 54323 1 0000000000000000 100 0 0 10 0
`
	inodeMap := map[string]string{
		"54321": "node",
		"54322": "vite",
	}

	// 0x0BB8 = 3000, 0x1435 = 5173
	entries := ParseProcNetContent(mockTcp, "TCP", true, inodeMap)
	if len(entries) != 2 {
		t.Fatalf("Expected 2 listening entries, got %d", len(entries))
	}

	if entries[0].Port != 3000 || entries[0].Process != "node" || entries[0].Address != "0.0.0.0" {
		t.Errorf("Entry 0 unexpected: %+v", entries[0])
	}
	if entries[1].Port != 5173 || entries[1].Process != "vite" || entries[1].Address != "127.0.0.1" {
		t.Errorf("Entry 1 unexpected: %+v", entries[1])
	}
}

func TestParseSSOutput(t *testing.T) {
	mockSS := `Netid State  Recv-Q Send-Q Local Address:Port  Peer Address:PortProcess
tcp   LISTEN 0      128          0.0.0.0:8080       0.0.0.0:*    users:(("go",pid=99,fd=3))
`
	entries := parseSSOutput(mockSS, true)
	if len(entries) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(entries))
	}
	if entries[0].Port != 8080 || entries[0].Process != "go" {
		t.Errorf("Unexpected ss entry: %+v", entries[0])
	}
}
