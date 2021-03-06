package util

import (
	"crypto/rsa"
	"github.com/shunr/passman_core/crypto"
	"golang.org/x/text/unicode/norm"
	"strings"
)

const VERSION_INFO string = "00"

func GenerateSalt(length int) []byte {
	return crypto.RandBytes(length)
}

func GenerateSecretKey(username string, key_length int) string {
	unambiguous_32_alphanumeric := "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"
	random_bytes := crypto.RandBytes(key_length)
	random_chars := make([]byte, key_length)
	for i := 0; i < key_length; i++ {
		random_chars[i] = unambiguous_32_alphanumeric[random_bytes[i]%32]
	}
	return strings.ToUpper(VERSION_INFO + (username + "XXXXXX")[:6] + string(random_chars))
}

func GenerateAsymmetricKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	return crypto.GenerateRSAKeyPair(4096)
}

func DeriveKeyFromMasterPasswordAndSecretKey(
	username string, master_password string, secret_key string, salt []byte) []byte {
	password_normalized := []byte(norm.NFD.String(strings.TrimSpace(master_password)))
	master_key_salt := crypto.HKDF(salt, []byte(username), []byte(VERSION_INFO))
	master_key_derived := crypto.Scrypt(password_normalized, master_key_salt)
	secret_key_derived := crypto.HKDF([]byte(secret_key), []byte(username), []byte(VERSION_INFO))

	derived_key := Bxor(master_key_derived, secret_key_derived)

	return derived_key
}
