package stravacommuterconfig

import (
  "fmt"
  "log"
  "gopkg.in/yaml.v2"
  "path/filepath"
  "io/ioutil"
)

type Config struct {
  AccessToken string `yaml:"access_token"`
}

func (config *Config) Read(pathString string) {
  // TODO: Move this into a function to DRY it. I'm doing the same thing twice.
  filename, _ := filepath.Abs(pathString)
  yamlFile, err := ioutil.ReadFile(filename)

  if err != nil {
    log.Fatalf("Error reading file %v: %v\n", filename, err)
  }

  err = yaml.Unmarshal(yamlFile, &config)

  if err != nil {
    log.Fatalf("Error parsing YAML from %v: %v\n", filename, err)
  }

  fmt.Printf("Using access token from %v: %v\n", filename, config.AccessToken)
}
