// Copyright 2018 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package ui // import "miniflux.app/ui"

import (
	"miniflux.app/config"
	"miniflux.app/oauth2"
)

func getOAuth2Manager() *oauth2.Manager {
	return oauth2.NewManager(
		config.Opts.OAuth2ClientID(),
		config.Opts.OAuth2ClientSecret(),
		config.Opts.OAuth2RedirectURL(),
	)
}
