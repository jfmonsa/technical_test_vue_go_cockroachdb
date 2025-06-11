package main

import (
	"strings"
	"time"
	"vue_go_cockroachdb/src/models"
)

func CalculateStockScore(s models.Stock) float64 {
	score := 0.0

	// 1. Profit potential
	if s.TargetFrom > 0 {
		potential := ((s.TargetTo - s.TargetFrom) / s.TargetFrom) * 100
		if potential > 0 {
			score += potential / 4 // More weight to upside
		} else {
			score += potential / 10 // Penalize downside, but less
		}
	}

	// 2. Recommended action
	action := strings.ToLower(s.Action)
	for _, a := range actionScore {
		if strings.Contains(action, a.Keyword) {
			score += a.Score
			break // only applies the first match
		}
	}

	// 3. Rating change
	from := normalizeRating(s.RatingFrom)
	to := normalizeRating(s.RatingTo)
	if from > 0 && to > 0 {
		diff := to - from
		score += diff * 2 // Weighs the rating change
	}

	// 4. Recent (more weight if it is from the last 3 days)
	var daysAgo float64
	if t, err := time.Parse(time.RFC3339, s.Time); err == nil {
		daysAgo = time.Since(t).Hours() / 24
	} else {
		daysAgo = 999 // If it cannot be parsed, it is assumed to be very old
	}
	if daysAgo < 1 {
		score += 1.5
	} else if daysAgo < 3 {
		score += 1
	} else if daysAgo < 7 {
		score += 0.5
	}

	return score
}

var ratingScore = []struct {
	Rating string
	Score  float64
}{
	{models.RatingNeutral, 5},
	{models.RatingUnchanged, 5},
	{models.RatingEqualWeight, 5},
	{models.RatingOutperform, 8},
	{models.RatingMarketPerform, 5},
	{models.RatingInLine, 5},
	{models.RatingHold, 4},
	{models.RatingBuy, 9},
	{models.RatingOverweight, 8},
	{models.RatingPositive, 8},
	{models.RatingMarketOutperform, 8},
	{models.RatingSectorOutperform, 8},
	{models.RatingStrongBuy, 10},
	{models.RatingSectorPerform, 5},
	{models.RatingUnderweight, 3},
	{models.RatingSell, 1},
	{models.RatingSpeculativeBuy, 7},
	{models.RatingSectorWeight, 5},
	{models.RatingOutperformer, 8},
	{models.RatingUnderperform, 2},
	{models.RatingPeerPerform, 5},
	{models.RatingSectorUnderperform, 2},
	{models.RatingAccumulate, 7},
	{models.RatingTopPick, 10},
	{models.RatingReduce, 2},
}

var actionScore = []struct {
	Keyword string
	Score   float64
}{
	{models.ActionUpgraded, 2},
	{models.ActionDowngraded, -2},
	{models.ActionInitiated, 1},
	{models.ActionReiterated, 0.5},
	{models.ActionTargetRaised, 1},
	{models.ActionTargetLowered, -1},
	{models.ActionTargetSet, 0.2},
}

// Assuming that ratings are always valid and normalized
// This function normalizes the rating to a numeric score.
func normalizeRating(rating string) float64 {
	for _, r := range ratingScore {
		if strings.EqualFold(r.Rating, rating) {
			return r.Score
		}
	}
	return 0
}
