package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/matiaskorhonen/oregon-go/oregonpi"
)

func main() {
	log.Println("Starting...")
	monitor, err := oregonpi.NewSensorMonitor(1, 0)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	readings := make(chan *oregonpi.SensorReading)
	monitor.ReadFromSensor(readings)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Println("Stopping...")
			monitor.Stop()
			close(readings)
		}
	}()

	for reading := range readings {
		log.Println(reading)
	}
}
