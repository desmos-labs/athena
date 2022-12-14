package endpoints

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RegisterRoutesList(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		routes := make([]string, len(r.Routes()))
		for i, route := range r.Routes() {
			routes[i] = route.Path
		}
		c.String(http.StatusOK, fmt.Sprintf("Available endpoints:\n%s", strings.Join(routes, "\n")))
	})
}
