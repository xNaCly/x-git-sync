package main

import (
	"fmt"
	"log"
	"os/exec"
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
	DebugLog(conf, fmt.Sprintf("parsed '%d' changed files...", len(res)))
	return res
}

func GitPull(conf Config){
	DebugLog(conf, "pulling changes from remote...")
	_, err := runCmd([]string{"git", "pull"})
	if err != nil {
		log.Println("[WARNING] pulling changes from remote failed: ", err)
		return
	}
	DebugLog(conf, "pulled changes from remote")
}

func GitRepoHasChanges(conf Config) bool {
	DebugLog(conf, "checking if repo has changes...");
    out, err := runCmd([]string{"git", "status", "-s"})
    return err == nil && len(out) != 0;
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
	log.Println("[INFO][PUSH]: pushed commits to remote...")
}

// makes a commit
func GitCommit(conf Config) {
	DebugLog(conf, "making commit...")
	commitContent := generateCommitContent(conf)
	log.Println("[INFO][COMMIT]:", strconv.Quote(strings.Join(commitContent, " ")))
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
	commitTime := time.Now().Format(conf.CommitTitleDateFormat)
	commitContent := conf.AutoCommitPrefix + commitTime
	commit := make([]string, 0)
	if conf.AddAffectedFiles {
		affectedFiles := gitAffectedFiles(conf)
		commitContent += "\n\n" + "Affected files:\n" + strings.Join(affectedFiles, "\n")
		commit = append(commit, strings.Split(conf.CommitCommand, " ")...)
	}
	commit = append(commit, commitContent)
	DebugLog(conf, fmt.Sprintf("generated commit content. (%s)", strconv.Quote(strings.Join(commit, " "))))
	return commit;
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
