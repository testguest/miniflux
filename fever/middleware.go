// Copyright 2018 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package fever // import "miniflux.app/fever"

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

func (m *middleware) serve(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.FormValue("api_key")
		if apiKey == "" {
			logger.Info("[Fever] No API key provided")
			json.OK(w, r, newAuthFailureResponse())
			return
		}

		user, err := m.store.UserByFeverToken(apiKey)
		if err != nil {
			logger.Error("[Fever] %v", err)
			json.OK(w, r, newAuthFailureResponse())
			return
		}

		if user == nil {
			logger.Info("[Fever] No user found with this API key")
			json.OK(w, r, newAuthFailureResponse())
			return
		}

		logger.Info("[Fever] User #%d is authenticated", user.ID)
		m.store.SetLastLogin(user.ID)

		ctx := r.Context()
		ctx = context.WithValue(ctx, request.UserIDContextKey, user.ID)
		ctx = context.WithValue(ctx, request.UserTimezoneContextKey, user.Timezone)
		ctx = context.WithValue(ctx, request.IsAdminUserContextKey, user.IsAdmin)
		ctx = context.WithValue(ctx, request.IsAuthenticatedContextKey, true)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
