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

  var file *os.File
  var err error
  var rowNum int = 0
  var content *bytes.Buffer

  if fileExists(path) {
    file, err = os.Open(path)

    if err != nil {
      fmt.Println("Error opening file", err)
      os.Exit(1)
    }

    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
      rowNum++
      /* #TODO use NewBuffer(scanner.Bytes())*/
      content = bytes.NewBufferString(scanner.Text())
      e.rows = append(e.rows, Row{rowNum, content, content.Len()})
    }

    if rowNum == 0 {
      content = bytes.NewBuffer(nil)
      e.rows = append(e.rows, Row{1, content, content.Len()})
      rowNum = 1
    }

    file.Close()
  } else {
    content = bytes.NewBuffer(nil)
    e.rows = append(e.rows, Row{1, content, content.Len()})
    rowNum = 1
  }

  e.totalRowsNum = rowNum
  e.setCursorXOffset()
  e.cursorXPos += e.cursorXOffset
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

  lineFormat := "%"+ strconv.Itoa(e.cursorXOffset-2) +"d|%s\x1b[38m\x1b[0K"

  for i := 0; i < e.totalRowsNum; i++ {
    row = e.rows[i]
    output += fmt.Sprintf(lineFormat, i+1, row.content.String())

    if i < e.totalRowsNum - 1 {
      output += "\r\n"
    }
  }

  x := strconv.Itoa(e.cursorXPos)
  y := strconv.Itoa(e.cursorYPos)

  output += "\r\n\r\nCursor x: " + x + " y: " +  y + " | "
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

  //fmt.Fprintf(os.Stderr, "%d %d %d %d\n", input[0],input[1],input[2],input[3])


  if input[1] == 0 && input[2] == 0 {
    key := input[0]
    if key == 3 {  /* ctrl-c */
      fmt.Println("\x1b[2J")
      terminal.Restore(0, e.termOldState)
      os.Exit(0)
    } else if key == 13 { /* enter */
      e.insertEmptyRow()
    } else if key == 17 { /* ctrl-q */
      fmt.Println("\x1b[2J")
      terminal.Restore(0, e.termOldState)
      os.Exit(0)
    } else if key == 19 { /* ctrl-s*/
      e.SaveChanges()
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


func (e *NanaoEditor) SaveChanges () {
  var outputLine string
  var file *os.File
  var err error
  var fileInfo os.FileInfo

  if !fileExists(e.filePath) {
    file, err = os.OpenFile(e.filePath, os.O_WRONLY|os.O_CREATE, 0644)
  } else {
    fileInfo, err = os.Stat(e.filePath)
    filePerms := os.FileMode(fileInfo.Mode())
    file, err = os.OpenFile(e.filePath, os.O_WRONLY | os.O_TRUNC, filePerms)

    if err != nil {
      fmt.Println("Error saving file")
      return
    }
  }

  for i := 0; i < e.totalRowsNum; i++ {
    outputLine = e.rows[i].content.String() + "\n"
    file.WriteString(outputLine)
  }

  file.Close()
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

  e.rows = rows
  e.totalRowsNum++
  e.setCursorXOffset()
  e.moveCursor(e.cursorXOffset, e.cursorYPos+1)
}


func (e *NanaoEditor) deleteRow () {
  /* #TODO Replace magic number with constant/variable */
  if e.cursorYPos == 1 {
    return
  }

  var rows []Row

  currRow := e.rows[e.cursorYPos-1]
  prevRow := e.rows[e.cursorYPos-2]
  currRowContent := currRow.content.Bytes()

  e.moveCursor(prevRow.content.Len()+e.cursorXOffset, e.cursorYPos-1)

  prevRow.content.Write(currRowContent)
  prevRow.size = prevRow.content.Len()
  rows = append(rows, e.rows[:e.cursorYPos]...)
  rows = append(rows, e.rows[e.cursorYPos+1:]...)

  e.rows = rows
  e.totalRowsNum--
  e.setCursorXOffset()
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
  if e.cursorXPos == e.cursorXOffset {
    e.deleteRow()
    return
  }

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
  e.cursorYPos = 1
  e.getWingowSize()
  e.isChanged = false
  e.termOldState, _ = terminal.MakeRaw(0)
  return e
}
