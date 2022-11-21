package middleware

import (
	"context"
	"github.com/ProninIgorr/alcomarket-nearst/pkg/helpers/uuid"
	"net/http"
)

func (m *Middleware) ContextRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := newContextWithRequestID(r.Context(), r)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func newContextWithRequestID(ctx context.Context, req *http.Request) context.Context {
	reqID := req.Header.Get("X-Request-ID")
	if reqID == "" {
		reqID = uuid.GenerateUUID()
	}

	return context.WithValue(ctx, "requestID", reqID)
}
