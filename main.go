package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/boltdb/bolt"
)

var (
	user        *string
	key         *string
	list_users  *bool
	interactive *bool
	ip          *string
)

// var clear map[string]func() //create a map for storing clear funcs

func init() {
	user = flag.String("user", "", "Username")
	key = flag.String("key", "", "Pub RSA key")
	ip = flag.String("ip", "", "IP address of the server")
	list_users = flag.Bool("list", false, "List all saved users. If set all other flags are ignored")
	interactive = flag.Bool("i", false, "Activate interactive mode. If set all other flags are ignored")
}

func main() {

	// Open the permit.db data file in the data directory.
	// It will be created if it doesn't exist.
	// TODO: Locate automatically the database file
	db, err := bolt.Open("~/.local/bin/.data/permit.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

	// Create bucket if it doesn't exist
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("DataBucket"))
		if err != nil {
			log.Fatalf("create bucket: %s", err)
		}
		return nil
	})

	flag.Parse()

	if *interactive == true {
		if err := interactive_mode(db); err != nil {
			log.Fatal(err)
		}
		return
	}

	// If true print all saved users and exit
	// TODO: Make the print prettier
	if *list_users == true {
		db.View(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
			b := tx.Bucket([]byte("DataBucket"))

			c := b.Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				fmt.Printf("USER: %s   KEY: %s\n", k, v)
			}
			return nil
		})
		return
	}

	if *user == "" && *key == "" || *ip == "" {
		fmt.Println("Usage: main.go")
		flag.PrintDefaults()
		return
	}

	// If *user and *key exist
	if *user != "" && *key != "" {
		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("DataBucket"))
			v := b.Get([]byte(*user))

			// Checks if the *key of the user equals the one stored in the db
			if string(v) != *key {

				// If the key is empty, the user didn't exist
				if string(v) == "" {
					fmt.Printf("Saving user...\n")

				} else { // Else, the *key stored is different from the inputed one // TODO: Promt the user to confirm *key update
					fmt.Printf("Updating user key...\n")
				}

				// The key is added to the DB and associated with the user
				if err := b.Put([]byte(*user), []byte(*key)); err != nil {
					log.Fatalf("Error: %s", err.Error())
				}
			}
			return nil
		})
		fmt.Printf("User saved successfully\n")
		// If *user exist but *key not
	} else if *user != "" && *key == "" {
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("DataBucket"))
			v := b.Get([]byte(*user))

			// If the key is empty the user doesn't exist, the program quits
			// because there is no supplied key to store
			// TODO: Promt user if they want to add the user by providing the key insted of quitting
			if string(v) == "" {
				log.Fatalf("The user %s doesn't exist", *user)
			}
			*key = string(v)
			return nil
		})
	}

	fmt.Printf("Checking if key already exists...\n")

	AddKey(*ip, *key)

	defer db.Close()
	return

}

func interactive_mode(db *bolt.DB) error {
	CallClear()

	answers := struct {
		Key          string
		Ip           string
		Confirmation string
	}{}

	err := survey.Ask(SimpleConnection, &answers)
	if err != nil {
		return err
	}

	CallClear()

	fmt.Printf("USER: %s   KEY: %s\n", answers.Key, answers.Ip)

	survey.AskOne(&prompt_confirmation, &answers.Confirmation, survey.WithValidator(survey.Required))

	if answers.Confirmation == "Yes" {
		AddKey(answers.Ip, answers.Key)
	} else {
		fmt.Printf("Key not added\n")
	}

	defer db.Close()
	return nil
}
