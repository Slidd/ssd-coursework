package auth

import (
	"context"
	"log"

	"golang.org/x/oauth2"

	oidc "github.com/coreos/go-oidc"
)

type Authenticator struct {
	Provider *oidc.Provider
	Config   oauth2.Config
	Ctx      context.Context
}

func NewAuthenticator() (*Authenticator, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, "https://dev-o6lnq6dg.eu.auth0.com/")
	if err != nil {
		log.Printf("failed to get provider: %v", err)
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     "UmOdRUfJnUUgNB5zHhlltXjTLqQJp5GM",                                 // put into env variable
		ClientSecret: "5W17SYrfStEKXnAfW39Tk-a4NfqNOXlElmIKaFoCqy0JEmk811j0cYSvrMuNwEoz", // put into env variable
		RedirectURL:  "http://localhost:3000/callback",                                   // update
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	return &Authenticator{
		Provider: provider,
		Config:   conf,
		Ctx:      ctx,
	}, nil
}
