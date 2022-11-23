package standard

// UtilityModule represents a module that contains utility method within it
type UtilityModule interface {
	GetDisplayName(userAddress string) string
}
