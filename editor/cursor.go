package editor

/* Cursor related functions */

import "strconv"


func (e *Editor) moveCursor(x, y int) {
  e.cursorXPos = x
  e.cursorYPos = y
}


func (e *Editor) moveCursorUp () {
  if e.cursorYPos == 1 {
    if e.rowsOffset > 0 {
      e.rowsOffset--
    }
  } else {
    e.cursorYPos--
  }

  e.boundCoursorRight()
}


func (e *Editor) moveCursorDown () {
  maxCursorYPos := e.screenRows - e.reservedRows

  if e.cursorYPos == maxCursorYPos {
    e.rowsOffset++
    e.cursorYPos = maxCursorYPos
  } else {
    e.cursorYPos++
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
  currRowSize := e.rows[e.GetCurrRowNum()].content.Len() + e.cursorXOffset

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
