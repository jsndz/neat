package commands

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/jsndz/neat/src/utils"
)

func Commit(message string) {
	var index utils.Path
	if indexContent, err := os.ReadFile(".neat/index"); err == nil {
		index = utils.ReadIndex(indexContent)
	}

	folders := GetFolders(index)
	treeSha := CreateTree(folders)
	parentCommitSha, exists := utils.GetParentCommit()

	var buf bytes.Buffer
	treeHex := hex.EncodeToString(treeSha)
	buf.WriteString(fmt.Sprintf("tree %s\n", treeHex))

	if exists {
		buf.WriteString(fmt.Sprintf("parent %s\n", parentCommitSha))
	}

	buf.WriteString(utils.AuthorLine() + "\n")
	buf.WriteString(utils.CommitterLine() + "\n\n")
	buf.WriteString(message)

	header := fmt.Sprintf("commit %d\x00", buf.Len())
	content := append([]byte(header), buf.Bytes()...)
	_, commitSha := utils.Sha1Hash(content)
	utils.WriteToObjects((commitSha), content)
	currentBranch := utils.GetCurrentBranch()
	branchPath := fmt.Sprintf(".neat/refs/heads/%s", currentBranch)
	os.WriteFile(branchPath, []byte(commitSha), 0644)

}

func CreateTree(folders utils.Folders) []byte {
	keys := make([]string, 0, len(folders))
	for key := range folders {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return strings.Count(keys[i], "/") > strings.Count(keys[j], "/")
	})
	treeSha := make(map[string][]byte)

	for _, folder := range keys {
		entries := folders[folder]

		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Name < entries[j].Name
		})

		var buff bytes.Buffer

		for _, e := range entries {
			var sha []byte

			if e.Mode == "040000" {
				childPath := e.Name
				if folder != "" {
					childPath = folder + "/" + e.Name
				}
				sha = treeSha[childPath]

			} else {
				sha = e.Sha
			}
			//  structure of tree entry
			// mode name \0sha -> everything in bytes (\0->00)
			entry := fmt.Sprintf("%s %s", e.Mode, e.Name)

			buff.Write([]byte(entry))
			buff.WriteByte(0)
			buff.Write(sha)

		}
		treeContent := buff.Bytes()
		// tree has header tree ,<size> 0
		// then all entries
		header := []byte(fmt.Sprintf("tree %d\x00", len(treeContent)))

		tree := append(header, treeContent...)

		shaBin, sha := utils.Sha1Hash(tree)
		utils.WriteToObjects(sha, tree)
		treeSha[folder] = shaBin
	}

	return treeSha[""]

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
