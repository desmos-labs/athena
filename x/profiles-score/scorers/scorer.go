package scorers

// Scorer represents a generic parses that gets data from an external application and converts it to a specific score
type Scorer interface {
	// SupportedApplications returns the list of applications that this scorer supports
	SupportedApplications() []string

	// RefreshScore refreshes the score for the user having the given address and username on the specified application
	RefreshScore(address string, username string, application string) error
}
