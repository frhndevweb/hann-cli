package environment

import (
	"bufio"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

// EnvironmentInfo holds detected system and runtime environment attributes.
type EnvironmentInfo struct {
	OS           string
	Architecture string
	Kernel       string
	Shell        string
	Hostname     string
	User         string
	WorkingDir   string
	IsTermux     bool
	IsPRoot      bool
}

// DetectEnvironment gathers all host and runtime environment properties.
func DetectEnvironment() *EnvironmentInfo {
	return &EnvironmentInfo{
		OS:           DetectOS(),
		Architecture: DetectArchitecture(),
		Kernel:       DetectKernel(),
		Shell:        DetectShell(),
		Hostname:     DetectHostname(),
		User:         DetectUser(),
		WorkingDir:   DetectWorkingDir(),
		IsTermux:     DetectTermux(),
		IsPRoot:      DetectPRoot(),
	}
}

// DetectOS attempts to read linux distribution release files.
func DetectOS() string {
	candidates := []string{"/etc/os-release", "/usr/lib/os-release", "/etc/lsb-release"}
	for _, path := range candidates {
		if content, err := os.ReadFile(path); err == nil {
			name := ParseOSRelease(string(content))
			if name != "" {
				return name
			}
		}
	}

	if issue, err := os.ReadFile("/etc/issue"); err == nil {
		firstLine := strings.Split(string(issue), "\n")[0]
		firstLine = strings.TrimSpace(strings.ReplaceAll(firstLine, "\\n", ""))
		firstLine = strings.TrimSpace(strings.ReplaceAll(firstLine, "\\l", ""))
		if firstLine != "" {
			return firstLine
		}
	}

	return runtime.GOOS
}

// ParseOSRelease parses KEY=VALUE os-release content and extracts distribution name.
func ParseOSRelease(content string) string {
	var prettyName, name, version string
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"'`)

		switch key {
		case "PRETTY_NAME":
			prettyName = val
		case "NAME":
			name = val
		case "VERSION":
			version = val
		}
	}

	if prettyName != "" {
		return prettyName
	}
	if name != "" && version != "" {
		return name + " " + version
	}
	if name != "" {
		return name
	}
	return ""
}

// DetectArchitecture returns the system machine architecture.
func DetectArchitecture() string {
	var uts syscall.Utsname
	if err := syscall.Uname(&uts); err == nil {
		var buf [65]byte
		for i, b := range uts.Machine {
			if b == 0 {
				break
			}
			buf[i] = byte(b)
		}
		mach := strings.TrimRight(string(buf[:]), "\x00")
		if mach != "" {
			return mach
		}
	}
	return runtime.GOARCH
}

// DetectKernel returns the operating system kernel release string.
func DetectKernel() string {
	if data, err := os.ReadFile("/proc/sys/kernel/osrelease"); err == nil {
		k := strings.TrimSpace(string(data))
		if k != "" {
			return "Linux " + k
		}
	}

	var uts syscall.Utsname
	if err := syscall.Uname(&uts); err == nil {
		var sysname, release [65]byte
		for i, b := range uts.Sysname {
			if b == 0 {
				break
			}
			sysname[i] = byte(b)
		}
		for i, b := range uts.Release {
			if b == 0 {
				break
			}
			release[i] = byte(b)
		}
		s := strings.TrimRight(string(sysname[:]), "\x00")
		r := strings.TrimRight(string(release[:]), "\x00")
		if s != "" || r != "" {
			return strings.TrimSpace(s + " " + r)
		}
	}

	return runtime.GOOS
}

// DetectShell inspects environment and process settings to find the current shell.
func DetectShell() string {
	if sh := os.Getenv("SHELL"); sh != "" {
		return filepath.Base(sh)
	}
	// Fallback to checking parent shell or common default
	return "bash"
}

// DetectHostname retrieves the system hostname.
func DetectHostname() string {
	if host, err := os.Hostname(); err == nil && host != "" {
		return host
	}
	if data, err := os.ReadFile("/proc/sys/kernel/hostname"); err == nil {
		h := strings.TrimSpace(string(data))
		if h != "" {
			return h
		}
	}
	return "localhost"
}

// DetectUser retrieves current username.
func DetectUser() string {
	if u, err := user.Current(); err == nil && u.Username != "" {
		return u.Username
	}
	if u := os.Getenv("USER"); u != "" {
		return u
	}
	if u := os.Getenv("LOGNAME"); u != "" {
		return u
	}
	return "unknown"
}

// DetectWorkingDir returns the current working directory safely.
func DetectWorkingDir() string {
	if cwd, err := os.Getwd(); err == nil {
		return cwd
	}
	if pwd := os.Getenv("PWD"); pwd != "" {
		return pwd
	}
	return "."
}

// DetectTermux determines if running within or under a Termux environment.
func DetectTermux() bool {
	if os.Getenv("TERMUX_VERSION") != "" || os.Getenv("TERMUX_MAIN_PACKAGE_FORMAT") != "" {
		return true
	}
	if strings.Contains(os.Getenv("PREFIX"), "com.termux") ||
		strings.Contains(os.Getenv("PATH"), "com.termux") ||
		strings.Contains(os.Getenv("LD_PRELOAD"), "termux") {
		return true
	}
	// Check standard Termux data paths
	if _, err := os.Stat("/data/data/com.termux"); err == nil {
		return true
	}
	return false
}

// DetectPRoot determines if running inside a PRoot or proot-distro container.
func DetectPRoot() bool {
	// Standard proot-distro environment variables
	if strings.EqualFold(os.Getenv("container"), "proot-distro") ||
		os.Getenv("PD_CONTAINER") != "" ||
		os.Getenv("PD_IMAGE") != "" ||
		os.Getenv("PROOT_L2S_DIR") != "" ||
		os.Getenv("PROOT_TMP_DIR") != "" ||
		os.Getenv("PROOT_FORCE_DISCARD") != "" {
		return true
	}

	// Kernel version signature often injected by proot-distro
	kernel := DetectKernel()
	if strings.Contains(strings.ToLower(kernel), "proot") {
		return true
	}

	// Check /proc/version
	if data, err := os.ReadFile("/proc/version"); err == nil {
		if strings.Contains(strings.ToLower(string(data)), "proot") {
			return true
		}
	}

	// Check /proc/1/cmdline
	if data, err := os.ReadFile("/proc/1/cmdline"); err == nil {
		if strings.Contains(strings.ToLower(string(data)), "proot") {
			return true
		}
	}

	return false
}
