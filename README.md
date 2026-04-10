# NetMap

**Visual Intelligence & Network Mapping Tool**

NetMap is een professionele CLI-tool voor developers en security specialisten om de structuur van een domein visueel in kaart te brengen. In plaats van een simpele lijst, bouwt NetMap een hiërarchische map van subdomeinen en endpoints met een sterke focus op **visual mapping** en **clarity**.

---

## 🔥 Quick Install (macOS & Linux)

Installeer de nieuwste versie met één commando:

```bash
curl -fsSL https://raw.githubusercontent.com/lucasenlucas/NetMap/main/install.sh | sh
```

*Na installatie kun je direct starten met `netmap --help`.*

---

## 🚀 Basis Gebruik

Breng de structuur van een website overzichtelijk in kaart:

```bash
netmap -d voorbeeld.nl
```

**Simpele uitleg:**
Brengt de structuur van een website overzichtelijk in kaart door visual mapping.

**Precieze uitleg:**
Start een OSINT-gebaseerde discovery (via Certificate Transparency logs) om subdomeinen te identificeren, en voert vervolgens gelijktijdige HTTP-requests uit om actieve endpoints te detecteren, te valideren en te structureren in een hiërarchische map.

---

## 🚩 Flags

### 1. Domein (`-d`, `--domain`)
*   **Simpele uitleg:** De website die je wilt analyseren.
*   **Precieze uitleg:** Definieert het root-domein voor de mapping. Wordt gebruikt als startpunt voor OSINT discovery en endpoint mapping.

### 2. Focus (`-f`, `--focus`)
*Opties: all, auth, admin, api*
*   **Simpele uitleg:** Laat alleen specifieke onderdelen zien, zoals login- of admin-pagina’s.
*   **Precieze uitleg:** Filtert nodes op basis van hun classificatie (EndpointType). Dit maakt het mogelijk om gericht delen van de infrastructuur te analyseren.

### 3. Output (`-o`, `--output`)
*Opties: text, json*
*   **Simpele uitleg:** Kies tussen een leesbare structuur of data voor andere tools.
*   **Precieze uitleg:** 
    *   `text` → Hiërarchische CLI visualisatie.
    *   `json` → Gestructureerde graph output (nodes + edges) voor frontend visualisatie.

### 4. Mode (`-m`, `--mode`)
*Opties: basic, advanced*
*   **Simpele uitleg:** Hoe diep NetMap analyseert.
*   **Precieze uitleg:** Bepaalt de intensiteit van discovery en validatie. `basic` is voor snelheid, `advanced` voor uitgebreide pad-detectie.

---

## 🧠 Intelligence Engine

| Onderdeel | Simpele uitleg | Precieze uitleg |
| :--- | :--- | :--- |
| **OSINT Discovery** | Zoekt openbare info over de site | Gebruikt Certificate Transparency logs (zoals crt.sh) als bron voor bekende subdomeinen. |
| **Live Probing** | Checkt welke onderdelen reageren | Voert gelijktijdige HTTP HEAD/GET requests uit om actieve endpoints te detecteren. |
| **Classification** | Labelt onderdelen | Classificeert endpoints via pattern matching (regex) in categorieën zoals auth, admin en api. |
| **Structuring** | Maakt er een overzicht van | Bouwt een hiërarchische graph (nodes + relaties) voor superieure visuele clarity. |

---

## 💡 Voorbeelden

```bash
# Alleen admin endpoints zien
netmap -d voorbeeld.nl -f admin

# Exporteren naar JSON voor weergave in websites
netmap -d voorbeeld.nl -o json > scan.json

# Uitgebreide scan met debug info
netmap -d voorbeeld.nl -m advanced -v
```

> **NetMap** maakt deel uit van het NET Ecosystem. Voor geautoriseerd gebruik en legitieme analyse.
