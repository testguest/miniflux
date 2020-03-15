// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package ui // import "miniflux.app/ui"

import (
	"net/http"

	"miniflux.app/http/response/html"
	"miniflux.app/http/request"
	"miniflux.app/http/route"
	"miniflux.app/ui/session"
	"miniflux.app/ui/view"
)

func (h *handler) showLoginPage(w http.ResponseWriter, r *http.Request) {
	if request.IsAuthenticated(r) {
		html.Redirect(w, r, route.Path(h.router, "unread"))
		return
	}

	sess := session.New(h.store, request.SessionID(r))
	view := view.New(h.tpl, r, sess)
	html.OK(w, r, view.Render("login"))
}
