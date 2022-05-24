package cmd

import (
	"os"

	"github.com/boltdb/bolt"
	"github.com/somedevv/permit-ssh/colors"
	"github.com/somedevv/permit-ssh/utils"
)

func AddLocal(db *bolt.DB, user, key, ip string) {
	// If user and key exist
	if user != "" && key != "" {
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

		// If user exist but key not
	} else if user != "" && key == "" {
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
	}

	if ip != "" {
		utils.AddKey(ip, key)
	}
	defer db.Close()
	os.Exit(0)
}
