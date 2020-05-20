package api

import (
	"fmt"
	"github.com/1Password/srp"
	"math/big"
	"passman_core/util"
)

const SECRET_KEY_LEN int = 26

func CreateAccount(username string, password string) error {
	secret_key := util.GenerateSecretKey(username, SECRET_KEY_LEN)
	master_unlock_salt := util.GenerateSalt(32)
	authentication_salt := util.GenerateSalt(32)

	// Generate key to decrypt
	master_unlock_key := util.DeriveKeyFromMasterPasswordAndSecretKey(
		username,
		password,
		secret_key,
		master_unlock_salt)

	// Generate x for SRP
	srp_x := util.DeriveKeyFromMasterPasswordAndSecretKey(
		username,
		password,
		secret_key,
		authentication_salt)

	private_key, public_key := util.GenerateAsymmetricKeyPair()

	x := new(big.Int)
	x.SetBytes(srp_x)

	srp_client := srp.NewSRPClient(srp.KnownGroups[srp.RFC5054Group4096], x, nil)
	if srp_client == nil {
		// TODO: return error here
		panic("couldn't setup client")
	}
	v, err := srp_client.Verifier()
	if err != nil {
		return err
	}

	fmt.Println(private_key, public_key)
	fmt.Println("MUK: ", master_unlock_key)
	fmt.Println("V: ", v)

	// TODO: send: username, display_name, auth_salt_hex, muk_salt_hex, auth_verifier_hex to server

	return nil
}
