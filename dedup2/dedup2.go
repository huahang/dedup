package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/huahang/dedup/v2/common"
)

func PrintUsage() {
	flag.PrintDefaults()
}

func main() {
	// parse flags
	base := flag.String(
		"base", "",
		"base sha1sum file",
	)
	target := flag.String(
		"target", "",
		"target sha1sum file",
	)
	flag.Parse()
	// open base file
	baseFile, err := os.Open(*base)
	if err != nil {
		flag.PrintDefaults()
		return
	}
	defer baseFile.Close()
	// open target file
	targetFile, err := os.Open(*target)
	if err != nil {
		flag.PrintDefaults()
		return
	}
	defer targetFile.Close()
	// scan base file
	base_map := make(map[string]string)
	common.ScanChecksumFile(baseFile, func(checksum string, path string) {
		base_map[checksum] = path
	})
	// scan target file
	common.ScanChecksumFile(targetFile, func(checksum string, path string) {
		if _, ok := base_map[checksum]; ok {
			fmt.Println("rm", "\""+path+"\"")
		}
	})
}
