package ai

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/syumai/tinyutil/httputil"
)

func Golden() string {
	c := New()

	var str strings.Builder
	str.WriteString(formatChart(c.Lookup("GC=F")))
	str.WriteString(formatChart(c.Lookup("SGC=F")))
	str.WriteString(formatChart(c.Lookup("SGU=F")))

	return str.String()
}

func JPYCNY() string {
	c := New()

	var str strings.Builder
	str.WriteString(formatChart(c.Lookup("JPYCNY=X")))
	return str.String()
}

func USDCNY() string {
	c := New()

	var str strings.Builder
	str.WriteString(formatChart(c.Lookup("CNY=X")))
	return str.String()
}

func formatChart(q *Chart, err error) string {
	if err != nil {
		return fmt.Sprintf(`
%s
`, escape(err.Error()))
	}
	return fmt.Sprintf(`
⭐ *%s*
*exchange*: %s
*current*: %s
*prev close price*: %s
*percentage*: %s%%
`,
		escape(q.ShortName()),
		escape(q.FullExchange()),
		escape(fmt.Sprint(q.Price())),
		escape(fmt.Sprint(q.PrevClosePrice())),
		escape(fmt.Sprint(((q.Price()-q.PrevClosePrice())/q.PrevClosePrice())*100)),
	)
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

func escape(str string) string {
	char := map[rune]bool{
		'_': true,
		'*': true,
		'[': true,
		']': true,
		'(': true,
		')': true,
		'~': true,
		'`': true,
		'>': true,
		'#': true,
		'+': true,
		'-': true,
		'=': true,
		'|': true,
		'{': true,
		'}': true,
		'.': true,
		'!': true,
	}

	s := &strings.Builder{}

	for _, c := range str {
		if char[c] {
			s.WriteRune('\\')
		}

		s.WriteRune(c)
	}

	return s.String()
}
