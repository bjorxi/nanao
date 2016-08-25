package editor

import "os"
import "fmt"
import "io/ioutil"
import "encoding/json"


func (e *Editor) ParseConf(path string) {
  configFile, err := ioutil.ReadFile(path)

  if err != nil {
    fmt.Fprintf(os.Stderr,"Can't open conf file %e", err.Error())
  }

  err = json.Unmarshal(configFile, e.conf)

  if err != nil {
    fmt.Fprintf(os.Stderr, "Can't Unmarshall %e", err.Error())
  }
}
