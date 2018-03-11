package main

import (
	"net/http"
	"log"
	"fmt"
    "database/sql"
    "encoding/json"
    "github.com/loganwilliams/adviceservisory/server/types"
    "github.com/loganwilliams/adviceservisory/server/env"
    "github.com/loganwilliams/adviceservisory/server/pretty"
    _ "github.com/lib/pq"
)

var db *sql.DB
var config *env.Config

func main() {
    initDb()
    types.DropRoutesTable(db)
    types.CreateAndPopulateRoutesTable(db)

    http.HandleFunc("/routes", allRoutesHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))

    defer db.Close()    
}

func allRoutesHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("allRoutesHandler")
    routes, err := types.ReadRoutes(db)
    if err != nil {
        log.Panic("Error querying routes", err)
    }

    response, err := json.Marshal(routes)
    if err != nil {
        log.Panic("Error marshalling json", err)
    }

    fmt.Fprintf(w, "%s", pretty.Json(string(response)))
}    

func initDb() {
    config := env.NewConfig()
    var err error
    dbInfo := fmt.Sprintf("host=%s user=%s "+
        "password=%s dbname=%s sslmode=disable",
        config.DB.Host, config.DB.Username, config.DB.Password, config.DB.Name)

    db, err = sql.Open(env.DB_DRIVER, dbInfo)
    if err != nil {
        panic(err)
    }
    err = db.Ping()
    if err != nil {
        panic(err)
    }
    log.Println("Successfully connected!")
}
