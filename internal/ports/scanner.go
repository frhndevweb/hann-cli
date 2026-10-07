package ports

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// PortEntry represents an active network port on the host.
type PortEntry struct {
	Port     int
	Protocol string
	Address  string
	State    string
	Process  string
	Inode    string
}

// ScannerOptions configures port scanning behavior.
type ScannerOptions struct {
	ListeningOnly bool
}

// ScanPorts scans listening/open ports using /proc/net first, with graceful tool fallbacks.
func ScanPorts(opts ScannerOptions) ([]PortEntry, error) {
	// 1. Primary approach: Direct Linux /proc/net reading
	entries, procErr := scanProcNet(opts)
	if procErr == nil && len(entries) > 0 {
		return sortAndDeduplicate(entries), nil
	}

	// 2. Fallback approach: ss
	if ssEntries, err := scanWithSS(opts); err == nil && len(ssEntries) > 0 {
		return sortAndDeduplicate(ssEntries), nil
	}

	// 3. Fallback approach: netstat
	if netstatEntries, err := scanWithNetstat(opts); err == nil && len(netstatEntries) > 0 {
		return sortAndDeduplicate(netstatEntries), nil
	}

	// 4. Fallback approach: lsof
	if lsofEntries, err := scanWithLsof(opts); err == nil && len(lsofEntries) > 0 {
		return sortAndDeduplicate(lsofEntries), nil
	}

	if procErr != nil {
		return nil, fmt.Errorf("unable to inspect sockets: %w", procErr)
	}

	return entries, nil
}

// scanProcNet parses /proc/net/{tcp,tcp6,udp,udp6}.
func scanProcNet(opts ScannerOptions) ([]PortEntry, error) {
	// Build map of socket inode to process name
	inodeMap := buildInodeProcessMap()

	var allEntries []PortEntry
	files := []struct {
		path     string
		protocol string
	}{
		{"/proc/net/tcp", "TCP"},
		{"/proc/net/tcp6", "TCP6"},
		{"/proc/net/udp", "UDP"},
		{"/proc/net/udp6", "UDP6"},
	}

	var lastErr error
	var successfullyReadFiles int

	for _, target := range files {
		content, err := os.ReadFile(target.path)
		if err != nil {
			lastErr = err
			continue
		}
		successfullyReadFiles++

		parsed := ParseProcNetContent(string(content), target.protocol, opts.ListeningOnly, inodeMap)
		allEntries = append(allEntries, parsed...)
	}

	if successfullyReadFiles == 0 && lastErr != nil {
		return nil, lastErr
	}

	return allEntries, nil
}

// ParseProcNetContent parses tabular /proc/net file data into PortEntry structs.
func ParseProcNetContent(content, protocol string, listeningOnly bool, inodeMap map[string]string) []PortEntry {
	var entries []PortEntry
	scanner := bufio.NewScanner(strings.NewReader(content))

	// Skip header line
	isHeader := true
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if isHeader {
			isHeader = false
			continue
		}
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}

		// fields[1] = local_address (hex_ip:hex_port)
		// fields[3] = st (hex socket state: 0A = TCP_LISTEN)
		// fields[9] = inode
		localAddr := fields[1]
		stateHex := fields[3]
		inode := fields[9]

		// 0A is TCP_LISTEN in hex
		isListening := stateHex == "0A" || (strings.HasPrefix(protocol, "UDP") && (stateHex == "07" || stateHex == "01"))
		if listeningOnly && !isListening {
			continue
		}

		ip, port, err := DecodeHexAddress(localAddr)
		if err != nil {
			continue
		}

		procName := "-"
		if inodeMap != nil && inode != "" && inode != "0" {
			if p, ok := inodeMap[inode]; ok && p != "" {
				procName = p
			}
		}

		stateName := "LISTEN"
		if !isListening {
			stateName = stateHex
		}

		entries = append(entries, PortEntry{
			Port:     port,
			Protocol: protocol,
			Address:  ip,
			State:    stateName,
			Process:  procName,
			Inode:    inode,
		})
	}

	return entries
}

// DecodeHexAddress parses IP:Port from Linux /proc/net format.
func DecodeHexAddress(hexAddr string) (string, int, error) {
	parts := strings.Split(hexAddr, ":")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid address format: %s", hexAddr)
	}

	ipHex := parts[0]
	portHex := parts[1]

	portUint, err := strconv.ParseUint(portHex, 16, 16)
	if err != nil {
		return "", 0, fmt.Errorf("invalid port hex: %w", err)
	}
	port := int(portUint)

	// IPv4 (8 hex characters)
	if len(ipHex) == 8 {
		b, err := hex.DecodeString(ipHex)
		if err != nil {
			return "", 0, err
		}
		// Little-endian order in /proc/net/tcp
		ip := fmt.Sprintf("%d.%d.%d.%d", b[3], b[2], b[1], b[0])
		return ip, port, nil
	}

	// IPv6 (32 hex characters)
	if len(ipHex) == 32 {
		ipBytes, err := hex.DecodeString(ipHex)
		if err != nil {
			return "", 0, err
		}
		// Linux writes 4 32-bit words in host byte order
		var fixedIP [16]byte
		for i := 0; i < 4; i++ {
			fixedIP[i*4+0] = ipBytes[i*4+3]
			fixedIP[i*4+1] = ipBytes[i*4+2]
			fixedIP[i*4+2] = ipBytes[i*4+1]
			fixedIP[i*4+3] = ipBytes[i*4+0]
		}
		ipNet := net.IP(fixedIP[:])
		if ipNet.IsUnspecified() {
			return "::", port, nil
		}
		return ipNet.String(), port, nil
	}

	return "", 0, fmt.Errorf("unsupported IP hex length: %d", len(ipHex))
}

// buildInodeProcessMap maps socket inodes to process names from /proc/[pid]/fd.
func buildInodeProcessMap() map[string]string {
	m := make(map[string]string)

	procEntries, err := os.ReadDir("/proc")
	if err != nil {
		return m
	}

	for _, entry := range procEntries {
		if !entry.IsDir() {
			continue
		}
		pid := entry.Name()
		// Only check numeric PID dirs
		if pid[0] < '0' || pid[0] > '9' {
			continue
		}

		fdDir := filepath.Join("/proc", pid, "fd")
		fds, err := os.ReadDir(fdDir)
		if err != nil {
			continue
		}

		// Read process command name once per PID
		var procName string
		commPath := filepath.Join("/proc", pid, "comm")
		if commData, err := os.ReadFile(commPath); err == nil {
			procName = strings.TrimSpace(string(commData))
		}

		for _, fd := range fds {
			link, err := os.Readlink(filepath.Join(fdDir, fd.Name()))
			if err != nil {
				continue
			}

			if strings.HasPrefix(link, "socket:[") && strings.HasSuffix(link, "]") {
				inode := link[8 : len(link)-1]
				if procName != "" {
					m[inode] = procName
				} else {
					m[inode] = pid
				}
			}
		}
	}

	return m
}

// scanWithSS uses `ss -tulnp` or `ss -tuln` as a fallback.
func scanWithSS(opts ScannerOptions) ([]PortEntry, error) {
	args := []string{"-tuln"}
	if cmdPath, err := exec.LookPath("ss"); err != nil {
		return nil, err
	} else if cmdPath == "" {
		return nil, fmt.Errorf("ss not found")
	}

	out, err := exec.Command("ss", append(args, "-p")...).Output()
	if err != nil {
		out, err = exec.Command("ss", args...).Output()
		if err != nil {
			return nil, err
		}
	}

	return parseSSOutput(string(out), opts.ListeningOnly), nil
}

// parseSSOutput parses output lines from the ss tool.
func parseSSOutput(output string, listeningOnly bool) []PortEntry {
	var entries []PortEntry
	scanner := bufio.NewScanner(strings.NewReader(output))
	procRe := regexp.MustCompile(`users:\(\("([^"]+)"`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "Netid") || line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		netid := strings.ToUpper(fields[0])
		state := fields[1]
		local := fields[4]

		if listeningOnly && state != "LISTEN" && !strings.Contains(netid, "UDP") {
			continue
		}

		host, portStr, err := net.SplitHostPort(local)
		if err != nil {
			idx := strings.LastIndex(local, ":")
			if idx > 0 {
				host = local[:idx]
				portStr = local[idx+1:]
			} else {
				continue
			}
		}

		port, err := strconv.Atoi(portStr)
		if err != nil {
			continue
		}

		process := "-"
		match := procRe.FindStringSubmatch(line)
		if len(match) > 1 {
			process = match[1]
		}

		entries = append(entries, PortEntry{
			Port:     port,
			Protocol: netid,
			Address:  host,
			State:    state,
			Process:  process,
		})
	}

	return entries
}

// scanWithNetstat parses output from `netstat -tuln`.
func scanWithNetstat(opts ScannerOptions) ([]PortEntry, error) {
	if _, err := exec.LookPath("netstat"); err != nil {
		return nil, err
	}

	out, err := exec.Command("netstat", "-tuln").Output()
	if err != nil {
		return nil, err
	}

	return parseNetstatOutput(string(out), opts.ListeningOnly), nil
}

func parseNetstatOutput(output string, listeningOnly bool) []PortEntry {
	var entries []PortEntry
	scanner := bufio.NewScanner(strings.NewReader(output))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "Active") || strings.HasPrefix(line, "Proto") || line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		proto := strings.ToUpper(fields[0])
		local := fields[3]
		state := ""
		if len(fields) >= 6 {
			state = fields[5]
		}

		if listeningOnly && state != "LISTEN" && !strings.Contains(proto, "UDP") {
			continue
		}

		idx := strings.LastIndex(local, ":")
		if idx == -1 {
			continue
		}
		host := local[:idx]
		portStr := local[idx+1:]
		port, err := strconv.Atoi(portStr)
		if err != nil {
			continue
		}

		entries = append(entries, PortEntry{
			Port:     port,
			Protocol: proto,
			Address:  host,
			State:    state,
			Process:  "-",
		})
	}

	return entries
}

// scanWithLsof parses output from `lsof -i -P -n`.
func scanWithLsof(opts ScannerOptions) ([]PortEntry, error) {
	if _, err := exec.LookPath("lsof"); err != nil {
		return nil, err
	}

	out, err := exec.Command("lsof", "-iTCP", "-sTCP:LISTEN", "-P", "-n").Output()
	if err != nil {
		return nil, err
	}

	var entries []PortEntry
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "COMMAND") || line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}
		command := fields[0]
		addrField := fields[8] // e.g. *:8080 or 127.0.0.1:3000

		idx := strings.LastIndex(addrField, ":")
		if idx == -1 {
			continue
		}
		host := addrField[:idx]
		portStr := addrField[idx+1:]
		port, err := strconv.Atoi(portStr)
		if err != nil {
			continue
		}

		entries = append(entries, PortEntry{
			Port:     port,
			Protocol: "TCP",
			Address:  host,
			State:    "LISTEN",
			Process:  command,
		})
	}

	return entries, nil
}

// sortAndDeduplicate sorts entries by port number and removes duplicates.
func sortAndDeduplicate(entries []PortEntry) []PortEntry {
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Port == entries[j].Port {
			return entries[i].Protocol < entries[j].Protocol
		}
		return entries[i].Port < entries[j].Port
	})

	var result []PortEntry
	seen := make(map[string]bool)
	for _, e := range entries {
		key := fmt.Sprintf("%d:%s:%s", e.Port, e.Protocol, e.Address)
		if !seen[key] {
			seen[key] = true
			result = append(result, e)
		}
	}

	return result
}
