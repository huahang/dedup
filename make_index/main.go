package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"

	bolt "go.etcd.io/bbolt"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	homeDir := usr.HomeDir
	db, err := bolt.Open(homeDir+"/index.db", 0644, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	counter := 0
	err = filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error walking through directory:", err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		hash := md5.New()
		limitedReader := io.LimitReader(file, 8192)
		if _, err := io.Copy(hash, limitedReader); err != nil {
			return err
		}
		hashInBytes := hash.Sum(nil)
		hashString := hex.EncodeToString(hashInBytes)
		err = db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("hash2path"))
			if err != nil {
				return err
			}
			err = b.Put([]byte(hashString), []byte(path))
			return err
		})
		if err != nil {
			panic(err)
		}
		counter++
		if counter%100 == 0 {
			log.Printf("Indexed %d files", counter)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error walking through directory:", err)
	}
}
