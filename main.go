package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	args := os.Args

	switch args[1] {
	case "init":
		initNeat()
	case "add":
		add(args[2])
	case "cat-file":
		cat(args[2], args[3])
	default:
		fmt.Println("Unknown command")
		return
	}
}

func initNeat() {
	if _, err := os.Stat(".neat"); err == nil {
		fmt.Println(".neat already exists")
		return
	}
	folders := []string{".neat", ".neat/refs", ".neat/objects"}
	for _, folder := range folders {
		if err := os.Mkdir(folder, 0755); err != nil {
			fmt.Println("Error in creating file", err)
			return
		}
	}

	headFileContent := []byte("ref: refs/heads/main\n")

	if err := os.WriteFile(".neat/HEAD", headFileContent, 0755); err != nil {
		fmt.Println("Error in writing HEAD", err)
		return
	}
	fmt.Println("Initialized neat repository")

}

func add(filename string) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error in getting working directory", err)
		return
	}
	path := filepath.Join(cwd, filename)
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file", err)
		return
	}
	size := len(content)
	header := []byte("blob " + strconv.Itoa(size) + "\x00")
	// \x00 is the null byte — a single byte with value 0.
	blob := append(header, content...)
	sha := sha1Hash(blob)
	dir := sha[:2]
	file := sha[2:]

	objectsDir := filepath.Join(cwd, ".neat", "objects")
	if err := os.MkdirAll(objectsDir, 0755); err != nil {
		fmt.Println("Error in creating objects", err)
		return
	}

	objDir := filepath.Join(objectsDir, dir)
	if err := os.MkdirAll(objDir, 0755); err != nil {
		fmt.Println("Error in creating object", err)
		return
	}

	blobPath := filepath.Join(objDir, file)
	if err := os.WriteFile(blobPath, blob, 0644); err != nil {
		fmt.Println("Error in writing to file", err)
		return
	}
	fmt.Println("Added object:", sha)

}

func sha1Hash(data []byte) string {
	h := sha1.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func cat() {

}
