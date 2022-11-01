package domain

import (
	"regexp"
	"time"

	"github.com/desmos-labs/djuno/v2/types"
)

var (
	httpRegEx = regexp.MustCompile("^https?://")
)

// --------------------------------------------------------------------------------------------------------------------

var (
	_ types.ProfileScoreDetails = &ScoreDetails{}
)

type ScoreDetails struct {
	CreatedAt time.Time `json:"created_at"`
}

func NewScoreDetails(createdAt time.Time) *ScoreDetails {
	return &ScoreDetails{
		CreatedAt: createdAt,
	}
}

// GetScore implements types.ProfileScoreDetails
func (d *ScoreDetails) GetScore() (score uint64) {
	// Base of 25 points
	score += 25

	domainAge := time.Since(d.CreatedAt)
	domainAgeYrs := domainAge.Nanoseconds() / types.Year.Nanoseconds()

	// 25 points for accounts older than 1 year
	if domainAgeYrs > 1 {
		score = 100
	}

	return score
}
