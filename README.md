# hann-cli

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20ARM64%20%7C%20AMD64-blue)](https://github.com/frhndevweb/hann-cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`hann-cli` adalah developer & system utility CLI super ringan buatan Farhan ("hann"). CLI ini gw rancang khusus buat environment Linux, terutama buat lo yang sering ngoding di Ubuntu via **Termux + proot-distro** di Android ARM64, tapi tetep jalan lancar jaya di server atau PC Linux biasa (x86_64 & aarch64).

100% local-first, minim dependensi, tanpa telemetri lebay, dan gak pake hardcoded info sistem abal-abal.

---

## Fitur Utama

- **Cek Lingkungan Sistem (`hann doctor`)**: Ngedeteksi distro OS, arsitektur CPU, kernel, shell, indikator Termux/PRoot, plus ngecek ketersediaan toolchain koding lo pake executable lookup real (`exec.LookPath`).
- **Metrik Sistem Real-Time (`hann system`)**: Baca langsung dari `/proc/cpuinfo`, `/proc/meminfo`, dan `/proc/uptime` tanpa API luar buat mantau core CPU, pemakaian RAM, uptime, dan info container lo.
- **Inspeksi Port Listening (`hann ports`)**: Decode tabel `/proc/net/tcp`, `/proc/net/tcp6`, `/proc/net/udp` bawaan Linux dan mapping inode socket ke nama proses. Lengkap dengan fallback cerdas (`ss`, `netstat`, `lsof`) kalau permission di environment lo dibatesin.
- **Discovery Project (`hann projects`)**: Nyari folder project kodingan di workspace lo (Go, Node.js, Rust, Python, PHP, Flutter, Java, Gradle, CMake, Git) secara teratur dengan batasan kedalaman folder yang aman.
- **Deteksi Akurat Termux & PRoot**: Gak perlu nebak-nebak, CLI ini langsung tau kalo lo lagi jalan di dalem Termux atau PRoot distro.
- **Terminal UI Rapi & Bersih**: Tampilan modern, adaptif warna (ngikut standar `NO_COLOR` dan `TERM=dumb`), plus fallback simbol Unicode/ASCII otomatis (`✓`/`[+]`, `✗`/`[-]`, `⚠`/`[!]`).

---

## Cara Install

### Persyaratan

- Go 1.22 ke atas udah kepasang di mesin lo
- Linux (x86_64, aarch64 / arm64)

### Install via Go Install

Kalo lo udah punya Go di sistem:
```bash
go install github.com/frhndevweb/hann-cli@latest
```

### Build Langsung dari Source

Lo tinggal clone repo ini terus build pake Makefile:

```bash
git clone https://github.com/frhndevweb/hann-cli.git
cd hann-cli
make build
make install
```

Binary bakal ke-generate di `bin/hann` dan langsung disalin ke `~/.local/bin/hann`.

---

## Cara Pake & Contoh Output

### 1. Menu Utama

Kalo lo jalanin `hann` tanpa argumen:

```bash
hann
```

Contoh tampilan yang bakal lo dapet:
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

### 2. Cek Kesehatan Sistem & Toolchain (`hann doctor`)

Jalanin perintah ini buat mastiin environment dan tool koding lo udah ready:

```bash
hann doctor
```

Contoh output di terminal lo:
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

### 3. Pantau Resource Sistem (`hann system`)

Buat ngeliat spesifikasi mesin, penggunaan memori real, dan uptime:

```bash
hann system
```

Contoh output:
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

### 4. Cek Port yang Lagi Buka (`hann ports`)

Kalo lo lagi jalanin server web atau API dan pengen tau port mana aja yang lagi listening:

```bash
hann ports
```

Contoh output:
```text
PORT    PROTOCOL    ADDRESS       PROCESS
3000    TCP         0.0.0.0       node
5173    TCP         127.0.0.1     vite
8080    TCP         0.0.0.0       go
```

Kalo lo mau filter khusus port yang statusnya listening:
```bash
hann ports --listening
```

> [!NOTE]
> Kalo lo jalanin di Android/PRoot tanpa akses root, Android SELinux sering ngeblok akses langsung ke `/proc/net/tcp`. Jangan panik, `hann-cli` bakal otomatis nyoba fallback (`ss`, `netstat`, `lsof`) dan ngasih panduan tanpa bikin terminal lo crash.

---

### 5. Temukan Project Kodingan (`hann projects`)

CLI ini bisa nge-scan folder workspace lo buat nyari project apa aja yang ada di sana:

```bash
hann projects
```

Secara default, perintah di atas bakal ngecek folder `~/projects`. Contoh output:
```text
Projects

~/projects

1.  api-sentinel         Go
2.  hann-cli             Go
3.  website              Node.js
4.  security-tool        Python
```

Kalo folder `~/projects` belum lo buat, CLI bakal ngasih panduan:
```text
Projects directory not found:
~/projects

Create it with:
mkdir -p ~/projects
```

Kalo lo mau nge-scan folder workspace lain, tinggal pake opsi `--path`:
```bash
hann projects --path ~/workspace
hann projects --path .
```

---

### 6. Cek Versi Lengkap (`hann version`)

```bash
hann version
# atau
hann --version
```

Contoh output:
```text
hann-cli

Version:    0.1.0
Commit:     80ad4eb
Built:      2026-10-07T15:39:27Z
Platform:   linux/arm64
Go runtime: go1.26.0
```

---

## Development & Testing

Buat lo yang pengen ngembangin atau ngetes kodenya langsung:

### Jalanin Unit Test

```bash
make test
# atau manual:
go test -v ./...
```

### Format & Linting Kode

```bash
make fmt
make lint
```

---

## Kompilasi Silang (Cross-Compilation)

Project ini murni pake Go standard library tanpa CGO, jadi lo gampang banget build buat arsitektur lain.

### Buat Linux ARM64 (aarch64 / Termux Android)

```bash
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o hann .
```

### Buat Linux AMD64 (x86_64 PC/Server)

```bash
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o hann .
```

---

## Keamanan & Privasi

- **Murni Local-First**: Semua pengecekan, pembacaan resource, dan scanning project dieksekusi 100% di mesin lo.
- **Nol Telemetri**: Gak ada data yang dikirim ke server mana pun. Gak ada tracking tersembunyi.
- **Scan Aman**: Scanner dibatesin hanya ke direktori yang lo tentuin dan otomatis ngeskip file rahasia kaya `.env`, file credentials, maupun folder sensitif lainnya.

---

## Roadmap Pengembangan

- [x] v0.1: Rilis awal (`doctor`, `system`, `ports`, `projects`, `version`)
- [ ] v0.2: Fitur killer proses server dev yang nyangkut di background Termux (`hann kill` / `hann ps`)
- [ ] v0.3: Pengujian latensi jaringan dan benchmark mirror repositori apt Termux/Ubuntu (`hann ping` / `hann mirrors`)
- [ ] v0.4: Helper optimasi storage dan cache di lingkungan PRoot

---

## Lisensi

Dirilis di bawah lisensi [MIT](LICENSE). Hak Cipta (c) 2026 Farhan (frhndevweb). Bebas lo pake, modifikasi, dan kembangin sesuka hati!
