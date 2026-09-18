package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const (
	requestId = "X-Request-ID"
)

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get(requestId)
			if id == "" {
				id = uuid.NewString()
			}

			w.Header().Add(requestId, id)
			ctx := context.WithValue(r.Context(), requestId, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}

func RequestIDFromContext(ctx context.Context) string {
	return ctx.Value("X-Request-ID").(string)
}
