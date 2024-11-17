package main

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"strings"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	homeDir := usr.HomeDir
	dupFile := homeDir + "/dup.txt"
	file, err := os.Open(dupFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		splits := strings.Split(line, ":")
		if len(splits) != 2 || splits[0] != "dup" {
			fmt.Println("Invalid line in dup.txt:", line)
			panic("Invalid line in dup.txt")
		}
		splits = strings.Split(splits[1], ",")
		if len(splits) != 2 {
			fmt.Println("Invalid line in dup.txt:", line)
			panic("Invalid line in dup.txt")
		}
		fmt.Printf("rm \"%s\"\n", splits[0])
	}
}
