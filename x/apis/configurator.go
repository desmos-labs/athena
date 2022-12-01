package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Configurator represents a function allowing to configure the Gin engine and the HTTP server.
// NOTE: Use this method only to configure either the router or the server, and NOT to register routes. Those will
// be registered using the various registrars set inside the module.
type Configurator func(router *gin.Engine, server *http.Server)
