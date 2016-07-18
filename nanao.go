package main


import (
  "os"
  "fmt"
  "./editor"
)

const version string = "0.0.0a"


func welcomeMessage () {
  fmt.Println("Nanao editor", version)
  fmt.Println("================")
  fmt.Println()
}


func main () {
  fileName := os.Args[1]

  e := editor.Init()
  e.Open(fileName)
  e.Edit()
}
