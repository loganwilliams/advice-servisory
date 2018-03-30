package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/loganwilliams/where-are-the-trains/server/transit_realtime"
)

func main() {
	f, _ := os.Open("gtfs_bad")
	gtfs, _ := ioutil.ReadAll(f)
	transit := &transit_realtime.FeedMessage{}

	if err := proto.Unmarshal(gtfs, transit); err != nil {
		log.Println("Failed to parse GTFS feed", err)
	}

	for _, entity := range transit.Entity {
		fmt.Println(entity)
	}

}
