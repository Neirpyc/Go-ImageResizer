package main

import (
	"crypto/sha512"
	"encoding/base64"
)

func sha512Str(src string) string {
	h := sha512.New()
	h.Write([]byte(src))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}
