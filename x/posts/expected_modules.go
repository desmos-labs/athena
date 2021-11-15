package posts

type ProfilesModule interface {
	UpdateProfiles(height int64, addresses []string) error
}
