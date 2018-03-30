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
    r.HandleFunc("/routes", app.AllRoutesHandler)
    r.HandleFunc("/route/{route_id}", app.RouteHandler)
    // r.HandleFunc("/route/{route_id}/{date}")
    // r.HandleFunc("/route/{route_id}/live")
    // r.HandleFunc("/route/{route_id}/live/geojson")

    r.HandleFunc("/trips", app.AllTripsHandler)
    r.HandleFunc("/trip/{trip_id}", app.TripUpdateHandler)
    // r.HandleFunc("/trip/{trip_id}/{date}")

    // r.HandleFunc("/history/{date}")
    // r.HandleFunc("/history/{date}/geojson")
    r.HandleFunc("/live", app.LiveUpdatesHandler)
    r.HandleFunc("/live/geojson", app.LiveUpdatesHandlerGJ)

    r.HandleFunc("/stops", app.AllStopsHandler)
    // r.HandleFunc("/stop/{stop_id}", app.StopHandler)
    r.HandleFunc("/station/{station_id}", app.StationHandler)

    http.Handle("/", r)

    log.Fatal(http.ListenAndServe(":8080", nil))

    defer app.DB.Close()
}
