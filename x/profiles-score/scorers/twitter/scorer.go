package twitter

import (
	"context"
	"fmt"
	"github.com/desmos-labs/djuno/v2/x/profiles-score/scorers"
	"github.com/forbole/juno/v3/types/config"
	"github.com/g8rswimmer/go-twitter/v2"
	"net/http"
)

var (
	_ scorers.Scorer = &Scorer{}
)

// Scorer represents a scorers.Scorer instance to score profiles based on their Twitter statistics
type Scorer struct {
	client *twitter.Client
}

// NewScorer returns a new Scorer instance
func NewScorer(junoCfg *config.Config) *Scorer {
	cfgBz, err := junoCfg.GetBytes()
	if err != nil {
		panic(err)
	}
	cfg, err := ParseConfig(cfgBz)
	if err != nil {
		panic(err)
	}

	return &Scorer{
		client: &twitter.Client{
			Authorizer: NewAuthorizer(cfg.Token),
			Client:     http.DefaultClient,
			Host:       "https://api.twitter.com",
		},
	}
}

// SupportedApplications implements scorers.Scorer
func (s *Scorer) SupportedApplications() []string {
	return []string{"twitter"}
}

// RefreshScore implements scorers.Scorer
func (s *Scorer) RefreshScore(address string, username string, application string) error {
	//TODO implement me
	panic("implement me")
}

// GetUser returns the Twitter id of the user having the given username
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
