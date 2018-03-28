package types

import (
  "time"
  "database/sql"
  "log"
)

// The current position of a train.
type TripUpdate struct {
  TripId string
  Route *Route
  StopId string
  Timestamp time.Time
  Direction int
}

func CreateTripUpdatesTable(db *sql.DB) error {
  mkTripUpdateTableStmt := `CREATE TABLE trip_updates(
    id serial primary key,
    trip_id varchar(100) references trips(id),
    stop varchar(100),
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

func (tu *TripUpdate) GetTrip() *Trip {
  return &Trip{
    Id: tu.TripId,
    Route: tu.Route,
    Direction: tu.Direction,
  }
}