package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path"
)

type Config struct {
	// will be inserted before the local date string in the commit title, default: "backup: "
	AutoCommitPrefix string `json:"auto_commit_prefix"`

	// TODO: implement
	//
	// CommitTitle string `json:"custom_commit_title"`

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
	CommitTitleDateFormat string `json:"commit_title_date_format"`

	// List filenames affected by the commit in the commit body, default: true
	AddAffectedFiles bool `json:"add_affected_files"`

	// time interval between backups (in s), default: 300
	BackupInterval int `json:"backup_interval"`

	// commit command, default: "git commit -m"
	CommitCommand string `json:"commit_cmd"`
}

// Loads and parses config from:
// - On Unix systems, $XDG_CONFIG_HOME or $HOME/.config
// - On Darwin, it returns $HOME/Library/Application Support
// - On Windows, it returns %AppData%
// - On Plan 9, it returns $home/lib
//
// config file location depends on os.UserConfigDir()
//
// if config is not found the fallback config is:
//
//	Config{
//      AutoCommitPrefix:      "backup: ",
//      BackupInterval:        300,
//      CommitCommand:         "git commit -m",
//      AddAffectedFiles:      true,
//      CommitTitleDateFormat: "2006-01-02 15:04:05",
//	}
func getConfig() Config {
	// all occuring errors are logged, but not treated like panics, due to the fact that a fallback config is provided
	fallbackConf := Config{
		AutoCommitPrefix:      "backup: ",
		BackupInterval:        300,
		CommitCommand:         "git commit -m",
		AddAffectedFiles:      true,
		CommitTitleDateFormat: "2006-01-02 15:04:05",
	}

	confDir, _ := os.UserConfigDir()

	confFile := path.Join(confDir, "gas.json")
	confContent, err := os.ReadFile(confFile)
	if err != nil {
		log.Println("[WARNING] gas config not found: ", err)
		log.Println("[INFO] using fallback config...")
		return fallbackConf
	}

	resConfig := Config{}

	err = json.Unmarshal(confContent, &resConfig)
	if err != nil {
		log.Println("[WARNING] couldn't parse config", err)
		log.Println("[INF] using fallback config...")
		return fallbackConf
	}
	return resConfig
}

func checkForGit() bool {
	_, err := exec.LookPath("git")
	return err == nil
}
