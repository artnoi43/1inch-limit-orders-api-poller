package get

import (
	"fmt"

	"github.com/artnoi43/1inch-limit-orders-api-poller/enums"
)

func getPath(chain enums.Chain, sw string) string {
	return fmt.Sprintf(
		enums.EndpointMap[sw], enums.ChainMap[chain],
	)
}

// Page returns the number of pages required by pagination
func getPages(count, limit int) int {
	return (count + limit - 1) / limit
}

// Page query returns the query string based on page and limit
func getQuery(page, limit int) string {
	return fmt.Sprintf(
		"page=%d&limit=%d&sortby=createDateTime",
		page,
		limit,
	)
}
