// Package apnscredentials validates shared console/IM P8 credentials.
package apnscredentials

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"regexp"
)

const MaxPrivateKeySize = 16 * 1024

var appleID = regexp.MustCompile(`^[A-Z0-9]{10}$`)

// ValidTopic validates a configured bundle, before appending a .voip suffix.
// It accepts only 1-100 ASCII letters, digits, periods and hyphens.
func ValidTopic(topic string) bool {
	if len(topic) == 0 || len(topic) > 100 {
		return false
	}
	for i := range topic {
		c := topic[i]
		if !(c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '.' || c == '-') {
			return false
		}
	}
	return true
}

// Validate accepts a single unencrypted PKCS#8 PEM ECDSA P-256 private key.
func Validate(privateKey []byte, keyID, teamID string) error {
	if !appleID.MatchString(keyID) || !appleID.MatchString(teamID) {
		return errors.New("P8 Key ID and Team ID must be 10 uppercase letters or digits")
	}
	return validateKey(privateKey)
}

func validateKey(privateKey []byte) error {
	invalid := errors.New("invalid P8 private key: expected PKCS#8 PEM ECDSA P-256, at most 16 KiB")
	if len(privateKey) == 0 || len(privateKey) > MaxPrivateKeySize {
		return invalid
	}
	input := bytes.TrimSpace(privateKey)
	if !bytes.HasPrefix(input, []byte("-----BEGIN PRIVATE KEY-----")) {
		return invalid
	}
	block, rest := pem.Decode(input)
	if block == nil || block.Type != "PRIVATE KEY" || len(block.Headers) != 0 || len(bytes.TrimSpace(rest)) != 0 {
		return invalid
	}
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return invalid
	}
	key, ok := parsed.(*ecdsa.PrivateKey)
	if !ok || key.Curve != elliptic.P256() || key.D == nil || key.D.Sign() <= 0 || key.D.Cmp(key.Params().N) >= 0 || !key.Curve.IsOnCurve(key.X, key.Y) {
		return invalid
	}
	x, y := key.Curve.ScalarBaseMult(key.D.Bytes())
	if x.Cmp(key.X) != 0 || y.Cmp(key.Y) != 0 {
		return invalid
	}
	return nil
}
