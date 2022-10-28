package twitter

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/forbole/juno/v3/types/config"
	"github.com/g8rswimmer/go-twitter/v2"

	"github.com/desmos-labs/djuno/v2/types"
	profilesscore "github.com/desmos-labs/djuno/v2/x/profiles-score"
)

var (
	_ profilesscore.Scorer = &Scorer{}
)

// Scorer represents a scorers.Scorer instance to score profiles based on their Twitter statistics
type Scorer struct {
	client *twitter.Client
}

// NewScorer returns a new Scorer instance
func NewScorer(junoCfg config.Config) *Scorer {
	cfgBz, err := junoCfg.GetBytes()
	if err != nil {
		panic(err)
	}
	cfg, err := ParseConfig(cfgBz)
	if err != nil {
		panic(err)
	}

	if cfg == nil {
		log.Info().Str("scorer", "twitter").Msg("no config set, skipping creation")
		return nil
	}

	return &Scorer{
		client: &twitter.Client{
			Authorizer: NewAuthorizer(cfg.Token),
			Client:     http.DefaultClient,
			Host:       "https://api.twitter.com",
		},
	}
}

// GetRateLimit implements Scorer
func (s *Scorer) GetRateLimit() *profilesscore.ScoreRateLimit {
	return profilesscore.NewScoreRateLimit(time.Minute*15, 900)
}

// GetScoreDetails implements Scorer
func (s *Scorer) GetScoreDetails(_ string, application string, username string) (types.ProfileScoreDetails, error) {
	if !strings.EqualFold(application, "twitter") {
		return nil, nil
	}

	user, err := s.GetUser(username)
	if err != nil {
		return nil, err
	}

	// Make sure the user exists
	if user == nil {
		return nil, nil
	}

	createdAt, err := time.Parse(time.RFC3339, user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return NewScoreDetails(
		createdAt,
		uint64(user.PublicMetrics.Followers),
		uint64(user.PublicMetrics.Following),
		uint64(user.PublicMetrics.Tweets),
		user.Verified,
	), nil
}

// GetUser returns the Twitter data of the user having the given username
func (s *Scorer) GetUser(username string) (*twitter.UserObj, error) {
	// Get the user details from the username
	res, err := s.client.UserNameLookup(context.Background(), []string{username}, twitter.UserLookupOpts{
		UserFields: []twitter.UserField{
			twitter.UserFieldID,
			twitter.UserFieldCreatedAt,
			twitter.UserFieldPublicMetrics,
			twitter.UserFieldVerified,
		},
	})
	if err != nil {
		return nil, err
	}

	if res.Raw.Errors != nil && hasNotFoundError(res.Raw.Errors) {
		return nil, fmt.Errorf("invalid Twitter username: %s", username)
	}

	return res.Raw.Users[0], nil
}
