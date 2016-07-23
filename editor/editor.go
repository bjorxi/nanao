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

  var rowNum int = 0
  var content *bytes.Buffer

  scanner := bufio.NewScanner(file)

  for scanner.Scan() {
    rowNum++
    content = bytes.NewBufferString(scanner.Text())
    e.rows = append(e.rows, Row{rowNum, content, content.Len()})
  }

  e.totalRowsNum = len(e.rows)

  file.Close()
}


func (e *NanaoEditor) Edit() {
  for {
    e.RefreshScreen()
    e.ProcessKeyboardInput()
  }
}


func (e *NanaoEditor) RefreshScreen() {
  var row Row

  output := "" /* #TODO replace string with bytes.Buffer */
  output += "\x1b[?25l" /* Hide cursor. */
  output += "\x1b[H" /* Go home. */
  numOfRows := len(e.rows)

  /* looks too complicated ?*/
  numOfRowsOffset := len(strconv.Itoa(numOfRows)) /* + 1 for the '|' */
  e.cursorXOffset = numOfRowsOffset + 2
  lineFormat := "%"+ strconv.Itoa(numOfRowsOffset) +"d|%s\x1b[38m\x1b[0K"

  for i := 0; i < numOfRows; i++ {
    row = e.rows[i]
    output += fmt.Sprintf(lineFormat, i+1, row.content.String())

    if i < numOfRows - 1 {
      output += "\r\n"
    }
  }

  x := strconv.Itoa(int(e.cursorXPos))
  y := strconv.Itoa(int(e.cursorYPos))

  output += "\r\nCursor x: " + x + " y: " +  y + " | "
  output += "lines: " + strconv.Itoa(e.totalRowsNum) + " | "
  output += "cursorXOffset: " + strconv.Itoa(e.cursorXOffset)
  output += "\r\nLine size " + strconv.Itoa(e.rows[e.cursorYPos-1].content.Len()) + "(" +
            strconv.Itoa(e.rows[e.cursorYPos-1].size) + ")"
  output += "\x1b["+y+";"+x+"f"

  output += "\x1b[?25h" /* Show cursor. */
  fmt.Printf("\x1b[2J")
  fmt.Printf("%s", output)
}


/**
 * 27 91 51 126 - backspace
 */
func (e *NanaoEditor) ProcessKeyboardInput() {
  var input []byte = make([]byte, 4)

  os.Stdin.Read(input)

  fmt.Fprintf(os.Stderr, "%d %d %d %d\n", input[0],input[1],input[2],input[3])


  if input[1] == 0 && input[2] == 0 {
    key := input[0]
    if key == 3 {  /* ctrl-c */
      fmt.Println("\x1b[2J")
      terminal.Restore(0, e.termOldState)
      os.Exit(0)
    } else if key == 13 { /* enter */
      e.insertEmptyRow()
    } else if key == 27 { /* ESC */
      return
    } else if key == 32 {
      e.insertChar(string(" "))
    } else if key >= 33 && key <= 126 {
      e.insertChar(string(input[0]))
    } else if key == 127 { /* DELETE */
      e.deleteChar()
    } else {
      return
    }
  }

  if input[0] == 27 {
    if input[1] == 91 {
      if input[2] == 68 {
        e.moveCursorLeft()
      } else if input[2] == 67 {
        e.moveCursorRight()
      } else if input[2] == 66 {
        e.moveCursorDown()
      } else if input[2] == 65 {
        e.moveCursorUp()
      } else {
        return
      }
    }
  }

  return
}

func (e *NanaoEditor) insertEmptyRow() {
  var rows []Row

  currRow := e.rows[e.cursorYPos-1]
  currRowContent := currRow.content.Bytes()
  sliceAt := e.cursorXPos - e.cursorXOffset

  newBuffer := bytes.NewBuffer(currRowContent[sliceAt:])
  newRow := Row{e.cursorYPos, newBuffer, newBuffer.Len()}

  e.rows[e.cursorYPos-1] = Row{e.cursorYPos-1, bytes.NewBuffer(currRowContent[:sliceAt]),
                               bytes.NewBuffer(currRowContent[:sliceAt]).Len()}

  rows = append(rows, e.rows[:e.cursorYPos]...)
  rows = append(rows, newRow)
  rows = append(rows, e.rows[e.cursorYPos:]...)

  /* looks too complicated ?*/
  numOfRows := len(rows)
  numOfRowsOffset := len(strconv.Itoa(numOfRows)) /* + 1 for the '|' */
  e.cursorXOffset = numOfRowsOffset + 2

  e.rows = rows
  e.totalRowsNum++
  e.moveCursor(e.cursorXOffset, e.cursorYPos+1)
}


func (e *NanaoEditor) insertChar (char string) {
  currRow := e.rows[e.cursorYPos-1]

  currRowContent := currRow.content.Bytes()
  newBuffer := bytes.NewBuffer(nil)

  newBuffer.Write(currRowContent[:e.cursorXPos-e.cursorXOffset])
  newBuffer.Write([]byte(char))
  newBuffer.Write(currRowContent[e.cursorXPos-e.cursorXOffset:])

  e.rows[e.cursorYPos-1].content = newBuffer
  e.rows[e.cursorYPos-1].size = newBuffer.Len()
  e.moveCursor(e.cursorXPos+1, e.cursorYPos)
}


func (e *NanaoEditor) deleteChar() {
  currRow := e.rows[e.cursorYPos-1]

  currRowContent := currRow.content.Bytes()
  newBuffer := bytes.NewBuffer(nil)

  if e.cursorXPos <= e.cursorXOffset {
    return
  }

  newBuffer.Write(currRowContent[:e.cursorXPos-e.cursorXOffset-1])
  newBuffer.Write(currRowContent[e.cursorXPos-e.cursorXOffset:])

  e.rows[e.cursorYPos-1].content = newBuffer
  e.rows[e.cursorYPos-1].size = newBuffer.Len()
  e.moveCursor(e.cursorXPos-1, e.cursorYPos)
}


func (e *NanaoEditor) moveCursor(x, y int) {
  e.cursorXPos = x
  e.cursorYPos = y
}


func (e *NanaoEditor) moveCursorUp () {
  if e.cursorYPos <= 1 {
    e.cursorYPos = 1
  } else {
    e.cursorYPos--
  }

  e.boundCoursorRight()
}


func (e *NanaoEditor) moveCursorDown () {
  e.cursorYPos++

  if e.cursorYPos >= e.totalRowsNum {
    e.cursorYPos = e.totalRowsNum
  }

  e.boundCoursorRight()
}


func (e *NanaoEditor) moveCursorLeft () {

  if e.cursorXPos <= e.cursorXOffset {
    e.cursorXPos = e.cursorXOffset
  } else {
    e.cursorXPos--
  }
}


func (e *NanaoEditor) moveCursorRight () {
  e.cursorXPos++
  e.boundCoursorRight()
}


func (e *NanaoEditor) boundCoursorRight () {
  currRowSize := e.rows[e.cursorYPos-1].size + e.cursorXOffset

  if e.cursorXPos >= currRowSize {
    e.cursorXPos = currRowSize
  }
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

  e.screenCols = int(ws.col)
  e.screenRows = int(ws.row)
}


func Init() Editor {
  e := &NanaoEditor{}
  e.cursorXOffset = 3
  e.cursorXPos = e.cursorXOffset
  e.cursorYPos = 1
  e.getWingowSize()
  e.isChanged = false
  e.termOldState, _ = terminal.MakeRaw(0)
  return e
}
