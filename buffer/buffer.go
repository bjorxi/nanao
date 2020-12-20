package buffer

import (
  "bufio"
  "bytes"
  "os"
  "fmt"
  "strconv"

  "github.com/bjorxi/nanao/util"
)


type Buffer struct {
	cursorXPos int /* cursor x position */
	cursorYPos int /* cursor y position */
  lineNumberDigits int // shows number of digits used to identify line numbers
  cursorXInitPos int
  maxVisibleRows int
  rowsOffset int
  reservedRows int
	isChanged bool /* Has a file been changed? */
	fileName string
	filePath string
	rows []Row /* File content */
}

func New() *Buffer {
  b := &Buffer{
    isChanged: false,
    cursorYPos: 0,
    cursorXInitPos: 6,
    lineNumberDigits: 3,
    reservedRows: 3,
  }

  b.cursorXPos = b.getLineMetaChars() + 1
  b.rows = append(b.rows, Row{1, bytes.NewBuffer(nil), 0})

  return b
}

// Returns a number of chars that used to display a line number + a separator
// Check the `GetVisibleContent` func to see the lineFormat
func (b *Buffer) getLineMetaChars() int {
  // we add 1 cause each line has a separator to separate line number from the actual content
  return b.lineNumberDigits + 1
}

func (b *Buffer) getRowEditIndex() int {
  return b.cursorXPos - b.getLineMetaChars() - 1
}


func NewFromFile (path string) *Buffer {
  b := &Buffer{
    isChanged: false,
    cursorXPos: 6,
    lineNumberDigits: 3,
    reservedRows: 3,
  }

  b.filePath = path

  var file *os.File
  var err error
  var content *bytes.Buffer
  rowNum := 0

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
    b.rows = append(b.rows, Row{rowNum, content, content.Len()})
  }

  if len(b.rows) == 0 {
    content = bytes.NewBuffer(nil)
    b.rows = append(b.rows, Row{1, content, content.Len()})
  }

  file.Close()

  return b
}


// #TODO create a file if it doesn't exist
func (b *Buffer) Open(path string) {
  if util.FileExists(path) {

  } else {

  }

  b.setCursorXOffset()
  b.cursorXPos += b.getLineMetaChars()
}


func (e *Buffer) SaveChanges () {
  // var outputLine string
  // var file *os.File
  // var err error
  // var fileInfo os.FileInfo
  //
  // if !util.FileExists(e.filePath) {
  //   file, err = os.OpenFile(e.filePath, os.O_WRONLY|os.O_CREATE, 0644)
  // } else {
  //   fileInfo, err = os.Stat(e.filePath)
  //   filePerms := os.FileMode(fileInfo.Mode())
  //   file, err = os.OpenFile(e.filePath, os.O_WRONLY | os.O_TRUNC, filePerms)
  //
  //   if err != nil {
  //     fmt.Println("Error saving file")
  //     return
  //   }
  // }
  //
  // for i := 0; i < e.totalRowsNum; i++ {
  //   outputLine = e.rows[i].content.String() + "\n"
  //   file.WriteString(outputLine)
  // }
  //
  // file.Close()
}


func (e *Buffer) InsertIndent () {
  // for i := 0; i < e.conf.TabSize; i++ {
  //   e.insertChar(" ")
  // }
}


func (b *Buffer) InsertEmptyRow() {
  var rows []Row

  currRow := b.rows[b.GetCurrRowNum()]
  currRowContent := currRow.content.Bytes()
  sliceAt := b.getRowEditIndex()

  newBuffer := bytes.NewBuffer(currRowContent[sliceAt:])
  newRow := Row{b.GetCurrRowNum(), newBuffer, newBuffer.Len()}

  b.rows[b.GetCurrRowNum()] = Row{b.cursorYPos-1, bytes.NewBuffer(currRowContent[:sliceAt]),
                               bytes.NewBuffer(currRowContent[:sliceAt]).Len()}

  rows = append(rows, b.rows[:b.GetCurrRowNum()+1]...)
  rows = append(rows, newRow)
  rows = append(rows, b.rows[b.GetCurrRowNum()+1:]...)

  b.rows = rows
  b.setCursorXOffset()
  b.MoveCursor(b.getLineMetaChars()+1, b.cursorYPos+1)
}


func (e *Buffer) GetCurrRowNum () int {
  return e.cursorYPos
}


func (e *Buffer) DeleteRow () {
  /* #TODO Replace magic number with constant/variable */
  // if e.cursorYPos == 1 && e.rowsOffset == 0 {
  //   return
  // }
  //
  // var rows []Row
  // currRowNum := e.GetCurrRowNum()
  // currRow := e.rows[currRowNum]
  // prevRowNum := e.GetCurrRowNum()-1
  // prevRow := e.rows[prevRowNum]
  // prevRowLen := prevRow.content.Len()
  // currRowContent := currRow.content.Bytes()
  //
  // e.MoveCursorUp()
  //
  // prevRow.content.Write(currRowContent)
  // prevRow.size = prevRow.content.Len()
  // rows = append(rows, e.rows[:currRowNum]...)
  // rows = append(rows, e.rows[currRowNum+1:]...)
  //
  // e.rows = rows
  // e.cursorXPos = prevRowLen + e.cursorXOffset
  // e.totalRowsNum--
  // e.setCursorXOffset()
}


func (b *Buffer) InsertChar (char string) {
  fmt.Fprintf(os.Stderr, "Buffer::InsertChar: %d, %d\n", b.cursorXPos, b.getLineMetaChars())

  b.rows[b.cursorYPos].insertChar(char, b.getRowEditIndex())
  b.MoveCursor(b.cursorXPos+1, b.cursorYPos)
}


func (e *Buffer) DeleteChar() {
  // if e.cursorXPos == e.cursorXOffset {
  //   e.DeleteRow()
  //   return
  // }
  //
  // currRow := e.rows[e.GetCurrRowNum()]
  //
  // currRowContent := currRow.content.Bytes()
  // newBuffer := bytes.NewBuffer(nil)
  //
  // if e.cursorXPos <= e.cursorXOffset {
  //   return
  // }
  //
  // newBuffer.Write(currRowContent[:e.cursorXPos-e.cursorXOffset-1])
  // newBuffer.Write(currRowContent[e.cursorXPos-e.cursorXOffset:])
  //
  // e.rows[e.GetCurrRowNum()].content = newBuffer
  // e.rows[e.GetCurrRowNum()].size = newBuffer.Len()
  // e.MoveCursor(e.cursorXPos-1, e.cursorYPos)
}

// Calculates max visible row Num
func (b *Buffer) getNumOfVisibleRows () int {
  visibleRows := b.rowsOffset + (b.maxVisibleRows - b.reservedRows)

  if visibleRows > len(b.rows) {
    visibleRows = len(b.rows)
  }

  return visibleRows
}


func (e *Buffer) GetVisibleContent () string {
  var row Row

  output := "" /* #TODO replace string with bytes.Buffer */
  // Move line format in the Buffer struct, or gloabl package space
  lineFormat := "%0"+ strconv.Itoa(e.lineNumberDigits) +"d|%s\x1b[38m\x1b[0K"
  maxScreenRows := e.getNumOfVisibleRows()

  for i := e.rowsOffset; i < maxScreenRows; i++ {
    row = e.rows[i]
    output += fmt.Sprintf(lineFormat, i+1, row.content.String())

    if i < maxScreenRows - 1 {
      output += "\r\n"
    }
  }

  return output
}

func (b *Buffer) GetCursorXPos() int {
  return b.cursorXPos
}

func (b *Buffer) GetCursorYPos() int {
  return b.cursorYPos
}

func (b *Buffer) SetMaxVisibleRows(num int) {
  b.maxVisibleRows = num
}
