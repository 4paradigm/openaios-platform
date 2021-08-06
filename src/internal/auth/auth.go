/*
 * Copyright © 2021 peizhaoyou <peizhaoyou@4paradigm.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package auth implements auth methods
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
		"oidc keycloak url like https://keycloak.pineapple.com:32443/auth/realms/develop")
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
