package types

import (
	"time"
)

type ProfileScore struct {
	DesmosAddress string
	Application   string
	Username      string
	Details       ProfileScoreDetails
	Timestamp     time.Time
}

func NewApplicationLinkScore(address string, application string, username string, details ProfileScoreDetails, timestamp time.Time) *ProfileScore {
	return &ProfileScore{
		DesmosAddress: address,
		Application:   application,
		Username:      username,
		Details:       details,
		Timestamp:     timestamp,
	}
}

type ProfileScoreDetails interface {
	GetScore() uint64
}
