package core

import (
  "database/sql"
  "fmt"
  "log"
  "time"

  "github.com/loganwilliams/adviceservisory/server/env"

  _ "github.com/lib/pq"
)

type AdviceServisory struct {
  DB     *sql.DB
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
