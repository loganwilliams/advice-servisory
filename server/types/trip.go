package types

import (
  "database/sql"
  "log"

  "github.com/lib/pq"
)

type Trip struct {
  Id        string `json:"id"`
  Route     *Route `json:"route"`
  Direction int    `json:"direction"`
}

const (
  North = iota
  South
)

var (
  tripReadAllStmt *sql.Stmt
  existsStmt      *sql.Stmt
  tripInsertStmt  *sql.Stmt
)

func CreateTripsTable(db *sql.DB) error {
  mkTripTableStmt := `CREATE TABLE IF NOT EXISTS trips(
    id varchar(100) primary key,
    route varchar(5) references routes(id),
    direction int)`

  _, err := db.Exec(mkTripTableStmt)
  if err != nil {
    log.Fatal("Error creating table `trips`", err)
  }

  return nil
}

func DropTripsTable(db *sql.DB) error {
  dpStmt := `DROP TABLE IF EXISTS trips CASCADE`
  _, err := db.Exec(dpStmt)
  if err != nil {
    log.Fatal("Error dropping table `trips`", err)
  }

  return nil
}

func (t *Trip) Insert(db *sql.DB) error {
  var err error

  if tripInsertStmt == nil {
    stmt := `INSERT INTO trips(
      id,
      route,
      direction
    ) VALUES ($1, $2, $3)`

    tripInsertStmt, err = db.Prepare(stmt)

    if err != nil {
      return err
    }
  }

  _, err = tripInsertStmt.Exec(t.Id, t.Route.Id, t.Direction)

  if err != nil {
    if err, ok := err.(*pq.Error); ok {
      if err.Code.Name() == "foreign_key_violation" {
        // okay so the deal is like this
        // sometimes the MTA invents new route names
        // so what we need to do is
        // add that new bullshit route
        log.Println("Inserting unexpected new route.")
        err := t.Route.Insert(db)

        if err != nil {
          log.Fatal("Could not insert new route", err)
        }

        return t.Insert(db)
      }
    }

  }

  return err
}

func (t *Trip) Upsert(db *sql.DB) error {
  exists, err := t.Exists(db)

  if exists {
    return nil
  }

  if err != nil {
    return err
  }

  // if we are here, we know that the trip doesn't exist yet
  return t.Insert(db)
}

func (t *Trip) Exists(db *sql.DB) (bool, error) {
  var err error

  if existsStmt == nil {
    stmt := `SELECT EXISTS(SELECT 1 FROM trips WHERE id=$1)`
    existsStmt, err = db.Prepare(stmt)

    if err != nil {
      log.Fatal("exists", err)
      return false, err
    }
  }

  var exists bool
  err = existsStmt.QueryRow(t.Id).Scan(&exists)

  if err != nil {
    return false, err
  }

  return exists, nil
}

func ReadTrips(db *sql.DB) ([]*Trip, error) {
  var (
    trips []*Trip = []*Trip{}
    err   error
  )

  // prepare statement if not already done so.
  if tripReadAllStmt == nil {
    stmt := `SELECT
              trips.id AS id, 
              direction, 
              routes.id AS route_id, 
              COALESCE(short_name, '') AS short_name, 
              COALESCE(name, '') AS name,
              COALESCE(description, '') AS description,
              COALESCE(type, 0) AS type,
              COALESCE(url, '') AS url,
              COALESCE(color, '') AS color
             FROM trips
             LEFT OUTER JOIN routes ON trips.route = routes.id`

    tripReadAllStmt, err = db.Prepare(stmt)
    if err != nil {
      log.Fatal("Error preparing statement: ", err)
    }
  }

  rows, err := tripReadAllStmt.Query()
  if err != nil {
    return trips, err
  }

  defer rows.Close()
  for rows.Next() {
    trip := &Trip{Route: &Route{}}

    err = rows.Scan(
      &trip.Id,
      &trip.Direction,
      &trip.Route.Id,
      &trip.Route.ShortName,
      &trip.Route.Name,
      &trip.Route.Description,
      &trip.Route.Type,
      &trip.Route.URL,
      &trip.Route.Color)

    if err != nil {
      return trips, err
    }

    // append scanned route into list of all trips
    trips = append(trips, trip)
  }

  if err := rows.Err(); err != nil {
    return trips, err
  }

  return trips, nil
}
