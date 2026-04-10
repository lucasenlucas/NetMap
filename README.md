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
| `-p, --pack`    | Discovery Pack (`standard`, `dns-extended`, `web-deep`, `api-focused`, `admin-stealth`, `full`, `ultra`) |
| `-w, --wordlist`| Path to custom local wordlist (.txt)   |
| `-f, --focus`   | Filter (`all`, `auth`, `admin`, `api`, `config`, `dev`, `dns`) |
| `-o, --output`  | `text` | `json`                        |
| `-m, --mode`    | `basic` | `advanced`                   |
| `-v, --verbose` | Debug output                           |

## Core Features
### ultra — Maximum Recon
```
netmap -d example.com -p ultra
```
The **Ultra Intelligence Engine** uses over 2500+ combined targets for deep infrastructure mapping.

### dns — Deep Infrastructure Recon
```
netmap -d example.com -f dns
```
NetMap now analyzes **CNAME, MX, and TXT** records to uncover hidden aliases, mail handlers, and cloud providers.

### Live Feedback
NetMap provides real-time terminal progress during high-volume scans:
`[~] Mapping: 12/50 hosts | 452/2500 endpoints | [Checking: /api/v2/config]`

## How it works
### OSINT Discovery
Uses Certificate Transparency logs (e.g. crt.sh)
→ finds known subdomains without directly targeting the server

### Live Probing & DNS Mapping
Performs concurrent HTTP requests and **Deep DNS Lookups**.
→ detects active endpoints and validates infrastructure records.

### Classification (Categorical Intelligence)
Endpoints are labeled automatically:
* `[AUTH]` → login / authentication
* `[ADMIN]` → admin panels
* `[API]` → API routes
* `[DNS]` → CNAME, MX, TXT records
* `[CONFIG]` → Sensitive configuration files (.env, settings)
* `[DEV]` → Development artifacts (.git, Dockerfiles)
* `[GENERAL]` → standard paths

## Author
Built by Lucas Mangroelal 
https://lucasmangroelal.nl

# ⚠️ Disclaimer
**This tool is for educational and authorized use only. Do not use it against systems without permission.**
