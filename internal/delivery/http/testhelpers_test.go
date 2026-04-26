package http

import (
	"context"
	"net/http"
)

func requestWithOwner(r *http.Request, userID string) *http.Request {
	ctx := context.WithValue(r.Context(), UserIDKey, userID)
	return r.WithContext(ctx)
}
