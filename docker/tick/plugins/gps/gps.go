package main

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)

	location := map[string]float64{
		"lat": 45.260860,
		"lng": 19.820638,
	}
	for {
		<-c

		if err := json.NewEncoder(os.Stdout).Encode(location); err != nil {
			panic(err)
		}
		location["lat"] += 0.1
		location["lng"] += 0.2
	}
}
