package commands

import (
	"fmt"
	"os"

	"github.com/jsndz/neat/src/utils"
)

func Switch(branch string) {
	headFileContent := fmt.Sprintf("ref: refs/heads/%s", branch)
	path := fmt.Sprintf(".neat/refs/heads/%s", branch)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Branch does not exist")
		return
	}

	commitSha := utils.GetBranchCommitSha(branch)
	if commitSha != "" {
		utils.RestoreFromCommit(commitSha)
	}

	if err := os.WriteFile(".neat/HEAD", []byte(headFileContent), 0644); err != nil {
		fmt.Println("Error in writing HEAD", err)
		return
	}
}
