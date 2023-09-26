package main

import (
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"
	"time"
	"unicode"
)

// x = SSN
// y = 10.7cm flux strength ($F_{10.7}$)
//
// Estimation function:
// y = 62.51 + 0.6422 * x
//
// From TK2018_F1 equation
// Table 2,
// Is the F10.7cm – Sunspot Number relation linear and stable?
// Frédéric Clette
// J. Space Weather Space Clim., 11 (2021) 2
// DOI: https://doi.org/10.1051/swsc/2020071
//
// Originally from:
// B. R. Tiwari and M. Kumar,
// The Solar Flux and Sunspot Number; A Long-Trend Analysis
// Int. Ann. Sci., vol. 5, no. 1, pp. 47–51, Jul. 2018.
// https://doi.org/10.21467/ias.5.1.47-51
//
// Note well: SSN should never be negative

func estimatedSSN(sfi float64) uint16 {
	essn := math.Round((50.0 * (100.0*sfi - 6251.0)) / 3211.0)
	if essn < 0 {
		return 0
	} else {
		return uint16(essn)
	}
}

// Get an HTTP(S) content from given URL

func getURLBody(url string) []byte {
	client := http.Client{
		// Set timeout to 10 seconds
		Timeout: time.Second * 10,
	}
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Print(err)
		return nil
	}
	request.Header.Set("User-Agent", "gossn")
	request.Close = true
	response, err := client.Do(request)
	if err != nil {
		log.Print(err)
		return nil
	}
	if response.Body != nil {
		defer response.Body.Close()
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
		return nil
	}
	return body
}

// Trim left spaces from a string

func trimLeftSpace(s string) string {
	return strings.TrimLeftFunc(s, unicode.IsSpace)
}
