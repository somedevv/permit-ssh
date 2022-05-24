package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"

	"github.com/boltdb/bolt"
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

func SaveKeyInLocalDB(db *bolt.DB, user, key, ip string) {
	// If user and key exist
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("DataBucket"))
		v := b.Get([]byte(user))

		// Checks if the key of the user equals the one stocolors.red in the db
		if string(v) != key {

			// If the key is empty, the user didn't exist
			// TODO: Prompt user for confirmation
			if string(v) == "" {
				colors.Yellow.Println("Saving user...")

			} else { // Else, the key stocolors.red is different from the inputed one // TODO: Prompt the user to confirm key update
				colors.Yellow.Println("Updating user key...")
			}

			// The key is added to the DB and associated with the user
			if err := b.Put([]byte(user), []byte(key)); err != nil {
				colors.Red.Printf("Error: %s\n", err.Error())
				os.Exit(1)
			}
		} else {
			colors.Red.Println("User already exists")
			os.Exit(0)
		}
		return nil
	})
	colors.Green.Println("User saved successfully")
}

func SearchUserInLocalDB(db *bolt.DB, user string) string {
	var key string
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("DataBucket"))
		v := b.Get([]byte(user))

		// If the key is empty the user doesn't exist, the program quits
		// because there is no supplied key to store
		// TODO: Prompt user if they want to add the user by providing the key insted of quitting
		if string(v) == "" {
			colors.Red.Printf("The user ")
			colors.WhiteBold.Print(user)
			colors.Red.Println(" doesn't exist")
			os.Exit(1)
		}
		key = string(v)
		return nil
	})
	return key
}
