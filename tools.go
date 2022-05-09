package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func CallClear() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func AddKey(ip string, key string) {
	// Checks if the key already exists
	cmd := exec.Command("ssh", ip, "grep -Fxq "+key+" .ssh/authorized_keys")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		// If the key doesn't exist, the key is added to the authorized_keys file
		fmt.Printf("Adding SSH key to server...\n")
		cmd = exec.Command("ssh", ip, "echo "+key+" >> .ssh/authorized_keys")
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			log.Fatalf("Error while adding key: %s", err)
		}
		fmt.Printf("Key added successfully\n")
	} else {
		// If the key already exists, the program quits
		fmt.Printf("Key already added to server\n")
	}
}
