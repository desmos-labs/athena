package x

import (
	"fmt"

	notificationscontext "github.com/desmos-labs/athena/x/notifications/context"
	notificationssender "github.com/desmos-labs/athena/x/notifications/sender"

	"github.com/desmos-labs/athena/database"
	"github.com/desmos-labs/athena/x/apis"
	"github.com/desmos-labs/athena/x/authz"
	contractsbuilder "github.com/desmos-labs/athena/x/contracts/builder"
	"github.com/desmos-labs/athena/x/feegrant"
	"github.com/desmos-labs/athena/x/notifications"
	notificationsbuilder "github.com/desmos-labs/athena/x/notifications/builder"
	standardnotificationsbuilder "github.com/desmos-labs/athena/x/notifications/builder/standard"
	"github.com/desmos-labs/athena/x/posts"
	"github.com/desmos-labs/athena/x/profiles"
	profilesscorebuilder "github.com/desmos-labs/athena/x/profiles-score/builder"
	"github.com/desmos-labs/athena/x/reactions"
	"github.com/desmos-labs/athena/x/relationships"
	"github.com/desmos-labs/athena/x/reports"
	"github.com/desmos-labs/athena/x/subspaces"

	"github.com/forbole/juno/v5/modules"
	"github.com/forbole/juno/v5/modules/registrar"
	"github.com/forbole/juno/v5/modules/telemetry"
	"github.com/forbole/juno/v5/node/remote"
)

type RegistrarOptions struct {
	NotificationsBuilderCreator notificationsbuilder.NotificationsBuilderCreator
	NotificationsSenderCreator  notificationssender.NotificationsSenderCreator
	APIsRegistrar               apis.Registrar
	APIsConfigurator            apis.Configurator
}

func (o RegistrarOptions) CreateNotificationsBuilder(context notificationscontext.Context) notificationsbuilder.NotificationsBuilder {
	if o.NotificationsBuilderCreator != nil {
		return o.NotificationsBuilderCreator(context)
	}
	return standardnotificationsbuilder.CreateNotificationsBuilder(context)
}

func (o RegistrarOptions) CreateNotificationsSender(context notificationscontext.Context) notificationssender.NotificationSender {
	if o.NotificationsSenderCreator != nil {
		return o.NotificationsSenderCreator(context)
	}
	return nil
}

func (o RegistrarOptions) GetAPIsRegistrar() apis.Registrar {
	if o.APIsRegistrar != nil {
		return o.APIsRegistrar
	}
	return apis.DefaultRegistrar
}

func (o RegistrarOptions) GetAPIsConfigurator() apis.Configurator {
	return o.APIsConfigurator
}

// --------------------------------------------------------------------------------------------------------------------

// ModulesRegistrar represents the modules.Registrar that allows to register all custom Athena modules
type ModulesRegistrar struct {
	options RegistrarOptions
}

// NewModulesRegistrar allows to build a new ModulesRegistrar instance
func NewModulesRegistrar() *ModulesRegistrar {
	return &ModulesRegistrar{}
}

// WithOptions sets the given option inside this registrar
func (r *ModulesRegistrar) WithOptions(options RegistrarOptions) *ModulesRegistrar {
	r.options = options
	return r
}

// BuildModules implements modules.Registrar
func (r *ModulesRegistrar) BuildModules(ctx registrar.Context) modules.Modules {
	cdc := ctx.EncodingConfig.Codec
	athenaDb := database.Cast(ctx.Database)

	remoteCfg, ok := ctx.JunoConfig.Node.Details.(*remote.Details)
	if !ok {
		panic(fmt.Errorf("cannot run Athena on local node"))
	}

	grpcConnection := remote.MustCreateGrpcConnection(remoteCfg.GRPC)

	// Juno modules
	telemetryModule := telemetry.NewModule(ctx.JunoConfig)

	// Athena modules
	apisModule := apis.NewModule(apis.NewContext(ctx, grpcConnection))
	if apisModule != nil {
		apisModule = apisModule.WithRegistrar(r.options.GetAPIsRegistrar())
		apisModule = apisModule.WithConfigurator(r.options.GetAPIsConfigurator())
	}

	authzModule := authz.NewModule(ctx.Proxy, cdc, athenaDb)
	contractsModule := contractsbuilder.BuildModule(ctx.JunoConfig, ctx.Proxy, grpcConnection, athenaDb)
	feegrantModule := feegrant.NewModule(ctx.Proxy, cdc, athenaDb)
	postsModule := posts.NewModule(ctx.Proxy, grpcConnection, cdc, athenaDb)
	profilesModule := profiles.NewModule(ctx.Proxy, grpcConnection, cdc, athenaDb)
	profilesScoreModule := profilesscorebuilder.BuildModule(ctx.JunoConfig, athenaDb)
	reactionsModule := reactions.NewModule(ctx.Proxy, grpcConnection, cdc, athenaDb)
	relationshipsModule := relationships.NewModule(profilesModule, grpcConnection, cdc, athenaDb)
	reportsModule := reports.NewModule(ctx.Proxy, grpcConnection, cdc, athenaDb)
	subspacesModule := subspaces.NewModule(ctx.Proxy, grpcConnection, cdc, athenaDb)

	context := notificationscontext.NewContext(ctx, ctx.Proxy, grpcConnection)
	notificationsModule := notifications.NewModule(ctx.JunoConfig, postsModule, reactionsModule, cdc, athenaDb)
	if notificationsModule != nil {
		notificationsModule = notificationsModule.
			WithNotificationsBuilder(r.options.CreateNotificationsBuilder(context)).
			WithNotificationSender(r.options.CreateNotificationsSender(context))
	}

	return []modules.Module{
		apisModule,
		authzModule,
		feegrantModule,
		profilesModule,
		relationshipsModule,
		subspacesModule,
		reportsModule,
		postsModule,
		reactionsModule,
		notificationsModule,
		telemetryModule,
		contractsModule,
		profilesScoreModule,
	}
}
