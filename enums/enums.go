package enums

import (
	"errors"
)

type Chain string
type ChainId int

const (
	ChainEthereum = "ethereum"
	ChainBsc      = "bsc"
	ChainPolygon  = "polygon"
)

var ChainMap = map[Chain]ChainId{
	ChainEthereum: 1,
	ChainBsc:      56,
	ChainPolygon:  137,
}

const (
	URL_BASE   = "https://limit-orders.1inch.io"
	PATH_ALL   = "/v2.0/%d/limit-order/all"
	PATH_COUNT = "/v2.0/%d/limit-order/count"
)

var (
	Endpoints = struct {
		All   string
		Count string
	}{
		All:   "all",
		Count: "count",
	}
	EndpointMap = map[string]string{
		Endpoints.All:   PATH_ALL,
		Endpoints.Count: PATH_COUNT,
	}
)

var (
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)
