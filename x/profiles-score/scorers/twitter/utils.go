package twitter

import (
	"github.com/g8rswimmer/go-twitter/v2"
)

const (
	NotFoundErrorType = "https://api.twitter.com/2/problems/resource-not-found"
)

// hasNotFoundError checks whether inside the given errors slice there is one NotFound error
func hasNotFoundError(errors []*twitter.ErrorObj) bool {
	for _, err := range errors {
		if err.Type == NotFoundErrorType {
			return true
		}
	}
	return false
}
