package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lggomez/cachecontrol-lite/middleware/cacheobject"
)

const (
	CacheControlHeader = "Cache-Control"
)

// WithCacheControl is a wrapper middleware for adding Cache-Control headers
func AddCacheControl(directives *cacheobject.ResponseCacheDirectives) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := c.Writer.Status()
		// Emit header only for valid 2xx requests
		if !c.IsAborted() && (status >= 200 && status <= 299) {
			c.Header(CacheControlHeader, directives.BuildResponseHeader())
		}
	}
}
