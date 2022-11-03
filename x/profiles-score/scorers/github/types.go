package github

import (
	"time"

	"github.com/desmos-labs/djuno/v2/x/profiles-score/scorers/utils"

	"github.com/desmos-labs/djuno/v2/types"
)

var (
	_ types.ProfileScoreDetails = &ScoreDetails{}
)

type ScoreDetails struct {
	CreatedAt               time.Time `json:"created_at"`
	FollowersCount          uint64    `json:"followers_count"`
	FollowingCount          uint64    `json:"following_count"`
	PublicRepositoriesCount uint64    `json:"public_repositories_count"`
	PublicGistsCount        uint64    `json:"public_gists_count"`
}

func NewScoreDetails(createdAt time.Time, followersCount, followingCount, publicReposCount, publicGistsCount uint64) *ScoreDetails {
	return &ScoreDetails{
		CreatedAt:               createdAt,
		FollowersCount:          followersCount,
		FollowingCount:          followingCount,
		PublicRepositoriesCount: publicReposCount,
		PublicGistsCount:        publicGistsCount,
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

	// 40 points for accounts with more than 100 followers, or with a followers:following ratio over 1
	isFamousAccount := d.FollowersCount > 100
	hasHigherThanAverageFollowersFollowingRatio := d.FollowingCount > 0 && d.FollowersCount/d.FollowingCount > 1
	if isFamousAccount || hasHigherThanAverageFollowersFollowingRatio {
		score += 40
	}

	// 10 points for accounts with more than 15 repositories
	if d.PublicRepositoriesCount > 15 {
		score += 10
	}

	return score
}
