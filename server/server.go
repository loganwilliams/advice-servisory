package main

import (
	"net/http"
	"log"
    "github.com/loganwilliams/adviceservisory/server/core"
    "time"
)

func main() {

    app := core.NewAdviceServisory()
    app.Setup()

    http.HandleFunc("/routes", app.AllRoutesHandler)
    http.HandleFunc("/trips", app.AllTripsHandler)

    // begin background polling ETL processes
    ticker := time.NewTicker(30 * time.Second)
    go app.Start(ticker)

    log.Fatal(http.ListenAndServe(":8080", nil))

    defer app.DB.Close()    
}


