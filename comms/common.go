package comms

import (
	"crypto/rand"
	"encoding/base64"
)

type HelloMsg struct {
	UserName string
}

// ClientIdMsg is sent from the server to the client in response to a HelloMsg. It provides
// the user with a unique ID.
type ClientIdMsg struct {
	UserID string
}

type PingMsg struct {
	UserID string
}

type PongMsg struct {
	UserID string
}

func generateRandomID(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Encode the random bytes using base64
	randomID := base64.RawURLEncoding.EncodeToString(randomBytes)

	return randomID, nil
}
