package main

import (
  "log"
  "net/http"
  "time"

  "github.com/gorilla/mux"
  "github.com/loganwilliams/adviceservisory/server/core"
)

func main() {
  app := core.NewAdviceServisory()
  app.Setup()

  // begin background polling ETL processes
  ticker := time.NewTicker(30 * time.Second)
  go app.Start(ticker)

  r := mux.NewRouter()
  r.HandleFunc("/api/routes", app.AllRoutesHandler)
  r.HandleFunc("/api/route/{route_id}", app.RouteHandler)
  // r.HandleFunc("/route/{route_id}/{date}")
  // r.HandleFunc("/route/{route_id}/live")
  // r.HandleFunc("/route/{route_id}/live/geojson")

  r.HandleFunc("/api/trips", app.AllTripsHandler)
  r.HandleFunc("/api/trip/{trip_id}", app.TripUpdateHandler)
  // r.HandleFunc("/trip/{trip_id}/{date}")

  // r.HandleFunc("/history/{date}")
  // r.HandleFunc("/history/{date}/geojson")
  r.HandleFunc("/api/live", app.LiveUpdatesHandler)
  r.HandleFunc("/api/live/geojson", app.LiveUpdatesHandlerGJ)

  r.HandleFunc("/api/stops", app.AllStopsHandler)
  // r.HandleFunc("/stop/{stop_id}", app.StopHandler)
  r.HandleFunc("/api/station/{station_id}", app.StationHandler)
  // r.HandleFunc("/station/{station_id}/today")
  // r.HandleFunc("/station/{station_id}/{date}")

  http.Handle("/", r)

  log.Fatal(http.ListenAndServe(":8080", nil))

  defer app.DB.Close()
}
