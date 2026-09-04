package handlers

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/occult/pagode/pkg/routenames"
	"github.com/occult/pagode/pkg/tests"
	"github.com/stretchr/testify/require"
)

func TestAdmin__RequiresAdmin(t *testing.T) {
	anon := request(t)
	anon.setRoute(routenames.AdminDashboard).get().assertStatusCode(http.StatusUnauthorized)
	anon.setRoute(routenames.AdminEntityList("User")).get().assertStatusCode(http.StatusUnauthorized)

	member, err := tests.CreateUser(c.ORM)
	require.NoError(t, err)
	asMember := loginAs(t, member.Email)
	asMember.setRoute(routenames.AdminDashboard).get().assertStatusCode(http.StatusUnauthorized)
	asMember.setRoute(routenames.AdminEntityList("User")).get().assertStatusCode(http.StatusUnauthorized)

	admin, err := tests.CreateUser(c.ORM)
	require.NoError(t, err)
	_, err = admin.Update().SetAdmin(true).Save(context.Background())
	require.NoError(t, err)
	asAdmin := loginAs(t, admin.Email)
	asAdmin.setRoute(routenames.AdminDashboard).get().assertStatusCode(http.StatusOK)
	asAdmin.setRoute(routenames.AdminEntityList("User")).get().assertStatusCode(http.StatusOK)
}

// loginAs signs in through POST /user/login with the password tests.CreateUser assigns,
// returning a request whose cookie jar carries the session.
func loginAs(t *testing.T, email string) *httpRequest {
	r := request(t)
	loginURL, err := url.Parse(srv.URL + c.Web.Reverse(routenames.Login))
	require.NoError(t, err)

	resp, err := r.client.Get(loginURL.String())
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())

	var csrf string
	for _, ck := range r.client.Jar.Cookies(loginURL) {
		if ck.Name == "XSRF-TOKEN" {
			csrf = ck.Value
		}
	}
	require.NotEmpty(t, csrf, "XSRF-TOKEN cookie after GET /user/login")

	form := url.Values{"email": {email}, "password": {"password"}}
	req, err := http.NewRequest(http.MethodPost, loginURL.String(), strings.NewReader(form.Encode()))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-XSRF-TOKEN", csrf)
	resp, err = r.client.Do(req)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	require.Equal(t, http.StatusOK, resp.StatusCode, "login should redirect to the dashboard and load it")
	return r
}
