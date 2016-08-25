package editor

import "os"
import "fmt"
import "io/ioutil"
import "encoding/json"


func (e *Editor) ParseConf(path string) {
  fmt.Fprintf(os.Stderr,"Parsing %s\n", path)
  configFile, err := ioutil.ReadFile(path)

  if err != nil {
    fmt.Fprintf(os.Stderr,"Can't open conf file %e", err.Error())
  }

  fmt.Fprintf(os.Stderr, "%s\n", string(configFile))
  err = json.Unmarshal(configFile, e.conf)

  fmt.Fprintf(os.Stderr, "Tab size %d\n", e.conf.TabSize)

  if err != nil {
    fmt.Fprintf(os.Stderr, "Can't Unmarshall %e", err.Error())
  }

  configFile.Close()
}
