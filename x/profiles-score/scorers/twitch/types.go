package twitch

import (
	"time"

	"github.com/desmos-labs/djuno/v2/x/profiles-score/scorers/utils"

	"github.com/desmos-labs/djuno/v2/types"
)

var (
	_ types.ProfileScoreDetails = &ScoreDetails{}
)

type ScoreDetails struct {
	CreatedAt       time.Time `json:"created_at"`
	BroadcasterType string    `json:"broadcaster_type"`
}

func NewScoreDetails(createdAt time.Time, broadcasterType string) *ScoreDetails {
	return &ScoreDetails{
		CreatedAt:       createdAt,
		BroadcasterType: broadcasterType,
	}
}

// GetScore implements types.ProfileScoreDetails
func (d *ScoreDetails) GetScore() (score uint64) {
	// Base of 25 points
	score += 25

	// 25 points for accounts older than 1 year
	if utils.GetTimeSinceInYears(d.CreatedAt) > 1 {
		score += 25
	}

	// Max points for accounts that are either affiliate or partners
	if d.BroadcasterType != "" {
		score = 100
	}

	return score
}
