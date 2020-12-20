package editor

import (
  "github.com/bjorxi/nanao/buffer"
  "github.com/bjorxi/nanao/terminal"
)



type Editor struct {
  rows int /* Number of rows terminal has*/
  cols int /* Number of columns terminal has */

  // How many row reserved for displaying meta info
  reservedRowsTop int
  reservedRowsBottom int

  buffers []*buffer.Buffer
  bufferIndex int

  termRows int

  termOldState *terminal.State /* TODO: move it somewhere */
}
