package buffer

import (
  "bytes"
)

type Row struct {
  number int
  content *bytes.Buffer
  size int
}

// adds a single charachter to a row
func (r *Row) insertChar (char string, index int) {
  content := r.content.Bytes()
  newBuffer := bytes.NewBuffer(nil)

  newBuffer.Write(content[:index])
  newBuffer.Write([]byte(char))
  newBuffer.Write(content[index:])

  r.content = newBuffer
  r.size = newBuffer.Len()
}
