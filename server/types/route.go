package types

import (
  "database/sql"
  "log"
  "github.com/loganwilliams/adviceservisory/server/gtfsstatic"
  "fmt"
)

type Route struct {
  Id          int             `json:"id"`
  Code        string          `json:"code"`
  ShortName   string          `json:"shortname"`
  Name        string          `json:"name"`
  Description string          `json:"description"`
  Type        int             `json:"type"`
  URL         string          `json:"url"`
  Color       string          `json:"color"`
}

// package level globals for storing prepared sql statements
var (
  routeReadAllStmt *sql.Stmt
)

func CreateAndPopulateRoutesTable(db *sql.DB) error {
  // Create import table
  mkImportTableStmt := `CREATE SCHEMA IF NOT EXISTS import;
    CREATE TABLE import.routes(
      route_id text,
      agency_id text,
      route_short_name text,
      route_long_name text,
      route_desc text,
      route_type text,
      route_url text,
      route_color text,
      route_text_color text
    )`

  _, err := db.Exec(mkImportTableStmt)
  if err != nil {
    log.Fatal("Error creating import.routes", err)
  }

  routesLocation, err := gtfsstatic.RoutesLocation()
  if err != nil {
    log.Fatal("Error getting routelocation ", err)
  }

  // Copy from CSV
  // TODO: change this so its not an absolute path
  copyFromCSVStmt := fmt.Sprintf(`COPY import.routes FROM 
    '%s' WITH DELIMITER ',' HEADER CSV`, routesLocation)

  _, err = db.Exec(copyFromCSVStmt)
  if err != nil {
    log.Fatal("Error copying routes from CSV")
  }

  // Create routes table
  mkTableStmt := `CREATE TABLE routes(
    id serial primary key,
    code varchar(5),
    short_name varchar(5),
    name varchar(200),
    description text,
    type int,
    url varchar(200),
    color varchar(6)
  )`

  _, err = db.Exec(mkTableStmt)
  if err != nil {
    log.Fatal("Error creating table `routes`", err)
  }

  // insert into routes table
  insertStmt := `INSERT INTO routes(
    code,
    short_name,
    name,
    description,
    type,
    url,
    color
  ) SELECT
    import.routes.route_id AS code,
    import.routes.route_short_name as short_name,
    import.routes.route_long_name AS name,
    import.routes.route_desc AS description,
    import.routes.route_type::int AS type,
    import.routes.route_url AS url,
    import.routes.route_color AS color
  FROM import.routes`

  _, err = db.Exec(insertStmt)
  if err != nil {
    log.Fatal("Error inserting routes", err)
  }

  return nil
}

func DropRoutesTable(db *sql.DB) error {
  dropStmt := `DROP TABLE IF EXISTS routes CASCADE;
    DROP TABLE IF EXISTS import.routes CASCADE;`
  _, err := db.Exec(dropStmt)
  if err != nil {
    log.Fatal("Error dropping routes table", err)
  }

  return nil
}

// func SetupRoutesTable(db *sql.DB) error {

// }

func ReadRoutes(db *sql.DB) ([]*Route, error) {
  var (
    routes []*Route = []*Route{}
    err error
  )

  // prepare statement if not already done so.
  if routeReadAllStmt == nil {
    stmt := `SELECT id, code, short_name, name, description, type, url, COALESCE(color,'') AS color
             FROM routes`
    routeReadAllStmt, err = db.Prepare(stmt)
    if err != nil {
      log.Fatal("Error preparing statement: ", err)
    }
  }

  rows, err := routeReadAllStmt.Query()
  if err != nil {
    return routes, err
  }

  defer rows.Close()
  for rows.Next() {
    route := &Route{}
    err = rows.Scan(
      &route.Id,
      &route.Code,
      &route.ShortName,
      &route.Name,
      &route.Description,
      &route.Type,
      &route.URL,
      &route.Color,
    )
    if err != nil {
      return routes, err
    }

    // append scanned route into list of all routes
    routes = append(routes, route)
  }
  if err := rows.Err(); err != nil {
    return routes, err
  }

  return routes, nil
}