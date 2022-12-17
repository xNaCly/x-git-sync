package main

import "fmt"


func main() {
    conf := getConfig()
    fmt.Println(conf.AutoCommitPrefix)
}
