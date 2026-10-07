# hann-cli

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20ARM64%20%7C%20AMD64-blue)](https://github.com/frhndevweb/hann-cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`hann-cli` is a fast, lightweight developer and system utility CLI created by Farhan ("hann"). Built specifically for Linux environments, it features first-class detection and optimization for Ubuntu running inside **Termux + proot-distro** on Android ARM64 devices, as well as standard Linux servers and workstations (x86_64 and aarch64).

Local-first, dependency-light, zero telemetry, and zero hardcoded system metrics.

---

## Features

- **System Diagnostics (`hann doctor`)**: Checks operating system distribution, architecture, kernel, shell, Termux/PRoot runtime indicators, and verifies your developer toolchain using real executable discovery (`exec.LookPath`).
- **Real-Time System Metrics (`hann system`)**: Reads directly from `/proc/cpuinfo`, `/proc/meminfo`, and `/proc/uptime` without external APIs to display CPU cores, live memory consumption, uptime, and container metadata.
- **Port Inspection (`hann ports`)**: Decodes Linux `/proc/net/tcp`, `/proc/net/tcp6`, `/proc/net/udp` tables and maps socket inodes to process names, with graceful fallbacks (`ss`, `netstat`, `lsof`) when permissions are restricted.
- **Project Discovery (`hann projects`)**: Scans workspace directories up to depth limits to discover projects across Go, Node.js, Rust, Python, PHP, Flutter, Java, Gradle, CMake, and Git.
- **Termux & PRoot Detection**: Accurately recognizes when running inside Android Termux or a PRoot chroot/distro environment.
- **Clean Terminal UI**: Minimalist, human-readable output with adaptive colors (respects `NO_COLOR` and `TERM=dumb`) and automatic Unicode/ASCII fallback glyphs.

---

## Installation

### Prerequisites

- Go 1.22 or newer
- Linux (x86_64, aarch64 / arm64)

### Install via Go

```bash
go install github.com/frhndevweb/hann-cli@latest
```

### Build from Source

```bash
git clone https://github.com/frhndevweb/hann-cli.git
cd hann-cli
make build
make install
```

The compiled binary will be placed at `bin/hann` and copied to `~/.local/bin/hann`.

---

## Usage & Examples

### Overview

```bash
hann
```

Example output:
```text
hann-cli

Version:     0.1.0
Description: Developer & system utility for Linux (Termux & PRoot)
Platform:    linux/arm64

Available Commands
  doctor       Inspect Linux environment and verify developer toolchain
  system       Display real-time system, CPU, memory, and uptime metrics
  ports        List active listening ports and associated processes
  projects     Discover development projects in your workspace
  version      Show version, commit, and build environment info

Flags
  -h, --help   Display help for hann or any command
  -v, --version Display version information

Run hann [command] --help for detailed usage.
```

---

### Environment & Tool Diagnostics (`hann doctor`)

Verify system indicators and installed CLI developer packages:

```bash
hann doctor
```

Example output:
```text
Hann Doctor

Environment
✓ OS             Ubuntu 24.04 LTS
✓ Architecture   arm64
✓ Kernel         Linux 6.17.0-PRoot-Distro
✓ Shell          bash
✓ PRoot          detected
✓ Termux         detected

Tools
✓ git
✓ go
✓ node
✓ npm
✗ pnpm           not installed
✓ python
✓ python3
✓ curl
✓ wget
✗ docker         not installed
✓ gcc
✓ make

Summary
16 checks passed
2 checks failed
```

---

### Real-Time System Information (`hann system`)

Inspect live resource statistics directly from Linux `/proc`:

```bash
hann system
```

Example output:
```text
Hann System

System
OS            Ubuntu 24.04 LTS
Kernel        Linux 6.17.0-PRoot-Distro
Architecture  aarch64
Hostname      localhost
User          hann
Shell         bash

Resources
CPU           8 cores
Memory        3.72 GB / 7.81 GB (48%)
Uptime        2h 31m

Environment
PRoot         yes
Termux        yes

Directory
/home/hann/projects/hann-cli
```

---

### Port Inspection (`hann ports`)

List active listening ports and resolved process names:

```bash
hann ports
```

Example output:
```text
PORT    PROTOCOL    ADDRESS       PROCESS
3000    TCP         0.0.0.0       node
5173    TCP         127.0.0.1     vite
8080    TCP         0.0.0.0       go
```

To filter strictly listening sockets (enabled by default):
```bash
hann ports --listening
```

> [!NOTE]
> On unrooted Android/PRoot devices, Android SELinux may restrict direct reading of `/proc/net/tcp`. `hann-cli` will attempt fallback tools (`ss`, `netstat`, `lsof`) and reports human-friendly diagnostics without crashing.

---

### Workspace Project Discovery (`hann projects`)

Scan for development projects (defaulting to `~/projects`):

```bash
hann projects
```

Example output:
```text
Projects

~/projects

1.  api-sentinel         Go
2.  hann-cli             Go
3.  website              Node.js
4.  security-tool        Python
```

If the target directory does not exist:
```text
Projects directory not found:
~/projects

Create it with:
mkdir -p ~/projects
```

Scan a custom workspace path using `--path`:
```bash
hann projects --path ~/workspace
```

---

### Version Information (`hann version` or `hann --version`)

```bash
hann version
```

Example output:
```text
hann-cli

Version:    0.1.0
Commit:     a1b2c3d
Built:      2026-10-07T15:30:00Z
Platform:   linux/arm64
Go runtime: go1.26.0
```

---

## Development & Testing

### Running Tests

Run all unit tests across modules:

```bash
make test
```

Or using native Go tooling:

```bash
go test -v ./...
```

### Code Formatting and Linting

```bash
make fmt
make lint
```

---

## Cross-Compilation

`hann-cli` is written in standard Go with zero CGO dependencies, making cross-compilation straightforward.

### Target: Linux ARM64 (aarch64 / Termux)

```bash
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o hann .
```

### Target: Linux AMD64 (x86_64)

```bash
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o hann .
```

---

## Security & Privacy

- **100% Local-First**: All diagnostics, hardware measurements, and workspace scans are processed entirely on your local machine.
- **Zero Telemetry**: No metrics, tracking beacons, or network phone-homes are embedded.
- **Safe Scanning**: Scans are bounded to specified workspace directories and explicitly skip hidden secret files, `.env` files, and private credentials.

---

## Roadmap

- [x] v0.1: Initial release with `doctor`, `system`, `ports`, `projects`, `version`
- [ ] v0.2: Process tree inspector and kill utility for orphaned background dev servers
- [ ] v0.3: Network latency and interface diagnostic submodule (`hann ping` / `hann net`)
- [ ] v0.4: Package manager accelerator helpers for Termux apt mirrors

---

## License

MIT License. Copyright (c) 2026 Farhan (frhndevweb). See [LICENSE](LICENSE) for details.
