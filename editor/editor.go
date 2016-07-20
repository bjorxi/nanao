package editor

import "os"
import "fmt"
import "bufio"
import "bytes"
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

  var rowNum uint32 = 0
  var content *bytes.Buffer

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    rowNum++
    content = bytes.NewBufferString(scanner.Text())
    e.rows = append(e.rows, Row{rowNum, content, content.Len()})
  }

  file.Close()
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
  /* looks too complicated ?*/
  numOfRowsOffset := len(strconv.Itoa(numOfRows)) + 1 /* + 1 for the '|' */
  e.cursorXOffset = numOfRowsOffset + 2
  lineFormat := "%"+ strconv.Itoa(numOfRowsOffset) +"d|%s\x1b[38m\x1b[0K"

  for i := 0; i < numOfRows; i++ {
    row = e.rows[i]
    output += fmt.Sprintf(lineFormat, i+1, row.content.String())
    // output += strconv.Itoa(i+1) + "| " + row.content.String() + "\x1b[39m" + "\x1b[0K"

    if i < numOfRows - 1 {
      output += "\r\n"
    }
  }

  x := strconv.Itoa(int(e.cursorXPos))
  y := strconv.Itoa(int(e.cursorYPos))

  output += "\r\nCursor x: " + x + " y: " +  y + " | "
  output += "lines: " + strconv.Itoa(e.totalRowsNum) + " | "
  output += "cursorXOffset: " + strconv.Itoa(e.cursorXOffset)
  output += "\x1b["+y+";"+x+"f"
  output += "\x1b[?25h" /* Show cursor. */
  fmt.Printf("\x1b[2J")
  fmt.Printf("%s", output)
}


func (e *NanaoEditor) ProcessKeyPress() {
  var keyPress int

  fmt.Scanf("%c", &keyPress)
  fmt.Println("Key pressed", keyPress)
  switch keyPress {
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
    default:
      e.insertChar(keyPress)
  }
}


func (e *NanaoEditor) insertChar (char int) {
  currRow := e.rows[0]
  currRowContent := currRow.content.Bytes()

  newBuffer := bytes.NewBuffer(nil)

  newBuffer.Write(currRowContent[:e.cursorXPos])
  newBuffer.Write([]byte(strconv.Itoa(char)))
  newBuffer.Write(currRowContent[e.cursorXPos:])

  e.rows[e.cursorYPos].content = newBuffer
  e.rows[e.cursorYPos].size = newBuffer.Len()
  e.moveCursor(e.cursorXPos+1, e.cursorYPos)
}


func (e *NanaoEditor) moveCursor(x, y uint32) {
  e.cursorXPos = x
  e.cursorYPos = y
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
