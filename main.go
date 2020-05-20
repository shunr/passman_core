package main

import (
	"bufio"
	"fmt"
	"os"
	"passman_core/api"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')

	api.CreateAccount(username, password)
}
