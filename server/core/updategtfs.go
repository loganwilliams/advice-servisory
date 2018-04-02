package core

import (
  "bytes"
  "errors"
  "fmt"
  "log"
  "net/http"
  "strings"
  "time"

  "github.com/golang/protobuf/proto"
  "github.com/loganwilliams/adviceservisory/server/transit_realtime"
  "github.com/loganwilliams/adviceservisory/server/types"
)

// GetLiveTrains() returns a GeoJSON []byte object with the most recent position of all trains in the NYC Subway, as
// reported by the MTA's GTFS feed.
func (a *AdviceServisory) GetTripUpdates() []*types.TripUpdate {
  // The MTA has several different endpoints for different lines. My API key is in here, but the abuse potential
  // seems low enough that I'm okay with that.
  datafeeds := [](string){
    "http://datamine.mta.info/mta_esi.php?key=" + a.Config.MTA.Key + "&feed_id=1",  // 123456S
    "http://datamine.mta.info/mta_esi.php?key=" + a.Config.MTA.Key + "&feed_id=26", // ACE
    "http://datamine.mta.info/mta_esi.php?key=" + a.Config.MTA.Key + "&feed_id=16", // NQRW
    "http://datamine.mta.info/mta_esi.php?key=" + a.Config.MTA.Key + "&feed_id=21", // BDFM
    "http://datamine.mta.info/mta_esi.php?key=" + a.Config.MTA.Key + "&feed_id=2",  // L
    "http://datamine.mta.info/mta_esi.php?key=" + a.Config.MTA.Key + "&feed_id=11", // SIR
    "http://datamine.mta.info/mta_esi.php?key=" + a.Config.MTA.Key + "&feed_id=31", // G
    "http://datamine.mta.info/mta_esi.php?key=" + a.Config.MTA.Key + "&feed_id=36", // JZ
    "http://datamine.mta.info/mta_esi.php?key=" + a.Config.MTA.Key + "&feed_id=51", // 7
  }

  var updates []*types.TripUpdate

  for _, url := range datafeeds {
    transit, err := getGTFS(url, 3)

    if err != nil {
      log.Println("Error getting GTFS feed: ", err)
    } else {
      updates = append(updates, a.trainList(transit, time.Now())...)
    }
  }

  return updates
}

func (a *AdviceServisory) AddTripUpdates() error {
  updates := a.GetTripUpdates()

  for _, update := range updates {
    // make sure the trip exists in the DB
    err := update.Trip.Upsert(a.DB)

    if err != nil {
      return err
    }

    update.Insert(a.DB)
  }

  return nil
}

// Returns a list of trains from an unmarshalled protobuf that have had an update within
// 10 minutes of the time "now".
func (a *AdviceServisory) trainList(transit *transit_realtime.FeedMessage, now time.Time) (updates []*types.TripUpdate) {
  cutoff := now.Add(-10.0 * time.Minute)

  for _, entity := range transit.Entity {
    update, err := a.trainPositionFromTripUpdate(entity)

    if err == nil {
      // Only include trains that have moved in the last 10 minutes, are reporting times in the present/past
      // and have a line associated with them.
      if update.Timestamp.After(cutoff) && update.Timestamp.Before(now) {
        updates = append(updates, update)
      }
    }
  }

  return
}

// trainPositionFromTripUpdate takes a GTFS protobuf entity and returns a Train object. If there is no
// trip update in the GTFS entity, it returns an empty Train and an error.
func (a *AdviceServisory) trainPositionFromTripUpdate(entity *transit_realtime.FeedEntity) (*types.TripUpdate, error) {
  if entity.TripUpdate == nil {
    return &types.TripUpdate{}, errors.New("No trip update in entity.")
  }

  tripId := entity.GetTripUpdate().GetTrip().GetTripId()
  direction := directionFromId(tripId)

  routeId := entity.GetTripUpdate().GetTrip().GetRouteId()
  stopTimes := entity.GetTripUpdate().GetStopTimeUpdate()
  if stopTimes == nil {
    return nil, errors.New("No stop times list.")
  }
  timestamp := time.Unix(int64(stopTimes[0].GetArrival().GetTime()), 0)
  stopId := stopTimes[0].GetStopId()

  stop := &types.Stop{Id: stopId}
  stop.GetDetails(a.DB)

  // fmt.Println(json.Marshal(&stop))

  route := &types.Route{Id: routeId}

  progress, err := route.Measure(stop)

  if err != nil {
    return nil, err
  }

  update := &types.TripUpdate{
    Trip:      &types.Trip{Id: tripId, Route: route, Direction: direction},
    Stop:      stop,
    Timestamp: timestamp,
    Progress:  progress}

  return update, nil
}

// Using the Trip ID, return a direction.
func directionFromId(id string) int {
  idParts := strings.Split(id, ".")
  direction := string(idParts[len(idParts)-1][0])

  switch direction {
  case "N":
    return types.North
  default:
    return types.South
  }
}

// getGTFS downloads a GTFS url from the MTA and unmarshals the protobuf.
func getGTFS(url string, retries int) (*transit_realtime.FeedMessage, error) {
  if retries <= 0 {
    return nil, fmt.Errorf("giving up on url %q", url)
  }

  resp, err := http.Get(url)
  defer resp.Body.Close()

  if err != nil {
    fmt.Printf("failed to fetch for url %q", url)
  }

  buf := new(bytes.Buffer)
  buf.ReadFrom(resp.Body)
  gtfs := buf.Bytes()

  transit := &transit_realtime.FeedMessage{}
  if err := proto.Unmarshal(gtfs, transit); err != nil {
    log.Println("Failed to parse GTFS feed", err)
    time.Sleep(time.Second)
    return getGTFS(url, retries-1)
  }

  return transit, nil
}
