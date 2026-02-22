package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		ctx := r.Context()
		requestID := middleware.GetReqID(ctx)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		log.Printf("[%s] started %s %s", requestID, r.Method, r.URL.Path)

		next.ServeHTTP(ww, r)

		latency := time.Since(start)
		status := ww.Status()
		log.Printf("[%s] completed %v %s %s %d in %v", requestID, r.Method, r.URL.Path, r.RemoteAddr, status, latency)
	})
}
