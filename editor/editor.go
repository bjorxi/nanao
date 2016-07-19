package editor

import "os"
import "fmt"
import "bufio"
import "strconv"
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
  var row Row

  output := "" /* #TODO replace string with bytes.Buffer */
  output += "\x1b[?25l" /* Hide cursor. */
  output += "\x1b[H" /* Go home. */
  numOfRows := len(e.rows)

  for i := 0; i < numOfRows; i++ {
    row = e.rows[i]
    output += row.content + "\x1b[39m" + "\x1b[0K"

    if i < numOfRows - 1 {
      output += "\r\n"
    }
  }

  x := strconv.Itoa(int(e.cursorXPos))
  y := strconv.Itoa(int(e.cursorYPos))

  output +=  "\r\nmoving cursor x: " + x + " y: " +  y + "|| " + "x1b["+y+";"+x+"f"
  output += "\x1b["+y+";"+x+"f"
  output += "\x1b[?25h" /* Show cursor. */
  fmt.Printf("%s", output)
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
    return /* #TODO handle this cases more efficient, now it forces screenRefresh */
  case 68: /* left arrow */
    e.moveCursorLeft()
  case 67: /* right arrow */
    e.moveCursorRight()
  case 65: /* up arrow */
    e.moveCursorUp()
  case 66: /* down arrow */
    e.moveCursorDown()
  }
}


func (e *NanaoEditor) moveCursorUp () {
  if e.cursorYPos <= 0 {
    e.cursorYPos = 0
  } else {
    e.cursorYPos--
    // fmt.Println("\x1b[1A")
  }
}

func (e *NanaoEditor) moveCursorDown () {
  e.cursorYPos++
  // fmt.Println("\x1b[1B")
}

func (e *NanaoEditor) moveCursorLeft () {
  if e.cursorXPos <= 0 {
    e.cursorXPos = 0
  } else {
    e.cursorXPos--
  }
  // fmt.Println("\x1b[1D")
}

func (e *NanaoEditor) moveCursorRight () {
  e.cursorXPos++
  // fmt.Println("\x1b[1C")
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
  e.cursorXPos = 0
  e.cursorYPos = 0
  e.getWingowSize()
  e.isChanged = false
  e.termOldState, _ = terminal.MakeRaw(0)
  return e
}
