package siteservice

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAvailableSessions(t *testing.T) {

	siteService := NewService("MyCookieSecret", nil, nil, nil, "test", true)
	request := &http.Request{}

	session, err := siteService.GetSession(request, SessionForRegistration, "akey")
	assert.NoError(t, err)
	assert.NotNil(t, session)

	session, err = siteService.GetSession(request, SessionInteractive, "usersession")
	assert.NoError(t, err)
	assert.NotNil(t, session)

}
