package auth

import (
	"context"
	"flag"
	oidc "github.com/coreos/go-oidc/v3/oidc"
	"github.com/pkg/errors"
	"os"
)

var (
	keycloakUrl = flag.String("keycloak-url", os.Getenv("PINEAPPLE_OIDC_KEYCLOAK_URL"),
		"oidc keycloak url like")
	clientID = flag.String("oidc-client-id", os.Getenv("PINEAPPLE_OIDC_CLIENT_ID"),
		"oidc client id")
	verifier *oidc.IDTokenVerifier = nil
)

type IDTokenClaim struct {
	Iss               string
	Sub               string
	PreferredUserName string `json:"preferred_username"`
	Email             string
}

func InitAuth() error {
	if *keycloakUrl == "" {
		return errors.New("flag to variable keycloakUrl is not set")
	}
	if *clientID == "" {
		return errors.New("flag to variable clientID is not set")
	}

	provider, err := oidc.NewProvider(context.TODO(), *keycloakUrl)
	if err != nil {
		return err
	}

	oidcConfig := &oidc.Config{
		ClientID: *clientID,
	}
	verifier = provider.Verifier(oidcConfig)
	return nil
}

func Verify(key string) (*IDTokenClaim, error) {
	if verifier == nil {
		return nil, errors.New("verifier not initialized")
	}
	idToken, err := verifier.Verify(context.TODO(), key)
	if err != nil {
		return nil, errors.Wrap(err, "token verify failed")
	}

	var idTokenClaim IDTokenClaim
	if err := idToken.Claims(&idTokenClaim); err != nil {
		return nil, errors.Wrap(err, "IDToken claim cast failed")
	}
	return &idTokenClaim, nil
}
