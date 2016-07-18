package editor

import "os"
import "../terminal"

type Editor interface {
  Open(path string)
  GetFilePath() string
  RefreshScreen()
  Edit()

  moveCursorUp()
  moveCursorDown()
  moveCursorLeft()
  moveCursorRight()

  GetNumOfRows()
  ProcessKeyPress()
  getWingowSize()
}


type NanaoEditor struct {
  cursorXPos uint32 /* cursor x position */
  cursorYPos uint32 /* cursor y position */
  screenRows int32 /* Number of rows */
  screenCols int32 /* Number of columns */
  isChanged bool /* Has a file been changed? */
  fileName string
  filePath string
  file *os.File
  rows []Row
  termOldState *terminal.State
}


type Row struct {
  number uint32
  content string
  size int
}


type winsize struct {
  row    uint16
  col    uint16
  xpixel uint16
  ypixel uint16
}
