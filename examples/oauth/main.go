package main

import (
	"context"
	"encoding/base64"
	"log"

	"golang.org/x/oauth2"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
	v1 "github.com/juhokoskela/pipedrive-go/pipedrive/v1"
	v2 "github.com/juhokoskela/pipedrive-go/pipedrive/v2"
)

func main() {
	const (
		clientID     = "YOUR_CLIENT_ID"
		clientSecret = "YOUR_CLIENT_SECRET"
		redirectURI  = "YOUR_REDIRECT_URI"
		code         = "AUTHORIZATION_CODE"
	)

	ctx := context.Background()

	oauthClient, err := v1.NewClient(pipedrive.Config{})
	if err != nil {
		log.Fatal(err)
	}

	authURL, err := oauthClient.OAuth.Authorize(
		ctx,
		v1.WithOAuthClientID(clientID),
		v1.WithOAuthRedirectURI(redirectURI),
		v1.WithOAuthState("state123"),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("open: %s", authURL)

	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(clientID+":"+clientSecret))
	tokens, err := oauthClient.OAuth.GetTokens(
		ctx,
		v1.WithOAuthAuthorization(authHeader),
		v1.WithOAuthCode(code),
		v1.WithOAuthRedirectURI(redirectURI),
	)
	if err != nil {
		log.Fatal(err)
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: tokens.AccessToken})
	apiClient, err := v2.NewClient(pipedrive.Config{
		BaseURL: tokens.APIDomain + "/api/v2",
		Auth:    pipedrive.OAuth2Auth{TokenSource: ts},
	})
	if err != nil {
		log.Fatal(err)
	}

	pipelines, _, err := apiClient.Pipelines.List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("pipelines=%d", len(pipelines))
}
