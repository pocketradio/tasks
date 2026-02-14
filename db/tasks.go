package db

import (
	"encoding/binary"
	"time"

	"github.com/boltdb/bolt"
)

var taskBucket = []byte("tasks")
var db *bolt.DB

type Task struct {
	Key   int
	Value string
}

func Init(dbPath string) error { // path will be a file path like todo.db

	// create file if it does not exist :
	var err error // so that it doesnt create a local "db" instance and instead uses the global one defined above.
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	// no error so we update database :
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(taskBucket)
		return err // this will be nil if its a success
	})
}

// tx is a transaction object that bolt gives to talk to the DB inside the transaction.
// tx *bolt.Tx is a pointer to a bolt transaction.

func CreateTask(task string) (int, error) {
	var id int

	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(taskBucket)  // tx.bucket retrieves a bucket by name. if not exists, then nil.
		id64, _ := bucket.NextSequence() //returns auto incr integer for the bucket.
		// nextSequence errors are ignored here since failure would imply DB corruption. not necessary for this appln
		id = int(id64)
		key := itob(id)
		return bucket.Put(key, []byte(task))
	})

	if err != nil {
		return -1, err
	}

	return id, nil
}

func AllTasks() ([]Task, error) {

	var tasks []Task

	err := db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket(taskBucket)

		// iterating through all kv pairs :

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			tasks = append(tasks, Task{
				Key:   btoi(k),
				Value: string(v),
			})
		}

		return nil // nothing in the loop can fail.
	})

	// closure will succeed ie. no errors.
	// But if the db.view itself fails, it can return an error internally . so the below part handles that.
	if err != nil {
		return nil, err
	}

	return tasks, nil

}

func DeleteTask(key int) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		return b.Delete(itob(key)) // delete returns an error ( hover )
	})
}

func itob(id int) []byte {
	b := make([]byte, 8) // 8 bcos its storing a uint64 = 64 bits = 8 bytes
	binary.BigEndian.PutUint64(b, uint64(id))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
