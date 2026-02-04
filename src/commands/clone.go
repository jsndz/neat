package commands

import (
	"github.com/jsndz/neat/src/utils"
)

func Clone(url string) {
	utils.Download(url, "neat.zip")
	utils.Unzip("neat.zip", ".")
	commitSha := utils.GetCurrentBranchCommitSha()
	if commitSha != "" {
		utils.RestoreFromCommit(commitSha)
	}
}
