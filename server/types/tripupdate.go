package types

import (
  "time"
  "database/sql"
  "log"
  // "encoding/json"
  // "fmt"
)

// The current position of a train.
type TripUpdate struct {
  Id int              `json:"id"`
  Trip *Trip          `json:"trip"`
  Stop *Stop          `json:"stop"`
  Timestamp time.Time `json:"timestamp"`
}

var (
  readTripUpdateStmt *sql.Stmt
  readAllTripUpdatesStmt *sql.Stmt
  readTripUpdateFromStationStmt *sql.Stmt
)

func CreateTripUpdatesTable(db *sql.DB) error {
  mkTripUpdateTableStmt := `CREATE TABLE IF NOT EXISTS trip_updates(
    id serial primary key,
    trip_id varchar(100) references trips(id),
    stop varchar(10) references stops(id),
    timestamp timestamptz not null
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
  stmt := `INSERT INTO trip_updates(
      trip_id,
      stop,
      timestamp)
    SELECT CAST($1 AS VARCHAR), CAST($2 AS VARCHAR), $3
      WHERE
        NOT EXISTS (
          SELECT * FROM trip_updates WHERE
          trip_id = $1 AND stop = $2)`

  createStmt, err := db.Prepare(stmt)

  if err != nil {
    log.Fatal("error preparing statement,", err)
    return err
  }

  _, err = createStmt.Exec(tu.Trip.Id, tu.Stop.Id, tu.Timestamp)

  return err
}

func ReadAllUpdates(db *sql.DB) ([]*TripUpdate, error) {
  var (
    updates []*TripUpdate = []*TripUpdate{}
    err error
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

func (s *Stop) ReadUpdates(db *sql.DB) ([]*TripUpdate, error) {
  var (
    updates []*TripUpdate = []*TripUpdate{}
    err error
  )

  if readTripUpdateFromStationStmt == nil {
    stmt := `SELECT
              trip_updates.id AS id,
              trip_id,
              stop,
              timestamp,
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
            WHERE stops.station = $1`

    readTripUpdateFromStationStmt, err = db.Prepare(stmt)
    if err != nil {
      return updates, err
    }
  }

  rows, err := readTripUpdateFromStationStmt.Query(s.Station)

  if err != nil {
    return updates, err
  }

  defer rows.Close()

  for rows.Next() {
    update := &TripUpdate{Stop: &Stop{}, Trip: &Trip{Route: &Route{}}}

    err = rows.Scan(
      &update.Id,
      &update.Trip.Id,
      &update.Stop.Id,
      &update.Timestamp,
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

func (t *Trip) ReadUpdates(db *sql.DB) ([]*TripUpdate, error) {
  var (
    updates []*TripUpdate = []*TripUpdate{}
    err error
  )

  if readTripUpdateStmt == nil {
    stmt := `SELECT
              trip_updates.id AS id,
              trip_id,
              stop,
              timestamp,
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
            WHERE trip_id = $1`

    readTripUpdateStmt, err = db.Prepare(stmt)
    if err != nil {
      return updates, err
    }
  }

  rows, err := readTripUpdateStmt.Query(t.Id)

  if err != nil {
    return updates, err
  }

  defer rows.Close()

  for rows.Next() {
    update := &TripUpdate{Stop: &Stop{}, Trip: &Trip{Route: &Route{}}}

    err = rows.Scan(
      &update.Id,
      &update.Trip.Id,
      &update.Stop.Id,
      &update.Timestamp,
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
