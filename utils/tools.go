package utils

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"

	"github.com/somedevv/permit-ssh/colors"
)

func CallClear() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func AddKey(ip string, key string) {
	colors.Yellow.Printf("Checking if key already exists in server [%s]\n", ip)
	// Checks if the key already exists
	cmd := exec.Command("ssh", ip, "grep -Fxq '"+key+"' .ssh/authorized_keys")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		// If the key doesn't exist, the key is added to the authorized_keys file
		colors.Yellow.Println("Adding SSH key to server...")
		cmd = exec.Command("ssh", ip, "echo "+key+" >> .ssh/authorized_keys")
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			colors.Red.Printf("Error while adding key: %s", err)
			os.Exit(1)
		}
		colors.Green.Println("Key added successfully")
	} else {
		// If the key already exists, the program quits
		colors.Yellow.Println("Key already added to server")
	}
}

func DeleteKey(ip string, key string) {
	colors.Yellow.Printf("Checking if key exists in [%s]\n", ip)
	// Checks if the key already exists
	cmd := exec.Command("ssh", ip, "grep -Fxq '"+key+"' .ssh/authorized_keys")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err == nil {
		// If the key exists, it's removed from the authorized_keys file
		colors.Yellow.Println("Removing key...")
		cmd = exec.Command("ssh", ip, `sudo sed '\%`+key+"%"+" d' .ssh/authorized_keys > .ssh/authorized_keys.tmp && sudo mv .ssh/authorized_keys.tmp .ssh/authorized_keys")
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			colors.Red.Printf("Error while trying to remove the key: %s\n", err)
			os.Exit(1)
		}
		colors.Green.Println("Key removed successfully")
	} else {
		// If the key doesn't exist, the program quits
		colors.Yellow.Printf("Key not found in [%s]\n", ip)
	}
}

func PrintKeyandUser(k string, v string) {
	colors.Green.Print("USER: ")
	colors.White.Printf("%s   ", k)
	colors.Green.Print("KEY: ")
	colors.White.Printf("%s\n", v)
}

func PrintKeyandIP(k string, v string) {
	colors.Green.Print("KEY: ")
	colors.White.Printf("%s   ", k)
	colors.Green.Print("IP: ")
	colors.White.Printf("%s\n", v)
}

func CheckIPAddress(ip string) error {
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("IP Address: %s - Invalid", ip)
	} else {
		return nil
	}
}
