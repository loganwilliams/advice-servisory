package main

import (
	"net/http"
	"log"
    "time"

    "github.com/loganwilliams/adviceservisory/server/core"
    "github.com/gorilla/mux"
)

func main() {

    app := core.NewAdviceServisory()
    app.Setup()

    // begin background polling ETL processes
    ticker := time.NewTicker(30 * time.Second)
    go app.Start(ticker)

    r := mux.NewRouter()
    r.HandleFunc("/routes", app.AllRoutesHandler)

    r.HandleFunc("/trips", app.AllTripsHandler)
    r.HandleFunc("/trip/{trip_id}", app.TripUpdateHandler)

    r.HandleFunc("/updates", app.AllTripUpdatesHandler)
    r.HandleFunc("/live", app.LiveUpdatesHandler)
    r.HandleFunc("/live/geojson", app.LiveGeojsonHandler)

    r.HandleFunc("/stops", app.AllStopsHandler)
    r.HandleFunc("/station/{station_id}", app.StationHandler)

    //TODO add live endpoint
    http.Handle("/", r)
    
    log.Fatal(http.ListenAndServe(":8080", nil))

    defer app.DB.Close()    
}
