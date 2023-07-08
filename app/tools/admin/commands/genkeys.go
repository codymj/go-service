package commands

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// GenKeys creates a x509 public/private key pair for auth tokens.
func GenKeys() error {
	// generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// create file to hold private key
	privateFile, err := os.Create("private.pem")
	if err != nil {
		return fmt.Errorf("error creating private.pem: %w", err)
	}
	defer privateFile.Close()

	// construct pem block for the private key
	privateBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	// write the private key to file
	if err = pem.Encode(privateFile, &privateBlock); err != nil {
		return fmt.Errorf("error encoding to private file: %w", err)
	}

	// marshal the public key
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("error marshalling public key: %w", err)
	}

	// create file to hold public key
	publicFile, err := os.Create("public.pem")
	if err != nil {
		return fmt.Errorf("error creating public.pem: %w", err)
	}
	defer publicFile.Close()

	// construct pem block for the public key
	publicBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	// write the public key to file
	if err = pem.Encode(publicFile, &publicBlock); err != nil {
		return fmt.Errorf("error encoding to public file: %w", err)
	}

	fmt.Println("public and private keys generated successfully")

	return nil
}
