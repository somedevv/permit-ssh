package cmd

import (
	"os"

	"github.com/boltdb/bolt"
	"github.com/somedevv/permit-ssh/colors"
	"github.com/somedevv/permit-ssh/utils"
)

func Remove(db *bolt.DB, user, key, ip string) {
	if ip != "" && key != "" {
		utils.DeleteKey(ip, key)
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
		utils.DeleteKey(ip, key)
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
}
