<p align="center">
  <img src="https://github.com/netseries/lucas_cdn/blob/main/netmap_banner.png?raw=true" alt="NetMap Banner"/>
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
curl -fsSL https://raw.githubusercontent.com/netseries/NetMap/main/install.sh | sh
```
Run instantly:
```
netmap --help
```
### Go install
```
go install github.com/netseries/netmap/cmd/netmap@latest
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

## Power Usage & Examples

### 🚀 The Ultra Scan (Most Powerful Command)
Use this for the most complete mapping of a target. It combines the largest intelligence set (2500+ targets) with deep DNS reconstruction and brute-forcing.
```bash
netmap -d example.com -p ultra -m advanced
```
* **What happens?**: Scans 2500+ targets, performs high-volume DNS brute-forcing, and retrieves all CNAME/MX/TXT records.

### 🏢 Infrastructure Only
Filter out all web noise and view pure DNS architecture and service providers.
```bash
netmap -d example.com -f dns
```
* **What happens?**: Displays only infrastructure nodes like mail handlers (MX) and CDN aliases (CNAME).

### 🔓 Security Audit
Focus specifically on sensitive configurations and development artifacts.
```bash
netmap -d example.com -p full -f config
```
* **What happens?**: Searches specifically for `.env`, `wp-config`, `.git` folders, and other high-risk points.

### 🤖 Automation / JSON Export
Export the full graph to JSON for use in external visualization tools or analysis pipelines.
```bash
netmap -d example.com -o json > map.json
```

## How it works
### OSINT Discovery
Uses Certificate Transparency logs (e.g., crt.sh)
→ finds known subdomains without directly targeting the server.

### Live Probing & DNS Mapping
Performs concurrent HTTP requests and **Deep DNS Lookups**.
→ detects active endpoints and validates infrastructure records in real-time.

### Live Progress Tracking
NetMap provides real-time terminal feedback during high-volume scans:
`[~] Mapping: 12/50 hosts | 452/2500 endpoints | [Checking: /api/v2/config]`

### Classification (Categorical Intelligence)
Endpoints are labeled automatically based on their nature:
* `[AUTH]` → login / authentication
* `[ADMIN]` → administrative panels
* `[API]` → API routes
* `[DNS]` → Infrastructure records (CNAME, MX, TXT)
* `[CONFIG]` → Sensitive configuration files (.env, settings)
* `[DEV]` → Development artifacts (.git, Dockerfiles)
* `[GENERAL]` → Standard web paths

## Author
Built by **Netseries Team** | [Netseries.dev](https://netseries.dev)
Contact: **team@netseries.dev**

# ⚠️ Disclaimer
**This tool is for educational and authorized use only. Do not use it against systems without permission.**
