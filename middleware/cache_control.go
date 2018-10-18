package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lggomez/cachecontrol-lite/middleware/cacheobject"
	"strconv"
)

const (
	CacheControlHeader = "Cache-Control"
	StatusCodeHeader   = "Status-Code"
)

// WithCacheControl is a wrapper middleware for adding Cache-Control headers
func AddCacheControl(directives *cacheobject.ResponseCacheDirectives) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate response state and build header
		status := 0
		statusHeader := c.GetHeader(StatusCodeHeader)
		if statusHeader != "" {
			status, _ = strconv.Atoi(statusHeader)
		}
		// Only cache valid 2xx requests
		if !c.IsAborted() && (status >= 200 && status <= 299) {
			c.Header(CacheControlHeader, directives.BuildResponseHeader())
		}
	}
}
