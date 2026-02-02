package commands

import (
	"fmt"
	"strings"

	"github.com/jsndz/neat/src/utils"
)

func Commit(args []string) {

}

func CreateTree() {

}

func GetFolders(files utils.Path) utils.Folders {
	// file path -> sha
	folders := make(map[string][]utils.TreeEntry)

	for p, sha := range files {
		parts := strings.Split(p, "/")
		folder := strings.Join(parts[:len(parts)-1], "/")
		file := parts[len(parts)-1]
		folders[folder] = append(folders[folder],
			utils.TreeEntry{
				Name: file,
				Sha:  sha,
				Mode: "100644",
			})
	}
	// folders folder ->  [...files]
	for folder := range folders {
		parts := strings.Split(folder, "/")
		parent := strings.Join(parts[:len(parts)-1], "/")
		child := parts[len(parts)-1]
		exists := false

		for _, e := range folders[parent] {
			if e.Name == child && e.Mode == "040000" {
				exists = true
				break
			}
		}
		if !exists {
			folders[parent] = append(folders[parent],
				utils.TreeEntry{
					Name: child,
					Mode: "040000",
				},
			)
		}
	}
	fmt.Print(folders)
	return folders
}
