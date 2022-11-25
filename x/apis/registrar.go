package apis

import (
	"github.com/forbole/juno/v4/modules/registrar"
	"github.com/gin-gonic/gin"

	"github.com/desmos-labs/djuno/v2/x/apis/endpoints"
)

// Registrar represents a function that allows registering API endpoints
type Registrar func(ctx registrar.Context, router *gin.Engine) error

// CombinedRegistrar returns a new Registrar combining the given API registrars together
func CombinedRegistrar(registrars ...Registrar) Registrar {
	return func(ctx registrar.Context, router *gin.Engine) error {
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
func DefaultRegistrar(_ registrar.Context, router *gin.Engine) error {
	endpoints.RegisterRoutesList(router)
	return nil
}
