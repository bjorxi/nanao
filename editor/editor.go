package editor

import (
  "os"
  "fmt"
  "strconv"

  "github.com/bjorxi/nanao/buffer"
  "github.com/bjorxi/nanao/terminal"
  "github.com/bjorxi/nanao/util"
)


func (e *Editor) Edit() {
  for {
    e.Draw()
    e.ProcessKeyboardInput()
  }
}

func (e *Editor) Draw() {
  x := strconv.Itoa(e.getCurrentBuffer().GetCursorXPos())
  y := strconv.Itoa(e.getCurrentBuffer().GetCursorYPos()+1)
  fmt.Fprintf(os.Stderr, "Editor::Draw | cursorYPos %s; CursorXPos %s\n", y, x)

  output := "" /* #TODO replace string with bytes.Buffer */

  output += "\x1b[?25l" /* Hide cursor. */
  output += "\x1b[H" /* Go home. */
  output += e.getCurrentBuffer().GetVisibleContent()
  output += "\x1b["+y+";"+x+"f" /* Set cursor position */
  output += "\x1b[?25h" /* Show cursor. */

  fmt.Printf("\x1b[2J")
  fmt.Printf("%s", output)
}


/**
 * 27 91 51 126 - backspace
 */
func (e *Editor) ProcessKeyboardInput() {
  var input []byte = make([]byte, 4)

  os.Stdin.Read(input)

  //fmt.Fprintf(os.Stderr, "%d %d %d %d\n", input[0],input[1],input[2],input[3])

  buffer := e.getCurrentBuffer()

  if input[1] == 0 && input[2] == 0 {
    key := input[0]
    if key == 3 {  /* ctrl-c */
      fmt.Println("\x1b[2J")
      terminal.Restore(0, e.termOldState)
      os.Exit(0)
    } else if key == 13 { /* enter */
      buffer.InsertEmptyRow()
    } else if key == 17 { /* ctrl-q */
      fmt.Println("\x1b[2J")
      terminal.Restore(0, e.termOldState)
      os.Exit(0)
    } else if key == 19 { /* ctrl-s*/
      buffer.SaveChanges()
    } else if key == 27 { /* ESC */
      return
    } else if key == 32 {
      buffer.InsertChar(string(" "))
    } else if key >= 33 && key <= 126 {
      buffer.InsertChar(string(input[0]))
    } else if key == 127 { /* DELETE */
      buffer.DeleteChar()
    } else if key == 9 { /* TAB */
      buffer.InsertIndent()
    } else {
      return
    }
  }

  if input[0] == 27 {
    if input[1] == 91 {
      if input[2] == 68 {
        buffer.MoveCursorLeft()
      } else if input[2] == 67 {
        buffer.MoveCursorRight()
      } else if input[2] == 66 {
        buffer.MoveCursorDown()
      } else if input[2] == 65 {
        buffer.MoveCursorUp()
      } else if input[2] == 72 {
        // buffer.MoveToLineStart()
      } else if input[2] == 70 {
        // buffer.MoveToLineEnd()
      } else {
        return
      }
    }
  }

  return
}


func (e *Editor) nextBuffer() *buffer.Buffer {
  e.bufferIndex -= 1

  if e.bufferIndex >= len(e.buffers) {
    e.bufferIndex = len(e.buffers) - 1
  }

  return e.buffers[e.bufferIndex]
}

func (e *Editor) prevBuffer() *buffer.Buffer {
  e.bufferIndex -= 1

  if e.bufferIndex < 0 {
    e.bufferIndex = 0
  }

  return e.buffers[e.bufferIndex]
}

func (e *Editor) getCurrentBuffer() *buffer.Buffer {
  return e.buffers[e.bufferIndex]
}

func (e *Editor) Open(path string) {
  b := buffer.New()
  b.SetMaxVisibleRows(e.termRows)
  e.buffers = append(e.buffers, b)
}

func New() *Editor {
  ws := util.NewWindowSize() // terminal window size

  e := &Editor{
    reservedRowsTop: 1,
    reservedRowsBottom: 1,
    termRows: ws.GetRowsInt(),
  }

  e.termOldState, _ = terminal.MakeRaw(0)
  return e
}
