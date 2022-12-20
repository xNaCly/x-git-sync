package main

import (
	"encoding/json"
	"log"
	"os"
	"path"
)

type Config struct {
    // will be inserted before the local date string in the commit title, default: "backup: "
	AutoCommitPrefix string `json:"auto_commit_prefix"`

    // TODO: implement
    //
	// CommitTitle string `json:"custom_commit_title"`

    // specifies the date format which the date will be formated as, default: "01-02-2006 15:04:05"
    CommitTitleDateFormat string `json:"commit_title_date_format"`

    // List filenames affected by the commit in the commit body, default: true
    AddAffectedFiles bool `json:"add_affected_files"`

    // time interval between backups (in s), default: 300
	BackupInterval int `json:"backup_interval"`

    // commit command, default: "git commit -m"
    CommitCommand string `json:"commit_cmd"`
}

// Loads and parses config from $HOME/.git_auto_sync.json
//
//
// config file location depends on os.UserConfigDir()
//
// if config is not found the fallback config is:
//
//	Config{
//	    AutoCommitPrefix: "backup:"
//	    BackupInterval: 300
//	}
func getConfig() Config {
    // all occuring errors are logged, but not treated like panics, due to the fact that a fallback config is provided
    fallbackConf := Config{
	    AutoCommitPrefix: "backup: ",
	    BackupInterval: 300,
        CommitCommand: "git commit -m",
        AddAffectedFiles: true,
        CommitTitleDateFormat: "2006-01-02 15:04:05",
	}

    confDir, _ := os.UserConfigDir()

    confFile := path.Join(confDir, ".git_auto_sync.json")
    confContent, err := os.ReadFile(confFile)
    if err != nil {
        log.Println("[ERR]", err)
        log.Println("[INF] using fallback config")
        return fallbackConf
    }

    resConfig := Config{}

    err = json.Unmarshal(confContent, &resConfig)
    if err != nil {
        log.Println("[ERR]", err)
        log.Println("[INF] using fallback config")
        return fallbackConf
    } 
    return resConfig
}

