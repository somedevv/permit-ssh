package cmd

import (
	"os"

	"github.com/boltdb/bolt"
	"github.com/somedevv/permit-ssh/utils"
)

func List(db *bolt.DB) {
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
