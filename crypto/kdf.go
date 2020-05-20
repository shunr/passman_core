package crypto

import (
	"crypto/sha256"
	"io"

	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/scrypt"
)

const KEY_LENGTH int = 32

func HKDF(secret []byte, salt []byte, info []byte) []byte {
	hash_fn := sha256.New
	hash := hkdf.New(hash_fn, secret, salt, info)
	key := make([]byte, KEY_LENGTH)
	if _, err := io.ReadFull(hash, key); err != nil {
		panic(err)
	}
	return key
}

func Scrypt(secret []byte, salt []byte) []byte {
	dk, err := scrypt.Key(secret, salt, 1<<15, 8, 1, KEY_LENGTH)
	if err != nil {
		panic(err)
	}
	return dk
}
