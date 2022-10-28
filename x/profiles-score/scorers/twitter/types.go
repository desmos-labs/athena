package twitter

import (
	"fmt"
	"net/http"
	"time"

	"github.com/g8rswimmer/go-twitter/v2"

	"github.com/desmos-labs/djuno/v2/types"
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
	CreatedAt      time.Time `json:"created_at"`
	FollowersCount uint64    `json:"followers_count"`
	FollowingCount uint64    `json:"following_count"`
	TweetsCount    uint64    `json:"tweets_count"`
	Verified       bool      `json:"verified"`
}

func NewScoreDetails(createdAt time.Time, followersCount, followingCount, tweetsCount uint64, verified bool) *ScoreDetails {
	return &ScoreDetails{
		CreatedAt:      createdAt,
		FollowersCount: followersCount,
		FollowingCount: followingCount,
		TweetsCount:    tweetsCount,
		Verified:       verified,
	}
}

// GetScore implements types.ProfileScoreDetails
func (d *ScoreDetails) GetScore() (score uint64) {
	if d.Verified {
		return 100
	}

	accountAge := time.Since(d.CreatedAt)
	accountAgeYrs := uint64(accountAge.Nanoseconds() / types.Year.Nanoseconds())

	// Base of 25 points
	score += 25

	// 25 points for accounts older than 1 year
	if accountAgeYrs > 1 {
		score += 25
	}

	isFamousAccount := d.FollowersCount > 1000
	hasHigherThanAverageFollowersFollowingRatio := d.FollowingCount > 0 && d.FollowersCount/d.FollowingCount > 1
	if isFamousAccount || hasHigherThanAverageFollowersFollowingRatio {
		score += 40
	}

	// Average user performs ~960 tweets per year
	hasHigherThanAverageTweetsPerYear := accountAgeYrs > 0 && d.TweetsCount/accountAgeYrs > 900
	if hasHigherThanAverageTweetsPerYear {
		score += 10
	}

	return score
}
