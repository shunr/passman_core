package api

import (
	"context"
	"fmt"
	"github.com/1Password/srp"
	pb "github.com/shunr/passman_core/proto"
	"github.com/shunr/passman_core/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
	"log"
	"math/big"
	"time"
)

const SERVER_HOST_OVERRIDE string = "x.test.youtube.com"
const SECRET_KEY_LEN int = 26

type PassmanClient struct {
	tls         bool
	ca_file     string
	server_addr string
	client      *pb.PassmanClient
	conn        *grpc.ClientConn
}

func NewPassmanClient(tls bool, ca_file string, server_addr string) PassmanClient {
	var opts []grpc.DialOption
	if tls {
		if ca_file == "" {
			ca_file = testdata.Path("ca.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(ca_file, SERVER_HOST_OVERRIDE)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(server_addr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	client := pb.NewPassmanClient(conn)
	return PassmanClient{tls, ca_file, server_addr, &client, conn}
}

func (c *PassmanClient) Close() {
	c.conn.Close()
}

func (c *PassmanClient) CreateAccount(username string, password string) error {
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
	v, err := srp_client.Verifier()
	if err != nil {
		return err
	}

	fmt.Println(private_key, public_key)
	fmt.Println("MUK: ", master_unlock_key)
	fmt.Println("V: ", v)

	// TODO: send: username, display_name, auth_salt_hex, muk_salt_hex, auth_verifier_hex to server

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// TODO: Timeout constant
	defer cancel()
	request := pb.CreateAccountRequest{Username: username} // TODO: Populate this object
	response, err := (*c.client).CreateAccount(ctx, &request)
	if err != nil {
		log.Fatalf("%v.CreateAccount(_) = _, %v: ", c.client, err)
	}
	log.Println(response)

	return nil
}
