package utils

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func RestoreFromCommit(commitSha string) {
	// get the tree from the commit
	CleanWorkingDir()
	treeSha := GetTreeSha(commitSha)
	treeHex := hex.EncodeToString(treeSha)
	RestoreTree(treeHex, ".")
}

func GetTreeSha(commitSha string) []byte {
	content, err := ReadObject(commitSha)

	if err != nil {
		fmt.Println(err)
	}
	body := StripHeader(content)
	line := strings.Split(string(body), "\n")[0]

	parts := strings.Split(line, " ")
	treeHex := parts[1]

	treeSha, _ := hex.DecodeString(treeHex)
	return (treeSha)
}

func RestoreTree(treeSha string, base string) {
	content, err := ReadObject(treeSha)
	if err != nil {
		fmt.Println(err)
	}
	body := StripHeader(content)
	entries := ParseTree(body)
	for _, entry := range entries {
		if entry.Mode == "100644" {
			shaHex := hex.EncodeToString(entry.Sha)

			content, err := ReadObject((shaHex))
			if err != nil {
				fmt.Println("Error in Reading object")
			}
			body := StripHeader(content)

			fp := filepath.Join(base, entry.Name)
			os.MkdirAll(filepath.Dir(fp), 0755)

			if err = os.WriteFile(fp, body, 0644); err != nil {
				fmt.Println("Error in writing to file ", err)
			}
		} else if entry.Mode == "040000" {

			childShaHex := hex.EncodeToString(entry.Sha)
			childBase := filepath.Join(base, entry.Name)

			os.MkdirAll(childBase, 0755)
			RestoreTree(childShaHex, childBase)

		}
	}

}

func ParseTree(tree []byte) []TreeEntry {
	var entries []TreeEntry
	i := 0
	for i < len(tree) {
		space := bytes.IndexByte(tree[i:], ' ')
		mode := string(tree[i : i+space])
		i = i + space + 1

		null := bytes.IndexByte(tree[i:], 0)
		name := string(tree[i : i+null])
		i = i + null + 1

		sha := make([]byte, 20)
		copy(sha, tree[i:i+20])
		i = i + 20
		entries = append(entries, TreeEntry{
			Name: name,
			Mode: mode,
			Sha:  sha,
		})
	}
	return entries
}

func StripHeader(data []byte) []byte {
	nullIndex := bytes.IndexByte(data, 0)
	return data[nullIndex+1:]
}

func CleanWorkingDir() {
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if path == "." || strings.HasPrefix(path, ".neat") {
			return nil
		}
		os.RemoveAll(path)
		return nil
	})
}
