package main

import (
  // "github.com/strava/go.strava"
  // "time"
  // "os"
  "fmt"
  "log"
  "gopkg.in/yaml.v2"
  "path/filepath"
  "io/ioutil"
)

// https://www.strava.com/oauth/authorize?client_id=7724&response_type=code&redirect_uri=http://localhost:3000&scope=write

// curl -X POST https://www.strava.com/oauth/token \
//      -F client_id=7724 \
//      -F client_secret=ceb5f0a3e4fd38e1bbbe18e4c1443f7abcd09785 \
//      -F code=70b48fab7ede3742952e39954115f04f61f90f65

// type Activity struct {
//   Name string
//   Duration int
//   Description string
//   GearId string `yaml:"gear_id"`
//   Distance float64
//   IsCommute bool `yaml:"commute"`
// }

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

func main() {
  config := Config{}
  config.Read("./config.yml")
  fmt.Printf("Got the access token %v\n", config.AccessToken)

  // filename, _ := filepath.Abs("./activity.yml")
  // yamlFile, err := ioutil.ReadFile(filename)

  // if err != nil {
  //   log.Fatalf("Error reading file %v: %v\n", filename, err)
  // }

  // activity := Activity{}
  // err = yaml.Unmarshal(yamlFile, &activity)

  // if err != nil {
  //   log.Fatalf("Error parsing YAML from %v: %v\n", filename, err)
  // }

  // fmt.Printf("Uploading activity from %v\n", filename)
  // fmt.Printf("Uploading\n%v\n\n", activity)

  // client := strava.NewClient(config.AccessToken)
  // service := strava.NewActivitiesService(client)

  // // TODO: Parse the start time from the YAML
  // startTime := time.Now()

  // stravaActivity, err := service.Create(activity.Name, strava.ActivityTypes.Ride, startTime, activity.Duration).
  //                                 Description(activity.Description).
  //                                 Distance(activity.Distance).
  //                                 Do()

  // if (err != nil) {
  //   fmt.Fprintf(os.Stderr, "error: %v\n", err)
  //   os.Exit(1)
  // }

  // _, err = service.Update(stravaActivity.Id).
  //   Commute(activity.IsCommute).
  //   Gear(activity.GearId).
  //   Do()

  // if (err != nil) {
  //   fmt.Fprintf(os.Stderr, "error: %v\n", err)
  //   os.Exit(1)
  // }

  // fmt.Printf("Activity created https://strava.com/activities/%v\n", stravaActivity.Id)
}
