package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
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
	count := 0
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
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("hash2path"))
			if b == nil {
				return nil
			}
			value := b.Get([]byte(hashString))
			if value != nil {
				fmt.Printf("dup:%s,%s\n", path, value)
			}
			return nil
		})
		count++
		return nil
	})
	fmt.Printf("%d files proccessed\n", count)
	if err != nil {
		fmt.Println("Error walking through directory:", err)
	}
}
