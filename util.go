package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

type Config struct {
	// specifies the format of the commit message, default: "backup: %date%"
	// currently supportes:
	// - %date%: the date of the commit, formatted as specified in Commit_date
	CommitFormat string `json:"commit_format"`

	// specifies the date format which the date will be formatted as, default: "2006-01-02 15:04:05"
	//
	//  - 2006 for the year, 06 would only be the last two integer
	//  - 01 for the month
	//  - 02 for the day
	//  - 15 for the hour (24-hour format), 05 for 12-hour format
	//  - 04 for the minute
	//  - 05 for the second
	//
	// time formatting in go is weird, see docs:
	//
	// https://www.digitalocean.com/community/tutorials/how-to-use-dates-and-times-in-go
	CommitDate string `json:"commit_date"`

	// List filenames affected by the commit in the commit body, default: true
	AddAffectedFiles bool `json:"add_affected_files"`

	// time interval between backups (in s), default: 300
	BackupInterval int `json:"backup_interval"`

	// commit command, default: "git commit -m"
	CommitCommand string `json:"commit_cmd"`

	// enables debug mode (verbose logging, extra infos, etc.), default: false
	DebugMode bool `json:"debug"`

	// enable pulling from remote on start, default: true
	PullOnStart bool `json:"pull_on_start"`
}

// Loads and parses config from:
// - On Unix systems, $XDG_CONFIG_HOME or $HOME/.config
// - On Darwin, it returns $HOME/Library/Application Support
// - On Windows, it returns %AppData%
// - On Plan 9, it returns $home/lib
//
// config file location depends on os.UserConfigDir()
func getConfig() Config {
	// all occuring errors are logged, but not treated like panics, due to the fact that a fallback config is provided
	fallbackConf := Config{
		CommitFormat:     "backup: %date%",
		BackupInterval:   300,
		CommitCommand:    "git commit -m",
		AddAffectedFiles: true,
		CommitDate:       "2006-01-02 15:04:05",
		DebugMode:        false,
		PullOnStart:      true,
	}

	confDir, _ := os.UserConfigDir()

	confFile := path.Join(confDir, "xgs.json")
	confContent, err := os.ReadFile(confFile)
	if err != nil {
		log.Println("[WARNING] xgs config not found: ", err)
		log.Println("using fallback config...")
		return fallbackConf
	}

	resConfig := Config{}

	err = json.Unmarshal(confContent, &resConfig)
	if err != nil {
		log.Println("[WARNING] couldn't parse config", err)
		log.Println("using fallback config...")
		return fallbackConf
	}
	return resConfig
}

func CheckForGit(conf Config) bool {
	DebugLog(conf, "checking for git executable in path...")
	_, err := exec.LookPath("git")
	return err == nil
}

func DebugLog(conf Config, msg string) {
	if conf.DebugMode {
		log.Println("[DEBUG]", msg)
	}
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
