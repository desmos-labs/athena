package standard

import (
	"github.com/desmos-labs/djuno/v2/database"
	notificationsbuilder "github.com/desmos-labs/djuno/v2/x/notifications/builder"
	notificationscontext "github.com/desmos-labs/djuno/v2/x/notifications/context"
	"github.com/desmos-labs/djuno/v2/x/profiles"
)

// CreateNotificationsBuilder returns the default NotificationsBuilderCreator implementation
func CreateNotificationsBuilder(context notificationscontext.Context) notificationsbuilder.NotificationsBuilder {
	db := database.Cast(context.Database)
	utilityModule := profiles.NewModule(context.Node, context.GRPCConnection, context.EncodingConfig.Marshaler, db)
	return NewDefaultBuilder(utilityModule)
}

// -------------------------------------------------------------------------------------------------------------------

var (
	_ notificationsbuilder.NotificationsBuilder = &Builder{}
)

// Builder represents the default NotificationsBuilder implementation
type Builder struct {
	utilityModule UtilityModule
}

func NewDefaultBuilder(utilityModule UtilityModule) *Builder {
	return &Builder{
		utilityModule: utilityModule,
	}
}

func (d *Builder) Reactions() notificationsbuilder.ReactionsNotificationsBuilder {
	return NewDefaultReactionsNotificationsBuilder(d.utilityModule)
}

func (d *Builder) Posts() notificationsbuilder.PostsNotificationsBuilder {
	return NewDefaultPostsNotificationsBuilder(d.utilityModule)
}

func (d *Builder) Relationships() notificationsbuilder.RelationshipsNotificationsBuilder {
	return NewDefaultRelationshipsNotificationsBuilder(d.utilityModule)
}
