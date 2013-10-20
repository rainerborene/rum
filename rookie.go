package rookie

import (
	"code.google.com/p/go.crypto/pbkdf2"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"strings"
)

type rookie struct {
	SecretKeyBase    []byte
	CookieSalt       []byte
	CookieSaltLength int
	Iterations       int
}

func New(secretKeyBase string) *rookie {
	return &rookie{
		Iterations:       1000,
		SecretKeyBase:    []byte(secretKeyBase),
		CookieSalt:       []byte("encrypted cookie"),
		CookieSaltLength: 64,
	}
}

func (r *rookie) key() []byte {
	return pbkdf2.Key(r.SecretKeyBase, r.CookieSalt, r.Iterations,
		r.CookieSaltLength, sha1.New)
}

func (r *rookie) Decode(cookie string) ([]byte, error) {
	raw, err := base64.StdEncoding.DecodeString(cookie)
	parts := strings.Split(string(raw), "--")
	data, err := base64.StdEncoding.DecodeString(parts[0])
	iv, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(r.key()[:32])
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(data, data)
	return data, nil
}
