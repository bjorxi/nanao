package editor

import "fmt"
import "io/ioutil"
import "encoding/json"


func (e *Editor) ParseConf(path string) {
  configFile, err := ioutil.ReadFile(path)

  if err != nil {
    fmt.Println("Can't open conf file", err.Error())
  }

  err = json.Unmarshal(configFile, e.conf)

  if err != nil {
    fmt.Println("Can't Unmarshall", err.Error())
  }

  // configFile.Close()
}
