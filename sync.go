package main

import (
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// gets the changed files using `git status -s`, prefixes files with the correct verb for the given status code, trims filenames
func gitAffectedFiles() []string {
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
	return res
}

// adds all changes to the staged area
func GitAdd() {
	_, err := runCmd([]string{"git", "add", "-A"})
	if err != nil {
		log.Println("[WARNING]", err)
	}
}

// pushes commits to remote
func GitPush() {
	out, err := runCmd([]string{"git", "push"})
	if err != nil {
		log.Println("[WARNING]", err)
	}
	log.Println("[INFO][PUSH]:", out)
}

// makes a commit
func GitCommit(conf Config) bool {
	commitContent := generateCommitContent(conf)
	log.Println("[INFO][COMMIT]:\n", strings.Join(commitContent, " "))
	_, err := runCmd(commitContent)
	return err == nil
}

// generates the commit content depending on the configuration made by the user in the Config:
//
// the commit consists of:
// - the commit prefix (AutoCommitPrefix)
// - the current datetime formated according to CommitTitleDateFormat
// - the affected files if AddAffectedFiles is true
func generateCommitContent(conf Config) []string {
	commitTime := time.Now().Format(conf.CommitTitleDateFormat)
	commitContent := conf.AutoCommitPrefix + commitTime
	commit := make([]string, 0)
	if conf.AddAffectedFiles {
		affectedFiles := gitAffectedFiles()
		commitContent += "\n\n" + "Affected files:\n" + strings.Join(affectedFiles, "\n")
		commit = append(commit, strings.Split(conf.CommitCommand, " ")...)
	}
	return append(commit, commitContent)
}

// executes command, trims output and returns it
func runCmd(cmd []string) (val string, err error) {
	command := exec.Command(cmd[0], cmd[1:]...)
	out, err := command.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
