package cmd

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/somedevv/permit-ssh/colors"
	"github.com/somedevv/permit-ssh/utils"
)

func addKeyToServer(ip, key string) {
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

func AddWithIP(ip, key string) {
	addKeyToServer(ip, key)
}

func AddWithAWS(profile, region, instance, key string) {
	if profile == "" && region == "" {
		colors.Red.Println("Error: At least AWS profile or region must be set")
		os.Exit(1)
	}
	ip := utils.GetAWSInstance(profile, region, instance)
	if ip == "" {
		colors.Red.Println("Error: No instance found")
		os.Exit(1)
	}
	addKeyToServer(ip, key)
}
