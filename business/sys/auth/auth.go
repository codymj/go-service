package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
)

// KeyLookup declares methods of behavior for looking up keys for jwt use.
type KeyLookup interface {
	PrivateKey(kid string) (*rsa.PrivateKey, error)
	PublicKey(kid string) (*rsa.PublicKey, error)
}

// Auth is used to authenticate clients.
type Auth struct {
	activeKid string
	keyLookup KeyLookup
	method    jwt.SigningMethod
	keyFunc   func(t *jwt.Token) (any, error)
	parser    jwt.Parser
}

// New creates an Auth to support authentication/authorization.
func New(activeKid string, lookup KeyLookup) (*Auth, error) {
	// activeKid represents the private key used to sign new tokens
	_, err := lookup.PrivateKey(activeKid)
	if err != nil {
		return nil, errors.New("active KID does not exist in store")
	}

	// get signing method
	method := jwt.GetSigningMethod("RS256")
	if method == nil {
		return nil, errors.New("invalid signing method")
	}

	// implement key function
	keyFunc := func(t *jwt.Token) (any, error) {
		kid, ok := t.Header["kid"]
		if !ok {
			return nil, errors.New("missing KID in token header")
		}
		kidId, ok := kid.(string)
		if !ok {
			return nil, errors.New("invalid KID in token header")
		}

		return lookup.PublicKey(kidId)
	}

	// create the token parser
	parser := jwt.NewParser(jwt.WithValidMethods([]string{"RS256"}))

	// set auth parameters
	return &Auth{
		activeKid: activeKid,
		keyLookup: lookup,
		method:    method,
		keyFunc:   keyFunc,
		parser:    *parser,
	}, nil
}

// GenerateToken generates a signed JWT token string representing user claims.
func (a *Auth) GenerateToken(claims Claims) (string, error) {
	// create token and set kid in header
	token := jwt.NewWithClaims(a.method, claims)
	token.Header["kid"] = a.activeKid

	// get private key
	privateKey, err := a.keyLookup.PrivateKey(a.activeKid)
	if err != nil {
		return "", errors.New("kid lookup failed")
	}

	// generate signing string
	str, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("error generating signing string: %w", err)
	}

	return str, nil
}

// ValidateToken creates the Claims that were used to generate a token.
func (a *Auth) ValidateToken(tokenStr string) (Claims, error) {
	// parse the claims
	var c Claims
	token, err := a.parser.ParseWithClaims(tokenStr, &c, a.keyFunc)
	if err != nil {
		return Claims{}, fmt.Errorf("error parsing token: %w", err)
	}
	if !token.Valid {
		return Claims{}, errors.New("invalid token string")
	}

	return c, nil
}
