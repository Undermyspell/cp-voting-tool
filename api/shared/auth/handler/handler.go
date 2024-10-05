package authhandler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"voting/internal/env"
	"voting/shared/auth/jwks"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

var (
	clientID     = ""
	clientSecret = ""
	tenant       = ""
	redirectURL  = ""
	oauthConfig  = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "",
		Scopes:       []string{},
		Endpoint:     microsoft.AzureADEndpoint(""),
	}
)

func InitOAuthConfig() {
	clientID = env.Env.AuthAzureClientId
	clientSecret = env.Env.AuthAzureClientSecret
	tenant = env.Env.AuthAzureTenantId
	redirectURL = env.Env.AuthRedirectUrl
	oauthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"openid User.Read email profile"},
		Endpoint:     microsoft.AzureADEndpoint(tenant),
	}
}

func Login(c *gin.Context) {
	session := sessions.Default(c)

	state := generateRandomString()
	nonce := generateRandomString()
	session.Set("oauth_state", state)
	session.Set("oauth_nonce", nonce)
	session.Save()

	url := oauthConfig.AuthCodeURL(state, oauth2.SetAuthURLParam("nonce", nonce), oauth2.AccessTypeOffline)
	c.Redirect(http.StatusFound, url)
}

func LoginCallback(c *gin.Context) {
	session := sessions.Default(c)

	state := c.Query("state")
	code := c.Query("code")

	storedState := session.Get("oauth_state").(string)
	storedNonce := session.Get("oauth_nonce").(string)

	if code == "" || state == "" {
		c.String(http.StatusBadRequest, "Code or state missing")
		return
	}

	if state != storedState {
		c.String(http.StatusBadRequest, "State parameter mismatch")
		return
	}

	token, err := oauthConfig.Exchange(context.Background(), code)

	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to exchange token: %s", err.Error())
		return
	}

	idToken := token.Extra("id_token")
	parsedIdToken, err := jwt.Parse(idToken.(string), jwks.GetProvider().GetKeyFunc())

	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to validate token")
		return
	}

	nonce := parsedIdToken.Claims.(jwt.MapClaims)["nonce"]

	if storedNonce != nonce {
		c.String(http.StatusBadRequest, "Nonce does not match")
		return
	}

	exp := parsedIdToken.Claims.(jwt.MapClaims)["exp"].(float64)
	session.Delete("oauth_state")
	session.Set("authenticated", true)
	session.Set("token", idToken)
	session.Set("tokenExpiry", int64(exp))
	session.Save()

	c.Redirect(http.StatusFound, "/")
}

func generateRandomString() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)
}
