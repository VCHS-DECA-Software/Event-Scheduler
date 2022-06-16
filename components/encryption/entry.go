package encryption

import (
	"crypto/rand"
)

func Initialize() error {
	signingKey = make([]byte, 16)
	_, err := rand.Read(signingKey)
	if err != nil {
		return err
	}
	return nil
}
