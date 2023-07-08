package commands

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/codymj/go-service/business/sys/database"
	"github.com/codymj/go-service/foundation/keystore"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog"
	"io"
	"os"
	"time"
)

// GenToken generates jwt token
func GenToken(logger *zerolog.Logger, db database.Config, ks keystore.AuthCfg) error {
	// open key file
	f, err := os.Open(fmt.Sprintf("%s%s.pem", ks.KeysFolder, ks.ActiveKid))
	if err != nil {
		return err
	}
	defer f.Close()

	// read pem
	privatePem, err := io.ReadAll(io.LimitReader(f, 1024*1024))
	if err != nil {
		return err
	}

	// parse key
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePem)
	if err != nil {
		return err
	}

	// define a set of claims for generating the jwt
	claims := struct {
		jwt.RegisteredClaims
		Roles []string
	}{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:  "go-service project",
			Subject: "12345",
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(8760 * time.Hour).UTC(),
			},
			IssuedAt: &jwt.NumericDate{
				Time: time.Now().UTC(),
			},
		},
		Roles: []string{"admin"},
	}

	// generate token and set key id
	token := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), claims)
	token.Header["kid"] = ks.ActiveKid

	tokenStr, err := token.SignedString(privateKey)
	if err != nil {
		return err
	}
	fmt.Println("===== TOKEN BEGIN =====")
	fmt.Println(tokenStr)
	fmt.Println("===== TOKEN END =====")
	fmt.Println()

	// marshal the public key
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("error marshalling public key: %w", err)
	}

	// construct pem block for the public key
	publicBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	// write the public key to file
	if err = pem.Encode(os.Stdout, &publicBlock); err != nil {
		return fmt.Errorf("error encoding to public file: %w", err)
	}

	// create token parser
	parser := jwt.NewParser(jwt.WithValidMethods([]string{"RS256"}))

	// key function
	keyFunc := func(t *jwt.Token) (any, error) {
		kid, ok := t.Header["kid"]
		if !ok {
			return nil, nil
		}
		kidId, ok := kid.(string)
		if !ok {
			return nil, nil
		}
		fmt.Printf("KID: %v\n", kidId)
		return &privateKey.PublicKey, nil
	}

	// token
	var parsedClaims struct {
		jwt.RegisteredClaims
		Roles []string
	}
	parsedToken, err := parser.ParseWithClaims(tokenStr, &parsedClaims, keyFunc)
	if err != nil {
		return err
	}
	if !parsedToken.Valid {
		return errors.New("invalid token")
	}
	fmt.Println("token validated")
	fmt.Println()

	return nil
}
