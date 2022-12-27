package main

import (
	"log"
	"time"
)

func main() {
	conf := getConfig()

	if !checkForGit() {
		log.Fatalln("[FATAL ERROR] 'git' executable not found, gas requires git to work properly - exiting.")
	}

	for true {
		GitAdd()
		GitCommit(conf)
		GitPush()
		time.Sleep(time.Duration(conf.BackupInterval) * time.Second)
	}
}
