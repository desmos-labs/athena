package notifications

import (
	"fmt"
)

// getDisplayName returns the name to be displayed for the user having the given address
func (m *Module) getDisplayName(userAddress string) string {
	profile, err := m.profilesModule.GetUserProfile(userAddress)
	if err != nil || profile == nil {
		return fmt.Sprintf("%[1]s...%[2]s", userAddress[:9], userAddress[len(userAddress)-5:])
	}

	switch {
	case profile.Nickname != "":
		return fmt.Sprintf("%[1]s (@%[2]s)", profile.Nickname, profile.DTag)

	default:
		return fmt.Sprintf("@%[1]s", profile.DTag)
	}
}
