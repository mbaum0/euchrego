package comms

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/fatih/color"
)

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

func logError(format string, args ...interface{}) {
	format = fmt.Sprintf("  %s\n", format)
	color.New(color.FgRed).Printf(format, args...)
}

func logWarn(format string, args ...interface{}) {
	format = fmt.Sprintf("  %s\n", format)
	color.New(color.FgYellow).Printf(format, args...)
}

func logInfo(format string, args ...interface{}) {
	format = fmt.Sprintf("  %s\n", format)
	color.New(color.FgBlue).Printf(format, args...)
}

func logSuccess(format string, args ...interface{}) {
	format = fmt.Sprintf("  %s\n", format)
	color.New(color.FgGreen).Printf(format, args...)
}
