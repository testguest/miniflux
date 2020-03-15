// Copyright 2018 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package api // import "miniflux.app/api"

import (
	"context"
	"net/http"

	"miniflux.app/http/request"
	"miniflux.app/http/response/json"
	"miniflux.app/logger"
	"miniflux.app/storage"
)

type middleware struct {
	store *storage.Storage
}

func newMiddleware(s *storage.Storage) *middleware {
	return &middleware{s}
}

// BasicAuth handles HTTP basic authentication.
func (m *middleware) serve(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		clientIP := request.ClientIP(r)
		username, password, authOK := r.BasicAuth()
		if !authOK {
			logger.Debug("[API] No authentication headers sent")
			json.Unauthorized(w, r)
			return
		}

		if err := m.store.CheckPassword(username, password); err != nil {
			logger.Error("[API] [ClientIP=%s] Invalid username or password: %s", clientIP, username)
			json.Unauthorized(w, r)
			return
		}

		user, err := m.store.UserByUsername(username)
		if err != nil {
			logger.Error("[API] %v", err)
			json.ServerError(w, r, err)
			return
		}

		if user == nil {
			logger.Error("[API] [ClientIP=%s] User not found: %s", clientIP, username)
			json.Unauthorized(w, r)
			return
		}

		logger.Info("[API] User authenticated: %s", username)
		m.store.SetLastLogin(user.ID)

		ctx := r.Context()
		ctx = context.WithValue(ctx, request.UserIDContextKey, user.ID)
		ctx = context.WithValue(ctx, request.UserTimezoneContextKey, user.Timezone)
		ctx = context.WithValue(ctx, request.IsAdminUserContextKey, user.IsAdmin)
		ctx = context.WithValue(ctx, request.IsAuthenticatedContextKey, true)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
