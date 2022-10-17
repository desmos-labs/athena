package standard

import (
	notificationsbuilder "github.com/desmos-labs/djuno/v2/x/notifications/builder"
)

// Creator returns the default NotificationsBuilderCreator implementation
func Creator(module notificationsbuilder.UtilityModule) notificationsbuilder.NotificationsBuilder {
	return NewDefaultBuilder(module)
}

// -------------------------------------------------------------------------------------------------------------------

var (
	_ notificationsbuilder.NotificationsBuilder = &DefaultBuilder{}
)

// DefaultBuilder represents the default NotificationsBuilder implementation
type DefaultBuilder struct {
	m notificationsbuilder.UtilityModule
}

func NewDefaultBuilder(utilityModule notificationsbuilder.UtilityModule) *DefaultBuilder {
	return &DefaultBuilder{
		m: utilityModule,
	}
}

func (d *DefaultBuilder) Reactions() notificationsbuilder.ReactionsNotificationsBuilder {
	return NewDefaultReactionsNotificationsBuilder(d.m)
}

func (d *DefaultBuilder) Posts() notificationsbuilder.PostsNotificationsBuilder {
	return NewDefaultPostsNotificationsBuilder(d.m)
}

func (d *DefaultBuilder) Relationships() notificationsbuilder.RelationshipsNotificationsBuilder {
	return NewDefaultRelationshipsNotificationsBuilder(d.m)
}
