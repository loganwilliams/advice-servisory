package core

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "time"

  "github.com/gorilla/mux"
  "github.com/loganwilliams/adviceservisory/server/pretty"
  "github.com/loganwilliams/adviceservisory/server/types"

  _ "github.com/lib/pq"
)

// TODO refactor this to reduce boilerplate
// TODO redo setHeaders as middleware

func (a *AdviceServisory) AllRoutesHandler(w http.ResponseWriter, r *http.Request) {
  w = setHeaders(w)
  routes, err := types.ReadRoutes(a.DB)

  if err != nil {
    log.Panic("Error querying routes", err)
  }

  response, err := json.Marshal(routes)
  if err != nil {
    log.Panic("Error marshalling json", err)
  }

  fmt.Fprintf(w, "%s", pretty.Json(string(response)))
}

func (a *AdviceServisory) RouteHandler(w http.ResponseWriter, r *http.Request) {
  w = setHeaders(w)
  vars := mux.Vars(r)
  route := &types.Route{Id: vars["route_id"]}
  updates, err := route.ReadUpdates(a.DB)

  if err != nil {
    log.Panic("Error querying updates", err)
  }

  response, err := json.Marshal(updates)
  if err != nil {
    log.Panic("Error marshalling json", err)
  }

  fmt.Fprintf(w, "%s", pretty.Json(string(response)))
}

func (a *AdviceServisory) AllTripsHandler(w http.ResponseWriter, r *http.Request) {
  w = setHeaders(w)
  trips, err := types.ReadTrips(a.DB)

  if err != nil {
    log.Panic("Error querying trips", err)
  }

  response, err := json.Marshal(trips)

  if err != nil {
    log.Panic("Error marshalling json", err)
  }

  fmt.Fprintf(w, "%s", pretty.Json(string(response)))
}

func (a *AdviceServisory) TripUpdateHandler(w http.ResponseWriter, r *http.Request) {
  w = setHeaders(w)
  vars := mux.Vars(r)
  trip := &types.Trip{Id: vars["trip_id"]}
  updates, err := trip.ReadUpdates(a.DB)

  if err != nil {
    log.Panic("Error querying update", err)
  }

  response, err := json.Marshal(updates)

  if err != nil {
    log.Panic("Error marshalling json", err)
  }

  fmt.Fprintf(w, "%s", pretty.Json(string(response)))
}

func (a *AdviceServisory) AllStopsHandler(w http.ResponseWriter, r *http.Request) {
  w = setHeaders(w)
  stops, err := types.ReadAllStops(a.DB)

  if err != nil {
    log.Panic("Error querying update", err)
  }

  response, err := json.Marshal(stops)

  if err != nil {
    log.Panic("Error marshalling json", err)
  }

  fmt.Fprintf(w, "%s", pretty.Json(string(response)))
}

func (a *AdviceServisory) StationHandler(w http.ResponseWriter, r *http.Request) {
  w = setHeaders(w)
  vars := mux.Vars(r)
  stop := &types.Stop{Station: vars["station_id"]}
  updates, err := stop.ReadUpdates(a.DB)

  if err != nil {
    log.Panic("Error querying update", err)
  }

  response, err := json.Marshal(updates)

  if err != nil {
    log.Panic("Error marshalling json", err)
  }

  fmt.Fprintf(w, "%s", pretty.Json(string(response)))
}

func (a *AdviceServisory) StationHandlerDate(w http.ResponseWriter, r *http.Request) {
  w = setHeaders(w)
  vars := mux.Vars(r)
  stop := &types.Stop{Station: vars["station_id"]}
  date, err := time.Parse("2006-01-02", vars["date"])
  updates, err := stop.ReadUpdatesAtDate(a.DB, date)

  if err != nil {
    log.Panic("Error querying update", err)
  }

  response, err := json.Marshal(updates)

  if err != nil {
    log.Panic("Error marshalling json", err)
  }

  fmt.Fprintf(w, "%s", pretty.Json(string(response)))
}

func (a *AdviceServisory) TwoStationHandler(w http.ResponseWriter, r *http.Request) {
  w = setHeaders(w)
  vars := mux.Vars(r)
  stop1 := &types.Stop{Station: vars["station_id1"]}
  stop2 := &types.Stop{Station: vars["station_id2"]}
  updates, err := stop1.UpdatesWithStop(a.DB, stop2)

  if err != nil {
    log.Panic("Error querying update", err)
  }

  response, err := json.Marshal(updates)

  if err != nil {
    log.Panic("Error marshalling json", err)
  }

  fmt.Fprintf(w, "%s", pretty.Json(string(response)))
}

func (a *AdviceServisory) StopHandler(w http.ResponseWriter, r *http.Request) {
  w = setHeaders(w)
  vars := mux.Vars(r)
  stop := &types.Stop{Id: vars["stop_id"]}
  updates, err := stop.ReadUpdates(a.DB)

  if err != nil {
    log.Panic("Error querying update", err)
  }

  response, err := json.Marshal(updates)

  if err != nil {
    log.Panic("Error marshalling json", err)
  }

  fmt.Fprintf(w, "%s", pretty.Json(string(response)))
}

func (a *AdviceServisory) StopHandlerDate(w http.ResponseWriter, r *http.Request) {
  w = setHeaders(w)
  vars := mux.Vars(r)
  stop := &types.Stop{Id: vars["stop_id"]}
  date, err := time.Parse("2006-01-02", vars["date"])

  updates, err := stop.ReadUpdatesAtDate(a.DB, date)

  if err != nil {
    log.Panic("Error querying update", err)
  }

  response, err := json.Marshal(updates)

  if err != nil {
    log.Panic("Error marshalling json", err)
  }

  fmt.Fprintf(w, "%s", pretty.Json(string(response)))
}

func (a *AdviceServisory) LiveUpdatesHandler(w http.ResponseWriter, r *http.Request) {
  w = setHeaders(w)
  updates, err := types.LiveUpdates(a.DB)

  if err != nil {
    log.Panic("Error querying update", err)
  }

  response, err := json.Marshal(updates)

  if err != nil {
    log.Panic("Error marshalling json", err)
  }

  fmt.Fprintf(w, "%s", pretty.Json(string(response)))
}

func (a *AdviceServisory) LiveUpdatesHandlerGJ(w http.ResponseWriter, r *http.Request) {
  w = setHeaders(w)
  updates, err := types.LiveUpdates(a.DB)
  geojson := types.MakeGeoJSON(updates)

  if err != nil {
    log.Panic("Error querying update", err)
  }

  response, err := json.Marshal(geojson)

  if err != nil {
    log.Panic("Error marshalling json", err)
  }

  fmt.Fprintf(w, "%s", pretty.Json(string(response)))
}

func setHeaders(w http.ResponseWriter) http.ResponseWriter {
  // Send the correct headers to enable CORS
  w.Header().Set("Content-Type", "text/json; charset=ascii")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")

  return w
}
