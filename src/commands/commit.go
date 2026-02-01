package commands

import (
	"strings"

	"github.com/jsndz/neat/src/utils"
)

func Commit(args []string) {

}

func createTree() {

}

func GetFolders(files utils.Path) {
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
}
