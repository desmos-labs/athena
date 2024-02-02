package profilesscore

import (
	"time"

	"github.com/desmos-labs/athena/v2/types"
)

type ScoreRateLimit struct {
	Duration  time.Duration
	RateLimit uint64
}

func NewScoreRateLimit(duration time.Duration, rateLimit uint64) *ScoreRateLimit {
	return &ScoreRateLimit{
		Duration:  duration,
		RateLimit: rateLimit,
	}
}

// Scorer represents a generic parses that gets data from an external application and converts it to a specific score
type Scorer interface {
	// GetRateLimit returns the rate limit for this scorer, if any
	GetRateLimit() *ScoreRateLimit

	// GetScoreDetails returns the score details for the user having the given address
	// and username on the specified application
	GetScoreDetails(address string, application string, username string) (types.ProfileScoreDetails, error)
}

type Scorers []Scorer

func (s Scorers) GetRateLimit() *ScoreRateLimit {
	// Get the highest duration
	var highestDuration time.Duration
	for _, scorer := range s {
		scorerRateLimit := scorer.GetRateLimit()
		if scorerRateLimit == nil {
			continue
		}

		// Get the maximum time frame size
		if highestDuration == 0 || scorerRateLimit.Duration > highestDuration {
			highestDuration = scorerRateLimit.Duration
		}
	}

	// Get the lowest rate limit
	var lowestRateLimit uint64
	for _, scorer := range s {
		scorerRateLimit := scorer.GetRateLimit()
		if scorerRateLimit == nil {
			continue
		}

		// Convert the scorer rate to the highest duration
		scorerConvertedRate := scorerRateLimit.RateLimit * uint64(highestDuration.Nanoseconds()/scorerRateLimit.Duration.Nanoseconds())
		if lowestRateLimit == 0 || scorerConvertedRate < lowestRateLimit {
			lowestRateLimit = scorerConvertedRate
		}
	}

	return NewScoreRateLimit(highestDuration, lowestRateLimit)
}
