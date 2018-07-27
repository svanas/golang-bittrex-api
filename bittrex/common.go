package bittrex

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// BaseURL represents the base URL for all requests
const (
	BaseURL = "https://bittrex.com/Api/v2.0"
)

// Auth represents the auth credentials to authenticate to the Bittrex API:
//
// It consists of a set of a private and a public key.
type Auth struct {
	ApiKey    string // The public key to connect to bittrex API.
	ApiSecret string // The private key to connect to bittrex API.
}

var (
	// defaultClient represents the default configuration for HTTP requests to the API.
	defaultClient = http.Client{
		Timeout: time.Second * 30,
	}
	// client represents the actual configuration for HTTP requests to the API.
	client = defaultClient
)

// SetCustomHTTPClient sets a custom client for requests.
func SetCustomHTTPClient(value http.Client) {
	client = value
}

// apiCall performs a generic API call.
func apiCall(req *http.Request) (*json.RawMessage, error) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Add("Cache-Control", "no-store")
	req.Header.Add("Cache-Control", "must-revalidate")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Status Code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ret response
	err = json.Unmarshal(content, &ret)
	if err != nil {
		return nil, err
	}

	if ret.Success == false {
		return nil, fmt.Errorf("Error Response: %s", ret.Message)
	}

	return ret.Result, nil
}

// publicCall performs a call to the public bittrex API.
//
// It does not need API Keys.
func publicCall(group, command string, params map[string]string) (*json.RawMessage, error) {
	url := fmt.Sprintf("%s/pub/%s/%s", BaseURL, group, command)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if params != nil {
		qry := req.URL.Query()
		for key, value := range params {
			qry.Set(key, value)
		}
		req.URL.RawQuery = qry.Encode()
	}

	return apiCall(req)
}

// authCall performs a call to the private bittrex API.
//
// It needs an Auth struct to be passed with valid Keys.
func authCall(group, command string, params map[string]string, auth *Auth) (*json.RawMessage, error) {
	if auth.ApiKey == "" || auth.ApiSecret == "" {
		return nil, errors.New("Cannot perform private api request without authentication keys")
	}

	url := fmt.Sprintf("%s/key/%s/%s", BaseURL, group, command)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	nonce := time.Now().UnixNano()

	qry := req.URL.Query()
	if params != nil { // Add them to query string
		for key, value := range params {
			qry.Set(key, value)
		}
	}
	qry.Set("apikey", auth.ApiKey)
	qry.Set("nonce", fmt.Sprintf("%d", nonce))
	req.URL.RawQuery = qry.Encode()

	mac := hmac.New(sha512.New, []byte(auth.ApiSecret))
	_, err = mac.Write([]byte(req.URL.String()))
	sig := hex.EncodeToString(mac.Sum(nil))
	req.Header.Add("apisign", sig)

	return apiCall(req)
}
