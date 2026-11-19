package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

const pbkdf2Iterations = 600_000

// Sealed is a secret encrypted with AES-256-GCM under a key derived from a
// passphrase. The passphrase is never stored.
type Sealed struct {
	Salt  []byte `json:"salt"`
	Nonce []byte `json:"nonce"`
	Data  []byte `json:"data"`
}

func deriveKey(pass string, salt []byte) ([]byte, error) {
	return pbkdf2.Key(sha256.New, pass, salt, pbkdf2Iterations, 32)
}

func Seal(plain, pass string) (*Sealed, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	key, err := deriveKey(pass, salt)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	return &Sealed{Salt: salt, Nonce: nonce, Data: gcm.Seal(nil, nonce, []byte(plain), nil)}, nil
}

func (s *Sealed) Open(pass string) (string, error) {
	key, err := deriveKey(pass, s.Salt)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	plain, err := gcm.Open(nil, s.Nonce, s.Data, nil)
	if err != nil {
		return "", errors.New("wrong passphrase (or the stored secret is corrupt)")
	}
	return string(plain), nil
}

var cachedPass string

// Passphrase returns $ABCTLX_PASSPHRASE, else asks on the terminal (once per
// run). confirm asks twice, for setting a new one.
func Passphrase(confirm bool) (string, error) {
	if cachedPass != "" {
		return cachedPass, nil
	}
	if v := os.Getenv("ABCTLX_PASSPHRASE"); v != "" {
		cachedPass = v
		return v, nil
	}
	p, err := ReadHidden("Passphrase: ")
	if err != nil {
		return "", err
	}
	if p == "" {
		return "", errors.New("empty passphrase")
	}
	if confirm {
		again, err := ReadHidden("Repeat passphrase: ")
		if err != nil {
			return "", err
		}
		if again != p {
			return "", errors.New("passphrases differ")
		}
	}
	cachedPass = p
	return p, nil
}

// ReadHidden prompts on the terminal without echoing.
func ReadHidden(prompt string) (string, error) {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		return "", fmt.Errorf("%q needs a terminal (or set $ABCTLX_PASSPHRASE / $ABCTLX_CLIENT_SECRET)", strings.TrimSuffix(prompt, ": "))
	}
	fmt.Fprint(os.Stderr, prompt)
	b, err := term.ReadPassword(fd)
	fmt.Fprintln(os.Stderr)
	return string(b), err
}
