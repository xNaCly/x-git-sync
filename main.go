package main

import (
	"log"
	"os"
	"time"
)

var devMode = false

func main() {
	conf := getConfig()

	if !checkForGit() {
		log.Fatalln("[FATAL ERROR] 'git' executable not found, gas requires git to work properly - exiting.")
	}

	if len(os.Args) > 1 && os.Args[1] == "--dev" {
		devMode = true
	}

	if devMode {
		log.Println(generateCommitContent(conf))
	} else {
		for true {
			GitAdd()
			GitCommit(conf)
			GitPush()
			time.Sleep(time.Duration(conf.BackupInterval) * time.Second)
		}
	}
}
