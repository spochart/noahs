package main

import (
	"fmt"
	"os"

	"github.com/spochart/noahs/cmd"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "latest-remotes":
		cmd.RunGitLog(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %q\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: noahs <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  latest-remotes    Show the last commit pushed to remote main/master for each repo")
}
