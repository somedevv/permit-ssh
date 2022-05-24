package cmd

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/boltdb/bolt"
	"github.com/somedevv/permit-ssh/colors"
	"github.com/somedevv/permit-ssh/utils"
)

func ListLocal(db *bolt.DB) {
	// TODO: Make the print prettier
	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("DataBucket"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			utils.PrintKeyandUser(string(k), string(v))
		}
		return nil
	})
	defer db.Close()
	os.Exit(0)
}

func ListAWS(profile, region string) {
	var cmd *exec.Cmd

	if region == "" {
		cmd = exec.Command("aws", "ec2", "describe-instances", "--query", `Reservations[*].Instances[*].{IP:PrivateIpAddress,Name:Tags[?Key=='Name']|[0].Value}`, "--output", "table", "--profile", profile)
	} else {
		cmd = exec.Command("aws", "ec2", "describe-instances", "--query", `Reservations[*].Instances[*].{IP:PrivateIpAddress,Name:Tags[?Key=='Name']|[0].Value}`, "--output", "table", "--profile", profile, "--region", region)
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		// If the key already exists, the program quits
		colors.Red.Println("Error:", err.Error())
		os.Exit(1)
	}
	print(string(out.Bytes()))
}
