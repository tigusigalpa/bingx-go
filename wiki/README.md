# Wiki Documentation

This folder contains all the Wiki pages for bingx-go. The docs are written in a conversational, human-friendly style — no corporate jargon, just practical guidance.

## What's Here

### Getting Started
- **[Home.md](Home.md)** — Main landing page with quick navigation
- **[Getting-Started.md](Getting-Started.md)** — Installation, API keys, first steps
- **[API-v3-Features.md](API-v3-Features.md)** — All the new v3 goodies

### Core Guides
- **[Market-Service.md](Market-Service.md)** — Prices, charts, order books
- **[Account-Service.md](Account-Service.md)** — Balance, positions, risk
- **[Trade-Service.md](Trade-Service.md)** — Orders, TWAP, position management

### Coming Soon
- **Wallet-Service.md** — Deposits, withdrawals, transfers
- **WebSocket-Service.md** — Real-time streaming
- **Error-Handling.md** — Dealing with problems gracefully
- **Coin-M-Futures.md** — Coin-margined contracts

## Publishing to GitHub Wiki

The easiest way:

```bash
# Clone the wiki repo
git clone https://github.com/tigusigalpa/bingx-go.wiki.git

# Copy the files over
cp wiki/*.md bingx-go.wiki/

# Push it
cd bingx-go.wiki
git add .
git commit -m "Update docs"
git push
```

Or just copy-paste each file manually through the GitHub web interface.

## Writing Style

These docs follow a few principles:

- **Be conversational** — Write like you're explaining to a friend
- **Show, don't tell** — Code examples > walls of text
- **Be practical** — Real-world examples that people can actually use
- **Skip the fluff** — No "In this section we will discuss..."

## Status

| Page | Done? |
|------|-------|
| Home | ✅ |
| Getting Started | ✅ |
| API v3 Features | ✅ |
| Market Service | ✅ |
| Account Service | ✅ |
| Trade Service | ✅ |
| Wallet Service | � |
| WebSocket Service | � |
| Error Handling | � |
| Coin-M Futures | � |

---

**Questions?** [Open an issue](https://github.com/tigusigalpa/bingx-go/issues)
