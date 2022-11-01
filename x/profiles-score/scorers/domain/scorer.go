package domain

import (
	"time"

	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"

	"net/http"
	"strings"

	"github.com/desmos-labs/djuno/v2/types"
	profilesscore "github.com/desmos-labs/djuno/v2/x/profiles-score"
)

var (
	_ profilesscore.Scorer = &Scorer{}
)

type Scorer struct {
	client http.Client
}

// NewScorer returns a new Scorer instance
func NewScorer() *Scorer {
	return &Scorer{
		client: http.Client{},
	}
}

// GetRateLimit implements Scorer
func (s *Scorer) GetRateLimit() *profilesscore.ScoreRateLimit {
	return profilesscore.NewScoreRateLimit(time.Second, 10)
}

// GetScoreDetails implements Scorer
func (s *Scorer) GetScoreDetails(_ string, application string, username string) (types.ProfileScoreDetails, error) {
	if !strings.EqualFold(application, "domain") {
		return nil, nil
	}

	info, err := s.GetDomainInfo(username)
	if err != nil {
		return nil, err
	}

	creationDate, err := time.Parse(time.RFC3339, info.Domain.CreatedDate)
	if err != nil {
		return nil, err
	}

	return NewScoreDetails(
		creationDate,
	), nil
}

// GetDomainInfo gets the WhoIs information about the given domain
func (s *Scorer) GetDomainInfo(domain string) (*whoisparser.WhoisInfo, error) {
	domain = httpRegEx.ReplaceAllString(domain, "")

	data, err := whois.Whois(domain)
	if err != nil {
		return nil, err
	}

	domainInfo, err := whoisparser.Parse(data)
	if err != nil {
		return nil, err
	}
	return &domainInfo, nil
}
