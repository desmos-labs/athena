package apis

import (
	"github.com/forbole/juno/v5/modules/registrar"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"github.com/desmos-labs/djuno/v2/x/apis/endpoints"
	"github.com/desmos-labs/djuno/v2/x/apis/endpoints/links"
)

// Context contains all the useful data that might be used when registering an API handler
type Context struct {
	registrar.Context
	GRPCConnection *grpc.ClientConn
}

func NewContext(ctx registrar.Context, grpcConnection *grpc.ClientConn) Context {
	return Context{
		Context:        ctx,
		GRPCConnection: grpcConnection,
	}
}

// Registrar represents a function that allows registering API endpoints
type Registrar func(ctx Context, router *gin.Engine) error

// CombinedRegistrar returns a new Registrar combining the given API registrars together
func CombinedRegistrar(registrars ...Registrar) Registrar {
	return func(ctx Context, router *gin.Engine) error {
		for _, registar := range registrars {
			err := registar(ctx, router)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// DefaultRegistrar returns the default API registrar
func DefaultRegistrar(ctx Context, router *gin.Engine) error {
	endpoints.RegisterRoutesList(router)
	links.RegisterRoutes(router, links.NewHandler(ctx.JunoConfig))
	return nil
}
