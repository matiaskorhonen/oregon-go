#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <time.h>
#include <iostream>
#include <string>
#include "RCSwitch.h"
#include "RcOok.h"
#include "Sensor.h"

void test() {
  int RXPIN = 1;
  int TXPIN = 0;

  if(wiringPiSetup() == -1)
    return;

  RCSwitch *rc = new RCSwitch(RXPIN,TXPIN);

  while (1)
  {
    if (rc->OokAvailable())
    {
      char message[100];

      rc->getOokCode(message);
      printf("%s\n",message);

      Sensor *s = Sensor::getRightSensor(message);
      if (s!= NULL)
      {
        std::cout << "Name : " << s->getSensorName() << "\n";
        printf("Temp : %f\n",s->getTemperature());
        printf("Humidity : %f\n",s->getHumidity());
        printf("Channel : %d\n",s->getChannel());
        printf("\n");
      }
      delete s;
    }
    delay(500);
  }
}
