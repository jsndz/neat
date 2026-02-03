package commands

import (
	"fmt"
	"os"

	"github.com/jsndz/neat/src/utils"
)

func Branch(name string) {
	path := fmt.Sprintf(".neat/refs/heads/%s", name)

	if _, err := os.Stat(path); err == nil {
		fmt.Println("Branch already exists")
		return
	}
	sha := utils.GetCurrentBranchCommitSha()

	os.WriteFile(path, []byte(sha), 0644)
}
