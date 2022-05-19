package utils

import (
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
