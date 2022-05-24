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

func GetAWSInstance(profile, region, instance string) string {
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
