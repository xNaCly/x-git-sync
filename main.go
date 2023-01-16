package main

import (
	"log"
	"os"
	"time"
)

func main() {
	// disables adding, commitin and pushing, only logs the generated commit message
	devMode := false;
	conf := getConfig()

	if !checkForGit(conf) {
		log.Fatalln("[FATAL ERROR] 'git' executable not found, gas requires git to work properly - exiting.")
	}

	if len(os.Args) > 1 && os.Args[1] == "--debug"{
		conf.DebugMode = true
	} else if len(os.Args) > 1 && os.Args[1] == "--dev" {
		devMode = true
	}

	if conf.DebugMode {
		DebugLog(conf, "Debug mode enabled");
	}

	if devMode {
		conf.DebugMode = true;
		DebugLog(conf, "Dev mode enabled, automatically enabled debug mode, adding, committing and pushing will be disabled.");
		generateCommitContent(conf)
		os.Exit(0)
	}
	
	if conf.PullOnStart {
		GitPull(conf)
	}

	log.Println("[INFO] Watching for changes...")

	for true {
		if GitRepoHasChanges(conf) {
			GitAdd(conf)
			GitCommit(conf)
			GitPush(conf)
		} else {
			log.Println("[INFO] No changes to commit, waiting for next iteration...")
		}
		time.Sleep(time.Duration(conf.BackupInterval) * time.Second)
	}
}
