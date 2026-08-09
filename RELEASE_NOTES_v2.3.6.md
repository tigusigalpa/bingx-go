# Release Notes v2.3.6 - Typed Book Tickers

**Release Date:** August 9, 2026

This patch adds precision-preserving typed book ticker helpers for BingX futures and spot markets while retaining the existing raw API methods.

---

## Added

- `Market().GetBookTickerData(*string) (*BookTicker, error)` for the futures response nested at `data.book_ticker`.
  - Maps `bid_price`, `bid_qty`, `ask_price`, `ask_qty`, `lastUpdateId`, and `time` to decimal-string fields.
  - Preserves JSON numeric tokens without converting monetary values through `float64`.
  - Also accepts string-encoded numeric fields.
- `Market().GetSpotBookTickerData(*string) (*SpotBookTicker, error)` for BingX spot's verified `data` array response at `/openApi/spot/v1/ticker/bookTicker`.
  - Maps spot `bidVolume` and `askVolume` fields to `BidQuantity` and `AskQuantity`.

## Compatibility

- `Market().GetBookTicker(*string) (map[string]interface{}, error)` is unchanged and continues to return the raw futures envelope.
- `Market().GetSpotBookTicker(*string) (map[string]interface{}, error)` is unchanged and continues to return its raw response.
- Both typed helpers reuse the existing request, signing, timeout, and API-error handling pipeline.

## Documentation

Documentation and examples now recommend `GetBookTickerData()` for futures quotes and describe the current `data.book_ticker.bid_price` and `ask_price` fields. Raw methods are explicitly documented for consumers needing the untouched exchange payload.

## Tests

Fixture-based tests cover the current futures payload, exact decimal-string preservation for JSON numbers and string numeric fields, malformed payloads, API error mapping, spot decoding, and unchanged raw futures behavior.

---

**Full Changelog**: https://github.com/tigusigalpa/bingx-go/compare/v2.3.5...v2.3.6
