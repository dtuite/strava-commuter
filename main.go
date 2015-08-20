package main

import (
  // "github.com/strava/go.strava"
  "regexp"
  "time"
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

func check(e error) {
  if e != nil {
    panic(e)
  }
}

var iso8601, _ = regexp.Compile("\\d\\d\\d\\d-\\d\\d-\\d\\dT\\d\\d:\\d\\d:\\d\\dZ")
const layout = "2006-01-02T15:04:05Z0700" // ISO 8601

func buildReplacer(distance time.Duration) func([]byte) []byte {
  return func(pointTime []byte) []byte {
    pointTimeTime, err := time.Parse(layout, string(pointTime))
    check(err)
    // fmt.Printf("Parsed time: %v -- %v -- %v\n", string(pointTime), pointTimeTime)
    pointTimeTime = pointTimeTime.Add(distance)
    return []byte(pointTimeTime.Format(layout))
  }
}

func durationSinceTemplateFinish(gpx []byte, now time.Time) time.Duration {
  gpxTimes := iso8601.FindAll(gpx, -1)
  gpxFinishTimeString := string(gpxTimes[len(gpxTimes) - 1])
  gpxFinishTime, _ := time.Parse(layout, gpxFinishTimeString)
  return now.Sub(gpxFinishTime)
}

func ReplaceFile(filepath string, newFilepath string, newFinishTime time.Time) {
  template, _ := ioutil.ReadFile(filepath)

  distance := durationSinceTemplateFinish(template, newFinishTime)
  template = iso8601.ReplaceAllFunc(template, buildReplacer(distance))

  err := ioutil.WriteFile(newFilepath, template, 0644)
  check(err)
}

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

  ReplaceFile("./to-work.gpx", "fixed-to-work.gpx", time.Now())

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
