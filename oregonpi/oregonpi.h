#ifndef OREGONPI_H
#define OREGONPI_H

#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

struct SensorReading {
  double temperature;
  double humidity;
  char* name;
  int sensor_type;
  int low_battery;
  int channel;
};

void* rc_switch_create(int RXPIN, int TXPIN);
struct SensorReading rc_read_from_sensor(void* rc_switch);
void rc_release(void* rc_switch);

#ifdef __cplusplus
}
#endif
#endif
