package types

import (
  "bufio"
  "encoding/csv"
  "io"
  "os"
  "strconv"
  "strings"
  "sync"
  "errors"

  "github.com/loganwilliams/adviceservisory/server/gtfsstatic"
  "github.com/paulmach/go.geo"
)

type LineCache struct {
  lines map[string]*geo.Path
  sync.Mutex
}

var c *LineCache

func (r *Route) Measure(s *Stop) (float64, error) {
  routeCode := r.Id

  if routeCode == "6X" {
    routeCode = "6"
  } else if routeCode == "7X" {
    routeCode = "7"
  } else if routeCode == "W" {
    routeCode = "N"
  }

  path, err := getLinePath(routeCode)

  if err != nil {
    return -1.0, err
  }

  if path == nil {
    return -1.0, errors.New("nil path")
  }

  p := projectPoint(float64(s.Latitude), float64(s.Longitude))

  return path.Measure(p), nil
}

func getLinePath(r string) (*geo.Path, error) {
  if c == nil {
    c = &LineCache{}
  }

  c.Lock()
  defer c.Unlock()

  if c.lines != nil {
    return c.lines[r], nil
  }

  c.lines = make(map[string]*geo.Path)

  linesLocation, err := gtfsstatic.LinesLocation()

  if err != nil {
    return nil, err
  }

  f, err := os.Open(linesLocation)
  defer f.Close()

  if err != nil {
    return nil, err
  }

  reader := csv.NewReader(bufio.NewReader(f))

  line, err := reader.Read()

  prevRouteCode := "1"
  prevIndex := -1
  var routeCode, direction string
  var path *geo.Path
  sizes := make(map[string]int)
  sizes["1"] = 0

  path = geo.NewPath()

  for {
    line, err = reader.Read()

    if err == io.EOF {
      break
    } else if err != nil {
      return nil, err
    }

    shapeId := line[0]
    latitude, _ := strconv.ParseFloat(line[1], 64)
    longitude, _ := strconv.ParseFloat(line[2], 64)
    index, _ := strconv.ParseInt(line[3], 0, 0)

    if prevIndex+1 != int(index) {
      // fmt.Println(prevRouteCode, prevIndex)
      if direction == "N" && prevIndex > sizes[prevRouteCode] {
        c.lines[prevRouteCode] = path
        sizes[prevRouteCode] = prevIndex

        path = geo.NewPath()
      } else {
        path = geo.NewPath()
      }
    }

    splitString := strings.Split(shapeId, ".")
    routeCode = splitString[0]
    direction = string(splitString[len(splitString)-1][0])

    prevRouteCode = routeCode
    prevIndex = int(index)

    path.Push(projectPoint(latitude, longitude))
  }

  c.lines[routeCode] = path

  return c.lines[r], nil
}

func projectPoint(latitude, longitude float64) *geo.Point {
  // We only care about this projection for distance calculations.
  // The magic nuumber below scales the lat and lon to have the same
  // relative size.
  // return geo.NewPointFromLatLng(latitude, longitude)
  return geo.NewPoint(latitude*111, longitude*0.75798863709*111)
}
