package main

import (
	"log"
	"os/exec"
	"strings"
	"time"
)

// gets the changed files using `git status -s`, prefixes files with the correct verb for the given status code, trims filenames
func gitAffectedFiles() []string {
    out, _ := runCmd([]string{"git", "status", "-s"})
    r := strings.Split(out, "\n")
    res := make([]string, 0)
    for _, file := range r {
        if len(file) == 0 {
            break
        }
        file = strings.TrimSpace(file)
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
        res = append(res, change+": "+strings.TrimSpace(file[1:]))
    }
    return res
}

// adds all changes to the staged area
func GitAdd() {
    _, err := runCmd([]string{"git", "add", "-A"})
    if err != nil {
        log.Println("[ERR]", err)
    }
}

// pushes commits to remote
func GitPush() {
    out, err := runCmd([]string{"git", "push"})
    if err != nil {
        log.Println("[ERR]", err)
    }
    log.Println("[INF][PUSH]", out)
}

// makes a commit depending on the configuration made by the user in the Config:
// 
// the commit consists of:
// - the commit prefix (AutoCommitPrefix)
// - the current datetime formated according to CommitTitleDateFormat
// - the affected files if AddAffectedFiles is true
func GitCommit(conf Config) bool {
    commitTime := time.Now().Format(conf.CommitTitleDateFormat)
    commitTitle := conf.AutoCommitPrefix + commitTime
    commit := ""
    if conf.AddAffectedFiles {
        affectedFiles := gitAffectedFiles()
        commit = commitTitle + "\n" + "Affected files:\n" + strings.Join(affectedFiles, "\n")
    } else {
        commit = commitTitle
    }
    log.Println("[INF][COMMIT]", strings.ReplaceAll(commit, "\n", "\\n"))
    // TODO: implement running this command
    return false
}

// executes command, trims output and returns it
func runCmd(cmd []string) (val string, err error) {
    inp := append(cmd)
    command := exec.Command(inp[0], inp[1:]...)
    out, err := command.CombinedOutput()
    if err != nil {
        return "", err
    }
    return strings.TrimSpace(string(out)), nil
}
