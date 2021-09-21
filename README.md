# cachecontrol-lite
Simple middleware for Cache-Control headers on responses for gingonic and chi routers

### Basic Usage

#### Import package

```go
import "github.com/lggomez/cachecontrol-lite/middleware"
import "github.com/lggomez/cachecontrol-lite/middleware/cacheobject"
```

For chi and gingonic respectively there are specific packages:

```go
import "github.com/lggomez/cachecontrol-lite/middleware/gin"
import "github.com/lggomez/cachecontrol-lite/middleware/chi"
```

#### gingonic - Adding the middleware on your router
```go
func mapUrlsToControllers(router *gin.Engine) {
    /* ... */
    defaultCacheControlConfig := &cacheobject.ResponseCacheDirectives{MaxAge: cacheobject.DeltaSeconds((time.Hour * 24).Seconds())}

    router.GET("/foo/:id",
        controller.Get,
        middleware.AddCacheControl(defaultCacheControlConfig),
    )
    /* ... */
}
```

#### chi - Adding the middleware on your router
The middleware wraps an `http.HandlerFunc`, so you can use your controller (or driver, or whatever conforms to its interface) in the following way:

```go
func mapUrlsToControllers(router *gin.Engine) {
    /* ... */
	r := chi.NewRouter()
	
    r.Get("/", WithCacheControl(controller.Handle, &cacheobject.ResponseCacheDirectives{
        MaxAge: cacheobject.DeltaSeconds((time.Hour * 24).Seconds()),
    }))
    /* ... */
}
```