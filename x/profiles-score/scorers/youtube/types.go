package youtube

import (
	"fmt"
	"net/http"
	"time"

	"github.com/desmos-labs/athena/x/profiles-score/scorers/utils"

	"github.com/g8rswimmer/go-twitter/v2"

	"github.com/desmos-labs/athena/types"
)

var (
	_ twitter.Authorizer = &Authorizer{}
)

// Authorizer implements twitter.Authorizer to authorize Twitter requests
type Authorizer struct {
	Token string
}

// NewAuthorizer implements a new Authorizer instance
func NewAuthorizer(token string) *Authorizer {
	return &Authorizer{
		Token: token,
	}
}

// Add implements twitter.Authorizer
func (a *Authorizer) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

// --------------------------------------------------------------------------------------------------------------------

var (
	_ types.ProfileScoreDetails = &ScoreDetails{}
)

type ScoreDetails struct {
	LinkedAt         *time.Time `json:"linked_at,omitempty"`
	SubscribersCount uint64     `json:"followers_count"`
}

func NewScoreDetails(linkedAt *time.Time, subscribersCount uint64) *ScoreDetails {
	return &ScoreDetails{
		LinkedAt:         linkedAt,
		SubscribersCount: subscribersCount,
	}
}

// GetScore implements types.ProfileScoreDetails
func (d *ScoreDetails) GetScore() (score uint64) {
	// 84.2% of YouTube channels have less than 1,000 subscribers
	if d.SubscribersCount > 1000 {
		return 100
	}

	// Base of 25 points
	score += 25

	// 25 points if public
	if d.LinkedAt != nil {
		score += 25

		// 25 points if older than one year
		if utils.GetTimeSinceInYears(*d.LinkedAt) > 1 {
			score += 25
		}
	}

	return score
}
