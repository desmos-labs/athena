package apis

import (
	"github.com/gin-gonic/gin"

	"github.com/desmos-labs/djuno/v2/x/apis/endpoints"
)

// Registrar represents a function that allows registering API endpoints
type Registrar func(router *gin.Engine) error

// CombinedRegistrar returns a new Registrar combining the given API registrars together
func CombinedRegistrar(registrars ...Registrar) Registrar {
	return func(router *gin.Engine) error {
		for _, registar := range registrars {
			err := registar(router)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// DefaultRegistrar returns the default API registrar
func DefaultRegistrar() Registrar {
	return func(router *gin.Engine) error {
		endpoints.RegisterRoutesList(router)
		return nil
	}
}
