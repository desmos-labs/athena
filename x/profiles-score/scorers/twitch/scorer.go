package twitch

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/forbole/juno/v5/types/config"
	"github.com/nicklaw5/helix"

	"github.com/desmos-labs/athena/v2/types"
	profilesscore "github.com/desmos-labs/athena/v2/x/profiles-score"
)

var (
	_ profilesscore.Scorer = &Scorer{}
)

type Scorer struct {
	client *helix.Client
}

// NewScorer returns a new Scorer instance
func NewScorer(junoCfg config.Config) profilesscore.Scorer {
	cfgBz, err := junoCfg.GetBytes()
	if err != nil {
		panic(err)
	}

	cfg, err := UnmarshalConfig(cfgBz)
	if err != nil {
		panic(err)
	}

	if cfg == nil {
		return nil
	}

	client, err := helix.NewClient(&helix.Options{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
	})
	if err != nil {
		panic(err)
	}

	// Set the app access token
	token, err := client.RequestAppAccessToken(nil)
	if err != nil {
		panic(err)
	}
	client.SetAppAccessToken(token.Data.AccessToken)

	return &Scorer{
		client: client,
	}
}

// GetRateLimit implements Scorer
func (s *Scorer) GetRateLimit() *profilesscore.ScoreRateLimit {
	return profilesscore.NewScoreRateLimit(time.Minute, 800)
}

// GetScoreDetails implements Scorer
func (s *Scorer) GetScoreDetails(_ string, application string, username string) (types.ProfileScoreDetails, error) {
	if !strings.EqualFold(application, "twitch") {
		return nil, nil
	}

	user, err := s.GetUser(username)
	if err != nil {
		return nil, err
	}

	return NewScoreDetails(user.CreatedAt.Time, user.BroadcasterType), nil
}

// GetUser returns the information for the given user
func (s *Scorer) GetUser(username string) (*helix.User, error) {
	res, err := s.client.GetUsers(&helix.UsersParams{
		Logins: []string{username},
	})
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(res.Error)
	}

	if len(res.Data.Users) == 0 {
		return nil, nil
	}

	return &res.Data.Users[0], nil
}
