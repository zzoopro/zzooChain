package db

import (
	"fmt"

	"github.com/zzoopro/zzoocoin/utils"
	bolt "go.etcd.io/bbolt"
)

const (
	DB_NAME = "blockchain"
	DATA_BUCKET = "data"
	BLOCKS_BUCKET = "blocks"
	DATA_BUCKET_KEY = "blockchain_data"
)

var db *bolt.DB

func getDbName() string {
	return fmt.Sprintf("%s_%s.db", DB_NAME, utils.GetPort())
}

func DB() *bolt.DB {
	if db == nil {
		utils.GetPort()
		dbPointer, err := bolt.Open(getDbName(), 0600, nil)
		utils.HandleErr(err)
		db = dbPointer
		err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(DATA_BUCKET))			
			_, err = tx.CreateBucketIfNotExists([]byte(BLOCKS_BUCKET))			
			return err
		})
		utils.HandleErr(err)
	}
	return db
}

func Close() {
	DB().Close()
}

func SaveBlock(hash string, data []byte) {
	// fmt.Printf("Saving block: %s\n data: %b", hash, data)
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS_BUCKET))
		err := bucket.Put([]byte(hash), data)
		return err
	})	
	utils.HandleErr(err)
}

func SaveBlockchain(data []byte){
 	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(DATA_BUCKET))
		err := bucket.Put([]byte(DATA_BUCKET_KEY), data)
		return err
	})
	utils.HandleErr(err)
}

func ChainData() []byte {
	var data []byte
	DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(DATA_BUCKET))
		data = bucket.Get([]byte(DATA_BUCKET_KEY))
		return nil
	})
	return data
}

func Block(hash string) []byte {
	var data []byte
	DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS_BUCKET))
		data = bucket.Get([]byte(hash))
		return nil
	})
	return data
}

func EmptyBlocks() {
	DB().Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte(BLOCKS_BUCKET))
		_, err := tx.CreateBucket([]byte(BLOCKS_BUCKET))
		utils.HandleErr(err)
		return nil
	})
}
