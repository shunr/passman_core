package main

import (
	"bufio"
	"fmt"
	"github.com/shunr/passman_core/api"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')

	client := api.NewPassmanClient(true, "", "localhost:13222")
	client.CreateAccount(username, password)
	defer client.Close()
}
