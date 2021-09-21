package chi

import (
	"net/http"

	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/lggomez/cachecontrol-lite/middleware"
	"github.com/lggomez/cachecontrol-lite/middleware/cacheobject"
)


func WithCacheControl(handler http.HandlerFunc, directives *cacheobject.ResponseCacheDirectives) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ww := chiMiddleware.NewWrapResponseWriter(w, r.ProtoMajor)

		handler.ServeHTTP(ww, r)

		status := ww.Status()
		// Emit header only for valid 2xx requests
		if status >= 200 && status <= 299 {
			if dir := directives.BuildResponseHeader(); dir != "" {
				ww.Header().Add(middleware.CacheControlHeader, dir)
			}
		}
	}
}