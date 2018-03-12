package types

import (
  "database/sql"
  "log"
  "github.com/loganwilliams/adviceservisory/server/gtfsstatic"
  "fmt"
)

type Stop struct {
	Id		int	`json:"id"`
	Code	int	`json:"code"`
	Name	string `json:"name"`
	Latitude float `json:"latitude"`
	Longitude float `json:"longitude"`
	Type int `json:"type"`
	URL string `json:"url"`
}

