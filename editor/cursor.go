package editor

/* Cursor related functions */

import "strconv"


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
  currRowSize := e.rows[e.cursorYPos-1].content.Len() + e.cursorXOffset

  if e.cursorXPos >= currRowSize {
    e.cursorXPos = currRowSize
  }
}


func (e *NanaoEditor) setCursorXOffset () {
  /* looks too complicated ?*/
  e.totalRowsNum = len(e.rows)
  numOfRowsOffset := len(strconv.Itoa(e.totalRowsNum)) /* + 1 for the '|' */
  e.cursorXOffset = numOfRowsOffset + 2
}
