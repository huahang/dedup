package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/huahang/dedup/v2/common"
)

func main() {
	target := flag.String(
		"target", "",
		"target sha1sum file",
	)
	flag.Parse()
	// open target file
	targetFile, err := os.Open(*target)
	if err != nil {
		flag.PrintDefaults()
		return
	}
	defer targetFile.Close()
	// scan target file
	checksum_map := make(map[string]map[string][]string)
	common.ScanChecksumFile(targetFile, func(checksum string, path string) {
		if strings.Contains(path, "@") {
			return
		}
		if strings.HasSuffix(path, ".xmp") {
			return
		}
		dir := filepath.Dir(path)
		if checksum_map[checksum] == nil {
			checksum_map[checksum] = make(map[string][]string)
		}
		checksum_map[checksum][dir] = append(checksum_map[checksum][dir], path)
	})
	// dedup
	for _, dir_map := range checksum_map {
		for _, dups := range dir_map {
			sort.Sort(sort.Reverse(sort.StringSlice(dups)))
			fmt.Println("# Potential duplicates:")
			for i := 0; i < len(dups); i++ {
				if i == 0 {
					fmt.Println("# keep", dups[i])
					continue
				}
				fmt.Println("rm", "\""+dups[i]+"\"")
			}
		}
	}
}
