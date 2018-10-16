# cachecontrol-lite
Simple gingonic middleware for Cache-Control headers on responses

### Basic Usage

#### Import package

```go
import "github.com/lggomez/cachecontrol-lite/middleware"
import "github.com/lggomez/cachecontrol-lite/cacheobject"
```

#### Add the middleware on your router
```go
func mapUrlsToControllers(router *gin.Engine) {
    /* ... */
    defaultCacheControlConfig := &cacheobject.ResponseCacheDirectives{MaxAge: cacheobject.DeltaSeconds(time.Hour * 24)}

    router.GET("/foo/:id",
        controller.Get,
        cachecontrol.AddCacheControl(defaultCacheControlConfig),
    )
    /* ... */
}
```