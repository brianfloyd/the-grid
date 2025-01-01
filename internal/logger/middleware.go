package logger

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
)

func Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		t1 := time.Now()
		defer func() {
			InfoArgs(r.Context(), "%s %s - %d, %d bytes written in %s", r.Method, r.URL, ww.Status(), ww.BytesWritten(), time.Since(t1))
		}()
		next.ServeHTTP(ww, r)
	}
	return http.HandlerFunc(fn)
}
