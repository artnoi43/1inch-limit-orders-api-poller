package get

import (
	"fmt"

	"github.com/artnoi43/1inch-limit-orders-api-poller/enums"
)

func getPath(chain enums.Chain, sw string) string {
	id, ok := enums.ChainMap[chain]
	if !ok {
		fmt.Println("Invalid chain", chain)
		// panic("invalid chain")
	}
	var s string
	switch sw {
	case "all":
		s = fmt.Sprintf(enums.PATH_ALL, id)
	case "count":
		s = fmt.Sprintf(enums.PATH_COUNT, id)
	default:
		panic("unexpected switch")
	}
	return s
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
