package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	args := os.Args

	switch args[1] {
	case "init":
		initNeat()
	case "add":
		add(args[2])
	case "cat-file":
		cat(args[2:])
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
		if err := os.MkdirAll(folder, 0755); err != nil {
			fmt.Println("Error in creating file", err)
			return
		}
	}

	headFileContent := []byte("ref: refs/heads/main\n")

	if err := os.WriteFile(".neat/HEAD", headFileContent, 0644); err != nil {
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
	header := []byte(fmt.Sprintf("blob %d\x00", size))
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

	if _, err := os.Stat(blobPath); err == nil {
		fmt.Println("Object already exists:", sha)
		return
	}

	var buff bytes.Buffer

	zw := zlib.NewWriter(&buff)
	_, err = zw.Write(blob)
	if err != nil {
		return
	}
	zw.Close()

	if err := os.WriteFile(blobPath, buff.Bytes(), 0644); err != nil {
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

func cat(args []string) {
	pretty := false
	sha := ""
	for _, arg := range args {
		switch arg {
		case "-p":
			pretty = true

		default:
			if sha != "" {
				fmt.Println("Error: multiple object IDs provided")
				return
			}
			sha = arg
		}
	}
	if sha == "" {
		fmt.Println("Error: object ID required")
		return
	}

	if len(sha) < 3 {
		fmt.Println("Error: invalid object ID")
		return
	}
	dir := sha[:2]
	file := sha[2:]
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error in getting working directory", err)
		return
	}
	path := filepath.Join(cwd, ".neat", "objects", dir, file)
	compressed, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file", err)
		return
	}
	zr, err := zlib.NewReader(bytes.NewReader(compressed))
	if err != nil {
		fmt.Println("Error decompressing file", err)
		return
	}
	defer zr.Close()

	var out bytes.Buffer
	_, err = io.Copy(&out, zr)
	if err != nil {
		fmt.Println("Error decompressing file", err)
		return
	}
	if pretty {
		fmt.Println(makePretty(out.Bytes()))
		return
	}
	fmt.Println((out.String()))

}

func makePretty(raw []byte) string {
	nullIndex := bytes.IndexByte(raw, 0)
	if nullIndex == -1 {
		fmt.Println("Error Couldn't find nullIndex, invalid object format")
		return ""
	}
	content := raw[nullIndex+1:]

	return string(content)
}
