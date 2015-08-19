package main

import (
  "io/ioutil"
  "regexp"
  "time"
  "fmt"
)

func check(e error) {
  if e != nil {
    panic(e)
  }
}

func RewriteGPX(data []byte) {

}

func main() {
  data, _ := ioutil.ReadFile("./to-work.gpx")
  r, _ := regexp.Compile("\\d\\d\\d\\d-\\d\\d-\\d\\d")

  year, month, day := time.Now().Date()
  data = r.ReplaceAll(data, []byte(fmt.Sprintf("%d-%02d-%02d", year, month, (day - 1))))
  err := ioutil.WriteFile("./fixed-to-work.gpx", data, 0644)
  check(err)
}
