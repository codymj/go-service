package password

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
)

// HashPassword generates a password hash.
func (s *service) HashPassword(password string) (string, error) {
	// Generate salt.
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Hash.
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		s.cfg.Time,
		s.cfg.Memory,
		s.cfg.Threads,
		s.cfg.KeyLength,
	)

	// base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Format full password hash.
	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	full := fmt.Sprintf(
		format,
		argon2.Version,
		s.cfg.Memory,
		s.cfg.Time,
		s.cfg.Threads,
		b64Salt,
		b64Hash,
	)

	return full, nil
}
