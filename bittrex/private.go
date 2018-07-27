package bittrex

import (
	"encoding/json"
	"strconv"
	"time"
)

const (
	ORDER_TYPE_LIMIT  = "LIMIT"
	ORDER_TYPE_MARKET = "MARKET"
)

const (
	GTC = "GOOD_TIL_CANCELLED"
	IOC = "IMMEDIATE_OR_CANCEL"
	FOK = "FILL_OR_KILL"
)

const (
	ConditionNone        = "NONE"
	LessThanOrEqualTo    = "LESS_THAN"
	GreaterThanOrEqualTo = "GREATER_THAN"
)

func GetBalance(auth *Auth, currencyName string) (*json.RawMessage, error) {
	now := time.Now().UnixNano()

	var params = map[string]string{
		"currencyName": currencyName,
		"_":            strconv.FormatInt(now, 10),
	}

	result, err := authCall("balance", "GetBalance", params, auth)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func TradeSell(auth *Auth, marketName string, orderType string, quantity float64, rate float64, timeInEffect string, conditionType string, target float64) (*Order, error) {
	now := time.Now().UnixNano()

	var params = map[string]string{
		"marketName":    marketName,
		"orderType":     orderType,
		"quantity":      strconv.FormatFloat(quantity, 'f', -1, 64),
		"rate":          strconv.FormatFloat(rate, 'f', -1, 64),
		"timeInEffect":  timeInEffect,
		"conditionType": conditionType,
		"target":        strconv.FormatFloat(target, 'f', -1, 64),
		"_":             strconv.FormatInt(now, 10),
	}

	raw, err := authCall("market", "TradeSell", params, auth)
	if err != nil {
		return nil, err
	}

	var result Order
	err = json.Unmarshal(*raw, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
