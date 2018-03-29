package core

import (
  "database/sql"
  "net/http"
  "log"
  "encoding/json"
  "fmt"
  "time"

  "github.com/loganwilliams/adviceservisory/server/env"
  "github.com/loganwilliams/adviceservisory/server/types"
  "github.com/loganwilliams/adviceservisory/server/pretty"
  "github.com/gorilla/mux"

  _ "github.com/lib/pq"
)

type AdviceServisory struct {
  DB *sql.DB
  Config *env.Config
}

func NewAdviceServisory() *AdviceServisory {
  a := &AdviceServisory{}

  a.InitDb()

  return a
}

func (a *AdviceServisory) Start(ticker *time.Ticker) {
  a.AddTripUpdates()

  // loop infinitely with a ticker

  for t := range ticker.C {
    fmt.Println("Updating trips", t)
    a.AddTripUpdates()
  } 

}

func (a *AdviceServisory) AllRoutesHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("allRoutesHandler")
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

func (a *AdviceServisory) AllTripUpdatesHandler(w http.ResponseWriter, r *http.Request) {
  updates, err := types.ReadAllUpdates(a.DB)

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


func (a *AdviceServisory) InitDb() {
    a.Config = env.NewConfig()
    var err error
    dbInfo := fmt.Sprintf("host=%s user=%s "+
        "password=%s dbname=%s sslmode=disable",
        a.Config.DB.Host, a.Config.DB.Username, a.Config.DB.Password, a.Config.DB.Name)

    a.DB, err = sql.Open(env.DB_DRIVER, dbInfo)
    if err != nil {
        panic(err)
    }
    err = a.DB.Ping()
    if err != nil {
        panic(err)
    }
    log.Println("Successfully connected!")
}