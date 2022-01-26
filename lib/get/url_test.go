package get

import (
	"fmt"
	"net/url"
	"path"
	"testing"

	"github.com/artnoi43/1inch-limit-orders-api-poller/enums"
)

func TestPath(t *testing.T) {
	expectedAll := "/v2.0/1/limit-order/all"
	actualAll := getPath(enums.ChainEthereum, "all")
	if actualAll != expectedAll {
		t.Logf(
			"Expected: %s\nActual: %s\n",
			expectedAll,
			actualAll,
		)
		t.Fatal("results not matched")
	}
	expectedCount := "/v2.0/56/limit-order/count"
	actualCount := getPath(enums.ChainBsc, "count")
	if actualCount != expectedCount {
		t.Logf(
			"Expected: %s\nActual: %s\n",
			expectedCount,
			actualCount,
		)
		t.Fatal("results not matched")
	}
}

func TestQuery(t *testing.T) {
	type param struct {
		pages int
		limit int
	}
	tests := map[param]string{
		{
			pages: 69,
			limit: 400,
		}: "page=69&limit=400&sortby=createDateTime",
		{
			pages: 10,
			limit: 169,
		}: "page=10&limit=169&sortby=createDateTime",
	}
	for p, expected := range tests {
		actual := getQuery(p.pages, p.limit)
		if actual != expected {
			t.Logf(
				"Expected: %s\nActual: %s\n",
				expected,
				actual,
			)
			t.Fatal("results not matched")
		}
	}
}

// Test the whole URL construction from count and limit
func TestCountURL(t *testing.T) {
	tests := map[[2]int]string{}
	limit := 500
	// Map count and query string
	for i := 0; i < 3; i++ {
		inputs := [2]int{i * 1000, limit}
		page := getPages(inputs[0], inputs[1])
		tests[inputs] = fmt.Sprintf(
			"page=%d&limit=%d&sortby=createDateTime",
			page, limit,
		)
	}

	for inputs, expected := range tests {
		actualPages := getPages(inputs[0], inputs[1])
		actualQuery := getQuery(actualPages, limit)
		if actualQuery != expected {
			t.Logf(
				"Expected: %s\nActual: %s\n",
				expected,
				actualQuery,
			)
			t.Fatal("results not matched")
		}
	}
}

func TestURL(t *testing.T) {
	type foo struct {
		endpoint string
		chain    enums.Chain
		count    int
		limit    int
	}
	tests := map[foo]string{
		{
			endpoint: "all",
			chain:    enums.ChainEthereum,
			count:    3000,
			limit:    500,
		}: "https://limit-orders.1inch.io/v2.0/1/limit-order/all?page=6&limit=500&sortby=createDateTime",
		{
			endpoint: "all",
			chain:    enums.ChainBsc,
			count:    2401,
			limit:    100,
		}: "https://limit-orders.1inch.io/v2.0/56/limit-order/all?page=25&limit=100&sortby=createDateTime",
		{
			endpoint: "count",
			chain:    enums.ChainPolygon,
			count:    69,
			limit:    69,
		}: "https://limit-orders.1inch.io/v2.0/137/limit-order/count?page=1&limit=69&sortby=createDateTime",
	}

	for input, expected := range tests {
		u, err := url.Parse(enums.URL_BASE)
		if err != nil {
			t.Fatalf("failed to parse URL: %v\n", err.Error())
		}
		u.Path = path.Join(u.Path, getPath(input.chain, input.endpoint))
		page := getPages(input.count, input.limit)
		u.RawQuery = getQuery(page, input.limit)
		actual := u.String()
		if actual != expected {
			t.Logf(
				"Expected: %s\nActual: %s\n",
				expected,
				actual,
			)
			t.Fatal("results not matched")
		}
	}
}
