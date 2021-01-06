package http

import "time"

// TrackRequest is the base struct to storage data to do the request in Gorreios endpoint
type TrackRequest struct {
	Code string
}

// TrackResponse is the base struct to storage API response from Gorreios
type TrackResponse struct {
	Codigo     string    `json:"codigo"`
	Host       string    `json:"host"`
	Time       float64   `json:"time"`
	Quantidade int       `json:"quantidade"`
	Servico    string    `json:"servico"`
	Ultimo     time.Time `json:"ultimo"`
	Eventos    []Events  `json:"eventos"`
}

// Events is the base struct to storage all events from TrackResponse
type Events struct {
	Data      string   `json:"data"`
	Hora      string   `json:"hora"`
	Local     string   `json:"local"`
	Status    string   `json:"status"`
	SubStatus []string `json:"subStatus"`
}
