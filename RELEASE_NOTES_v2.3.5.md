# Release Notes v2.3.5 - Spot Account Wallet Overview

**Release Date:** August 8, 2026

BingX retired the `/openApi/wallets/v1/capital/fundBalance` endpoint used by `SpotAccount().GetFundBalance()`, which now returns `100400 this api is not exist`. This release removes that broken call and exposes the supported `GET /openApi/account/v1/allAccountBalance` overview endpoint as `GetAccountOverview()`.

---

## 🐛 Bug Fixes

### Spot Account Wallet Overview

- **Removed** the retired endpoint `/openApi/wallets/v1/capital/fundBalance` from `SpotAccountService.GetFundBalance()`.
- `GetFundBalance()` is now marked `Deprecated`. It no longer makes an HTTP request and returns the local error:
  ```
  GetFundBalance is retired; use GetAccountOverview instead
  ```
  It will be removed entirely in the next major release.
- Added `SpotAccountService.GetAccountOverview(accountType *string)`:
  - Calls `GET /openApi/account/v1/allAccountBalance`.
  - `accountType` is optional. When `nil`, the API returns all wallet types.
  - When provided, passes `accountType=<value>` in the query.
- Added account type string constants in `services/spotaccount.go`:
  ```go
  AccountTypeSpotFund     = "sopt"
  AccountTypeStdFutures   = "stdFutures"
  AccountTypeCoinMPerp    = "coinMPerp"
  AccountTypeUSDTMPerp    = "USDTMPerp"
  AccountTypeCopyTrading  = "copyTrading"
  AccountTypeGrid         = "grid"
  AccountTypeWealth       = "eran"
  AccountTypeC2C          = "c2c"
  ```

---

## 🧪 Testing

New regression tests in `services/spotaccount_test.go` using `httptest.Server`:

- `GetAccountOverview(nil)` calls exactly `/openApi/account/v1/allAccountBalance`.
- `GetAccountOverview(&AccountTypeUSDTMPerp)` produces `?accountType=USDTMPerp`.
- Signature is computed from the raw canonical query string and placed last.
- `GetFundBalance()` returns the retirement error and does not make an HTTP call.

---

## 🔄 Migration Guide

Replace any call to `client.SpotAccount().GetFundBalance()` with `client.SpotAccount().GetAccountOverview(nil)`:

```go
// Before (retired, returns an error)
balance, err := client.SpotAccount().GetFundBalance()

// After — overview of all wallet types
overview, err := client.SpotAccount().GetAccountOverview(nil)

// After — only the USDT-M perpetual wallet
usdtM := services.AccountTypeUSDTMPerp
usdtOverview, err := client.SpotAccount().GetAccountOverview(&usdtM)
```

---

## 📚 Documentation

- `README.md` Spot Account Service section updated to recommend `GetAccountOverview` and document `AccountType*` constants.
- `examples/basic/main.go` and `examples/v3_features/main.go` updated to demonstrate `GetAccountOverview` usage.
- `CHANGELOG.md` updated with the v2.3.5 release notes.

---

**Full Changelog**: https://github.com/tigusigalpa/bingx-go/compare/v2.3.4...v2.3.5
