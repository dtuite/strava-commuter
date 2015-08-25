package main

import (
  "github.com/strava/go.strava"
  "regexp"
  "time"
  "log"
  "gopkg.in/yaml.v2"
  "path/filepath"
  "io/ioutil"
  "os"
)

// https://www.strava.com/oauth/authorize?client_id=7724&response_type=code&redirect_uri=http://localhost:3000&scope=write

// curl -X POST https://www.strava.com/oauth/token \
//      -F client_id=7724 \
//      -F client_secret=ceb5f0a3e4fd38e1bbbe18e4c1443f7abcd09785 \
//      -F code=70b48fab7ede3742952e39954115f04f61f90f65

func check(e error) {
  if e != nil {
    log.Fatal(e)
  }
}

var iso8601, _ = regexp.Compile("\\d\\d\\d\\d-\\d\\d-\\d\\dT\\d\\d:\\d\\d:\\d\\dZ")
const layout = "2006-01-02T15:04:05Z0700" // ISO 8601

func buildReplacer(distance time.Duration) func([]byte) []byte {
  return func(pointTime []byte) []byte {
    pointTimeTime, err := time.Parse(layout, string(pointTime))
    check(err)
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

func ReplaceFile(readFilepath string, writeFilepath string, newFinishTime time.Time) {
  template, _ := ioutil.ReadFile(readFilepath)

  distance := durationSinceTemplateFinish(template, newFinishTime)
  template = iso8601.ReplaceAllFunc(template, buildReplacer(distance))

  err := ioutil.WriteFile(writeFilepath, template, 0644)
  check(err)
}

type Config struct {
  AccessToken string `yaml:"access_token"`
  BikeGearId string `yaml:"bike_gear_id"`
  DefaultActivityDescription string `yaml:"default_activity_description"`
}

func (config *Config) Read(pathString string) {
  filename, _ := filepath.Abs(pathString)
  yamlFile, err := ioutil.ReadFile(filename)
  check(err)
  err = yaml.Unmarshal(yamlFile, &config)
  check(err)
}

func UploadGPX(config Config, gpxFilepath string) *strava.ActivityDetailed {
  client := strava.NewClient(config.AccessToken)
	uploadService := strava.NewUploadsService(client)
  activityService := strava.NewActivitiesService(client)

  fileReader, err := os.Open(gpxFilepath)
  check(err)

  log.Printf("About to upload file\n")

  upload, err := uploadService.
		Create(strava.FileDataTypes.GPX, "fixed-to-work.gpx", fileReader).
    ActivityType(strava.ActivityTypes.Ride).
    Name("Test Upload 1").
    Description(config.DefaultActivityDescription).
		Do()

  check(err)
  log.Printf("File uploaded.\n")

  time.Sleep(5 * time.Second)

  uploadSummary, err := uploadService.Get(upload.Id).Do()
  log.Printf("Got upload summary for upload ID: %v. %v\n", uploadSummary.Id, uploadSummary.ActivityId)

  activity, err := activityService.Update(uploadSummary.ActivityId).
    Commute(true).
    Gear(config.BikeGearId).
    Do()

  return activity
}

func main() {
  templateFile, _ := filepath.Abs("./to-work.gpx")
  targetFile, _ := filepath.Abs("./fixed-to-work.gpx")
  configFile, _ := filepath.Abs("./config.yml")

  config := Config{}
  config.Read(configFile)
  log.Printf("Using access token: %v\n", config.AccessToken)

  ReplaceFile(templateFile, targetFile, time.Now())

  log.Printf("Template times replaced\n")

  activity := UploadGPX(config, targetFile)

  log.Printf("Activity created successfully ID: %v\n", activity.Id)
}
