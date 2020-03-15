// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package form // import "miniflux.app/ui/form"

import (
	"net/http"

	"miniflux.app/errors"
)

// AuthForm represents the authentication form.
type AuthForm struct {
	Username string
	Password string
}

// Validate makes sure the form values are valid.
func (a AuthForm) Validate() error {
	if a.Username == "" || a.Password == "" {
		return errors.NewLocalizedError("error.fields_mandatory")
	}

	return nil
}

// NewAuthForm returns a new AuthForm.
func NewAuthForm(r *http.Request) *AuthForm {
	return &AuthForm{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}
}
