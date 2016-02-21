// Copyright 2012 The KidStuff Authors.
// Copyright 2016 the ItsYou.online developers
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sessions

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type TokenGetSeter interface {
	GetToken(req *http.Request, name string) (string, error)
	SetToken(rw http.ResponseWriter, name, value string, options *sessions.Options)
}

type CookieToken struct{}

func (c *CookieToken) GetToken(req *http.Request, name string) (string, error) {
	cook, err := req.Cookie(name)
	if err != nil {
		return "", err
	}

	return cook.Value, nil
}

func (c *CookieToken) SetToken(rw http.ResponseWriter, name, value string,
	options *sessions.Options) {
	http.SetCookie(rw, sessions.NewCookie(name, value, options))
}
