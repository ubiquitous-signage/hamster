package main

import (
	"time"

	"github.com/ubiquitous-signage/hamster/lectures"
	"github.com/ubiquitous-signage/hamster/newsletters"
	"github.com/ubiquitous-signage/hamster/portal"
	"github.com/ubiquitous-signage/hamster/schedules"
	"github.com/ubiquitous-signage/hamster/server"
	"github.com/ubiquitous-signage/hamster/trains"
	// "github.com/ubiquitous-signage/hamster/weather"
)

func main() {
	go server.Run()
	go lectures.Run()
	go newsletters.Run()
	go portal.Run()
	go schedules.Run()
	go trains.Run()
	for {
		time.Sleep(1 * time.Second)
	}
}
