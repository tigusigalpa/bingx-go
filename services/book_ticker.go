package services

import (
	"encoding/json"
	"fmt"
	"regexp"
)

var decimalPattern = regexp.MustCompile(`^-?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+-]?\d+)?$`)

type BookTicker struct {
	Symbol       string
	BidPrice     string
	BidQuantity  string
	AskPrice     string
	AskQuantity  string
	LastUpdateID string
	Timestamp    string
}

type SpotBookTicker struct {
	Symbol      string
	BidPrice    string
	BidQuantity string
	AskPrice    string
	AskQuantity string
	Timestamp   string
}

func (s *MarketService) GetBookTickerData(symbol *string) (*BookTicker, error) {
	params := bookTickerParams(symbol)
	var response struct {
		Data json.RawMessage `json:"data"`
	}
	if err := s.client.RequestJSON("GET", "/openApi/swap/v2/quote/bookTicker", params, &response); err != nil {
		return nil, err
	}
	if len(response.Data) == 0 || string(response.Data) == "null" {
		return nil, fmt.Errorf("BingX book ticker response is missing data")
	}

	var data struct {
		BookTicker json.RawMessage `json:"book_ticker"`
	}
	if err := json.Unmarshal(response.Data, &data); err != nil {
		return nil, fmt.Errorf("BingX book ticker response has malformed data: %w", err)
	}
	if len(data.BookTicker) == 0 || string(data.BookTicker) == "null" {
		return nil, fmt.Errorf("BingX book ticker response is missing data.book_ticker")
	}

	var ticker BookTicker
	if err := json.Unmarshal(data.BookTicker, &ticker); err != nil {
		return nil, fmt.Errorf("BingX book ticker response has malformed data.book_ticker: %w", err)
	}
	return &ticker, nil
}

func (s *MarketService) GetSpotBookTickerData(symbol *string) (*SpotBookTicker, error) {
	params := bookTickerParams(symbol)
	var response struct {
		Data json.RawMessage `json:"data"`
	}
	if err := s.client.RequestJSON("GET", "/openApi/spot/v1/ticker/bookTicker", params, &response); err != nil {
		return nil, err
	}
	if len(response.Data) == 0 || string(response.Data) == "null" {
		return nil, fmt.Errorf("BingX spot book ticker response is missing data")
	}

	var tickers []SpotBookTicker
	if err := json.Unmarshal(response.Data, &tickers); err != nil {
		return nil, fmt.Errorf("BingX spot book ticker response has malformed data: %w", err)
	}
	if len(tickers) == 0 {
		return nil, fmt.Errorf("BingX spot book ticker response contains no ticker")
	}
	return &tickers[0], nil
}

func (t *BookTicker) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	var err error
	if t.Symbol, err = requiredString(fields, "symbol"); err != nil {
		return err
	}
	if t.BidPrice, err = requiredDecimal(fields, "bid_price"); err != nil {
		return err
	}
	if t.BidQuantity, err = requiredDecimal(fields, "bid_qty"); err != nil {
		return err
	}
	if t.AskPrice, err = requiredDecimal(fields, "ask_price"); err != nil {
		return err
	}
	if t.AskQuantity, err = requiredDecimal(fields, "ask_qty"); err != nil {
		return err
	}
	if t.LastUpdateID, err = requiredDecimal(fields, "lastUpdateId"); err != nil {
		return err
	}
	if t.Timestamp, err = requiredDecimal(fields, "time"); err != nil {
		return err
	}
	return nil
}

func (t *SpotBookTicker) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	var err error
	if t.Symbol, err = requiredString(fields, "symbol"); err != nil {
		return err
	}
	if t.BidPrice, err = requiredDecimal(fields, "bidPrice"); err != nil {
		return err
	}
	if t.BidQuantity, err = requiredDecimal(fields, "bidVolume"); err != nil {
		return err
	}
	if t.AskPrice, err = requiredDecimal(fields, "askPrice"); err != nil {
		return err
	}
	if t.AskQuantity, err = requiredDecimal(fields, "askVolume"); err != nil {
		return err
	}
	if t.Timestamp, err = requiredDecimal(fields, "time"); err != nil {
		return err
	}
	return nil
}

func bookTickerParams(symbol *string) map[string]interface{} {
	params := map[string]interface{}{}
	if symbol != nil {
		params["symbol"] = *symbol
	}
	return params
}

func requiredString(fields map[string]json.RawMessage, name string) (string, error) {
	raw, ok := fields[name]
	if !ok {
		return "", fmt.Errorf("missing %s", name)
	}
	if string(raw) == "null" {
		return "", fmt.Errorf("%s must be a string", name)
	}

	var value string
	if err := json.Unmarshal(raw, &value); err != nil {
		return "", fmt.Errorf("%s must be a string", name)
	}
	return value, nil
}

func requiredDecimal(fields map[string]json.RawMessage, name string) (string, error) {
	raw, ok := fields[name]
	if !ok {
		return "", fmt.Errorf("missing %s", name)
	}
	if string(raw) == "null" {
		return "", fmt.Errorf("%s must be a JSON number or string", name)
	}

	var stringValue string
	if err := json.Unmarshal(raw, &stringValue); err == nil {
		if !decimalPattern.MatchString(stringValue) {
			return "", fmt.Errorf("%s must be a decimal string", name)
		}
		return stringValue, nil
	}

	var number json.Number
	if err := json.Unmarshal(raw, &number); err != nil {
		return "", fmt.Errorf("%s must be a JSON number or string", name)
	}
	return number.String(), nil
}
