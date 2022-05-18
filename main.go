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
	remove      *flaggy.Subcommand
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
	remove = flaggy.NewSubcommand("remove")
	remove.String(&user, "u", "user", "The user to remove")
	remove.String(&key, "k", "key", "The key to remove")
	remove.String(&ip, "ip", "address", "The IP of the server to remove the user")
	flaggy.AttachSubcommand(remove, 1)

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

	if remove.Used {
		cmd.Remove(db, user, key, ip)
	}

	if add.Used {
		cmd.Add(db, user, key, ip)
	}

	defer db.Close()
	os.Exit(0)
}
