<p align="center">
  <img src="https://github.com/lucasenlucas/lucas_cdn/blob/main/Scherm%C2%ADafbeelding%202026-04-10%20om%2021.09.22.png?raw=true" alt="NetMap Banner"/>
</p>

<p align="center">
  <strong>NetMap</strong> — Visual Network Mapping CLI
</p>

<p align="center">
  Structure • Insight • Clarity — see how systems are built
</p>

<p align="center">
  <strong>Understand the system. Not just the endpoints.</strong>
</p>

---

## What is NetMap?

NetMap is a **visual mapping CLI** that transforms domains, subdomains, and endpoints into a clear, structured overview.

No endless lists. No guessing.

Just:
```
netmap -d example.com
```
→ and instantly see how everything connects.
Built as part of the **NET Ecosystem** alongside NetScope and NetForce.

## Authorized Use Only
NetMap is designed strictly for:
* Systems you own
* Systems you have explicit permission to analyze
This tool performs network requests and structure mapping. Use responsibly.

## Quick Install
### macOS & Linux
```
curl -fsSL https://raw.githubusercontent.com/lucasenlucas/NetMap/main/install.sh | sh
```
Run instantly:
```
netmap --help
```
### Go install
```
go install github.com/lucasenlucas/netmap/cmd/netmap@latest
```

### Manual build
```
git clone https://github.com/lucasenlucas/NetMap.git
cd NetMap
go mod tidy
go build -o netmap ./cmd/netmap
./netmap -d example.com
```

## CLI Philosophy
Same system as NetScope & NetForce → no learning curve.
```
netmap -d <domain> [options]
```

## Core Flags
| Flag            | Description                            |
| --------------- | -------------------------------------- |
| `-d, --domain`  | Target domain                          |
| `-p, --pack`    | Discovery Pack (`standard`, `dns-extended`, `web-deep`, `api-focused`, `admin-stealth`) |
| `-w, --wordlist`| Path to custom local wordlist (.txt)   |
| `-f, --focus`   | Filter (`all`, `auth`, `admin`, `api`) |
| `-o, --output`  | `text` | `json`                        |
| `-m, --mode`    | `basic` | `advanced`                   |
| `-v, --verbose` | Debug output                           |

## Core Features
map - Full Structure Mapping
```
netmap -d example.com
```

### focus — Targeted View
```
netmap -d example.com -f auth
netmap -d example.com -f admin
```
### export — Data Output
```
netmap -d example.com -o json > map.json
```

### advanced — Deeper Mapping
```
netmap -d example.com -m advanced
```

## How it works
### OSINT Discovery
Uses Certificate Transparency logs (e.g. crt.sh)
→ finds known subdomains without directly targeting the server

### Live Probing
Performs concurrent HTTP requests
→ detects active endpoints and validates responses

### Classification
Endpoints are labeled automatically:
* [AUTH] → login / authentication
* [ADMIN] → admin panels
* [API] → API routes
* [GENERAL] → standard paths

### Structuring
All data is transformed into a hierarchical map
→ from raw data → to clear structure

### Output Example
```
╔══════════════════════════════════════════╗
║           NetMap — Mapping               ║
╚══════════════════════════════════════════╝

example.com
├── api.example.com
│   ├── /v1/users [API]
│   └── /auth/login [AUTH]
├── admin.example.com
│   └── /login [ADMIN]
└── www.example.com
    └── /dashboard

Summary:
- Subdomains:          3
- Endpoints:           4
- High-interest nodes: 2
```

### Why NetMap?
Most tools give you:
* raw data
* messy output
* no structure
NetMap gives you:
* clarity
* structure
* visual understanding

### Part of the NET Ecosystem
* NetScope → discovery
* NetForce → performance testing
* NetIntel → analysis
* NetMap → visualization
Together: → a complete network analysis toolkit

### Roadmap
* Interactive web visualization
* Node color highlighting (risk-based)
* NetScope integration
* NetIntel overlay
* Graph export (nodes + edges)

## Author
Built by Lucas Mangroelal 
https://lucasmangroelal.nl

## ❤️ Support
* ⭐ Star the repo
* 🔗 Share it
* 🛠️ Contribute

# ⚠️ Disclaimer
**This tool is for educational and authorized use only. Do not use it against systems without permission.**
