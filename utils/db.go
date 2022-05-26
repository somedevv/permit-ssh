package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/somedevv/permit-ssh/colors"
)

func SetupLocalDB() *bolt.DB {
	// Open the permit.db data file in the data directory.
	// It will be created if it doesn't exist.
	// TODO: Locate automatically the database file
	db, err := bolt.Open(os.Getenv("HOME")+"/.local/bin/.permit_data/users.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		colors.Red.Println(err)
		os.Exit(1)
	}
	// Create bucket if it doesn't exist
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("DataBucket"))
		if err != nil {
			colors.Red.Println("create bucket: %s", err)
			os.Exit(1)
		}
		return nil
	})
	return db
}

func SaveKeyInLocalDB(db *bolt.DB, user, key, ip string) {
	// If user and key exist
	if err := db.Update(func(tx *bolt.Tx) error {
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
			return fmt.Errorf("User already exists, so it will not be added")
		}
		return nil
	}); err != nil {
		colors.Red.Println(err.Error())
		return
	}
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

func RemoveKeyFromLocalDB(db *bolt.DB, user, key string) {
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
}
