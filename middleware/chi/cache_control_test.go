package chi

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/lggomez/cachecontrol-lite/middleware/cacheobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testWriter struct {
	statusCode int
}

func (cw *testWriter) Handle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(cw.statusCode)
	w.Write([]byte("{\"foo\":\"bar\"}"))
}

func TestRequestCacheControl_2XX(t *testing.T) {
	recorder := httptest.NewRecorder()
	_ = recorder.Header()

	r := chi.NewRouter()
	tw := &testWriter{http.StatusOK}
	r.Get("/", WithCacheControl(tw.Handle, &cacheobject.ResponseCacheDirectives{
		MaxAge: cacheobject.DeltaSeconds((time.Hour * 24).Seconds()),
	}))

	req, _ := http.NewRequest("GET", "/", nil)

	r.ServeHTTP(recorder, req)
	res := recorder.Result()

	require.NotNil(t, res)
	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.Equal(t, recorder.Header().Get("Cache-Control"), "max-age=86400")

	data, _ := ioutil.ReadAll(recorder.Body)
	assert.Equal(t, string(data), "{\"foo\":\"bar\"}")
}

func TestRequestCacheControl_2XX_EmptyDirective(t *testing.T) {
	recorder := httptest.NewRecorder()
	_ = recorder.Header()

	r := chi.NewRouter()
	tw := &testWriter{http.StatusOK}
	r.Get("/", WithCacheControl(tw.Handle, &cacheobject.ResponseCacheDirectives{}))

	req, _ := http.NewRequest("GET", "/", nil)

	r.ServeHTTP(recorder, req)
	res := recorder.Result()

	require.NotNil(t, res)
	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.Equal(t, recorder.Header().Get("Cache-Control"), "")

	data, _ := ioutil.ReadAll(recorder.Body)
	assert.Equal(t, string(data), "{\"foo\":\"bar\"}")
}

func TestRequestCacheControl_400(t *testing.T) {
	recorder := httptest.NewRecorder()
	_ = recorder.Header()

	r := chi.NewRouter()
	tw := &testWriter{http.StatusTeapot}
	r.Get("/", WithCacheControl(tw.Handle, &cacheobject.ResponseCacheDirectives{
		MaxAge: cacheobject.DeltaSeconds((time.Hour * 24).Seconds()),
	}))

	req, _ := http.NewRequest("GET", "/", nil)

	r.ServeHTTP(recorder, req)
	res := recorder.Result()

	require.NotNil(t, res)
	assert.Equal(t, res.StatusCode, http.StatusTeapot)
	assert.Equal(t, recorder.Header().Get("Cache-Control"), "")

	data, _ := ioutil.ReadAll(recorder.Body)
	assert.Equal(t, string(data), "{\"foo\":\"bar\"}")
}
