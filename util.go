package main

import (
	"encoding/json"
	"log"
	"os"
	"path"
)

type Config struct {
	// will be inserted before the local datestring in the commit title
	AutoCommitPrefix string `json:"auto_commit_prefix"`
	// replaces the default commit title
	CustomCommitTitle string `json:"custom_commit_title"`
	// time interval between backups (in s)
	BackupInterval int `json:"backup_interval"`
}

// Generates a new commit message based on the users configuration:
//
// 1. By default the commit title will be formated like so: "[Config.AutoCommitPrefix] yyyy-mm-dd HH:MM:SS"
//
// 2. if [Config.CustomCommitMsg] is set the commit msg will be exactly the content specified in it: [Config.CustomCommitMsg]
func getCommitTitle() string {
	return "here should be the time"
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
	    AutoCommitPrefix: "backup:",
	    BackupInterval: 300,
	}

    confDir, err := os.UserConfigDir()
    if err != nil {
        log.Println("[ERR]", err)
    }

    confFile := path.Join(confDir, ".git_auto_sync.json")
    confContent, err := os.ReadFile(confFile)
    if err != nil {
        log.Println("[ERR]", err)
    }

    resConfig := Config{}

    err = json.Unmarshal(confContent, &resConfig)
    if err != nil {
        log.Println("[ERR]", err)
    } else {
        return resConfig
    }

    return fallbackConf
}
