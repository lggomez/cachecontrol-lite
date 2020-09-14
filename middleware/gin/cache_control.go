package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/lggomez/cachecontrol-lite/middleware"
	"github.com/lggomez/cachecontrol-lite/middleware/cacheobject"
)

// WithCacheControl is a wrapper middleware for adding Cache-Control headers
func WithCacheControl(directives *cacheobject.ResponseCacheDirectives) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := c.Writer.Status()
		// Emit header only for valid 2xx requests
		if !c.IsAborted() && (status >= 200 && status <= 299) {
			c.Header(middleware.CacheControlHeader, directives.BuildResponseHeader())
		}
	}
}
