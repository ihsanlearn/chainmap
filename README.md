# Chainmap

**Chainmap** is a modular Nmap workflow tool designed to "Do one thing and do it well, then pipe it." usage philosophy. It is optimized for chaining with tools like `naabu` or `httpx`.

## Core Philosophy

1.  **Input Module**: Accepts inputs from stdin or file (supports `IP:PORT` format).
2.  **Parser Module**: Groups inputs by unique IP to minimize Nmap execution overhead (e.g., `1.1.1.1:80` and `1.1.1.1:443` becomes `1.1.1.1 -p 80,443`).
3.  **Concurrency Engine**: Uses a Worker Pool to manage load and stability.
4.  **Execution Module**: Runs optimized Nmap scans.

## Features

- **Smart Grouping**: Automatically groups ports for the same IP.
- **Worker Pool**: Control concurrency with `-threads`.
- **Safety**: Built-in timeout management (`-timeout`).
- **Flexible Input**: Handles `IP`, `Domain`, and `IP:PORT` formats seamlessly.

## Installation

```bash
go install github.com/ihsanlearn/chainmap/cmd/chainmap@latest
```

## Usage

### Basic Usage

**Scan a list of targets (IP:PORT or IP):**

```bash
chainmap -l targets.txt
```

**Pipe from other tools:**

```bash
cat targets.txt | chainmap -threads 10
```

**Custom Nmap Flags**

```bash
chainmap -l targets.txt -n "-sC -sV"
```

### Options

```
Flags:
INPUT:
   -l, -list string    Input file containing list of IPs/IP:PORT
   -t, -target string  Single target IP

CONFIGURATION:
   -c, -threads int        Number of concurrent threads (default 5)
   -T, -timeout int        Timeout in minutes (default 10)
   -n, -nmap-flags string  Nmap flags to use (default "-sV -T3 --version-intensity 5 -Pn -n")

OPTIMIZATION:
   -s, -silent           Silent mode
   -V, -version          Display application version
```

## Requirements

- **Nmap**: Ensure `nmap` is installed and available in your system's PATH.
