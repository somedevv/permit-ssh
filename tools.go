package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
)

func CallClear() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func AddKey(ip string, key string) {
	Yellow.Printf("Checking if key already exists in server [%s]\n", ip)
	// Checks if the key already exists
	cmd := exec.Command("ssh", ip, "grep -Fxq "+key+" .ssh/authorized_keys")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		// If the key doesn't exist, the key is added to the authorized_keys file
		Yellow.Printf("Adding SSH key to server...\n")
		cmd = exec.Command("ssh", ip, "echo "+key+" >> .ssh/authorized_keys")
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			Red.Printf("Error while adding key: %s", err)
			os.Exit(1)
		}
		Green.Printf("Key added successfully\n")
	} else {
		// If the key already exists, the program quits
		Yellow.Printf("Key already added to server\n")
	}
}

func PrintKeyandUser(k string, v string) {
	Green.Print("USER: ")
	White.Printf("%s   ", k)
	Green.Print("KEY: ")
	White.Printf("%s\n", v)
}

func PrintKeyandIP(k string, v string) {
	Green.Print("KEY: ")
	White.Printf("%s   ", k)
	Green.Print("IP: ")
	White.Printf("%s\n", v)
}

func CheckIPAddress(ip string) error {
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("IP Address: %s - Invalid", ip)
	} else {
		return nil
	}
}
