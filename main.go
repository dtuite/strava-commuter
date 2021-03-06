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
  "flag"
  "fmt"
  "strings"
  "os/user"
)

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
  DefaultActivityDuration int `yaml:"default_activity_duration"`
  DefaultActivityDistance float64 `yaml:"default_activity_distance"`
  DefaultActivityIsPrivate bool `yaml:"default_activity_is_private"`
  DefaultActivityIsCommute bool `yaml:"default_activity_is_commute"`
}

func (config *Config) Read(pathString string) {
  filename, _ := filepath.Abs(pathString)
  yamlFile, err := ioutil.ReadFile(filename)
  check(err)
  err = yaml.Unmarshal(yamlFile, &config)
  check(err)
}

func UploadGPX(config Config, gpxFilename string) *strava.ActivityDetailed {
  client := strava.NewClient(config.AccessToken)
	uploadService := strava.NewUploadsService(client)
  activityService := strava.NewActivitiesService(client)

  activityName := strings.Replace(gpxFilename, "-", " ", -1)
  activityName = strings.TrimSuffix(activityName, ".gpx")
  activityName = strings.Title(activityName)

  targetFile, _ := filepath.Abs(gpxFilename)
  fileReader, err := os.Open(targetFile)
  check(err)

  log.Printf("About to upload file\n")

  upload, err := uploadService.
		Create(strava.FileDataTypes.GPX, gpxFilename, fileReader).
    ActivityType(strava.ActivityTypes.Ride).
    Name(activityName).
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

func CreateManual(config Config, activityName string, finishTime time.Time) *strava.ActivityDetailed {
  client := strava.NewClient(config.AccessToken)
  activityService := strava.NewActivitiesService(client)

  durationString, err := time.ParseDuration(fmt.Sprintf("-%vs",config.DefaultActivityDuration))
  check(err)
  activityStartTime := finishTime.Add(durationString)

  activity, err := activityService.Create(activityName, strava.ActivityTypes.Ride, activityStartTime, config.DefaultActivityDuration).
    Description(config.DefaultActivityDescription).
    Distance(config.DefaultActivityDistance).
    Do()

  check(err)

  activity, err = activityService.Update(activity.Id).
    Private(config.DefaultActivityIsPrivate).
    Commute(config.DefaultActivityIsCommute).
    Gear(config.BikeGearId).
    Do()

  check(err)

  return activity
}

func main() {
  // gpxFileFlag := flag.String("route", "to-work.gpx", "The path to the GPX file to repeat")
  configFileFlag := flag.String("config", ".strava-commuter/config.yml", "The path to the config file in your home directory.")
  finishTimeFlag := flag.String("finish-time", time.Now().Format("15:04"), "The time you finished the activity. 24 hour, zero padded.")
  finishDateFlag := flag.String("finish-date", time.Now().Format("2006-01-02"), "The date you finished the activity. Zero padded.")
  activityName := flag.String("name", "To Work", "The activity name")

  flag.Parse()

  usr, err := user.Current()
  check(err)

  // templateFile, _ := filepath.Abs(*gpxFileFlag)
  configFile := filepath.Join(usr.HomeDir, *configFileFlag)

  config := Config{}
  config.Read(configFile)
  log.Printf("Using access token: %v\n", config.AccessToken)

  finishDateTime, err := time.Parse(layout, fmt.Sprintf("%vT%v:00Z", *finishDateFlag, *finishTimeFlag))
  check(err)

  // ReplaceFile(templateFile, fmt.Sprintf("./fixed-%v", *gpxFileFlag), finishDateTime)

  // log.Printf("Template times replaced\n")

  // activity := UploadGPX(config, *gpxFileFlag)

  activity := CreateManual(config, *activityName, finishDateTime)

  log.Printf("Activity created successfully.\n")
  log.Printf("https://strava.com/activities/%v\n", activity.Id)
}
