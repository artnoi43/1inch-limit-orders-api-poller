package get

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/artnoi43/1inch-limit-orders-api-poller/enums"
	"github.com/artnoi43/1inch-limit-orders-api-poller/lib/types"
	"github.com/pkg/errors"
)

// GetAndParse fetches data from the HTTP endpoint,
// and parses response's JSON body into the interface,
// and lastly calls afterGet.
func getAndParse(u string, v interface{}) error {
	resp, err := http.Get(u)
	if err != nil {
		return errors.Wrap(
			err,
			fmt.Sprintf(
				"failed to get url: %s",
				u,
			),
		)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusTooManyRequests:
		return enums.ErrRateLimitExceeded
	default:
		return errors.New("non-200 status code")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(
			err,
			fmt.Sprintf(
				"failed to read body: %s",
				resp.Body,
			),
		)
	}
	if err := json.Unmarshal(body, &v); err != nil {
		return errors.Wrap(
			err,
			fmt.Sprintf(
				"failed to unmarshal: %s",
				body,
			),
		)
	}
	if err := afterGet(resp, v); err != nil {
		return errors.Wrap(err, "afterGet failed")
	}
	return nil
}

// afterGet is called by GetAndParse after it's done with unmarshaling.
// afterGet uses a type switch to determine what to do with the data,
// since values from GetAndParse can be of any types.
func afterGet(resp *http.Response, v interface{}) error {
	switch v := v.(type) {
	// Add timestamps (response) to limit orders
	case *[]*types.LimitOrder:
		orders := *v
		for _, order := range orders {
			t, err := time.Parse(
				time.RFC1123,
				resp.Header.Get("Date"),
			)
			if err != nil {
				return err
			}
			// TZ
			location, _ := time.LoadLocation("Asia/Bangkok")
			order.ResponseTime = t.In(location)
		}
	}
	return nil
}

// GetCount uses GetAndParse to get 1inch limit order count.
// It is called by PollCount.
func GetCount(chain enums.Chain) (int, error) {
	// Prepare URL
	u, _ := url.Parse(enums.URL_BASE)
	u.Path = path.Join(
		u.Path, getPath(chain, "count"),
	)
	// Get and parse count
	var count struct {
		Count int `json:"count"`
	}
	if err := getAndParse(
		u.String(),
		&count,
	); err != nil {
		return 0, err
	}

	return count.Count, nil
}

// GetPage deals with 1inch pagination
// by modifying query string based pn page and limit.
// It also sends orders received to orderChan channel.
func getPage(
	chain enums.Chain,
	page int,
	limit int,
	orderChan chan<- []*types.LimitOrder,
) error {
	u, err := url.Parse(enums.URL_BASE)
	if err != nil {
		return err
	}
	// Prepare URLs
	u.Path = path.Join(
		u.Path, getPath(chain, "all"),
	)
	u.RawQuery = getQuery(
		page, limit,
	)
	// Get and parse, then send orders
	var orders []*types.LimitOrder
	if err := getAndParse(
		u.String(),
		&orders,
	); err != nil {
		return err
	}
	orderChan <- orders
	return nil
}
