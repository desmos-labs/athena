package x

import (
	"fmt"

	notificationscontext "github.com/desmos-labs/djuno/v2/x/notifications/context"

	"github.com/desmos-labs/djuno/v2/database"
	"github.com/desmos-labs/djuno/v2/x/apis"
	"github.com/desmos-labs/djuno/v2/x/authz"
	contractsbuilder "github.com/desmos-labs/djuno/v2/x/contracts/builder"
	"github.com/desmos-labs/djuno/v2/x/feegrant"
	"github.com/desmos-labs/djuno/v2/x/fees"
	"github.com/desmos-labs/djuno/v2/x/notifications"
	notificationsbuilder "github.com/desmos-labs/djuno/v2/x/notifications/builder"
	standardnotificationsbuilder "github.com/desmos-labs/djuno/v2/x/notifications/builder/standard"
	messagebuilder "github.com/desmos-labs/djuno/v2/x/notifications/message-builder"
	topicfirebasemessagebuilder "github.com/desmos-labs/djuno/v2/x/notifications/message-builder/topic"
	"github.com/desmos-labs/djuno/v2/x/posts"
	"github.com/desmos-labs/djuno/v2/x/profiles"
	profilesscorebuilder "github.com/desmos-labs/djuno/v2/x/profiles-score/builder"
	"github.com/desmos-labs/djuno/v2/x/reactions"
	"github.com/desmos-labs/djuno/v2/x/relationships"
	"github.com/desmos-labs/djuno/v2/x/reports"
	"github.com/desmos-labs/djuno/v2/x/subspaces"

	"github.com/forbole/juno/v4/modules"
	"github.com/forbole/juno/v4/modules/registrar"
	"github.com/forbole/juno/v4/modules/telemetry"
	"github.com/forbole/juno/v4/node/remote"
)

type RegistrarOptions struct {
	NotificationsBuilderCreator   notificationsbuilder.NotificationsBuilderCreator
	FirebaseMessageBuilderCreator messagebuilder.FirebaseMessageBuilderCreator
	APIsRegistrar                 apis.Registrar
	APIsConfigurator              apis.Configurator
}

func (o RegistrarOptions) CreateNotificationsBuilder(context notificationscontext.Context) notificationsbuilder.NotificationsBuilder {
	if o.NotificationsBuilderCreator != nil {
		return o.NotificationsBuilderCreator(context)
	}
	return standardnotificationsbuilder.CreateNotificationsBuilder(context)
}

func (o RegistrarOptions) CreateFirebaseMessageBuilder(context notificationscontext.Context) messagebuilder.FirebaseMessageBuilder {
	if o.FirebaseMessageBuilderCreator != nil {
		return o.FirebaseMessageBuilderCreator(context)
	}
	return topicfirebasemessagebuilder.CreateMessageBuilder(context)
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

// ModulesRegistrar represents the modules.Registrar that allows to register all custom DJuno modules
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
	cdc := ctx.EncodingConfig.Marshaler
	djunoDb := database.Cast(ctx.Database)

	remoteCfg, ok := ctx.JunoConfig.Node.Details.(*remote.Details)
	if !ok {
		panic(fmt.Errorf("cannot run DJuno on local node"))
	}

	grpcConnection := remote.MustCreateGrpcConnection(remoteCfg.GRPC)

	// Juno modules
	telemetryModule := telemetry.NewModule(ctx.JunoConfig)

	// DJuno modules
	apisModule := apis.NewModule(apis.NewContext(ctx, grpcConnection))
	if apisModule != nil {
		apisModule = apisModule.WithRegistrar(r.options.GetAPIsRegistrar())
		apisModule = apisModule.WithConfigurator(r.options.GetAPIsConfigurator())
	}

	authzModule := authz.NewModule(ctx.Proxy, cdc, djunoDb)
	contractsModule := contractsbuilder.BuildModule(ctx.JunoConfig, ctx.Proxy, grpcConnection, djunoDb)
	feegrantModule := feegrant.NewModule(ctx.Proxy, cdc, djunoDb)
	feesModule := fees.NewModule(ctx.Proxy, grpcConnection, cdc, djunoDb)
	postsModule := posts.NewModule(ctx.Proxy, grpcConnection, cdc, djunoDb)
	profilesModule := profiles.NewModule(ctx.Proxy, grpcConnection, cdc, djunoDb)
	profilesScoreModule := profilesscorebuilder.BuildModule(ctx.JunoConfig, djunoDb)
	reactionsModule := reactions.NewModule(ctx.Proxy, grpcConnection, cdc, djunoDb)
	relationshipsModule := relationships.NewModule(profilesModule, grpcConnection, cdc, djunoDb)
	reportsModule := reports.NewModule(ctx.Proxy, grpcConnection, cdc, djunoDb)
	subspacesModule := subspaces.NewModule(ctx.Proxy, grpcConnection, cdc, djunoDb)

	context := notificationscontext.NewContext(ctx, ctx.Proxy, grpcConnection)
	notificationsModule := notifications.NewModule(ctx.JunoConfig, postsModule, reactionsModule, cdc, djunoDb)
	if notificationsModule != nil {
		notificationsModule = notificationsModule.WithNotificationsBuilder(r.options.CreateNotificationsBuilder(context))
		notificationsModule = notificationsModule.WithFirebaseMessageBuilder(r.options.CreateFirebaseMessageBuilder(context))
	}

	return []modules.Module{
		apisModule,
		authzModule,
		feegrantModule,
		feesModule,
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
