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

	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		id := req.Header.Get(requestId)
		if id == "" {
			id = uuid.NewString()
		}

		res.Header().Add(requestId, id)
		ctx := context.WithValue(req.Context(), "requestCtxId", id)
		next.ServeHTTP(res, req.WithContext(ctx))

	})

}
