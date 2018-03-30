package types

import (
  "database/sql"
  "fmt"
  "log"

  "github.com/loganwilliams/adviceservisory/server/gtfsstatic"
)

type Stop struct {
  Id        string  `json:"id"`
  Name      string  `json:"name"`
  Station   string  `json:"station"`
  Latitude  float32 `json:"latitude"`
  Longitude float32 `json:"longitude"`
}

var (
  stopReadAllStmt *sql.Stmt
)

func CreateAndPopulateStopsTable(db *sql.DB) error {
  mkImportTableStmt := `CREATE SCHEMA IF NOT EXISTS impor;
    CREATE TABLE import.stops(
      stop_id text,
      stop_code text,
      stop_name text,
      stop_desc text,
      stop_lat text,
      stop_lon text,
      zone_id text,
      stop_url text,
      location_type text,
      parent_station text
      )`

  _, err := db.Exec(mkImportTableStmt)

  if err != nil {
    log.Fatal("Error creating import.stops", err)
  }

  stopsLocation, err := gtfsstatic.StopsLocation()

  if err != nil {
    log.Fatal("Error getting stopsLocation", err)
  }

  copyFromCSVStmt := fmt.Sprintf(`COPY import.stops FROM 
  '%s' WITH DELIMITER ',' HEADER CSV`, stopsLocation)

  _, err = db.Exec(copyFromCSVStmt)
  if err != nil {
    log.Fatal("Error copying routes from CSV")
  }

  mkTableStmt := `CREATE TABLE stops(
    id varchar(10) primary key,
    name varchar(100),
    station varchar(10),
    latitude float,
    longitude float
    )`

  _, err = db.Exec(mkTableStmt)
  if err != nil {
    log.Fatal("Error creating table `stops`", err)
  }

  // insert into routes table
  insertStmt := `INSERT INTO stops(
    id,
    name,
    station,
    latitude,
    longitude
  ) SELECT
    import.stops.stop_id AS id,
    import.stops.stop_name AS name,
    import.stops.parent_station AS station,
    import.stops.stop_lat::float AS latitude,
    import.stops.stop_lon::float AS longitude
  FROM import.stops
  WHERE import.stops.location_type = '0'`

  _, err = db.Exec(insertStmt)
  if err != nil {
    log.Fatal("Error inserting stops", err)
  }

  return nil
}

func DropStopsTable(db *sql.DB) error {
  dropStmt := `DROP TABLE IF EXISTS stops CASCADE;
    DROP TABLE IF EXISTS import.stops CASCADE;
    DROP TABLE IF EXISTS stations CASCADE;`
  _, err := db.Exec(dropStmt)
  if err != nil {
    log.Fatal("Error dropping stops table", err)
  }

  return nil
}

func ReadAllStops(db *sql.DB) ([]*Stop, error) {
  var (
    stops []*Stop = []*Stop{}
    err   error
  )

  // prepare statement if not already done so.
  if stopReadAllStmt == nil {
    stmt := `SELECT
              stops.id AS id,
              stops.name AS name,
              stops.latitude AS latitude,
              stops.longitude AS longitude,
              stops.station AS station
             FROM stops`

    stopReadAllStmt, err = db.Prepare(stmt)
    if err != nil {
      log.Fatal("Error preparing statement: ", err)
    }
  }

  rows, err := stopReadAllStmt.Query()
  if err != nil {
    return stops, err
  }

  defer rows.Close()
  for rows.Next() {
    stop := &Stop{}

    err = rows.Scan(
      &stop.Id,
      &stop.Name,
      &stop.Latitude,
      &stop.Longitude,
      &stop.Station)

    if err != nil {
      return stops, err
    }

    stops = append(stops, stop)
  }

  if err := rows.Err(); err != nil {
    return stops, err
  }

  return stops, nil
}
