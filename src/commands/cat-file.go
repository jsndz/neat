package commands

import (
	"fmt"

	"github.com/jsndz/neat/src/flags"
	"github.com/jsndz/neat/src/utils"
)

func CatFile(args []string) {
	pretty := false
	sha := ""

	for _, arg := range args {
		switch arg {
		case "-p":
			pretty = true
		default:
			if sha != "" {
				fmt.Println("Error: multiple object IDs provided")
				return
			}
			sha = arg
		}
	}

	if sha == "" || len(sha) < 3 {
		fmt.Println("Invalid object ID")
		return
	}

	raw, err := utils.ReadObject(sha)
	if err != nil {
		fmt.Println("Error reading object:", err)
		return
	}

	if pretty {
		fmt.Print(flags.MakePretty(raw))
		return
	}

	fmt.Print(string(raw))
}
