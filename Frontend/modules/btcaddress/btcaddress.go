// Package btcaddress makes sure the Bitcoin address for the Ripple donations
// account on Keybase is up-to-date.
// This shows a default key if the system is not set up, or otherwise refreshes
// the address every 30 minutes (using redis to hold the cache).
// The wallet with the name "Ripple" will be used for getting the address.
package btcaddress

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"strings"

	"gopkg.in/redis.v5"
)

// DefaultAddress is the address used if the package is not set up to refresh
// from coinbase.
const DefaultAddress = "1CKGzZqrVwKoXwHEWpTobWqqtkYuqoFNro"

// AccountName is the name the account on coinbase must have to be considered as
// the actual account.
const AccountName = "Ripple"

// Configuration variables.
var (
	Redis     *redis.Client
	APIKey    string
	APISecret string
)

// Get retrieves the Bitcoin address, using the following methods:
//
// First the address is requested from Redis, and if the key is available from
// redis then that is used.
//
// After that, the key is requested from coinbase, and the account ID is found
// if not already saved in Redis.
func Get() string {
	v := Redis.Get("hanayo:btcaddress").Val()
	if v != "" {
		return v
	}
	if APIKey == "" || APISecret == "" {
		return DefaultAddress
	}

	a, err := getFromCoinbase()
	if err != nil {
		fmt.Println(err)
		return DefaultAddress
	}

	Redis.Set("hanayo:btcaddress", a, time.Minute*30)

	return a
}

type account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type accountsResp struct {
	Data []account `json:"data"`
}

type address struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type createAddressResp struct {
	Data address `json:"data"`
}

type errors struct {
	Errors []struct {
		ID      string `json:"id"`
		Message string `json:"message"`
	} `json:"errors"`
}

const coinbaseAPIDate = "2017-01-23"
const coinbaseAPIBase = "https://api.coinbase.com/v2/"

var accountID string

func getFromCoinbase() (string, error) {
	if accountID == "" {
		var accs accountsResp
		err := req("accounts", nil, &accs)
		if err != nil {
			return "", err
		}
		for _, el := range accs.Data {
			if el.Name == AccountName {
				accountID = el.ID
				break
			}
		}
	}

	var addrResp createAddressResp
	err := req("accounts/"+accountID+"/addresses", struct{}{}, &addrResp)
	if err != nil {
		return "", err
	}

	return addrResp.Data.Address, nil
}

func req(endpoint string, data interface{}, out interface{}) error {
	var dataJSON string
	method := "GET"
	if data != nil {
		method = "POST"
		// sort of a dirty hack, but oh well
		djr, err := json.Marshal(data)
		if err != nil {
			return err
		}
		dataJSON = string(djr)
	}

	req, err := http.NewRequest(method, coinbaseAPIBase+endpoint, strings.NewReader(dataJSON))
	if err != nil {
		return err
	}

	if data != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	// Not usind .Header.Add because we want to make sure the names are kept
	// uppercase.
	ts := time.Now().Unix()
	req.Header["CB-VERSION"] = []string{coinbaseAPIDate}
	req.Header["CB-ACCESS-KEY"] = []string{APIKey}
	req.Header["CB-ACCESS-TIMESTAMP"] = []string{fmt.Sprintf("%d", ts)}

	sig := fmt.Sprintf(
		"%d%s/v2/%s%s",
		ts,
		method,
		endpoint,
		string(dataJSON),
	)

	hm := hmac.New(
		sha256.New,
		[]byte(APISecret),
	)
	hm.Write([]byte(sig))

	req.Header["CB-ACCESS-SIGN"] = []string{
		fmt.Sprintf("%x", hm.Sum(nil)),
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		var errs errors
		json.NewDecoder(resp.Body).Decode(&errs)
		return fmt.Errorf("btcaddress: response status code is %d - %v", resp.StatusCode, errs)
	}

	err = json.NewDecoder(resp.Body).Decode(&out)
	return err
}
