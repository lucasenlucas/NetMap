# NetMap

NetMap is a professional visual intelligence and mapping CLI tool built in Go. It turns discovered network domains, subdomains, and endpoints into clean, highly readable tree maps directly in your terminal. 

Designed for red teams, security researchers, and systems administrators, NetMap maps the structure of target environments lightning-fast via automated OSINT (Certificate Transparency logs) and deeply concurrent active HTTP checking.

---

## ⚡ Quick Install (macOS & Linux)

Paste this single command in your terminal to instantly install the latest version:

```bash
curl -fsSL https://raw.githubusercontent.com/lucasenlucas/NetMap/main/install.sh | sh
```

*(You might be prompted for your password to install into `/usr/local/bin`)*

## 🪟 Windows
Go to [Releases](https://github.com/lucasenlucas/NetMap/releases) and download the `.exe` directly.

---

## 📖 Help Schema & Usage

If you ever need a quick reminder of how to run NetMap, simply run:
```bash
netmap --help
```

### Advanced Examples

Map an entire root domain implicitly traversing discovered endpoints:
```bash
netmap -d example.com
```

Focus the renderer dynamically to strictly print protected, highly-valuable endpoints like **admin**, **login**, and **API** pathways:
```bash
netmap -d target.com -f auth
netmap -d target.com -f admin
netmap -d target.com -f api
```

Incorporate NetMap into extensive backend scanning workflows by generating purely deterministic JSON blobs (structured perfectly for React/D3 pipelines):
```bash
netmap -d example.com -o json > report.json
```

---

## 🛠 Features

- **Concurrent Engine**: Blazing fast mapping sweeping up to hundreds of possible HTTP paths simultaneously without saturating standard routing architectures.
- **Categorical Intelligence**: Automatically attributes endpoint signatures (Regex tagging) into recognizable groups. 
- **Zero Binary Dependency**: Standalone pure golang logic. Zero extra scripts required.
- **Beautiful Graph Output**: Premium hierarchical formatting tailored exclusively to terminal purists. Employs mathematically padded trees layout formats. 
- **Safe & Covert OSINT**: Native `crt.sh` querying logic grabs domain trees without explicitly triggering loud recursive DNS checks.

---

> ⚠️ For authorized testing only. NetMap is part of the NET Ecosystem Toolkit.
