package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jsndz/neat/src/utils"
)

func Add(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	header := []byte(fmt.Sprintf("blob %d\x00", len(content)))
	blob := append(header, content...)

	sha := utils.Sha1Hash(blob)

	file, objDir := utils.ObjectPath(sha)
	if err := utils.EnsureDir(objDir); err != nil {
		fmt.Println("Error creating object dir:", err)
		return
	}

	blobPath := filepath.Join(objDir, file)

	if _, err := os.Stat(blobPath); err == nil {
		fmt.Println("Object already exists:", sha)
		return
	}

	compressed, err := utils.Compress(blob)
	if err != nil {
		fmt.Println("Compression error:", err)
		return
	}

	if err := os.WriteFile(blobPath, compressed, 0644); err != nil {
		fmt.Println("Write error:", err)
		return
	}

	fmt.Println("Added object:", sha)
}
