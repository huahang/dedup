package common

import (
	"bufio"
	"log"
	"os"
	"regexp"
)

type ChecksumCallback func(checksum string, path string)

func ScanChecksumFile(file *os.File, callback ChecksumCallback) {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		splits := regexp.MustCompile(`\s+`).Split(scanner.Text(), 2)
		if len(splits) != 2 {
			log.Fatal("invalid line: ", line)
		}
		callback(splits[0], splits[1])
	}
}
