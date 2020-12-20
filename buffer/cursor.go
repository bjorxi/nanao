package buffer

/* Cursor related functions */

import (
  // "strconv"
  "fmt"
  "os"
)


func (b *Buffer) MoveCursor(x, y int) {
  b.cursorXPos = x
  b.cursorYPos = y
  fmt.Fprintf(os.Stderr, "cursorYPos %d\n", b.cursorYPos)

  maxCursorYPos := b.maxVisibleRows - b.reservedRows
  fmt.Fprintf(os.Stderr, "%d := %d - %d\n", maxCursorYPos, b.maxVisibleRows, b.reservedRows)

  if b.cursorYPos > maxCursorYPos {
    b.rowsOffset++
    b.cursorYPos = maxCursorYPos
  }
}


func (e *Buffer) MoveCursorUp () {
  if e.cursorYPos == 1 {
    if e.rowsOffset > 0 {
      e.rowsOffset--
    }
  } else {
    e.cursorYPos--
  }

  e.boundCoursorRight()
}


func (b *Buffer) MoveCursorDown () {
  if b.cursorYPos >= len(b.rows) {
    return
  }

  if b.cursorYPos >= b.maxVisibleRows {
    b.rowsOffset++
    b.cursorYPos = b.maxVisibleRows
  } else {
    b.cursorYPos++
  }

  b.boundCoursorRight()
}


func (b *Buffer) MoveCursorLeft () {
  b.cursorXPos--

  if b.cursorXPos <= b.getLineMetaChars() {
    b.cursorXPos = b.getLineMetaChars() + 1
  }
}


func (e *Buffer) MoveCursorRight () {
  e.cursorXPos++
  e.boundCoursorRight()
}


func (e *Buffer) boundCoursorRight () {
  currRowSize := e.rows[e.GetCurrRowNum()].content.Len() + e.getLineMetaChars() + 1

  if e.cursorXPos >= currRowSize {
    e.cursorXPos = currRowSize
  }
}


func (b *Buffer) setCursorXOffset () {
  // TODO replace this with a config option
  // minCursorXOffset := 4
  // /* looks too complicated ?*/
  // numOfRowsOffset := len(strconv.Itoa(len(b.rows))) /* + 1 for the '|' */
  //
  // if numOfRowsOffset < minCursorXOffset {
  //   b.cursorXOffset = minCursorXOffset
  // } else {
  //   b.cursorXOffset = numOfRowsOffset + 2
  // }

  // return b.getLineMetaChars()
}
