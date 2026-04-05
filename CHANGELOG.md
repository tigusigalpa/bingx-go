# Changelog

All notable changes to this project will be documented in this file.

## [1.4.0] - 2026-04-05

### 🎉 BingX API v3 Support

This major release brings full compatibility with BingX API v3, adding 20+ new methods and enhanced functionality across all services.

### Added

#### Trade Service - New Order Types & Methods
- **Order Type Constants**: Added support for new order types
  - `OrderTypeTriggerLimit` - TRIGGER_LIMIT orders
  - `OrderTypeTrailingStopMarket` - TRAILING_STOP_MARKET orders
  - `OrderTypeTrailingTPSL` - TRAILING_TP_SL orders
  
- **One-Click Reverse Position**: `OneClickReversePosition(symbol, recvWindow)` - Instantly reverse position direction
- **Auto Margin Addition**: `SetAutoAddMargin(symbol, positionSide, autoAddMargin, recvWindow)` - Hedge mode automatic margin addition
- **Multi-Assets Mode**: 
  - `SwitchMultiAssetsMode(multiAssetsMargin, recvWindow)` - Enable/disable multi-assets margin mode
  - `GetMultiAssetsMode(recvWindow)` - Query current multi-assets mode status
  - `GetMultiAssetsRules(recvWindow)` - Get multi-assets trading rules
  - `GetMultiAssetsMargin(recvWindow)` - Query multi-assets margin details
  
- **TWAP Orders** (Time-Weighted Average Price):
  - `PlaceTWAPOrder(params)` - Create TWAP order for large position building
  - `CancelTWAPOrder(orderId, recvWindow)` - Cancel active TWAP order
  - `GetTWAPOrder(orderId, recvWindow)` - Query specific TWAP order
  - `GetTWAPOrders(symbol, status, startTime, endTime, limit, recvWindow)` - List TWAP orders with filters

#### Market Service - Enhanced Market Data
- `GetOpenInterest(symbol)` - Get current open interest for symbol
- `GetOpenInterestHistory(symbol, period, limit, startTime, endTime)` - Historical open interest data
- `GetFundingRateInfo(symbol)` - Current funding rate information
- `GetBookTicker(symbol)` - Best bid/ask prices (futures)
- `GetSpotBookTicker(symbol)` - Best bid/ask prices (spot)
- `GetIndexPrice(symbol)` - Current index price
- `GetTickerPrice(symbol)` - Latest ticker price data

#### Account Service - Advanced Account Management
- `GetPositionRisk(symbol, recvWindow)` - Detailed position risk metrics
- `GetIncomeHistory(symbol, incomeType, startTime, endTime, limit, recvWindow)` - Income/PnL history with filtering
- `GetCommissionHistory(symbol, startTime, endTime, limit, recvWindow)` - Trading commission history
- `GetForceOrders(symbol, autoCloseType, startTime, endTime, limit, recvWindow)` - Liquidation/force close orders
- `GetPositionMode(recvWindow)` - Query hedge/one-way position mode
- `SetPositionMode(dualSidePosition, recvWindow)` - Switch between hedge and one-way mode

#### Coin-M Futures Service
- `GetPositionRisk(symbol, recvWindow)` - Position risk for Coin-M contracts
- `GetIncomeHistory(...)` - Income history for Coin-M
- `GetFundingRateHistory(symbol, limit)` - Historical funding rates
- `GetMarkPrice(symbol)` - Mark price for Coin-M
- `GetIndexPrice(symbol)` - Index price for Coin-M
- `GetRecentTrades(symbol, limit)` - Recent trades data

### Changed

#### HTTP Client
- Enhanced error handling with additional v3 error codes:
  - Added `100412` for authentication errors
  - Added `100429` for rate limit errors
  - Added `200002` for insufficient balance scenarios
  - Added explicit handling for success code `0`

#### Coin-M Trade Service
- Updated `CreateOrder` endpoint from v1 to v2: `/openApi/swap/v2/trade/order`

#### WebSocket
- Confirmed v3 compatibility for WebSocket URLs
- Added documentation comments for endpoint URLs

### Documentation

- Updated README.md with v3 API reference
- Increased method counts: 240+ total API methods
- Updated order types documentation
- Changed API documentation link to v3: https://bingx-api.github.io/docs-v3/
- Enhanced service descriptions with new capabilities

### API Compatibility

**Breaking Changes**: None - All changes are backward compatible additions

**New Capabilities**:
- TWAP order execution for institutional-grade trading
- Multi-assets margin mode for portfolio margin
- One-click position reversal for quick strategy changes
- Enhanced risk management with position risk metrics
- Comprehensive income and commission tracking

**Tested Against**: BingX API v3 (April 2026)

---

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
