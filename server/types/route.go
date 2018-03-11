package types

import (
  "database/sql"
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
      log.Fatal("Error preparing statement", err)
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