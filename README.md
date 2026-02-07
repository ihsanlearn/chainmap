<h1 align="center">
  Chainmap
  <br>
</h1>

<h4 align="center">A modular, high-performance Nmap workflow orchestrator designed for automation pipelines.</h4>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#installation">Installation</a> •
  <a href="#usage">Usage</a> •
  <a href="#scan-modes">Scan Modes</a> •
  <a href="#workflow-integration">Workflow Integration</a>
</p>

---

**Chainmap** is built on the philosophy of "Do one thing and do it well, then pipe it." It serves as the bridge between target discovery tools (like `subfinder`, `naabu`, `httpx`) and the deep scanning capabilities of Nmap.

Instead of sequentially scanning targets one by one, Chainmap intelligently groups inputs and utilizes a worker pool to execute multiple Nmap instances in parallel, significantly reducing scan times for large target lists.

## Features

- **Smart Target Parsing**: Automatically handles `IP`, `Domain`, and `IP:PORT` formats.
- **Intelligent Grouping**: Consolidates multiple ports for the same IP into a single Nmap command (e.g., `1.1.1.1:80` + `1.1.1.1:443` -> `nmap 1.1.1.1 -p 80,443`).
- **Concurrency Control**: Configurable worker pool to manage load and network stability.
- **Optimized Scan Modes**: Built-in presets for `Fast` triage and `Deep` inspection.
- **Unified Reporting**: Merges individual XML results into a single comprehensive report (XML & HTML).
- **Resilience**: Built-in timeout management to prevent stalled scans.

## Installation

Ensure you have **Go 1.21+** installed.

```bash
go install github.com/ihsanlearn/chainmap/cmd/chainmap@latest
```

### Dependencies

Chainmap requires `nmap` to be installed and available in your system's PATH.

- **Optional**: `xsltproc` is required for generating HTML reports.

## Usage

```bash
chainmap -h
```

### Basic Scans

**Scan a single target:**

```bash
sudo chainmap -t 192.168.1.10
```

**Scan a list of targets (File):**

```bash
sudo chainmap -l targets.txt
```

**Scan via Stdin (Pipeline):**

```bash
cat targets.txt | sudo chainmap
```

### Scan Modes

**Fast Mode (`-fast`)**
Ideal for quick triage. Uses aggressive timing, limited port range (Top 1000), and skips host discovery.
_Requires root privileges for SYN scan._

```bash
sudo chainmap -l targets.txt -fast
```

**Deep Mode (`-deep`)**
Comprehensive auditing. Enables version detection (`-sV`), default scripts (`-sC`), and OS detection.
_Requires root privileges for SYN scan._

```bash
sudo chainmap -l targets.txt -deep
```

### Advanced Configuration

| Flag              | Description                                | Default       |
| :---------------- | :----------------------------------------- | :------------ |
| `-c, -threads`    | Number of concurrent Nmap instances        | `5`           |
| `-T, -timeout`    | Timeout per scan in minutes                | `10`          |
| `-o, -output`     | Output file path (supports .xml and .html) | `results.xml` |
| `-n, -nmap-flags` | Custom Nmap flags (overrides modes)        | _Dynamic_     |
| `-s, -silent`     | Suppress standard output logs              | `false`       |

## Workflow Integration

Chainmap shines when integrated into bug bounty or pentest workflows.

**Example: Discovery to Scan Pipeline**

```bash
subfinder -d example.com | httpx -ports 80,443,8080 -ip | awk '{print $n}' | sudo chainmap -fast -o triage.html
```

## License

This project is licensed under the MIT License.
