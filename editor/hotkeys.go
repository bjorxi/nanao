package editor


// Moves cursor to the start of the line
// Keys:
//    Mac: fn + Left
func (e *Editor) moveToLineStart () {
  e.cursorXPos = e.cursorXOffset
}


// Moves cursor to the end of the line
// Keys:
//    Mac: fn + Right
func (e *Editor) moveToLineEnd () {
  currRow := e.rows[e.GetCurrRowNum()]
  currRowLen := currRow.content.Len()
  e.cursorXPos = e.cursorXOffset + currRowLen
}
