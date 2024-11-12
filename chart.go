package ai

import (
	"time"
)

type Chart struct {
	TickerData struct {
		Results []struct {
			Meta struct {
				Currency       string  `json:"currency"`
				Symbol         string  `json:"symbol"`
				Exchange       string  `json:"exchangeName"`
				FullExchange   string  `json:"fullExchangeName"`
				Instrument     string  `json:"instrumentType"`
				FirstTrade     int64   `json:"firstTradeDate"`
				MarketTime     int64   `json:"regularMarketTime"`
				GMTOffset      int     `json:"gmtoffset"`
				TimeZone       string  `json:"timezone"`
				TimeZoneName   string  `json:"exchangeTimezoneName"`
				MarketPrice    float64 `json:"regularMarketPrice"`
				PrevClosePrice float64 `json:"previousClose"`
				ShortName      string  `json:"shortName"`
			} `json:"meta"`
		} `json:"result"`
	} `json:"chart"`
}

func (c *Chart) empty() bool {
	if c == nil {
		return true
	}

	if len(c.TickerData.Results) == 0 {
		return true
	}

	return false
}

func (c *Chart) Currency() string {
	if c.empty() {
		return ""
	}

	return c.TickerData.Results[0].Meta.Currency
}

func (c *Chart) Symbol() string {
	if c.empty() {
		return ""
	}

	return c.TickerData.Results[0].Meta.Symbol
}

func (c *Chart) ShortName() string {
	return c.TickerData.Results[0].Meta.ShortName
}

func (c *Chart) FullExchange() string {
	return c.TickerData.Results[0].Meta.FullExchange
}

func (c *Chart) Exchange(symbol bool) string {
	if c.empty() {
		return ""
	}

	exSym := c.TickerData.Results[0].Meta.Exchange
	if symbol {
		return exSym
	}

	switch exSym {
	case "NYQ":
		return "New York Stock Exchange"
	case "NMS":
		return "National Market System"
	default:
		return exSym
	}
}

func (c *Chart) InstrumentType() string {
	if c.empty() {
		return ""
	}

	return c.InstrumentType()
}

func (c *Chart) FirstTrade() time.Time {
	if c.empty() {
		return time.Time{}
	}

	return time.Unix(c.TickerData.Results[0].Meta.FirstTrade, 0)
}

func (c *Chart) MarketTime() time.Time {
	if c.empty() {
		return time.Time{}
	}

	return time.Unix(c.TickerData.Results[0].Meta.MarketTime, 0)
}

func (c *Chart) Price() float64 {
	if c.empty() {
		return 0
	}

	return c.TickerData.Results[0].Meta.MarketPrice
}

func (c *Chart) PrevClosePrice() float64 {
	if c.empty() {
		return 0
	}

	return c.TickerData.Results[0].Meta.PrevClosePrice
}

func (c *Chart) Delta() float64 {
	if c.empty() {
		return 0
	}

	return c.Price() - c.PrevClosePrice()
}

func (c *Chart) DeltaPerc() float64 {
	if c.empty() {
		return 0
	}

	d := c.Delta()
	p := d / c.PrevClosePrice()
	return p * 100
}
