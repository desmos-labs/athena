package github

import (
	"context"
	"time"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/forbole/juno/v3/types/config"
	"github.com/google/go-github/v48/github"
	"github.com/rs/zerolog/log"

	"github.com/desmos-labs/djuno/v2/types"
	profilesscore "github.com/desmos-labs/djuno/v2/x/profiles-score"

	"net/http"
	"strings"
)

var (
	_ profilesscore.Scorer = &Scorer{}
)

type Scorer struct {
	client *github.Client
}

// NewScorer returns a new Scorer instance
func NewScorer(junoCfg config.Config) *Scorer {
	cfgBz, err := junoCfg.GetBytes()
	if err != nil {
		panic(err)
	}
	cfg, err := UnmarshalConfig(cfgBz)
	if err != nil {
		panic(err)
	}

	if cfg == nil {
		log.Info().Str("scorer", "github").Msg("no config set, skipping creation")
		return nil
	}

	// Wrap the shared transport for use with the integration ID 1 authenticating with installation ID 99.
	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, cfg.AppID, cfg.InstallationID, cfg.PrivateKeyFilePath)
	if err != nil {
		panic(err)
	}

	return &Scorer{
		client: github.NewClient(&http.Client{Transport: itr}),
	}
}

// GetRateLimit implements Scorer
func (s *Scorer) GetRateLimit() *profilesscore.ScoreRateLimit {
	return profilesscore.NewScoreRateLimit(time.Hour, 5000)
}

// GetScoreDetails implements Scorer
func (s *Scorer) GetScoreDetails(_ string, application string, username string) (types.ProfileScoreDetails, error) {
	if !strings.EqualFold(application, "github") {
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

	return NewScoreDetails(
		user.GetCreatedAt().UTC(),
		uint64(user.GetFollowers()),
		uint64(user.GetFollowing()),
		uint64(user.GetPublicRepos()),
		uint64(user.GetPrivateGists()),
	), nil
}

// GetUser returns the GitHub data of the user having the given username
func (s *Scorer) GetUser(username string) (*github.User, error) {
	user, res, err := s.client.Users.Get(context.Background(), username)
	if res.Response.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	return user, err
}
