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

func Decrypt(privateKey string, secret string, writer io.Writer) error {
	identity, err := age.ParseX25519Identity(privateKey)
	if err != nil {
		return fmt.Errorf("Failed to parse private key %q: %v", privateKey, err)
	}

	s, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return fmt.Errorf("Cannot decode base64 secret %v", err)
	}
	out := &bytes.Buffer{}
	f := strings.NewReader(string(s))

	r, err := age.Decrypt(f, identity)
	if err != nil {
		return fmt.Errorf("Failed to open encrypted file: %v", err)
	}
	if _, err := io.Copy(out, r); err != nil {
		return fmt.Errorf("Failed to read encrypted file: %v", err)
	}

	fmt.Fprintf(writer, "%q\n", out.Bytes())
	return nil
}

func Encrypt(publicKey string, secret string, writer io.Writer) error {
	recipient, err := age.ParseX25519Recipient(publicKey)
	if err != nil {
		return fmt.Errorf("Failed to parse public key %q: %v", publicKey, err)
	}

	buf := &bytes.Buffer{}
	// armorWriter := armor.NewWriter(buf)

	w, err := age.Encrypt(buf, recipient)
	if err != nil {
		return fmt.Errorf("Failed to create encrypted file: %v", err)
	}
	defer w.Close()

	if _, err := io.WriteString(w, secret); err != nil {
		return fmt.Errorf("Failed to write to encrypted file: %v", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("Failed to close encrypted file: %v", err)
	}

	result := base64.StdEncoding.EncodeToString(buf.Bytes())

	fmt.Fprintf(writer, "%s\n", result)
	return nil
}
