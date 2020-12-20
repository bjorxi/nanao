package config

import "os"
import "fmt"
import "io/ioutil"
import "encoding/json"

type Config struct {
  ShowTabs bool `json:"show_tabs"`
  ShowStatusLine bool `json:"show_status_line"`
  TabSize int `json:"tab_size"`
  TranslateTabsToSpaces bool `json:"translate_tabs_to_spaces"`
}

var _conf *Config

func New() *Config {
  if _conf != nil {
    return _conf
  }

  _conf := &Config{
    ShowTabs: true,
    ShowStatusLine: true,
    TabSize: 2,
    TranslateTabsToSpaces: true,
  }

  return _conf
}


func Parse(path string) *Config {
  conf := &Config{}

  configFile, err := ioutil.ReadFile(path)

  if err != nil {
    fmt.Fprintf(os.Stderr,"Can't open conf file %e", err.Error())
  }

  err = json.Unmarshal(configFile, conf)

  if err != nil {
    fmt.Fprintf(os.Stderr, "Can't Unmarshall %e", err.Error())
  }

  return conf;
}
