package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"

	"github.com/somedevv/permit-ssh/colors"
	"github.com/somedevv/permit-ssh/models"
)

func CallClear() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func AddKey(ip, key string) {
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

func DeleteKey(ip, key string) {
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

func PrintKeyandUser(k, v string) {
	colors.Green.Print("USER: ")
	colors.White.Printf("%s   ", k)
	colors.Green.Print("KEY: ")
	colors.White.Printf("%s\n", v)
}

func PrintKeyandIP(k, v string) {
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

func GetAWSInstance(profile, region, instance, key string) string {
	var cmd *exec.Cmd
	var data []models.EC2Instance

	if region == "" {
		cmd = exec.Command("aws", "ec2", "describe-instances", "--query", `Reservations[*].Instances[].{IP:PrivateIpAddress,Name:Tags[?Key=='Name']|[0].Value}`, "--output", "json", "--profile", profile)
	} else {
		cmd = exec.Command("aws", "ec2", "describe-instances", "--query", `Reservations[*].Instances[].{IP:PrivateIpAddress,Name:Tags[?Key=='Name']|[0].Value}`, "--output", "json", "--profile", profile, "--region", region)
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		// If the key already exists, the program quits
		colors.Red.Println("Error:", err.Error())
		os.Exit(1)
	}

	if err := json.Unmarshal(out.Bytes(), &data); err != nil {
		colors.Red.Println("Error:", err.Error())
		os.Exit(1)
	}

	for _, v := range data {
		if v.Name == instance {
			return v.IP
		}
	}
	return ""
}
