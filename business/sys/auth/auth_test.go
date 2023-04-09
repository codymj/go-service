package auth_test

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/codymj/go-service/business/sys/auth"
	"github.com/golang-jwt/jwt/v4"
	"testing"
	"time"
)

const (
	success = "\u2713"
	failed  = "\u2717"
)

type keyStore struct {
	pk *rsa.PrivateKey
}

func (k *keyStore) PrivateKey(_ string) (*rsa.PrivateKey, error) {
	return k.pk, nil
}

func (k *keyStore) PublicKey(_ string) (*rsa.PublicKey, error) {
	return &k.pk.PublicKey, nil
}

// =============================================================================

func TestAuth(t *testing.T) {
	t.Log("given the need to be able to authenticate and authorize access")
	{
		testId := 0
		t.Logf("\ttest %d:\twhen handling a single user", testId)
		{
			msg := "should be able to create private key"
			const keyId = "1b24502a-4781-47cb-99c2-3403c23bedac"
			privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				t.Fatalf("\t%s\ttest %d:\t%s: %v", failed, testId, msg, err)
			}
			t.Logf("\t%s\ttest %d:\t%s", success, testId, msg)

			msg = "should be able to create an authenticator"
			authenticator, err := auth.New(keyId, &keyStore{pk: privateKey})
			if err != nil {
				t.Fatalf("\t%s\ttest %d:\t%s: %v", failed, testId, msg, err)
			}
			t.Logf("\t%s\ttest %d:\t%s", success, testId, msg)

			msg = "should be able to generate a JWT"
			claims := auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:  "service project",
					Subject: "1b24502a-4781-47cb-99c2-3403c23bedac",
					ExpiresAt: &jwt.NumericDate{
						Time: time.Now().Add(time.Hour).UTC(),
					},
					IssuedAt: &jwt.NumericDate{
						Time: time.Now().UTC(),
					},
				},
				Roles: []string{auth.RoleAdmin},
			}
			token, err := authenticator.GenerateToken(claims)
			if err != nil {
				t.Fatalf("\t%s\ttest %d:\t%s: %v", failed, testId, msg, err)
			}
			t.Logf("\t%s\ttest %d:\t%s", success, testId, msg)

			msg = "should be able to parse claims"
			parsedClaims, err := authenticator.ValidateToken(token)
			if err != nil {
				t.Fatalf("\t%s\ttest %d:\t%s: %v", failed, testId, msg, err)
			}
			t.Logf("\t%s\ttest %d:\t%s", success, testId, msg)

			msg = "should have expected number of roles"
			if exp, got := len(claims.Roles), len(parsedClaims.Roles); exp != got {
				t.Logf("\t\ttest %d:\texp: %d", testId, exp)
				t.Logf("\t\ttest %d:\tgot: %d", testId, got)
				t.Logf("\t%s\ttest %d:\t%s: %v", failed, testId, msg, err)
			}

			msg = "should have expected roles"
			if exp, got := claims.Roles[0], parsedClaims.Roles[0]; exp != got {
				t.Logf("\t\ttest %d:\texp: %s", testId, exp)
				t.Logf("\t\ttest %d:\tgot: %s", testId, got)
				t.Logf("\t%s\ttest %d:\t%s: %v", failed, testId, msg, err)
			}
		}
	}
}
