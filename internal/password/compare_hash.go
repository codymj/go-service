package password

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strings"
)

// CompareHash is used to compare a user-inputed password to a hash.
func (s *service) CompareHash(password, hash string) (bool, error) {
	// Split hash into parts.
	parts := strings.Split(hash, "$")

	// Scan parameters.
	cfg := &Config{}
	_, err := fmt.Sscanf(
		parts[3],
		"m=%d,t=%d,p=%d",
		&cfg.Memory,
		&cfg.Time,
		&cfg.Threads,
	)
	if err != nil {
		return false, err
	}

	// Extract salt.
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	// Extract decoded hash and length.
	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}
	cfg.KeyLength = uint32(len(decodedHash))

	// Compare hashes.
	comparisonHash := argon2.IDKey(
		[]byte(password),
		salt,
		s.cfg.Time,
		s.cfg.Memory,
		s.cfg.Threads,
		s.cfg.KeyLength,
	)

	return subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1, nil
}
