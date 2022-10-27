package profilesscore

import (
	"time"

	"github.com/desmos-labs/djuno/v2/types"
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
	var rateLimit *ScoreRateLimit
	for _, scorer := range s {
		scorerRateLimit := scorer.GetRateLimit()
		if rateLimit == nil {
			rateLimit = scorerRateLimit
			continue
		}
		if scorerRateLimit == nil {
			continue
		}

		// Get the minimum amount of request per time frame
		if rateLimit.RateLimit > scorerRateLimit.RateLimit {
			rateLimit.RateLimit = scorerRateLimit.RateLimit
		}

		// Get the maximum time frame size
		if rateLimit.Duration < scorerRateLimit.Duration {
			rateLimit.Duration = scorerRateLimit.Duration
		}
	}
	return rateLimit
}
