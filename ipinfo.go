package ipinfo

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2026 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"errors"
	"fmt"
	"strings"

	"github.com/essentialkaos/ek/v14/req"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Info contains basic information about IP address
type Info struct {
	IP            string `json:"ip"`             // IP address being queried
	ASN           string `json:"asn"`            // Autonomous System Number
	ASName        string `json:"as_name"`        // Organization name
	ASDomain      string `json:"as_domain"`      // Organization's domain name
	Country       string `json:"country"`        // Full country name
	CountryCode   string `json:"country_code"`   // Two-letter country code (ISO 3166-1 alpha-2)
	Continent     string `json:"continent"`      // Full continent name
	ContinentCode string `json:"continent_code"` // Two-letter continent code
}

// ////////////////////////////////////////////////////////////////////////////////// //

// apiError is API error
type apiError struct {
	Error string `json:"error"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

// API is URL of IPInfo Lite API
var API = "https://api.ipinfo.io/lite"

// ////////////////////////////////////////////////////////////////////////////////// //

// Token is API token
var Token string

// ////////////////////////////////////////////////////////////////////////////////// //

var (
	ErrEmptyToken = errors.New("Token is empty")
	ErrEmptyIP    = errors.New("IP is empty")
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Get gets info about IP from IPInfo Lite API
func Get(ip string) (*Info, error) {
	switch {
	case Token == "":
		return nil, ErrEmptyToken
	case ip == "":
		return nil, ErrEmptyIP
	case !isLooksLikeIP(ip):
		return nil, fmt.Errorf("IP %q is not valid", ip)
	}

	resp, err := req.Request{
		URL:    API + "/" + ip,
		Query:  req.Query{"token": Token},
		Accept: req.CONTENT_TYPE_JSON,
	}.Get()

	if err != nil {
		return nil, fmt.Errorf("Can't send request to IPinfo API: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != req.STATUS_OK {
		errData := &apiError{}
		err = resp.JSON(errData)

		if err != nil {
			return nil, fmt.Errorf("IPinfo API returned non-ok status code (%d)", resp.StatusCode)
		}

		return nil, fmt.Errorf("IPinfo API returned error: %s", errData.Error)
	}

	info := &Info{}
	err = resp.JSON(info)

	if err != nil {
		return nil, fmt.Errorf("Can't decode IPinfo API response: %w", err)
	}

	return info, nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// isLooksLikeIP is very simple IP validation
func isLooksLikeIP(ip string) bool {
	switch {
	case strings.Count(ip, ".") != 3:
		return false
	case strings.Trim(ip, "01234567890.") != "":
		return false
	}

	return true
}
