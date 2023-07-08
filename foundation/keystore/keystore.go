package keystore

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"io/fs"
	"path"
	"strings"
	"sync"
)

// AuthCfg for setting up key store
type AuthCfg struct {
	KeysFolder string
	ActiveKid  string
}

// KeyStore represents an in-memory store implementation of the KeyStorer
// interface for use with the auth package.
type KeyStore struct {
	mu    sync.RWMutex
	store map[string]*rsa.PrivateKey
}

// New constructs an empty KeyStore.
func New() *KeyStore {
	return &KeyStore{
		store: make(map[string]*rsa.PrivateKey),
	}
}

// NewMap constructs a KeyStore with an initial set of keys.
func NewMap(store map[string]*rsa.PrivateKey) *KeyStore {
	return &KeyStore{
		store: store,
	}
}

// NewFS constructs a KeyStore based on a set of PEM files rooted inside a
// directory. The name of each PEM file will be used as the key ID.
func NewFS(fsys fs.FS) (*KeyStore, error) {
	// create key store
	ks := KeyStore{
		store: make(map[string]*rsa.PrivateKey),
	}

	// function to search for key files
	fn := func(fn string, dir fs.DirEntry, err error) error {
		// return on error during directory walk
		if err != nil {
			return fmt.Errorf("walkdir failure: %w", err)
		}

		// return if directory
		if dir.IsDir() {
			return nil
		}

		// return if there are no .pem files
		if path.Ext(fn) != ".pem" {
			return nil
		}

		// open file
		f, err := fsys.Open(fn)
		if err != nil {
			return fmt.Errorf("error opening key file: %w", err)
		}
		defer f.Close()

		// read file
		privatePem, err := io.ReadAll(io.LimitReader(f, 1024*1024))
		if err != nil {
			return fmt.Errorf("error reading pem file: %w", err)
		}

		// parse file
		privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePem)
		if err != nil {
			return fmt.Errorf("error parsing pem file: %w", err)
		}

		// store key
		ks.store[strings.TrimSuffix(dir.Name(), ".pem")] = privateKey

		return nil
	}

	if err := fs.WalkDir(fsys, ".", fn); err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}

	return &ks, nil
}

// Add a private key and combination kid to the store.
func (ks *KeyStore) Add(privateKey *rsa.PrivateKey, kid string) {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	ks.store[kid] = privateKey
}

// Remove a private key from key store.
func (ks *KeyStore) Remove(kid string) {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	delete(ks.store, kid)
}

// PrivateKey searches the key store for a given kid and returns private key.
func (ks *KeyStore) PrivateKey(kid string) (*rsa.PrivateKey, error) {
	ks.mu.RLock()
	defer ks.mu.RUnlock()

	privateKey, found := ks.store[kid]
	if !found {
		return nil, errors.New("kid lookup failed")
	}

	return privateKey, nil
}

// PublicKey searches the key store for a given kid and returns the public key.
func (ks *KeyStore) PublicKey(kid string) (*rsa.PublicKey, error) {
	ks.mu.RLock()
	defer ks.mu.RUnlock()

	privateKey, found := ks.store[kid]
	if !found {
		return nil, errors.New("kid lookup failed")
	}

	return &privateKey.PublicKey, nil
}
