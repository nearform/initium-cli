package secrets

import (
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestGenerateKeys(t *testing.T) {
	keys, err := GenerateKeys()
	secretKeyPrefix := "AGE-SECRET-KEY-"
	publicKeyPrefix := "age1"

	if err != nil {
		t.Error(err)
	}

	if !strings.HasPrefix(keys.Private, secretKeyPrefix) {
		t.Errorf("Secret key doesn't start with %s", secretKeyPrefix)
	}

	if !strings.HasPrefix(keys.Public, publicKeyPrefix) {
		t.Errorf("Public key doesn't start with %s", publicKeyPrefix)
	}
}

func TestDecrypt(t *testing.T) {
	keys, err := GenerateKeys()
	if err != nil {
		t.Error(err)
	}

	expected := "simple"

	secret, err := Encrypt(keys.Public, expected)
	if err != nil {
		t.Error(err)
	}

	plain, err := Decrypt(keys.Private, secret)
	if err != nil {
		t.Error(err)
	}

	assert.Assert(t, plain == expected, "Expected %q, got %q", expected, plain)
}
