package httpService

import (
	"encoding/json"
	"net/http"
	"net/url"
	http2 "rastreamento_gorreios/models/http"
	"time"
)

var client = &http.Client{
	Timeout: 15 * time.Second,
}

const (
	// Defined from API
	userData  = "teste"
	tokenData = "1abcd00b2731640e886fb41a8a9671ad1434c599dbaa0a0de9a5aa619f29a83f"
	urlAPI    = "https://api.linketrack.com/track/json"
)

// MakeRequest is responsible to make the request and get track data from Gorreios API
func MakeRequest(t *http2.TrackRequest) (tRes http2.TrackResponse, err error) {
	urlReq, err := url.Parse(urlAPI)
	if err != nil {
		return
	}

	urlReq.RawQuery = url.Values{
		"user":   []string{userData},
		"token":  []string{tokenData},
		"codigo": []string{t.Code},
	}.Encode()

	res, err := client.Get(urlReq.String())
	if err != nil {
		return
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&tRes)
	if err != nil {
		return
	}

	return
}
