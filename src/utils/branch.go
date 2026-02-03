package utils

import (
	"fmt"
	"os"
	"strings"
)

func GetCurrentBranch() string {
	content, err := os.ReadFile(".neat/HEAD")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return ""
	}

	line := strings.TrimSpace(string(content))

	parts := strings.Split(line, " ")

	if len(parts) != 2 {
		return ""
	}

	refPath := parts[1]
	refParts := strings.Split(refPath, "/")

	return refParts[len(refParts)-1]
}

func GetCurrentBranchCommitSha() string {
	currentBranch := GetCurrentBranch()
	branchPath := fmt.Sprintf(".neat/refs/heads/%s", currentBranch)

	content, err := os.ReadFile(branchPath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return ""
	}

	return strings.TrimSpace(string(content))

}
