package main


func main() {
    conf := getConfig()

    // for true {
        GitAdd()
        GitCommit(conf)
        GitPush()
        // time.Sleep(time.Duration(conf.BackupInterval) * time.Second)
    // }
}
