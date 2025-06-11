package models

// not stored in the database, but used to represent a stock recommendation
// for internal processing and transformation.
type Stock struct {
	Ticker     string  `json:"ticker"`
	Company    string  `json:"company"`
	Brokerage  string  `json:"brokerage"`
	Action     string  `json:"action"`
	RatingFrom string  `json:"rating_from"`
	RatingTo   string  `json:"rating_to"`
	TargetFrom float64 `json:"target_from"`
	TargetTo   float64 `json:"target_to"`
	Time       string  `json:"time"`
}

// Represents a stock recommendation with its details in the database (stocks table).
type StockWithScore struct {
	Stock
	RecommendationScore float64 `json:"recommendation_score,omitempty"`
}

// Constants for stock ratings to avoid magic strings in the code.
// These are the expected values for the `rating_from` and `rating_to` fields in the Stock model.
// Note: This values can be verified using: `SELECT DISTINCT rating_from FROM stocks;` against our db
const (
	RatingNeutral            = "Neutral"
	RatingUnchanged          = "Unchanged"
	RatingEqualWeight        = "Equal Weight"
	RatingOutperform         = "Outperform"
	RatingMarketPerform      = "Market Perform"
	RatingInLine             = "In-Line"
	RatingHold               = "Hold"
	RatingBuy                = "Buy"
	RatingOverweight         = "Overweight"
	RatingPositive           = "Positive"
	RatingMarketOutperform   = "Market Outperform"
	RatingSectorOutperform   = "Sector Outperform"
	RatingStrongBuy          = "Strong-Buy"
	RatingSectorPerform      = "Sector Perform"
	RatingUnderweight        = "Underweight"
	RatingSell               = "Sell"
	RatingSpeculativeBuy     = "Speculative Buy"
	RatingSectorWeight       = "Sector Weight"
	RatingOutperformer       = "Outperformer"
	RatingUnderperform       = "Underperform"
	RatingPeerPerform        = "Peer Perform"
	RatingSectorUnderperform = "Sector Underperform"
	RatingAccumulate         = "Accumulate"
	RatingTopPick            = "Top Pick"
	RatingReduce             = "Reduce"
)

// Constants for stock actions to avoid magic strings in the code.
// These are the expected values for the `action` field in the Stock model.
// Note: This values can be verified using: `SELECT DISTINCT action FROM stocks;` against our db
const (
	ActionUpgraded      = "upgraded"
	ActionDowngraded    = "downgraded"
	ActionInitiated     = "initiated"
	ActionReiterated    = "reiterated"
	ActionTargetRaised  = "target raised"
	ActionTargetLowered = "target lowered"
	ActionTargetSet     = "target set"
)
