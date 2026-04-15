package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func RunGitLog(args []string) {
	fs := flag.NewFlagSet("git-log", flag.ExitOnError)
	root := fs.String("dir", ".", "root directory to search for git repos")
	fs.Parse(args)

	repos, err := findGitRepos(*root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if len(repos) == 0 {
		fmt.Println("no git repositories found under", *root)
		return
	}

	for _, repo := range repos {
		printRemoteLog(repo)
	}
}

func findGitRepos(root string) ([]string, error) {
	var repos []string

	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		path := filepath.Join(root, e.Name())
		if _, err := os.Stat(filepath.Join(path, ".git")); err == nil {
			repos = append(repos, path)
		}
	}

	return repos, nil
}

func getRemotes(repoPath string) ([]string, error) {
	out, err := exec.Command("git", "-C", repoPath, "remote").Output()
	if err != nil {
		return nil, err
	}

	var remotes []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line != "" {
			remotes = append(remotes, line)
		}
	}
	return remotes, nil
}

func remoteHead(repoPath string, remote string) string {
	for _, branch := range []string{"main", "master"} {
		ref := remote + "/" + branch
		out, err := exec.Command(
			"git", "-C", repoPath,
			"rev-parse", "--verify", "--quiet", ref,
		).Output()
		if err == nil && len(out) > 0 {
			return ref
		}
	}
	return ""
}

func printRemoteLog(repoPath string) {
	repoName := filepath.Base(repoPath)

	remotes, err := getRemotes(repoPath)
	if err != nil || len(remotes) == 0 {
		fmt.Printf("%-40s  (no remotes configured)\n", repoName)
		return
	}

	exec.Command("git", "-C", repoPath, "fetch", "--all", "--quiet").Run()

	printed := false
	for _, remote := range remotes {
		branch := remoteHead(repoPath, remote)
		if branch == "" {
			fmt.Printf("%-40s  %-20s  (no main/master branch found)\n", repoName, remote)
			continue
		}

		out, err := exec.Command(
			"git", "-C", repoPath,
			"log", branch,
			"-1",
			"--pretty=format:%h  %ai  %an  %s",
		).Output()

		if err != nil || len(out) == 0 {
			fmt.Printf("%-40s  %-20s  (could not read log)\n", repoName, remote)
			continue
		}

		fmt.Printf("%-40s  %-20s  %s\n",
			repoName,
			remote,
			strings.TrimSpace(string(out)),
		)
		printed = true
	}

	if !printed {
		fmt.Printf("%-40s  (no remote commits found)\n", repoName)
	}
}
