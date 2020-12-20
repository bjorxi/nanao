package util

import (
  "syscall"
  "unsafe"
)

type WindowSize struct {
  row    uint16
  col    uint16
  xpixel uint16
  ypixel uint16
}


func NewWindowSize() *WindowSize {
  ws := &WindowSize{}
  retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
      uintptr(syscall.Stdin),
      uintptr(syscall.TIOCGWINSZ),
      uintptr(unsafe.Pointer(ws)))

  if int(retCode) == -1 {
      panic(errno)
  }

  return ws
}

func (ws *WindowSize) GetRowsInt() int {
  return int(ws.col)
}
