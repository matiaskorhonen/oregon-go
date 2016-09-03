# oregon-go

Read data from Oregon Scientific wireless (433MHz) sensors and push the readings to AWS IoT. Still a work in progress.

## Dependencies

### Hardware

* A Raspberry PI (test with a [Raspberry Pi 2 Model B][])
* A 433 MHz receiver (tested with a [Quasar DSQAM-RX3-1][])
  * INPUT GPIO 1 (See [wiringPi pins][])
* A compatible Oregon Scientific wireless sensor (See `Sensor.cpp`)

### Software

* Tested on [Raspbian Jessie][Raspbian]
* [wiringPi][]
  * See the [wiringPi instructions][] for more info.

    ```sh
    sudo apt-get update
    sudo apt-get upgrade
    sudo apt-get install git-core
    git clone git://git.drogon.net/wiringPi
    cd wiringPi
    ./build
    ```
  * wiringPi is released under the GNU Lesser Public License version 3.

## License

Licensed under GPLv3 due to the required dependencies. See the LICENSE file for details.

Based on the Disk19.com/Paul Pinault [rfrpi][] project (GPLv3), [modified by Emilio Peña][OregonPi].

Includes code from Suat Özgür's [RCSwitch][] project (GNU Lesser General Public License 2.1).



[Raspberry Pi 2 Model B]: https://www.raspberrypi.org/products/raspberry-pi-2-model-b/
[Quasar DSQAM-RX3-1]: http://www.quasaruk.co.uk/acatalog/info_QAM_RX3_433.html
[wiringPi]: http://wiringpi.com/
[wiringPi Pins]: https://projects.drogon.net/raspberry-pi/wiringpi/pins/
[wiringPi instructions]: http://wiringpi.com/download-and-install/
[rfrpi]: https://bitbucket.org/disk_91-admin/rfrpi
[RCSwitch]: https://github.com/sui77/rc-switch
[OregonPi]: https://github.com/1000io/OregonPi
[Raspbian]: https://www.raspberrypi.org/downloads/raspbian/
