package editor

import "os"


/**
 * Checks if a file exists
 */
func fileExists (path string) bool {
  if _, err := os.Stat(path); os.IsNotExist(err) {
    return false
  }

  return true
}
