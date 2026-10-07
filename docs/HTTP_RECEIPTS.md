# Bounded HTTP receipts

All SDK HTTP methods, including map-based methods and error responses, have a default body limit of 16 MiB. Configure a positive `int64` limit with `bingx.WithMaxResponseBytes`. Zero, negative values and `math.MaxInt64` fail with `http.ErrInvalidResponseLimit` before I/O. `NewBaseHTTPClient` retains its original signature; direct HTTP users can use `NewBaseHTTPClientWithOptions`.

```go
package example

import (
    "context"
    "errors"
    "net/http"
    "time"

    bingx "github.com/tigusigalpa/bingx-go/v2"
    bxhttp "github.com/tigusigalpa/bingx-go/v2/http"
)

func Fetch(ctx context.Context, apiKey, apiSecret string) (*bxhttp.RawResponse, error) {
    transport := http.DefaultTransport.(*http.Transport).Clone()
    transport.MaxResponseHeaderBytes = 64 << 10
    client := bingx.NewClient(apiKey, apiSecret,
        bingx.WithHTTPClient(&http.Client{Transport: transport, Timeout: 10 * time.Second}),
        bingx.WithMaxResponseBytes(1 << 20),
    )
    receipt, err := client.Market().GetKlinesRaw(ctx, "BTC-USDT", "1m", 100, nil, nil)
    if err != nil {
        var responseError *bxhttp.ResponseError
        if errors.As(err, &responseError) {
            // Complete bounded evidence is available even for rejected responses.
            archive(responseError.Receipt())
        }
        return receipt, err
    }
    archive(receipt)
    return receipt, nil
}
```

`archive` is application-owned. Configure the injected client and its transport before use. The SDK copies the `http.Client` value, shares its transport/jar, and does not mutate `http.DefaultTransport`. Transport resources and idle-connection cleanup remain caller-owned. A nil injected client uses the SDK's 30-second default. A client with no timeout requires an application deadline.

## Evidence and failure behavior

The SDK reads at most `limit+1` body bytes and closes every response body. Allocation is proportional to the configured limit rather than the advertised length or stream length; copying and JSON parsing also consume memory within that bounded scope. The limit applies to bytes delivered by `Response.Body`, including transparently decompressed HTTP bodies. For encoded bytes, configure a transport with compression disabled and handle that encoding yourself. Transport/header allocations are separate: configure `MaxResponseHeaderBytes` on your transport if needed.

Only a completed read produces a receipt. Accessors return `BodyBytes()` (copied exact bytes), `CompletedAt()` (local time at read completion), `StatusCode()`, `Headers()`, `Route()` and `Selectors()`. Completion time is not provider event time. The exported `Body` and `RetrievedAt` fields remain available for compatibility; changing them cannot alter accessor evidence. Do not mutate legacy fields concurrently.

Headers are allowlisted to validated Date, Retry-After, Content-Type (media type only), Content-Length and Content-Encoding. Selectors are allowlisted to symbol, interval, period, limit, fromId, startTime and endTime, with restricted values. Route contains only the original SDK request path. Cookies, Location, authorization headers, API keys, signing timestamp and signature are excluded. Receipts and SDK errors omit sensitive payloads and signed URLs from ordinary `%v`, `%+v` and `%#v` diagnostics. Explicitly reading payloads or unwrapping errors is for evidence processing, not automatic logging.

| Condition | Raw result | Stable classification |
| --- | --- | --- |
| Body exceeds limit | nil receipt | `*http.ResponseTooLargeError`, `errors.Is(err, http.ErrResponseTooLarge)` |
| Truncated/failed read | nil receipt | `*http.IncompleteResponseError`, `errors.Is(err, http.ErrIncompleteResponse)` |
| Cancellation/deadline | nil receipt | `errors.Is(err, context.Canceled/DeadlineExceeded)`; body failures are also incomplete |
| Complete malformed JSON, HTTP failure or provider failure | non-nil receipt plus error | `*http.ResponseError`; `.Receipt()` returns independent evidence |

Treat **every non-nil error as failure**, even when evidence is present. Map-based methods retain their original provider exception types; raw errors wrap those exceptions for `errors.As`. No partial body is ever a successful raw result.

## Retry ownership

The SDK adds no retry loop and makes one `http.Client.Do` call per request. Every raw market wrapper passes through the supplied context. The application owns pagination, rate limits, finite retry count, an overall deadline and cancellable backoff. Do not retry canceled operations or malformed/oversized bodies unchanged. Use safe Retry-After evidence when appropriate. Standard/custom transport retries and HTTP redirects are part of the injected client's policy; they must also be bounded. Coordinate one retry owner rather than stacking retry loops.

## Historical source admission

The [documentation fixture](../services/testdata/contracts/swap-market.documentation.json) pins the full official docs-v3 source revision `1e8fae7f702525e99fdf7cf96b67070b51c33644` and the shorter AI reference revision `5fb44d121b7e10ef3493bb4de21fedf7e5c98ac6`. [Fixture classifications](../services/testdata/contracts/README.md) distinguish documentation and synthetic tests. No live exchange capture is claimed.

Historical trades: the official route is `/openApi/swap/v1/market/historicalTrades`. ID/time field types are documented, but pagination direction/inclusion and quantity units remain unproven. Accordingly no new historical-trades method or endpoint substitution is made. Legacy `GetAggregateTradesRaw` retains its existing route and is **not eligible for proven historical admission**.

Funding: the [full official reference](https://github.com/BingX-API/docs-v3/blob/1e8fae7f702525e99fdf7cf96b67070b51c33644/static/js/app.8bc50bc0c8a6a308c3e4.js) (module `kGAv`, clarification dated 2026-08-31) explicitly documents completed historical settlement records: `fundingRate` is the settled rate and `fundingTime` is the settlement's Unix millisecond timestamp. It excludes the current/estimated next-settlement rate and `nextFundingTime`. This resolves the conflicting current/next description in the shorter AI reference. `GetFundingRatesRaw` preserves the array response unchanged. Pagination order, boundary inclusivity and complete time-range coverage are not established; consumers must validate those separately rather than infer completeness from one successful receipt. No live behavior is asserted by this documentation fixture.

Candles: bounded transport receipts establish retrieval provenance, not finality. Consumers must separately verify/filter closure and their accepted provider payload contract before admission; SDK boundary tests do not establish it.
