package links

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/desmos-labs/djuno/v2/x/apis/utils"
)

func RegisterRoutes(r *gin.Engine, handler *Handler) {
	r.POST("/links", func(c *gin.Context) {
		// Get the request body
		body, err := c.GetRawData()
		if err != nil {
			utils.HandleError(c, err)
			return
		}

		// Parse the request
		req, err := handler.ParseGenerateLinkRequest(body)
		if err != nil {
			utils.HandleError(c, err)
			return
		}

		// Validate the request
		err = handler.ValidateLinkRequest(req)
		if err != nil {
			utils.HandleError(c, err)
			return
		}

		// Handle the request
		res, err := handler.HandleGenerateLinkRequest(req)
		if err != nil {
			utils.HandleError(c, err)
			return
		}

		// Return the response
		c.JSON(http.StatusOK, &res)
	})
}
