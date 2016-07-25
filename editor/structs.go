package editor

import "bytes"
import "../terminal"


type Editor interface {
  Open(path string)
  GetFilePath() string
  RefreshScreen()
  ProcessKeyboardInput()
  Edit()
  GetNumOfRows()
  SaveChanges()

  moveCursorUp()
  moveCursorDown()
  moveCursorLeft()
  moveCursorRight()
  boundCoursorRight()
  getWindowSize()
}


type NanaoEditor struct {
  cursorXPos int /* cursor x position */
  cursorYPos int /* cursor y position */
  cursorXOffset int
  screenRows int /* Number of rows */
  screenCols int /* Number of columns */
  rowsOffset int
  colsOffset int
  isChanged bool /* Has a file been changed? */
  fileName string
  filePath string
  rows []Row /* File content */
  totalRowsNum int
  termOldState *terminal.State /* TODO: move it somewhere */
}


type Row struct {
  number int
  content *bytes.Buffer
  size int
}


type winsize struct {
  row    uint16
  col    uint16
  xpixel uint16
  ypixel uint16
}
