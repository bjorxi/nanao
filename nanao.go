package main


import (
  "os"
  "fmt"
  "./editor"
)

const version string = "v0.0.1b"


func welcomeMessage () {
  fmt.Println("Nanao editor", version)
  fmt.Println("================")
  fmt.Println()
}


func main () {
  if len(os.Args) < 2 {
    fmt.Println("Usage: nanao </path/to/file>")
    return
  }

  if os.Args[1] == "-v" || os.Args[1] == "--version" {
    fmt.Println(version)
    return
  }

  if os.Args[1] == "-h" || os.Args[1] == "--help" {
    fmt.Println("Nanao editor", version, "\n")
    fmt.Println("Usage: nanao </path/to/file>\n")
    fmt.Println("Arguments:")
    fmt.Println("\t-v(--version) show the version")
    fmt.Println("\t-h(--help) show help\n")
    fmt.Println("Author: Nodari Lipartiya <nodari.lipartiya(at)gmail.com>")
    fmt.Println("Project url: https://github.com/leemalmac/nanao\n")
    return
  }

  fileName := os.Args[1]

  e := editor.Init()
  e.Open(fileName)
  e.Edit()
}
