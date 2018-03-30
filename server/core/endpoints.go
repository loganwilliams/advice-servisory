package core

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"

  "github.com/gorilla/mux"
  "github.com/loganwilliams/adviceservisory/server/pretty"
  "github.com/loganwilliams/adviceservisory/server/types"

  _ "github.com/lib/pq"
)

func (a *AdviceServisory) AllRoutesHandler(w http.ResponseWriter, r *http.Request) {
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

func (a *AdviceServisory) LiveUpdatesHandler(w http.ResponseWriter, r *http.Request) {
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
