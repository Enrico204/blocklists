package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func SHA256String(s string) string {
	var hash = sha256.New()
	hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}
