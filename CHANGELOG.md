# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased] - 2026-03-15

### Added

#### Market Service
- **Spot Symbols Endpoint Update**: Updated `GetSpotSymbols()` to use `/openApi/spot/v1/common/symbols` endpoint
  - Added support for `maxMarketNotional` field (maximum notional amount for single market order)
  - Added new symbol status value `29 = Pre-Delisted`
  - Full status values: 0=Offline, 1=Online, 5=Pre-open, 10=Accessed, 25=Suspended, 29=Pre-Delisted, 30=Delisted

- **Spot Klines v2 Endpoint**: Updated `GetSpotKlines()` to use `/openApi/spot/v2/market/kline`
  - Added optional `timeZone` parameter (0=UTC (default), 8=UTC+8)
  - Updated max limit from 1000 to 1440 records

#### Spot Account Service
- **Wallet Type Constants**: Added constants for wallet types
  ```go
  WalletTypeFund             = 1  // Fund Account
  WalletTypeStandardFutures  = 2  // Standard Futures Account
  WalletTypePerpetualFutures = 3  // Perpetual Futures Account
  WalletTypeSpot             = 4  // Spot Account (NEW)
  ```

- **Internal Transfer Update**: Updated `InternalTransfer()` method signature
  - Changed `walletType` parameter to `int` (use constants above)
  - Added `userAccountType` parameter (1=UID, 2=Phone number, 3=Email)
  - Added `userAccount` parameter
  - Added optional `callingCode` parameter (required when userAccountType=2)
  - Added optional `transferClientID` parameter (custom ID, max 100 chars)
  - Added optional `recvWindow` parameter

#### Sub-Account Service
- **Sub-Account Wallet Type Constants**: Added constants for sub-account wallet types
  ```go
  SubAccountWalletTypeFund             = 1  // Fund Account
  SubAccountWalletTypeStandardFutures  = 2  // Standard Futures Account
  SubAccountWalletTypePerpetualFutures = 3  // Perpetual Futures Account
  SubAccountWalletTypeSpot             = 15 // Spot Account (NEW)
  ```

- **Sub-Account Internal Transfer Update**: Updated `SubAccountInternalTransfer()` method
  - Changed `walletType` parameter to `int` (use constants above)
  - Added `userAccountType` parameter (1=UID, 2=Phone number, 3=Email)
  - Added `userAccount` parameter
  - Added optional `callingCode` parameter
  - Added optional `transferClientID` parameter
  - Added optional `recvWindow` parameter

- **New Method**: `SubMotherAccountAssetTransfer()` - Sub-Mother Account Asset Transfer Interface
  - Flexible asset transfer between parent and sub-accounts
  - Supports account types: 1=Funding, 2=Standard futures, 3=Perpetual U-based, 15=Spot
  - Only available to master account
  - Returns `tranId` (transfer record ID)

- **New Method**: `GetSubMotherAccountTransferableAmount()` - Query Sub-Mother Account Transferable Amount
  - Query supported coins and available transferable amounts
  - Only available to master account
  - Returns coins array with id, name, and availableAmount

- **New Method**: `GetSubMotherAccountTransferHistory()` - Query Sub-Mother Account Transfer History
  - Query transfer history between sub-accounts and parent account
  - Supports filtering by type, tranId, time range
  - Pagination support (pageID, pagingSize)
  - Only available to master account

### Changed
- Updated BingX API integration to support changes from December 2025 through February 2026
- Improved type safety with dedicated constants for wallet types
- Enhanced parameter validation across all updated methods

### API Compatibility
- Breaking changes in method signatures for `InternalTransfer()` and `SubAccountInternalTransfer()`
- New constants provide better type safety and code clarity
- All new parameters are optional with sensible defaults where applicable

---

## API Reference
For detailed API documentation, see: https://bingx-api.github.io/docs-v3/
