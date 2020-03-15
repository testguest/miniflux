// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package oauth2 // import "miniflux.app/oauth2"

import (
	"context"
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"
)

type googleProfile struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
}

type googleProvider struct {
	clientID     string
	clientSecret string
	redirectURL  string
}

func (g googleProvider) GetUserExtraKey() string {
	return "google_id"
}

func (g googleProvider) GetRedirectURL(state string) string {
	return g.config().AuthCodeURL(state)
}

func (g googleProvider) GetProfile(code string) (*Profile, error) {
	conf := g.config()
	ctx := context.Background()
	token, err := conf.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	client := conf.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user googleProfile
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&user); err != nil {
		return nil, fmt.Errorf("unable to unserialize google profile: %v", err)
	}

	profile := &Profile{Key: g.GetUserExtraKey(), ID: user.Sub, Username: user.Email}
	return profile, nil
}

func (g googleProvider) config() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  g.redirectURL,
		ClientID:     g.clientID,
		ClientSecret: g.clientSecret,
		Scopes:       []string{"email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}
}

func newGoogleProvider(clientID, clientSecret, redirectURL string) *googleProvider {
	return &googleProvider{clientID: clientID, clientSecret: clientSecret, redirectURL: redirectURL}
}
