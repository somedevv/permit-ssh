package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/boltdb/bolt"
	"github.com/somedevv/permit-ssh/colors"
	"github.com/somedevv/permit-ssh/utils"
)

var (
	// FLAGS
	user        *string
	key         *string
	list_users  *bool
	interactive *bool
	ip          *string
	delete      *bool
)

func init() {
	user = flag.String("user", "", "Username")
	key = flag.String("key", "", "Pub RSA key")
	ip = flag.String("ip", "", "IP address of the server")
	list_users = flag.Bool("list", false, "List all saved users. If set all other flags are ignocolors.red")
	interactive = flag.Bool("i", false, "Activate interactive mode. If set all other flags are ignocolors.red")
	delete = flag.Bool("del", false, "Delete a user or key. If IP is set, the user will be deleted from the server, otherwise, the user will be deleted from the database")
}

func main() {

	flag.Parse()

	// Open the permit.db data file in the data directory.
	// It will be created if it doesn't exist.
	// TODO: Locate automatically the database file
	db, err := bolt.Open(os.Getenv("HOME")+"/.local/bin/.data/permit.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
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

	if *interactive == true {
		if err := interactive_mode(db); err != nil {
			colors.Red.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// If true print all saved users and exit
	// TODO: Make the print prettier
	if *list_users == true {
		db.View(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
			b := tx.Bucket([]byte("DataBucket"))

			c := b.Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				utils.PrintKeyandUser(string(k), string(v))
			}
			return nil
		})
		os.Exit(0)
	}

	if *delete == true {
		if *ip != "" && *key != "" {
			utils.DeleteKey(*ip, *key)
			os.Exit(0)
		} else if *ip != "" && *user != "" {
			db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("DataBucket"))
				u := b.Get([]byte(*user))
				if string(u) == "" {
					colors.Red.Println("User not found")
					os.Exit(1)
				}
				return nil
			})
			utils.DeleteKey(*ip, *key)
			os.Exit(0)
		} else if *user != "" || *key != "" {
			db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("DataBucket"))
				if *user != "" {
					u := b.Get([]byte(*user))
					if string(u) == "" {
						colors.Red.Println("User not found")
						os.Exit(1)
					}
				} else {
					c := b.Cursor()
					for u, k := c.First(); k != nil; u, k = c.Next() {
						if string(k) == *key {
							*user = string(u)
							break
						}
					}
					if *user == "" {
						colors.Red.Println("User not found")
						os.Exit(1)
					}
				}
				err := b.Delete([]byte(*user))
				if err != nil {
					colors.Red.Println(err)
					os.Exit(1)
				}
				return nil
			})
			colors.Green.Println("User removed")
			os.Exit(0)
		} else if *user == "" && *key == "" {
			if *ip == "" {
				colors.Red.Println("You must specify a user or key, and/or IP address")
				os.Exit(1)
			}
			colors.Red.Println("You must specify a user or key")
			os.Exit(1)
		}
	}

	if *user == "" && *key == "" || *ip == "" && (*user == "" && *key == "") {
		fmt.Println("Usage: main.go")
		flag.PrintDefaults()
		os.Exit(0)
	}

	// If *user and *key exist
	if *user != "" && *key != "" {
		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("DataBucket"))
			v := b.Get([]byte(*user))

			// Checks if the *key of the user equals the one stocolors.red in the db
			if string(v) != *key {

				// If the key is empty, the user didn't exist
				// TODO: Prompt user for confirmation
				if string(v) == "" {
					colors.Yellow.Println("Saving user...")

				} else { // Else, the *key stocolors.red is different from the inputed one // TODO: Prompt the user to confirm *key update
					colors.Yellow.Println("Updating user key...")
				}

				// The key is added to the DB and associated with the user
				if err := b.Put([]byte(*user), []byte(*key)); err != nil {
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

		// If *user exist but *key not
	} else if *user != "" && *key == "" {
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("DataBucket"))
			v := b.Get([]byte(*user))

			// If the key is empty the user doesn't exist, the program quits
			// because there is no supplied key to store
			// TODO: Prompt user if they want to add the user by providing the key insted of quitting
			if string(v) == "" {
				colors.Red.Printf("The user ")
				colors.WhiteBold.Print(*user)
				colors.Red.Println(" doesn't exist")
				os.Exit(1)
			}
			*key = string(v)
			return nil
		})
	}

	if *ip != "" {
		utils.AddKey(*ip, *key)
	}

	defer db.Close()
	os.Exit(0)

}

func interactive_mode(db *bolt.DB) error {
	utils.CallClear()

	answers := struct {
		Key          string
		Ip           string
		Confirmation string
	}{}

	err := survey.Ask(utils.SimpleConnection, &answers)
	if err != nil {
		return err
	}

	utils.CallClear()

	utils.PrintKeyandIP(answers.Key, answers.Ip)
	survey.AskOne(&utils.Prompt_confirmation, &answers.Confirmation, survey.WithValidator(survey.Required))

	if answers.Confirmation == "Yes" {
		utils.AddKey(answers.Ip, answers.Key)
	} else {
		colors.Red.Println("Key not added")
	}

	defer db.Close()
	return nil
}
