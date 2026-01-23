package main

import (
	"fmt"
	"os"

	"github.com/jsndz/neat/src/commands"
)

func main() {
	args := os.Args

	switch args[1] {
	case "init":
		commands.InitNeat()
	case "add":
		commands.Add(args[2])
	case "cat-file":
		commands.CatFile(args[2:])
	case "hash-object":
		commands.HashObject(args[2:])
	default:
		fmt.Println("Unknown command")
		return
	}
}
