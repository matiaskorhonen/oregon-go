package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	"github.com/BurntSushi/toml"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iotdataplane"
	"github.com/matiaskorhonen/oregon-go/oregonpi"
)

// ConfiguredSensor is used to configure which sensors are being listened for
type ConfiguredSensor struct {
	SensorType      int    `toml:"type"`
	Channel         int    `toml:"channel"`
	ThingEndpoint   string `toml:"thing_endpoint"`
	SkipHumidity    bool   `toml:"skip_humidity"`
	SkipTemperature bool   `toml:"skip_temperature"`
}

// Config is used to confifigure the GPIO pins and what sensors to listen for
type Config struct {
	RXPin   int                `toml:"rx_pin"`
	TXPin   int                `toml:"tx_pin"`
	Sensors []ConfiguredSensor `toml:"sensors"`
}

type thingState struct {
	SensorName       string  `json:"sensorName"`
	SensorType       int     `json:"sensorType"`
	SensorChannel    int     `json:"sensorChannel"`
	SensorLowBattery bool    `json:"sensorLowBattery"`
	Temperature      float32 `json:"temperature"`
	Humidity         float32 `json:"humidity"`
}

var config Config

func init() {
	var help bool
	var configPath string

	flag.StringVar(&configPath, "config", "", "path to the config file")
	flag.BoolVar(&help, "help", false, "this help mesage")
	flag.Parse()

	if configPath == "" || help {
		flag.PrintDefaults()
		os.Exit(1)
	}

	tomlData, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	var conf Config
	_, err = toml.Decode(string(tomlData), &conf)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func main() {
	log.Println("Starting...")
	monitor, err := oregonpi.NewSensorMonitor(config.RXPin, config.TXPin)

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
		updateThingShadow(reading)
	}
}

func updateThingShadow(reading *oregonpi.SensorReading) {
	sess, err := session.NewSession()
	if err != nil {
		log.Println("Failed to create AWS session: ", err)
		return
	}

	svc := iotdataplane.New(sess, aws.NewConfig().WithEndpoint(os.Getenv("AWS_IOT_ENDPOINT")))

	reportedState := thingState{
		SensorName:       reading.Sensor.Name,
		SensorType:       reading.Sensor.Type,
		SensorChannel:    reading.Sensor.Channel,
		SensorLowBattery: reading.Sensor.LowBattery,
		Temperature:      reading.Temperature,
		Humidity:         reading.Humidity,
	}

	payload, err := json.Marshal(map[string]map[string]thingState{
		"state": map[string]thingState{
			"reported": reportedState,
		},
	})
	if err != nil {
		log.Println("Serialization error: ", err)
		return
	}

	log.Println("Updating Thing Shadow…")
	params := &iotdataplane.UpdateThingShadowInput{
		Payload:   payload,
		ThingName: aws.String("OutsideWeatherSensor"),
	}
	_, err = svc.UpdateThingShadow(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		log.Println(err.Error())
		return
	}

	log.Println("Thing Shadow updated")
}
