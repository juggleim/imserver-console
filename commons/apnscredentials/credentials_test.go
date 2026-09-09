package apnscredentials

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"strings"
	"sync"
	"testing"
)

func testKey(t *testing.T, curve elliptic.Curve) []byte {
	t.Helper()
	key, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
}

func TestValidate(t *testing.T) {
	valid := testKey(t, elliptic.P256())
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	rsaDER, err := x509.MarshalPKCS8PrivateKey(rsaKey)
	if err != nil {
		t.Fatal(err)
	}
	if err := Validate(valid, "ABC1234567", "XYZ1234567"); err != nil {
		t.Fatal(err)
	}
	for name, key := range map[string][]byte{
		"empty": nil, "malformed": []byte("secret-invalid-key"), "oversized": bytes.Repeat([]byte("x"), MaxPrivateKeySize+1),
		"p384": testKey(t, elliptic.P384()), "rsa": pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: rsaDER}),
		"multiple": append(bytes.Clone(valid), valid...), "prefix": append([]byte("garbage"), valid...),
		"wrong PEM": bytes.ReplaceAll(valid, []byte("PRIVATE KEY"), []byte("EC PRIVATE KEY")),
	} {
		t.Run(name, func(t *testing.T) {
			if err := Validate(key, "ABC1234567", "XYZ1234567"); err == nil {
				t.Fatal("accepted invalid key")
			}
		})
	}
	for _, id := range []string{"", "abc1234567", "ABCDEFGHI", "ABCDEFGHIJK", "ABC 234567", "ABCDEFGHI!"} {
		if Validate(valid, id, "XYZ1234567") == nil || Validate(valid, "ABC1234567", id) == nil {
			t.Fatal("accepted invalid metadata")
		}
	}
}

func TestValidTopic(t *testing.T) {
	for _, topic := range []string{"a", "com.Example-app123", strings.Repeat("a", 100)} {
		if !ValidTopic(topic) {
			t.Fatalf("valid bundle rejected: length %d", len(topic))
		}
	}
	for _, topic := range []string{"", "com.example app", "com/example", "com_example", "com.例", "com.example\n", " com.example", strings.Repeat("a", 101)} {
		if ValidTopic(topic) {
			t.Fatal("invalid bundle accepted")
		}
	}
}

func TestValidatePreservesCallerBytesConcurrently(t *testing.T) {
	key := append([]byte("\r\n  "), testKey(t, elliptic.P256())...)
	key = append(key, []byte(" \r\n")...)
	want := bytes.Clone(key)
	var wg sync.WaitGroup
	for range 16 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := Validate(key, "ABC1234567", "XYZ1234567"); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	if !bytes.Equal(key, want) {
		t.Fatal("validation mutated caller bytes")
	}
}
