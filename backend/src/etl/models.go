package main

import "vue_go_cockroachdb/src/models"

// APIResponse models the response from the external API containing
// a list of action recommendations and a pagination token.
//
// Example response:
//
//	{
//	  "items": [ ... ], // list of recommendations (APIRawItem)
//	  "next_page": "TRIN" // token for the next page of results
//	}
type APIResponse struct {
	Items    []APIRawItem `json:"items"`     // List of action recommendations.
	NextPage string       `json:"next_page"` // Token to get the next page of results.
}

// APIRawItem represents an action recommendation as returned by the external API.
// All fields are received as strings.
//
// Example item:
//
//	{
//	  "ticker": "MOMO",
//	  "company": "Hello Group",
//	  "brokerage": "Benchmark",
//	  "action": "reiterated by",
//	  "rating_from": "Buy",
//	  "rating_to": "Buy",
//	  "target_from": "$13.00",
//	  "target_to": "$13.00",
//	  "time": "2025-03-14T00:30:05.974622332Z"
//	}
//

type APIRawItem struct {
	Ticker     string `json:"ticker"`      // Stock symbol (e.g., "MOMO").
	Company    string `json:"company"`     // Company name (e.g., "Hello Group").
	Brokerage  string `json:"brokerage"`   // Name of the brokerage issuing the recommendation.
	Action     string `json:"action"`      // Action performed (e.g., "upgraded by", "reiterated by", etc).
	RatingFrom string `json:"rating_from"` // Previous rating (e.g., "Buy", "Neutral", etc).
	RatingTo   string `json:"rating_to"`   // New rating (e.g., "Buy", "Outperform", etc).
	TargetFrom string `json:"target_from"` // Previous target price (e.g., "$13.00").
	TargetTo   string `json:"target_to"`   // New target price (e.g., "$13.00").
	Time       string `json:"time"`        // Recommendation date and time in RFC3339 format.
}

// Verify if the response contains a valid rating in order to be inserted into the database.
func isValidRating(rating string) bool {
	return rating == models.RatingNeutral ||
		rating == models.RatingUnchanged ||
		rating == models.RatingEqualWeight ||
		rating == models.RatingOutperform ||
		rating == models.RatingMarketPerform ||
		rating == models.RatingInLine ||
		rating == models.RatingHold ||
		rating == models.RatingBuy ||
		rating == models.RatingOverweight ||
		rating == models.RatingPositive ||
		rating == models.RatingMarketOutperform ||
		rating == models.RatingSectorOutperform ||
		rating == models.RatingStrongBuy ||
		rating == models.RatingSectorPerform ||
		rating == models.RatingUnderweight ||
		rating == models.RatingSell ||
		rating == models.RatingSpeculativeBuy ||
		rating == models.RatingSectorWeight ||
		rating == models.RatingOutperformer ||
		rating == models.RatingUnderperform ||
		rating == models.RatingPeerPerform ||
		rating == models.RatingSectorUnderperform ||
		rating == models.RatingAccumulate ||
		rating == models.RatingTopPick ||
		rating == models.RatingReduce
}
