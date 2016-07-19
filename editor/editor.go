package editor

import "os"
import "fmt"
import "bufio"
import "unsafe"
import "syscall"
import "../terminal"


// #TODO create a file if it doesn't exist
func (e *NanaoEditor) Open(path string) {
  e.filePath = path
  file, err := os.Open(path)

  if err != nil {
    fmt.Println("Error opening file")
    os.Exit(1)
  }

  e.file = file

  var rowNum uint32 = 0
  var content string

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    rowNum++
    content = scanner.Text()
    e.rows = append(e.rows, Row{rowNum, content, len(content)})
  }

  e.file.Close()
}


func (e *NanaoEditor) Edit() {
  for {
    e.RefreshScreen()
    e.ProcessKeyPress()
  }
}


func (e *NanaoEditor) RefreshScreen() {
  fmt.Println()
  var row Row
  for i := 0; i < len(e.rows); i++ {
    row = e.rows[i]
    fmt.Printf("%d %s\n", row.number, row.content)
  }
}


func (e *NanaoEditor) ProcessKeyPress() {
  var keyPress int

  fmt.Scanf("%c", &keyPress)
  fmt.Println("Key pressed", keyPress)
  switch keyPress {
  default:
    fmt.Printf("%c", keyPress)
  case 3:
    fmt.Println("^C")
    terminal.Restore(0, e.termOldState)
    os.Exit(0)
  case 10: /* enter */
    fmt.Println()
  case 27, 91:
    return
  case 68: /* left arrow */
    fmt.Println("Left")
  case 67: /* right arrow */
    fmt.Println("Right")
  case 65: /* up arrow */
    fmt.Println("Up")
  case 66: /* down arrow */
    fmt.Println("Down")
  }
}


func (e *NanaoEditor) moveCursorUp () {
  e.cursorYPos++
}

func (e *NanaoEditor) moveCursorDown () {
  e.cursorYPos--
}

func (e *NanaoEditor) moveCursorLeft () {
  e.cursorXPos--
}

func (e *NanaoEditor) moveCursorRight () {
  e.cursorXPos++
}

func (e *NanaoEditor) GetNumOfRows() {

}


func (e *NanaoEditor) GetFilePath() string {
  return e.filePath
}


func (e NanaoEditor) getWingowSize() {
  ws := &winsize{}
  retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
      uintptr(syscall.Stdin),
      uintptr(syscall.TIOCGWINSZ),
      uintptr(unsafe.Pointer(ws)))

  if int(retCode) == -1 {
      panic(errno)
  }

  e.screenCols = int32(ws.col)
  e.screenRows = int32(ws.row)
}


func Init() Editor {
  e := &NanaoEditor{}
  e.getWingowSize()
  e.isChanged = false
  e.termOldState, _ = terminal.MakeRaw(0)
  return e
}
