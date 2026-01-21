package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args

	switch args[1] {
	case "init":
		initNeat()

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
		if err := os.Mkdir(folder, 0755); err != nil {
			fmt.Println("Error in creating file", err)
			return
		}
	}

	headFileContent := []byte("ref: refs/heads/main\n")

	if err := os.WriteFile(".neat/HEAD", headFileContent, 0755); err != nil {
		fmt.Println("Error in writing HEAD", err)
		return
	}
	fmt.Println("Initialized neat repository")

}
