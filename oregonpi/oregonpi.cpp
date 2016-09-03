#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <time.h>
#include <iostream>
#include <string>
#include "RCSwitch.h"
#include "RcOok.h"
#include "Sensor.h"
#include "oregonpi.h"

extern "C" {
  void * rc_switch_create(int RXPIN, int TXPIN) {
    if(wiringPiSetup() == -1)
      return NULL;

    return new RCSwitch(RXPIN,TXPIN);
  }

  struct SensorReading rc_read_from_sensor(void *rc_switch) {
    RCSwitch *rc = (RCSwitch *)rc_switch;

    while (1) {
      if (rc->OokAvailable()) {
        char message[100];

        rc->getOokCode(message);

        Sensor *s = Sensor::getRightSensor(message);
        if (s != NULL) {
          SensorReading sr;
          sr.temperature = s->getTemperature();
          sr.humidity = s->getHumidity();
          sr.name = strdup(s->getSensorName().c_str());
          sr.sensor_type = s->getSensType();
          sr.channel = s->getChannel();
          sr.low_battery = s->isBatteryLow() ? 1 : 0;
          delete s;
          return sr;
        }
        delete s;
      }
      delay(500);
    }
  }

  void rc_release(void *rc_switch) {
    delete (RCSwitch *)rc_switch;
  }
}
