package secrets

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"filippo.io/age"
)

type Keys struct {
	Private string
	Public  string
}

func GenerateKeys() (Keys, error) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		return Keys{}, err
	}

	keys := Keys{
		Private: identity.String(),
		Public:  identity.Recipient().String(),
	}
	return keys, nil
}

func Decrypt(privateKey string, secret string) (string, error) {
	identity, err := age.ParseX25519Identity(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key %q: %v", privateKey, err)
	}

	s, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("cannot decode base64 secret %v", err)
	}
	out := &bytes.Buffer{}
	f := strings.NewReader(string(s))

	r, err := age.Decrypt(f, identity)
	if err != nil {
		return "", fmt.Errorf("failed to open encrypted file: %v", err)
	}
	if _, err := io.Copy(out, r); err != nil {
		return "", fmt.Errorf("failed to read encrypted file: %v", err)
	}

	return out.String(), nil
}

func Encrypt(publicKey string, secret string) (string, error) {
	recipient, err := age.ParseX25519Recipient(publicKey)
	if err != nil {
		return "", fmt.Errorf("failed to parse public key %q: %v", publicKey, err)
	}

	buf := &bytes.Buffer{}
	// armorWriter := armor.NewWriter(buf)

	w, err := age.Encrypt(buf, recipient)
	if err != nil {
		return "", fmt.Errorf("failed to create encrypted file: %v", err)
	}
	defer w.Close()

	if _, err := io.WriteString(w, secret); err != nil {
		return "", fmt.Errorf("failed to write to encrypted file: %v", err)
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("failed to close encrypted file: %v", err)
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
