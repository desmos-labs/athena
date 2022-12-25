package profilesscore

import (
	"time"

	"github.com/rs/zerolog/log"

	"github.com/desmos-labs/djuno/v2/types"
)

// RefreshApplicationLinksScores reads all the applications links stored inside the database and refreshes their scores
func (m *Module) RefreshApplicationLinksScores() error {
	applicationLinks, err := m.db.GetApplicationLinkInfos()
	if err != nil {
		return err
	}
	if applicationLinks == nil {
		return nil
	}

	// Split the application links based on the scores rate limit
	scorers := m.GetScorers()
	chunks, waitTime := getApplicationLinksChunksAndWaitTime(applicationLinks, scorers.GetRateLimit())
	for i, chunk := range chunks {
		for _, link := range chunk {
			for _, scorer := range scorers {
				// Get the score details for this link
				details, err := scorer.GetScoreDetails(link.User, link.Application, link.Username)
				if err != nil {
					log.Error().Err(err).Str("user", link.User).Str("application", link.Application).
						Str("username", link.Username).Msg("error while getting score details")
					continue
				}
				if details == nil {
					continue
				}

				// Save the score
				score := types.NewApplicationLinkScore(link.User, link.Application, link.Username, details, time.Now())
				err = m.db.SaveApplicationLinkScore(score)
				if err != nil {
					return err
				}
			}
		}

		if i > 0 && i < len(chunk)-1 {
			// Wait before processing the next chink
			time.Sleep(waitTime)
		}
	}

	return nil
}

// getApplicationLinksChunksAndWaitTime takes the given application links info slice, and splits them into multiple
// chunks based on the max amount of requests set inside the given rate limit. It also returns the duration that
// should be waited between each slice update based on the rate limit given
func getApplicationLinksChunksAndWaitTime(
	applicationLinks []types.ApplicationLinkInfo, rateLimit *ScoreRateLimit,
) ([][]types.ApplicationLinkInfo, time.Duration) {
	if rateLimit == nil {
		return [][]types.ApplicationLinkInfo{applicationLinks}, 0
	}

	var chunks [][]types.ApplicationLinkInfo
	var chunkSize = int(rateLimit.RateLimit)
	for {
		if len(applicationLinks) == 0 {
			break
		}

		// Necessary check to avoid slicing beyond slice capacity
		if len(applicationLinks) < chunkSize {
			chunkSize = len(applicationLinks)
		}

		chunks = append(chunks, applicationLinks[0:chunkSize])
		applicationLinks = applicationLinks[chunkSize:]
	}

	return chunks, rateLimit.Duration
}
