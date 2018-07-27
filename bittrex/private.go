package bittrex

import (
	"encoding/json"
	"strconv"
	"time"
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
