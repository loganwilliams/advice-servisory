package types

import (
  "database/sql"
  "fmt"
  "log"
  "time"

  "github.com/lib/pq"
  "github.com/paulmach/go.geojson"
)

// The current position of a train.
type TripUpdate struct {
  Id        int       `json:"id"`
  Trip      *Trip     `json:"trip"`
  Stop      *Stop     `json:"stop"`
  Timestamp time.Time `json:"timestamp"`
  Progress  float64   `json:"progress"`
}

var (
  readTripUpdateStmt            *sql.Stmt
  readAllTripUpdatesStmt        *sql.Stmt
  readTripUpdateFromStationStmt *sql.Stmt
  liveUpdatesStmt               *sql.Stmt
  readRouteTripUpdateStmt       *sql.Stmt
  createStmt                    *sql.Stmt
  readTripUpdateFromStopStmt    *sql.Stmt
  readUpdatesWithTwoStops       *sql.Stmt
)

func CreateTripUpdatesTable(db *sql.DB) error {
  mkTripUpdateTableStmt := `CREATE TABLE IF NOT EXISTS trip_updates(
    id serial primary key,
    code varchar(50) UNIQUE,
    trip_id varchar(100) references trips(id),
    stop varchar(10) references stops(id),
    timestamp timestamptz not null,
    progress float
    )`

  _, err := db.Exec(mkTripUpdateTableStmt)
  if err != nil {
    log.Fatal("Error creating table `trip_updates`", err)
  }

  return nil
}

func DropTripUpdatesTable(db *sql.DB) error {
  dpStmt := `DROP TABLE IF EXISTS trip_updates CASCADE`
  _, err := db.Exec(dpStmt)
  if err != nil {
    log.Fatal("Error dropping table `trip_updates`", err)
  }

  return nil
}

func (tu *TripUpdate) Insert(db *sql.DB) error {
  var err error

  if createStmt == nil {
    stmt := `INSERT INTO trip_updates(
        trip_id,
        code,
        stop,
        timestamp,
        progress)
      SELECT $1::VARCHAR(100), $2::VARCHAR(40), $3::VARCHAR(10), $4, $5
      WHERE NOT EXISTS (SELECT 1 FROM trip_updates WHERE code = $2)`

    createStmt, err = db.Prepare(stmt)

    if err != nil {
      log.Fatal("error preparing statement,", err)
      return err
    }
  }

  code := tu.Timestamp.Format("2006-01-02T15") + "_" + tu.Trip.Id + "_" + tu.Stop.Id

  _, err = createStmt.Exec(tu.Trip.Id, code, tu.Stop.Id, tu.Timestamp, tu.Progress)

  if err != nil {
    if err, ok := err.(*pq.Error); ok {
      if err.Code.Name() == "foreign_key_violation" {
        // okay so the deal is like this
        // sometimes the MTA invents new stop names
        // so what we need to do is
        // add that new bullshit stop
        log.Println("Inserting unexpected new stop.")
        err := tu.Stop.Insert(db)

        if err != nil {
          log.Fatal("Could not insert new stop", err)
        }

        return tu.Insert(db)
      }
    } else {
      fmt.Println(err)
    }

  }

  return err
}

func ReadAllUpdates(db *sql.DB) ([]*TripUpdate, error) {
  var (
    updates []*TripUpdate = []*TripUpdate{}
    err     error
  )

  if readAllTripUpdatesStmt == nil {
    stmt := `SELECT
              id,
              trip_id,
              stop,
              timestamp
            FROM trip_updates`

    readAllTripUpdatesStmt, err = db.Prepare(stmt)
    if err != nil {
      return updates, err
    }
  }

  rows, err := readAllTripUpdatesStmt.Query()

  if err != nil {
    return updates, err
  }

  defer rows.Close()

  for rows.Next() {
    update := &TripUpdate{Stop: &Stop{}, Trip: &Trip{}}

    err = rows.Scan(
      &update.Id,
      &update.Trip.Id,
      &update.Stop.Id,
      &update.Timestamp)

    if err != nil {
      return updates, err
    }

    updates = append(updates, update)
  }

  if err := rows.Err(); err != nil {
    return updates, err
  }

  return updates, nil
}

func (s1 *Stop) UpdatesWithStop(db *sql.DB, s2 *Stop) ([]*TripUpdate, error) {
  var (
    updates []*TripUpdate = []*TripUpdate{}
    err     error
  )

  if readUpdatesWithTwoStops == nil {
    stmt := `SELECT
              trip_updates.id AS id,
              trip_id,
              stop,
              timestamp,
              progress,
              direction,
              route,
              stops.name AS stop_name,
              stops.latitude AS latitude,
              stops.longitude AS longitude,
              stops.station AS station
            FROM trip_updates
            LEFT OUTER JOIN trips ON trip_updates.trip_id = trips.id
            LEFT OUTER JOIN stops ON trip_updates.stop = stops.id
            WHERE trip_id IN (
              SELECT
                trip_id
              FROM trip_updates
              LEFT OUTER JOIN stops ON trip_updates.stop = stops.id
              WHERE stops.station = $1
              AND timestamp > $3)
            AND trip_id IN (
              SELECT
                trip_id
              FROM trip_updates
              LEFT OUTER JOIN stops ON trip_updates.stop = stops.id
              WHERE stops.station = $2
              AND timestamp > $3)
            AND (stops.station = $1
              OR stops.station = $2)
            AND timestamp > $3
            ORDER BY trip_id, timestamp DESC`

    readUpdatesWithTwoStops, err = db.Prepare(stmt)
    if err != nil {
      return updates, err
    }
  }

  rows, err := readUpdatesWithTwoStops.Query(s1.Station, s2.Station, time.Now().Add(-22*time.Hour))

  if err != nil {
    return updates, err
  }

  defer rows.Close()

  return scanRowsShort(rows)
}

func (s *Stop) ReadUpdatesAtDate(db *sql.DB, time.Time *date) ([]*TripUpdate, error) {
  var (
    updates []*TripUpdate = []*TripUpdate{}
    err     error
  )

  // We might want to query based on the stop or the station,
  // depending on what values are set in the s Object.

  if readTripUpdateFromStationStmt == nil {
    stmt := `SELECT
              trip_updates.id AS id,
              trip_id,
              stop,
              timestamp,
              progress,
              direction,
              route,
              stops.name AS stop_name,
              stops.latitude AS latitude,
              stops.longitude AS longitude,
              stops.station AS station
            FROM trip_updates
            LEFT OUTER JOIN trips ON trip_updates.trip_id = trips.id
            LEFT OUTER JOIN stops ON trip_updates.stop = stops.id
            WHERE trip_id IN (
              SELECT
                trip_id
              FROM trip_updates
              LEFT OUTER JOIN stops ON trip_updates.stop = stops.id
              WHERE stops.station = $1
              AND timestamp >= $2
              AND timestamp < $3)
            AND timestamp >= $2
            AND timestamp < $3
            ORDER BY trip_id, timestamp DESC`

    readTripUpdateFromStationStmt, err = db.Prepare(stmt)

    if err != nil {
      return updates, err
    }
  }

  if readTripUpdateFromStopStmt == nil {
    stmt := `SELECT
              trip_updates.id AS id,
              trip_id,
              stop,
              timestamp,
              progress,
              direction,
              route,
              stops.name AS stop_name,
              stops.latitude AS latitude,
              stops.longitude AS longitude,
              stops.station AS station
            FROM trip_updates
            LEFT OUTER JOIN trips ON trip_updates.trip_id = trips.id
            LEFT OUTER JOIN stops ON trip_updates.stop = stops.id
            WHERE trip_id IN (
              SELECT
                trip_id
              FROM trip_updates
              WHERE trip_updates.stop = $1
              AND timestamp >= $2
              AND timestamp < $3)
            AND timestamp >= $2
            AND timestamp < $3
            ORDER BY trip_id, timestamp DESC`

    readTripUpdateFromStopStmt, err = db.Prepare(stmt)
    if err != nil {
      return updates, err
    }
  }

  query := readTripUpdateFromStopStmt
  queryVal := s.Id

  if s.Id == "" {
    fmt.Println(s.Station)
    query = readTripUpdateFromStationStmt
    queryVal = s.Station
  }

  startTime := date
  endTime := date.Add(24*time.Hour)

  rows, err := query.Query(queryVal, startTime, endTime)

  if err != nil {
    return updates, err
  }

  defer rows.Close()

  return scanRowsShort(rows)
}

func (s *Stop) ReadUpdates(db *sql.DB) ([]*TripUpdate, error) {
  var (
    updates []*TripUpdate = []*TripUpdate{}
    err     error
  )

  // We might want to query based on the stop or the station,
  // depending on what values are set in the s Object.

  if readTripUpdateFromStationStmt == nil {
    stmt := `SELECT
              trip_updates.id AS id,
              trip_id,
              stop,
              timestamp,
              progress,
              direction,
              route,
              stops.name AS stop_name,
              stops.latitude AS latitude,
              stops.longitude AS longitude,
              stops.station AS station
            FROM trip_updates
            LEFT OUTER JOIN trips ON trip_updates.trip_id = trips.id
            LEFT OUTER JOIN stops ON trip_updates.stop = stops.id
            WHERE trip_id IN (
              SELECT
                trip_id
              FROM trip_updates
              LEFT OUTER JOIN stops ON trip_updates.stop = stops.id
              WHERE stops.station = $1
              AND timestamp > $2)
            AND timestamp > $2
            ORDER BY trip_id, timestamp DESC`

    readTripUpdateFromStationStmt, err = db.Prepare(stmt)
    if err != nil {
      return updates, err
    }
  }

  if readTripUpdateFromStopStmt == nil {
    stmt := `SELECT
              trip_updates.id AS id,
              trip_id,
              stop,
              timestamp,
              progress,
              direction,
              route,
              stops.name AS stop_name,
              stops.latitude AS latitude,
              stops.longitude AS longitude,
              stops.station AS station
            FROM trip_updates
            LEFT OUTER JOIN trips ON trip_updates.trip_id = trips.id
            LEFT OUTER JOIN stops ON trip_updates.stop = stops.id
            WHERE trip_id IN (
              SELECT
                trip_id
              FROM trip_updates
              WHERE trip_updates.stop = $1
              AND timestamp > $2)
            AND timestamp > $2
            ORDER BY trip_id, timestamp DESC`

    readTripUpdateFromStopStmt, err = db.Prepare(stmt)
    if err != nil {
      return updates, err
    }
  }

  query := readTripUpdateFromStopStmt
  queryVal := s.Id

  if s.Id == "" {
    fmt.Println(s.Station)
    query = readTripUpdateFromStationStmt
    queryVal = s.Station
  }

  rows, err := query.Query(queryVal, time.Now().Add(-24*time.Hour))

  if err != nil {
    return updates, err
  }

  defer rows.Close()

  return scanRowsShort(rows)
}

func LiveUpdates(db *sql.DB) ([]*TripUpdate, error) {
  var (
    updates []*TripUpdate = []*TripUpdate{}
    err     error
  )

  if liveUpdatesStmt == nil {
    stmt := `SELECT DISTINCT ON (trip_id)
              trip_updates.id AS id,
              trip_id,
              stop, 
              timestamp, 
              progress,
              direction,
              routes.id AS route_id,
              COALESCE(short_name, '') AS short_name,
              COALESCE(routes.name, '') AS route_name,
              COALESCE(description, '') AS description,
              COALESCE(type, 0) AS type,
              COALESCE(url, '') AS url,
              COALESCE(color, '') AS color,
              stops.name AS stop_name,
              stops.latitude AS latitude,
              stops.longitude AS longitude,
              stops.station AS station
            FROM trip_updates
            LEFT OUTER JOIN trips ON trip_updates.trip_id = trips.id
            LEFT OUTER JOIN routes ON trips.route = routes.id
            LEFT OUTER JOIN stops ON trip_updates.stop = stops.id
            WHERE timestamp > $1
            ORDER BY trip_id, timestamp DESC`

    liveUpdatesStmt, err = db.Prepare(stmt)
    if err != nil {
      return updates, err
    }
  }

  rows, err := liveUpdatesStmt.Query(time.Now().Add(-20 * time.Minute))

  if err != nil {
    return updates, err
  }

  defer rows.Close()

  return scanRows(rows)
}

func (r *Route) ReadUpdates(db *sql.DB) ([]*TripUpdate, error) {
  var (
    updates []*TripUpdate = []*TripUpdate{}
    err     error
  )

  if readRouteTripUpdateStmt == nil {
    stmt := `SELECT
              trip_updates.id AS id,
              trip_id,
              stop,
              timestamp,
              progress,
              direction,
              route,
              stops.name AS stop_name,
              stops.latitude AS latitude,
              stops.longitude AS longitude,
              stops.station AS station
            FROM trip_updates
            LEFT OUTER JOIN trips ON trip_updates.trip_id = trips.id
            LEFT OUTER JOIN stops ON trip_updates.stop = stops.id
            WHERE trips.route = $1 AND timestamp > $2
            ORDER BY trip_id, timestamp DESC`

    readRouteTripUpdateStmt, err = db.Prepare(stmt)
    if err != nil {
      return updates, err
    }
  }

  rows, err := readRouteTripUpdateStmt.Query(r.Id, time.Now().Add(-24*time.Hour))

  if err != nil {
    return updates, err
  }

  defer rows.Close()

  return scanRowsShort(rows)
}

func (t *Trip) ReadUpdates(db *sql.DB) ([]*TripUpdate, error) {
  var (
    updates []*TripUpdate = []*TripUpdate{}
    err     error
  )

  if readTripUpdateStmt == nil {
    stmt := `SELECT
              trip_updates.id AS id,
              trip_id,
              stop,
              timestamp,
              progress,
              direction,
              routes.id AS route_id,
              COALESCE(short_name, '') AS short_name,
              COALESCE(routes.name, '') AS route_name,
              COALESCE(description, '') AS description,
              COALESCE(type, 0) AS type,
              COALESCE(url, '') AS url,
              COALESCE(color, '') AS color,
              stops.name AS stop_name,
              stops.latitude AS latitude,
              stops.longitude AS longitude,
              stops.station AS station
            FROM trip_updates
            LEFT OUTER JOIN trips ON trip_updates.trip_id = trips.id
            LEFT OUTER JOIN routes ON trips.route = routes.id
            LEFT OUTER JOIN stops ON trip_updates.stop = stops.id
            WHERE trip_id = $1 AND timestamp > $2
            ORDER BY timestamp DESC`

    readTripUpdateStmt, err = db.Prepare(stmt)
    if err != nil {
      return updates, err
    }
  }

  rows, err := readTripUpdateStmt.Query(t.Id, time.Now().Add(-24*time.Hour))

  if err != nil {
    return updates, err
  }

  defer rows.Close()

  return scanRows(rows)
}

func scanRows(rows *sql.Rows) ([]*TripUpdate, error) {
  var updates []*TripUpdate = []*TripUpdate{}

  for rows.Next() {
    update := &TripUpdate{Stop: &Stop{}, Trip: &Trip{Route: &Route{}}}

    err := rows.Scan(
      &update.Id,
      &update.Trip.Id,
      &update.Stop.Id,
      &update.Timestamp,
      &update.Progress,
      &update.Trip.Direction,
      &update.Trip.Route.Id,
      &update.Trip.Route.ShortName,
      &update.Trip.Route.Name,
      &update.Trip.Route.Description,
      &update.Trip.Route.Type,
      &update.Trip.Route.URL,
      &update.Trip.Route.Color,
      &update.Stop.Name,
      &update.Stop.Latitude,
      &update.Stop.Longitude,
      &update.Stop.Station)

    if err != nil {
      return updates, err
    }

    updates = append(updates, update)
  }

  if err := rows.Err(); err != nil {
    return updates, err
  }

  return updates, nil
}

func scanRowsShort(rows *sql.Rows) ([]*TripUpdate, error) {
  var updates []*TripUpdate = []*TripUpdate{}

  for rows.Next() {
    update := &TripUpdate{Stop: &Stop{}, Trip: &Trip{Route: &Route{}}}

    err := rows.Scan(
      &update.Id,
      &update.Trip.Id,
      &update.Stop.Id,
      &update.Timestamp,
      &update.Progress,
      &update.Trip.Direction,
      &update.Trip.Route.Id,
      &update.Stop.Name,
      &update.Stop.Latitude,
      &update.Stop.Longitude,
      &update.Stop.Station)

    if err != nil {
      return updates, err
    }

    updates = append(updates, update)
  }

  if err := rows.Err(); err != nil {
    return updates, err
  }

  return updates, nil
}

// MakeGeoJSON takes a list of TripUpdates objects and constructs a GeoJSON FeatureCollection.
func MakeGeoJSON(updates []*TripUpdate) *geojson.FeatureCollection {
  fc := geojson.NewFeatureCollection()

  for _, update := range updates {
    f := geojson.NewPointFeature([]float64{float64(update.Stop.Longitude), float64(update.Stop.Latitude)})
    f.SetProperty("trip", update.Trip.Id)
    f.SetProperty("stop", update.Stop.Id)
    f.SetProperty("line", update.Trip.Route.ShortName)
    f.SetProperty("color", update.Trip.Route.Color)
    f.SetProperty("time", update.Timestamp)
    f.SetProperty("direction", update.Trip.Direction)
    fc.AddFeature(f)
  }

  return fc
}
