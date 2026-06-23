<div align="center">

# 🛰️ tg-proxy-list — Free Telegram Proxy List (Verified MTProto Proxies)

**Fresh, auto-verified [Telegram](https://telegram.org) MTProto proxies — connection-tested, latency-ranked, geo-located, and updated every 6 hours.**
One click to connect. No app, no signup, no ads.

[![Verified proxies](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fraw.githubusercontent.com%2Fzakky8%2Ftg-proxy-list%2Fmain%2Fproxies.json&query=%24.count&label=verified%20proxies&color=22c55e&style=for-the-badge)](https://zakky8.github.io/tg-proxy-list/)
[![Last updated](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fraw.githubusercontent.com%2Fzakky8%2Ftg-proxy-list%2Fmain%2Fproxies.json&query=%24.generated_at_utc&label=updated&color=229ed9&style=for-the-badge)](https://zakky8.github.io/tg-proxy-list/)
[![Update proxies](https://github.com/zakky8/tg-proxy-list/actions/workflows/update.yml/badge.svg)](https://github.com/zakky8/tg-proxy-list/actions/workflows/update.yml)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg?style=for-the-badge)](LICENSE)

### 👉 **[Open the list & connect →](https://zakky8.github.io/tg-proxy-list/)** 👈

English · [Русский](README_RU.md) · [فارسی](README_FA.md) · [中文](README_CN.md)

</div>

---

## What this is

**tg-proxy-list is a free, continuously updated list of public Telegram MTProto proxies.** It helps you use Telegram where it is **blocked, throttled, or censored** — without installing anything.

Unlike most "free proxy" dumps, every proxy here is **actually verified** before it is published:

- ✅ **DNS resolved** and **TCP-connected** with **measured round-trip latency**
- ✅ **FakeTLS handshake** tested for `ee` proxies (a real protocol probe)
- 🌍 **Geo-located** to a country
- 📈 **Uptime** tracked across runs
- 🔁 **Re-checked every 6 hours** — dead proxies drop off automatically

> **Honest by design.** We never claim a proxy "works 100%". A proxy is labeled **`reachable`** (answered TCP) or **`handshake_ok`** (also completed a FakeTLS handshake). Public proxies can still go offline anytime — if one fails, pick another.

## Quick start

### Easiest — the website
**[zakky8.github.io/tg-proxy-list](https://zakky8.github.io/tg-proxy-list/)** → search/filter by country, sort by speed, tap **Connect**. Telegram opens and asks to enable the proxy. Done.

### From the raw list
Open **[`all_proxies.txt`](all_proxies.txt)** and tap any `https://t.me/proxy?...` link on your phone — Telegram handles the rest.

### Subscribe (always-fresh raw URLs)
```
# All proxies, fastest first (one t.me link per line)
https://raw.githubusercontent.com/zakky8/tg-proxy-list/main/all_proxies.txt

# Structured data (country, latency, type, uptime, status)
https://raw.githubusercontent.com/zakky8/tg-proxy-list/main/proxies.json

# Just your country, e.g. Iran / Russia / Germany
https://raw.githubusercontent.com/zakky8/tg-proxy-list/main/by_country/IR.txt
```

## Output formats

| File | What it is |
|------|------------|
| [`all_proxies.txt`](all_proxies.txt) | Flat list of `https://t.me/proxy?...` links, fastest first |
| [`proxies.json`](proxies.json) | Structured records: `server, port, secret, type, country, latency_ms, status, uptime_pct, last_checked_utc` |
| [`sorted_by_latency.txt`](sorted_by_latency.txt) | Human-readable table (latency · country · status · link) |
| [`by_country/`](by_country/) | One `.txt` per country (`IR.txt`, `RU.txt`, `DE.txt`, …) |

## How it works

```
collect  →  verify  →  geo  →  publish
```

1. **collect** — fetches candidates from many public proxy feeds (see [`sources.txt`](sources.txt)), parses `tg://`, `t.me/proxy`, and `host:port:secret` formats, and dedupes.
2. **verify** — a bounded worker pool resolves DNS, TCP-connects (measuring latency), and runs a FakeTLS handshake probe for `ee` secrets.
3. **geo** — maps each IP to a country using the public [DB-IP](https://db-ip.com) lite dataset.
4. **publish** — writes every output format and updates uptime history.

A GitHub Action ([`.github/workflows/update.yml`](.github/workflows/update.yml)) runs the whole pipeline every 6 hours and republishes the site.

## Run it yourself

Requires [Go](https://go.dev) 1.22+.

```bash
git clone https://github.com/zakky8/tg-proxy-list
cd tg-proxy-list
go run ./cmd/tgproxy            # full run
go run ./cmd/tgproxy --limit 100 --handshake-only   # quick / strict run
go test ./...                   # tests
```

Flags: `--sources`, `--concurrency`, `--timeout`, `--limit`, `--handshake-only`, `--geo`, `--out`, `--docs`.

Add your own feeds by editing [`sources.txt`](sources.txt) (one URL per line).

## FAQ

**What is an MTProto proxy?** Telegram's own proxy protocol. It relays *only* your Telegram traffic to help you reach Telegram when it is blocked — it doesn't touch the rest of your device.

**Is it safe?** Your messages stay protected by Telegram's existing encryption. A proxy operator only sees that you connect to Telegram. Prefer `handshake_ok` proxies and avoid sensitive logins over unknown proxies.

**How is this different from `SoliSpirit/mtproto`?** That project publishes a single unverified `.txt`. tg-proxy-list verifies reachability + latency, geo-locates, tracks uptime, ships multiple formats, and includes a one-click connect website.

## Disclaimer

These are **public** proxies aggregated from open sources, provided **as is** with no guarantee of availability, security, or performance. This project does not operate the proxies. Use lawfully and responsibly. Country data © [DB-IP](https://db-ip.com) ([CC-BY-4.0](https://creativecommons.org/licenses/by/4.0/)).

## License

[MIT](LICENSE) — free to use, modify, and distribute.

---

<div align="center">
<sub>Keywords: free telegram proxy · mtproto proxy · telegram proxy list 2026 · working telegram proxies · telegram proxy server · bypass telegram censorship · прокси для телеграм · پروکسی تلگرام</sub>
</div>
