package types

import (
	"time"
)

// This struct imitates 1st-level JSON key
type LimitOrder struct {
	Signature            string         `json:"signature"`
	OrderHash            string         `json:"orderHash"`
	CreateDateTimeStr    string         `json:"createDateTime"`
	CreateDateTime       time.Time      `json:"createDateTimeTime"`
	RemainingMakerAmount string         `json:"remainingMakingOrderAmount"`
	MakerBalance         string         `json:"makerBalance"`
	MakerAllowance       string         `json:"makerAllowance"`
	MakerRate            string         `json:"makerRate"`
	TakerRate            string         `json:"takerRate"`
	Data                 LimitOrderData `json:"data"`
	// ResponseTime may be relevant when determining which data should we use
	ResponseTime time.Time `json:"responseTime"`
}

type LimitOrderData struct {
	Maker          string `json:"taker"`
	MakingAmount   string `json:"makingAmount"`
	TakingAmount   string `json:"takingAmount"`
	MakerAsset     string `json:"makerAsset"`
	TakerAsset     string `json:"takerAsset"`
	GetMakerAmount string `json:"getMakerAmount"`
	GetTakerAmount string `json:"getTakerAmount"`
	MakerAssetData string `json:"makerAssetData"`
	TakerAssetData string `json:"takerAssetData"`
	Salt           string `json:"salt"`
	Permit         string `json:"permit"`
	Predicate      string `json:"predicate"`
	Interaction    string `json:"interaction"`
	Receiver       string `json:"receiver"`
	AllowedSender  string `json:"allowedSender"`
}

// ParseTime parses CreateDateTimeString to time.Time and save it to CreateDateTime
func (o *LimitOrder) ParseTime() error {
	t, err := time.Parse(time.RFC3339, o.CreateDateTimeStr)
	if err != nil {
		return err
	}
	o.CreateDateTime = t
	return nil
}
