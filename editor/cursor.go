package editor

/* Cursor related functions */

import "strconv"


func (e *Editor) moveCursor(x, y int) {
  e.cursorXPos = x
  e.cursorYPos = y
}


func (e *Editor) moveCursorUp () {
  if e.cursorYPos <= 1 {
    e.cursorYPos = 1
  } else {
    e.cursorYPos--
  }
  /* #TODO replace magic numbers with constan/variable */
  if e.cursorYPos <= 1 && e.rowsOffset > 0 {
    e.rowsOffset--
    e.cursorYPos = 1
  }

  e.boundCoursorRight()
}


func (e *Editor) moveCursorDown () {
  e.cursorYPos++

  if e.cursorYPos >= e.totalRowsNum {
    e.cursorYPos = e.totalRowsNum
  }

  if e.cursorYPos >= e.screenRows - e.reservedRows {
    e.rowsOffset++
    e.cursorYPos = e.screenRows - e.reservedRows
  }

  e.boundCoursorRight()
}


func (e *Editor) moveCursorLeft () {

  if e.cursorXPos <= e.cursorXOffset {
    e.cursorXPos = e.cursorXOffset
  } else {
    e.cursorXPos--
  }
}


func (e *Editor) moveCursorRight () {
  e.cursorXPos++
  e.boundCoursorRight()
}


func (e *Editor) boundCoursorRight () {
  currRowSize := e.rows[e.cursorYPos-1].content.Len() + e.cursorXOffset

  if e.cursorXPos >= currRowSize {
    e.cursorXPos = currRowSize
  }
}


func (e *Editor) setCursorXOffset () {
  /* looks too complicated ?*/
  e.totalRowsNum = len(e.rows)
  numOfRowsOffset := len(strconv.Itoa(e.totalRowsNum)) /* + 1 for the '|' */
  e.cursorXOffset = numOfRowsOffset + 2
}
