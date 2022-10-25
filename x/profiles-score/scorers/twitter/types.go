package twitter

import (
	"fmt"
	"github.com/g8rswimmer/go-twitter/v2"
	"net/http"
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
