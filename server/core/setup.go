package core

import (
  "github.com/loganwilliams/adviceservisory/server/types"
)

func (a *AdviceServisory) Setup() error {
  types.DropRoutesTable(a.DB)
  types.DropStopsTable(a.DB)
  types.DropTripsTable(a.DB)
  types.DropTripUpdatesTable(a.DB)
  types.CreateAndPopulateRoutesTable(a.DB)
  types.CreateAndPopulateStopsTable(a.DB)
  types.CreateTripsTable(a.DB)
  types.CreateTripUpdatesTable(a.DB)

  return nil
}