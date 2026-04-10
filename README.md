# NetMap

NetMap is a professional visual intelligence and mapping CLI tool built in Go. It turns discovered network domains, subdomains, and endpoints into clean, highly readable tree maps directly in your terminal.

## Quick Install (macOS & Linux)

Paste this single command in your terminal to instantly install the latest version:

```bash
curl -fsSL https://raw.githubusercontent.com/lucasenlucas/NetMap/main/install.sh | sh
```


## Windows
Go to [Releases](https://github.com/lucasenlucas/NetMap/releases) and download the `.exe` directly.

## Usage

Map an entire target natively in your terminal:
```bash
netmap -d example.com
```

Focus on specific high-interest groups (`auth`, `admin`, `api`, `all`):
```bash
netmap -d example.com -f auth
```

Output as structurally clean JSON for web pipelines:
```bash
netmap -d example.com -o json
```

## Features

- **Blazing Fast**: Compiled directly to machine code globally across multiple OS architectures with exactly zero external dependencies.
- **Categorical Intelligence**: Groups standard pathways accurately using pattern detection.
- **Structured Rendering**: Outputs elegant hierarchical visualizations right at the command line.

> NetMap is part of the NET Ecosystem.
