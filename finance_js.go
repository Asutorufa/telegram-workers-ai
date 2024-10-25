package ai

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/syumai/tinyutil/httputil"
)

func Golden() string {
	c := New()

	q, err := c.Lookup("GC=F")
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("exchange: %s\ncurrent: %f\nprev close price: %f\npercentage: %f%%", q.FullExchange(), q.Price(), q.PrevClosePrice(), ((q.Price()-q.PrevClosePrice())/q.PrevClosePrice())*100)
}

type client struct {
	httpClient *httputil.Client
}

const (
	yfurl = "https://query1.finance.yahoo.com"

	// do we need this?
	otherURL = "https://finance.yahoo.com/quote"
)

func New() *client {
	return &client{
		httpClient: httputil.DefaultClient,
	}
}

func (c *client) Lookup(ticker string) (*Chart, error) {
	yURL := fmt.Sprintf("%s/v8/finance/chart/%s", yfurl, ticker)

	request, reqErr := http.NewRequest(http.MethodGet, yURL, nil)
	if reqErr != nil {
		return nil, reqErr
	}

	request.Header.Set("User-Agent", "curl/7.68.0")

	response, reqErr := c.httpClient.Do(request)
	if reqErr != nil {
		return nil, reqErr
	}

	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("received status code %d", response.StatusCode)
	}

	var chart Chart
	if decErr := json.NewDecoder(response.Body).Decode(&chart); decErr != nil {
		return nil, decErr
	}

	return &chart, nil
}
