package editor

import "bytes"
import "../terminal"


type Editor struct {
  cursorXPos int /* cursor x position */
  cursorYPos int /* cursor y position */
  cursorMaxYPos int
  cursorMinYPost int
  cursorXOffset int
  reservedRows int
  statusLineRows int
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
  conf *Config
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


type Config struct {
  tabSize int `json:"tab_size"`
  translateTabsToSpaces bool `json:"translate_tabs_to_spaces"`
}
