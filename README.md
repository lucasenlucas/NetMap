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

## Power Usage & Examples

### 🚀 The Ultra Scan (Krachtigste Command)
Gebruik dit voor de meest complete mapping van een target. Het combineert de grootste intelligentie-set met diepe DNS-reconstructie.
```bash
netmap -d voorbeeld.nl -p ultra -m advanced
```
* **Wat gebeurt er?**: Scant 2500+ targets, voert DNS brute-forcing uit en haalt alle CNAME/MX/TXT records op.

### 🏢 Infrastructure Only
Filter alle web-ruis weg en bekijk puur de DNS-architectuur en providers.
```bash
netmap -d voorbeeld.nl -f dns
```
* **Wat gebeurt er?**: Toont alleen de infrastructuur-nodes zoals mailservers en CDN-aliassen.

### 🔓 Security Audit
Focus op gevoelige configuraties en development artifacts.
```bash
netmap -d voorbeeld.nl -p full -f config
```
* **Wat gebeurt er?**: Zoekt specifiek naar `.env`, `wp-config`, `.git` mappen en andere risico-punten.

### 🤖 Automation / JSON Export
Exporteer de volledige graaf naar JSON voor gebruik in externe visualisatie-tools of pipelines.
```bash
netmap -d voorbeeld.nl -o json > map.json
```

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

## Author
Built by Lucas Mangroelal 
https://lucasmangroelal.nl

# ⚠️ Disclaimer
**This tool is for educational and authorized use only. Do not use it against systems without permission.**
