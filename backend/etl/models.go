package main

import "time"

// StockItem represents a stock recommendation record
// after it has been transformed and is ready to be inserted into the database.
type StockItem struct {
	Ticker     string
	Company    string
	Brokerage  string
	Action     string
	RatingFrom string
	RatingTo   string
	TargetFrom float64
	TargetTo   float64
	Time       time.Time
}

// APIResponse models the structure of the API response
// containing a list of raw stock recommendation items and a pagination token.
type APIResponse struct {
	Items    []APIRawItem `json:"items"`
	NextPage string       `json:"next_page"`
}

// APIRawItem represents the raw stock recommendation data
// as received from the external API, with all fields as strings.
type APIRawItem struct {
	Ticker     string `json:"ticker"`
	Company    string `json:"company"`
	Brokerage  string `json:"brokerage"`
	Action     string `json:"action"`
	RatingFrom string `json:"rating_from"`
	RatingTo   string `json:"rating_to"`
	TargetFrom string `json:"target_from"`
	TargetTo   string `json:"target_to"`
	Time       string `json:"time"`
}
