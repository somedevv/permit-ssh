package cmd

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/somedevv/permit-ssh/colors"
	"github.com/somedevv/permit-ssh/utils"
)

func deleteKeyFromServer(ip, key string) {
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

func DeleteWithIP(ip, key string) {
	deleteKeyFromServer(ip, key)
}

func DeleteWithAWS(profile, region, instance, key string) {
	var ip string

	ip = utils.GetAWSInstance(profile, region, instance)
	if ip == "" {
		colors.Red.Println("Error: No instance found")
		os.Exit(1)
	}
	deleteKeyFromServer(ip, key)
}
