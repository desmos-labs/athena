package utils

import (
	"fmt"
	"net/http"
	"unicode"

	"github.com/gin-gonic/gin"
)

type HTTPError struct {
	StatusCode int
	Response   string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("status %d: %s", e.StatusCode, e.Response)
}

// WrapErr wraps the given error into a new one that contains the given status code and response
func WrapErr(statusCode int, res string) error {
	return &HTTPError{
		StatusCode: statusCode,
		Response:   res,
	}
}

func ucFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

// UnwrapErr unwraps the given error returning the status code and response
func UnwrapErr(err error) (statusCode int, res string) {
	if httpErr, ok := err.(*HTTPError); ok {
		return httpErr.StatusCode, httpErr.Response
	}
	return http.StatusInternalServerError, ucFirst(err.Error())
}

type errorJSONResponse struct {
	Error string `json:"error"`
}

// HandleError handles the given error by returning the proper response
func HandleError(c *gin.Context, err error) {
	statusCode, res := UnwrapErr(err)
	c.Abort()
	c.Error(err)
	c.JSON(statusCode, errorJSONResponse{Error: res})
}
