package cmd

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/boltdb/bolt"
	"github.com/somedevv/permit-ssh/colors"
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

func RemoveLocal(db *bolt.DB, user, key, ip string) {
	if ip != "" && key != "" {
		deleteKeyFromServer(ip, key)
		os.Exit(0)
	} else if ip != "" && user != "" {
		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("DataBucket"))
			k := b.Get([]byte(user))
			if string(k) == "" {
				colors.Red.Println("User not found")
				os.Exit(1)
			}
			key = string(k)
			return nil
		})
		deleteKeyFromServer(ip, key)
		os.Exit(0)
	} else if user != "" || key != "" {
		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("DataBucket"))
			if user != "" {
				u := b.Get([]byte(user))
				if string(u) == "" {
					colors.Red.Println("User not found")
					os.Exit(1)
				}
			} else {
				c := b.Cursor()
				for u, k := c.First(); k != nil; u, k = c.Next() {
					if string(k) == key {
						user = string(u)
						break
					}
				}
				if user == "" {
					colors.Red.Println("User not found")
					os.Exit(1)
				}
			}
			err := b.Delete([]byte(user))
			if err != nil {
				colors.Red.Println(err)
				os.Exit(1)
			}
			return nil
		})
		colors.Green.Println("User removed")
		os.Exit(0)
	} else if user == "" && key == "" {
		if ip == "" {
			colors.Red.Println("You must specify a user or key, and/or IP address")
			os.Exit(1)
		}
		colors.Red.Println("You must specify a user or key")
		os.Exit(1)
	}
	defer db.Close()
	os.Exit(0)
}
