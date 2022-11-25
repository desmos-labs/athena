package youtube

import (
	"context"
	"strings"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"

	"github.com/rs/zerolog/log"

	"github.com/forbole/juno/v4/types/config"

	"github.com/desmos-labs/djuno/v2/types"
	profilesscore "github.com/desmos-labs/djuno/v2/x/profiles-score"
)

var (
	_ profilesscore.Scorer = &Scorer{}
)

// Scorer represents a scorers.Scorer instance to score profiles based on their Twitter statistics
type Scorer struct {
	client *youtube.Service
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
		log.Info().Str("scorer", "youtube").Msg("no config set, skipping creation")
		return nil
	}

	service, err := youtube.NewService(context.Background(), option.WithAPIKey(cfg.APIKey))
	if err != nil {
		panic(err)
	}

	return &Scorer{
		client: service,
	}
}

// GetRateLimit implements Scorer
func (s *Scorer) GetRateLimit() *profilesscore.ScoreRateLimit {
	return profilesscore.NewScoreRateLimit(time.Hour*24, 10000)
}

// GetScoreDetails implements Scorer
func (s *Scorer) GetScoreDetails(_ string, application string, username string) (types.ProfileScoreDetails, error) {
	if !strings.EqualFold(application, "youtube") {
		return nil, nil
	}

	channel, err := s.GetChannel(username)
	if err != nil {
		return nil, err
	}

	var linkedTime *time.Time
	if channel.Status.IsLinked {
		date, err := time.Parse(time.RFC3339, channel.ContentOwnerDetails.TimeLinked)
		if err != nil {
			return nil, err
		}
		linkedTime = &date
	}

	return NewScoreDetails(
		linkedTime,
		channel.Statistics.SubscriberCount,
	), nil
}

// GetChannel returns the YouTube channel data of the user having the given username
func (s *Scorer) GetChannel(username string) (*youtube.Channel, error) {
	call := s.client.Channels.List([]string{"contentDetails", "contentOwnerDetails", "snippet", "statistics", "status"})
	call = call.ForUsername(username)
	res, err := call.Do()
	if err != nil {
		return nil, err
	}

	return res.Items[0], nil
}
