package main

import (
	"fmt"
	"os"

	"github.com/jsndz/neat/src/commands"
	"github.com/jsndz/neat/src/utils"
)

func main() {
	args := os.Args

	switch args[1] {
	case "init":
		commands.InitNeat()
	case "add":
		if args[2] == "." {
			commands.AddAll()
		} else {
			commands.Add(args[2])
		}
	case "cat-file":
		commands.CatFile(args[2:])
	case "hash-object":
		commands.HashObject(args[2:])
	case "ls-tree":
		commands.LsTree()
	case "commit":
		if len(args) >= 4 && args[2] == "-m" {
			commands.Commit(args[3])
		} else {
			fmt.Print("Enter the commit message: ")
			message := utils.ReadInput()

			if message == "" {
				fmt.Println("Empty commit message. Aborting.")
				return
			}

			commands.Commit(message)
		}
	case "branch":
		commands.Branch(args[2])
	case "switch":
		commands.Switch(args[2])
	case "checkout":
		commands.Switch(args[2])
	case "clone":
		commands.Clone(args[2])
	default:
		fmt.Println("Unknown command")
		return
	}
}
