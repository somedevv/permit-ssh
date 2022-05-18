package main

import (
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/integrii/flaggy"
	"github.com/somedevv/permit-ssh/cmd"
	"github.com/somedevv/permit-ssh/colors"
	"github.com/somedevv/permit-ssh/utils"
)

// Version of the program, set at buildtime with -ldflags "-X main.version=X"
var version = ""

var (
	// FLAG VARIABLES
	user string
	key  string
	ip   string

	// SUBCOMMANDS
	delete      *flaggy.Subcommand
	add         *flaggy.Subcommand
	list        *flaggy.Subcommand
	interactive *flaggy.Subcommand
)

func init() {
	//------META------//
	flaggy.SetName("permit")
	flaggy.SetDescription("Your own SSH key manager and friend, made by somedevv")
	flaggy.SetVersion(version)
	flaggy.DefaultParser.AdditionalHelpPrepend = "https://github.com/somedevv/permit-ssh"

	//---------SUBCOMMANDS---------//

	//------DELETE------//
	delete = flaggy.NewSubcommand("remove")
	delete.String(&user, "u", "user", "The user to delete")
	delete.String(&key, "k", "key", "The key to delete")
	delete.String(&ip, "ip", "address", "The IP of the server to delete the user")
	flaggy.AttachSubcommand(delete, 1)

	//------ADD------//
	add = flaggy.NewSubcommand("add")
	add.String(&user, "u", "user", "The user to add")
	add.String(&key, "k", "key", "The key to add")
	add.String(&ip, "ip", "address", "The IP of the server to add the user")
	flaggy.AttachSubcommand(add, 1)

	//------LIST------//
	list = flaggy.NewSubcommand("list")
	flaggy.AttachSubcommand(list, 1)

	//------INTERACTIVE------//
	interactive = flaggy.NewSubcommand("interactive")
	flaggy.AttachSubcommand(interactive, 1)

	//------PARSE------//
	flaggy.Parse()
}

func main() {

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

	if interactive.Used {
		if err := cmd.Interactive_mode(db); err != nil {
			colors.Red.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if list.Used {
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
		os.Exit(0)
	}

	if delete.Used {
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

	if add.Used {
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
	}

	defer db.Close()
	os.Exit(0)
}
