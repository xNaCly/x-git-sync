package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// gets the changed files using `git status -s`, prefixes files with the correct verb for the given status code, trims filenames
func gitAffectedFiles(conf Config) []string {
	DebugLog(conf, "parsing files and their state")
	out, _ := runCmd([]string{"git", "status", "-s"})
	r := strings.Split(out, "\n")
	res := make([]string, 0)
	for _, file := range r {
		file = strings.TrimSpace(file)
		if len(file) == 0 {
			break
		}
		change := ""
		switch file[0] {
		case 'M':
			change = "modified"
		case 'A':
			change = "added"
		case 'D':
			change = "deleted"
		case 'R':
			change = "renamed"
		case 'C':
			change = "copied"
		case 'U':
			change = "updated but unmerged"
		case '?':
			continue
		}
		if strings.Contains(file, "\"") {
			v, err := strconv.Unquote(strings.TrimSpace(file[1:]))
			if err != nil {
				log.Fatalln("[ERROR] couldn't parse encoded characters: ", err)
			}
			file = " " + v
		}
		res = append(res, strings.TrimSpace(file[1:])+" ("+change+")")
	}
	var c rune
	if len(res) > 1 {
		c = 's'
	}
	DebugLog(conf, fmt.Sprintf("parsed '%d' changed file%c...", len(res), c))
	return res
}

func GitPull(conf Config) {
	_, err := runCmd([]string{"git", "pull"})
	if err != nil {
		log.Println("[WARNING] pulling changes from remote failed: ", err)
		return
	}
	DebugLog(conf, "pulled changes from remote")
}

func GitRepoHasChanges(conf Config) bool {
	DebugLog(conf, "checking if repo has changes...")
	out, err := runCmd([]string{"git", "status", "-s"})
	return err == nil && len(out) != 0
}

func CheckIfGitRepo(conf Config) bool {
	DebugLog(conf, "checking if current directory is a git repository...")
	_, err := runCmd([]string{"git", "status", "-s"})
	return err == nil
}

// adds all changes to the staged area
func GitAdd(conf Config) {
	DebugLog(conf, "adding all changes to the staged area...")
	_, err := runCmd([]string{"git", "add", "-A"})
	if err != nil {
		log.Println("[WARNING] adding to staging area failed: ", err)
	}
	DebugLog(conf, "added changes")
}

// pushes commits to remote
func GitPush(conf Config) {
	DebugLog(conf, "pushing commits to remote...")
	_, err := runCmd([]string{"git", "push"})
	if err != nil {
		log.Println("[WARNING] push to remote failed: ", err)
		return
	}
	log.Println("pushed commits to remote...")
}

// makes a commit
func GitCommit(conf Config) {
	commitContent := generateCommitContent(conf)
	log.Println("new commit:", strconv.Quote(strings.Join(commitContent, " ")))
	_, err := runCmd(commitContent)
	if err != nil {
		log.Println("[WARNING] commiting failed: ", err)
		return
	}
	DebugLog(conf, "made commit")
}

// generates the commit content depending on the configuration made by the user in the Config:
//
// the commit consists of:
// - the commit prefix (AutoCommitPrefix)
// - the current datetime formated according to CommitTitleDateFormat
// - the affected files if AddAffectedFiles is true
func generateCommitContent(conf Config) []string {
	DebugLog(conf, "generating commit content...")
	commitTime := time.Now().Format(conf.CommitDate)
	commitContent := strings.ReplaceAll(conf.CommitFormat, "%date%", commitTime)
	commit := make([]string, 0)
	if conf.AddAffectedFiles {
		affectedFiles := gitAffectedFiles(conf)
		commitContent += "\n\n" + "Affected files:\n" + strings.Join(affectedFiles, "\n")
		commit = append(commit, strings.Split(conf.CommitCommand, " ")...)
	}
	commit = append(commit, commitContent)
	return commit
}
