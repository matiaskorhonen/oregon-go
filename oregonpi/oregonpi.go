package oregonpi

// #cgo CFLAGS: -O3
// #cgo CXXFLAGS: -O3
// #cgo LDFLAGS: -lwiringPi
// #include "oregonpi.h"
// #include <stdlib.h>
import "C"
import (
	"errors"
	"sync"
	"unsafe"
)

type Sensor struct {
	Name       string
	Type       int
	Channel    int
	LowBattery bool
}

type SensorReading struct {
	Temperature float32
	Humidity    float32
	Sensor
}

type SensorMonitor struct {
	rcSwitch  unsafe.Pointer
	active    bool
	mu        sync.Mutex // protects active
	terminate chan chan bool
}

func NewSensorMonitor(RXPIN, TXPIN int) (*SensorMonitor, error) {
	var rcSwitch unsafe.Pointer
	rcSwitch = C.rc_switch_create(C.int(RXPIN), C.int(TXPIN))
	if rcSwitch == nil {
		return nil, errors.New("oregonpi: wiringPiSetup failed")
	}
	return &SensorMonitor{
		rcSwitch:  rcSwitch,
		terminate: make(chan chan bool),
	}, nil
}

func (sm *SensorMonitor) ReadFromSensor(ch chan<- *SensorReading) {
	sm.mu.Lock()
	sm.active = true
	sm.mu.Unlock()

	go func() {
		for {
			select {
			case terminated := <-sm.terminate:
				C.rc_release(sm.rcSwitch)
				terminated <- true
				return
			default:
				ch <- sm.getReadingFromSensor()
			}
		}
	}()
}

func (sm *SensorMonitor) Stop() {
	sm.mu.Lock()
	active := sm.active
	sm.mu.Unlock()

	if active {
		terminated := make(chan bool)
		sm.terminate <- terminated
		<-terminated
	} else {
		C.rc_release(sm.rcSwitch)
	}
}

func (sm *SensorMonitor) getReadingFromSensor() *SensorReading {
	r := C.rc_read_from_sensor(sm.rcSwitch)
	reading := &SensorReading{
		Temperature: float32(r.temperature),
		Humidity:    float32(r.humidity),
		Sensor: Sensor{
			Name:       C.GoString(r.name),
			Type:       int(r.sensor_type),
			Channel:    int(r.channel),
			LowBattery: int(r.low_battery) != 0,
		},
	}
	C.free(unsafe.Pointer(r.name))
	return reading
}
