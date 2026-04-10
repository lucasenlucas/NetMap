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

## CLI Philosophy
Same system as NetScope & NetForce → no learning curve.
```
netmap -d <domain> [options]
```

## Core Flags
| Flag            | Description                            |
| --------------- | -------------------------------------- |
| `-d, --domain`  | Target domain                          |
| `-p, --pack`    | Discovery Pack (`standard`, `dns-extended`, `web-deep`, `api-focused`, `admin-stealth`, `full`) |
| `-w, --wordlist`| Path to custom local wordlist (.txt)   |
| `-f, --focus`   | Filter (`all`, `auth`, `admin`, `api`, `config`, `dev`) |
| `-o, --output`  | `text` | `json`                        |
| `-m, --mode`    | `basic` | `advanced`                   |
| `-v, --verbose` | Debug output                           |

## Core Features
### map — Full Structure Mapping
```
netmap -d example.com -p full
```

### focus — Targeted View
```
netmap -d example.com -f config
netmap -d example.com -f dev
```
### export — Data Output
```
netmap -d example.com -o json > map.json
```

## How it works
### OSINT Discovery
Uses Certificate Transparency logs (e.g. crt.sh)
→ finds known subdomains without directly targeting the server

### Live Probing
Performs concurrent HTTP requests
→ detects active endpoints and validates responses. The **Full Intelligence Engine** uses over 700+ combined targets.

### Classification (Categorical Intelligence)
Endpoints are labeled automatically:
* `[AUTH]` → login / authentication
* `[ADMIN]` → admin panels
* `[API]` → API routes
* `[CONFIG]` → Sensitive configuration files (.env, settings)
* `[DEV]` → Development artifacts (.git, Dockerfiles)
* `[GENERAL]` → standard paths

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
├── dev.example.com
│   └── /.env [CONFIG]
└── www.example.com
    └── /dashboard
```

## Author
Built by Lucas Mangroelal 
https://lucasmangroelal.nl

# ⚠️ Disclaimer
**This tool is for educational and authorized use only. Do not use it against systems without permission.**
