# ChainMetric: Sensor System

[![golang badge]][golang]&nbsp;
[![lines counter]][this repo]&nbsp;
[![commit activity badge]][repo commit activity]&nbsp;
[![hardware badge]][raspberry pi]&nbsp;
[![license badge]][license url]

## Overview

_**Chainmetric Sensor System**_, being an embedded IoT sensors-equipped device, designed to be compatible with a permissioned blockchain network based on Hyperledger Fabric stack. 

By leveraging highly concurrent engine driver such implementation is ideal for harvesting environmental surrounding conditions and further publishing them onto the distributed, immutable ledger, where such data can be validated by on-chain Smart Contacts against previously assigned requirements. 

The device itself is intended for deployment in the areas where assets requiring monitoring are stored or being delivered on. Thus, providing a general-purpose supply-chain monitoring solution.

## Supports

[![max44009 badge]][max44009]
[![si1145 badge]][si1145]

[![hdc1080 badge]][si1145]
[![dht22 badge]][dht22]

[![ccs811 badge]][ccs811]

[![bmp280 badge]][bmp280]
[![adxl345 badge]][adxl345]

[![max33102 badge]][max33102]

## Requirements
- Raspberry Pi 3/4/Zero

## Deployment

## Wrap up

## License

Licensed under the [Apache 2.0][license file].


[golang badge]: https://img.shields.io/badge/Code-Golang-informational?style=flat&logo=go&logoColor=white&color=6AD7E5
[lines counter]: https://img.shields.io/tokei/lines/github/timoth-y/chainmetric-sensorsys?color=teal&label=Lines
[commit activity badge]: https://img.shields.io/github/commit-activity/m/timoth-y/chainmetric-sensorsys?label=Commit%20activity&color=teal
[hardware badge]: https://img.shields.io/badge/Hardware-Raspberry%20Pi-informational?style=flat&logo=Raspberry%20Pi&color=953347
[license badge]: https://img.shields.io/badge/License-Apache%202.0-informational?style=flat&color=blue

[this repo]: https://github.com/timoth-y/chainmetric-sensorsys
[golang]: https://golang.org
[repo commit activity]: https://github.com/timoth-y/kicksware-api/graphs/commit-activity
[raspberry pi]: https://www.raspberrypi.org
[license url]: https://www.apache.org/licenses/LICENSE-2.0

[max44009 badge]: https://img.shields.io/badge/Luminosity-MAX44009-informational?style=for-the-badge&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGVuYWJsZS1iYWNrZ3JvdW5kPSJuZXcgMCAwIDI0IDI0IiBoZWlnaHQ9IjI0IiB2aWV3Qm94PSIwIDAgMjQgMjQiIHdpZHRoPSIyNCI+PGc+PHBhdGggZD0iTTAsMGgyNHYyNEgwVjB6IiBmaWxsPSJub25lIi8+PC9nPjxnPjxnPjxyZWN0IGhlaWdodD0iMyIgd2lkdGg9IjIiIHg9IjExIiB5PSIxOSIvPjxyZWN0IGhlaWdodD0iMiIgd2lkdGg9IjMiIHg9IjIiIHk9IjExIi8+PHJlY3QgaGVpZ2h0PSIyIiB3aWR0aD0iMyIgeD0iMTkiIHk9IjExIi8+PHJlY3QgaGVpZ2h0PSIzIiB0cmFuc2Zvcm09Im1hdHJpeCgwLjcwNzEgLTAuNzA3MSAwLjcwNzEgMC43MDcxIC03LjY2NjUgMTcuODAxNCkiIHdpZHRoPSIxLjk5IiB4PSIxNi42NiIgeT0iMTYuNjYiLz48cmVjdCBoZWlnaHQ9IjEuOTkiIHRyYW5zZm9ybT0ibWF0cml4KDAuNzA3MSAtMC43MDcxIDAuNzA3MSAwLjcwNzEgLTEwLjk3OTEgOS44MDQxKSIgd2lkdGg9IjMiIHg9IjQuODUiIHk9IjE3LjE2Ii8+PHBhdGggZD0iTTE1LDguMDJWM0g5djUuMDJDNy43OSw4Ljk0LDcsMTAuMzcsNywxMmMwLDIuNzYsMi4yNCw1LDUsNXM1LTIuMjQsNS01QzE3LDEwLjM3LDE2LjIxLDguOTQsMTUsOC4wMnogTTExLDVoMnYyLjEgQzEyLjY4LDcuMDQsMTIuMzQsNywxMiw3cy0wLjY4LDAuMDQtMSwwLjFWNXoiLz48L2c+PC9nPjwvc3ZnPg==&labelColor=FBE967&color=434343
[max44009]: https://datasheets.maximintegrated.com/en/ds/MAX44009.pdf

[si1145 badge]: https://img.shields.io/badge/Ambient%20light%20(UV%20index)-SI1145-informational?style=for-the-badge&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0iYmxhY2siIHdpZHRoPSIxNHB4IiBoZWlnaHQ9IjE0cHgiPjxwYXRoIGQ9Ik0wIDBoMjR2MjRIMHoiIGZpbGw9Im5vbmUiLz48cGF0aCBkPSJNMjAgMTUuMzFMMjMuMzEgMTIgMjAgOC42OVY0aC00LjY5TDEyIC42OSA4LjY5IDRINHY0LjY5TC42OSAxMiA0IDE1LjMxVjIwaDQuNjlMMTIgMjMuMzEgMTUuMzEgMjBIMjB2LTQuNjl6TTEyIDE4Yy0zLjMxIDAtNi0yLjY5LTYtNnMyLjY5LTYgNi02IDYgMi42OSA2IDYtMi42OSA2LTYgNnoiLz48L3N2Zz4=&labelColor=FBE967&color=434343
[si1145]: https://cdn-shop.adafruit.com/datasheets/Si1145-46-47.pdf

[hdc1080 badge]: https://img.shields.io/badge/Temperature%20&%20Humidity-HDC1080-informational?style=for-the-badge&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGNvbG9yPSJ3aGl0ZSIgaGVpZ2h0PSIxNCIgdmlld0JveD0iMCAwIDI0IDI0IiB3aWR0aD0iMTQiPjxwYXRoIGQ9Ik0wIDBoMjR2MjRIMHoiIGZpbGw9Im5vbmUiLz48cGF0aCBkPSJNMTUgMTNWNWMwLTEuNjYtMS4zNC0zLTMtM1M5IDMuMzQgOSA1djhjLTEuMjEuOTEtMiAyLjM3LTIgNCAwIDIuNzYgMi4yNCA1IDUgNXM1LTIuMjQgNS01YzAtMS42My0uNzktMy4wOS0yLTR6bS00LThjMC0uNTUuNDUtMSAxLTFzMSAuNDUgMSAxaC0xdjFoMXYyaC0xdjFoMXYyaC0yVjV6Ii8+PC9zdmc+&labelColor=72F5F5&color=434343
[hdc1080]: https://www.ti.com/lit/ds/symlink/hdc1080.pdf

[dht22 badge]: https://img.shields.io/badge/Temperature%20&%20Humidity-DHT11\22-informational?style=for-the-badge&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGNvbG9yPSJ3aGl0ZSIgaGVpZ2h0PSIxNCIgdmlld0JveD0iMCAwIDI0IDI0IiB3aWR0aD0iMTQiPjxwYXRoIGQ9Ik0wIDBoMjR2MjRIMHoiIGZpbGw9Im5vbmUiLz48cGF0aCBkPSJNMTUgMTNWNWMwLTEuNjYtMS4zNC0zLTMtM1M5IDMuMzQgOSA1djhjLTEuMjEuOTEtMiAyLjM3LTIgNCAwIDIuNzYgMi4yNCA1IDUgNXM1LTIuMjQgNS01YzAtMS42My0uNzktMy4wOS0yLTR6bS00LThjMC0uNTUuNDUtMSAxLTFzMSAuNDUgMSAxaC0xdjFoMXYyaC0xdjFoMXYyaC0yVjV6Ii8+PC9zdmc+&labelColor=72F5F5&color=434343
[dht22]: https://www.sparkfun.com/datasheets/Sensors/Temperature/DHT22.pdf

[ccs811 badge]: https://img.shields.io/badge/Air%20Quality%20(CO2,TVOC)-CCS811-informational?style=for-the-badge&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGVuYWJsZS1iYWNrZ3JvdW5kPSJuZXcgMCAwIDI0IDI0IiBoZWlnaHQ9IjI0IiB2aWV3Qm94PSIwIDAgMjQgMjQiIHdpZHRoPSIyNCI+PGc+PHBhdGggZD0iTTAsMGgyNHYyNEgwVjB6IiBmaWxsPSJub25lIi8+PC9nPjxnPjxnPjxwYXRoIGQ9Ik0xNC41LDE3YzAsMS42NS0xLjM1LDMtMywzcy0zLTEuMzUtMy0zaDJjMCwwLjU1LDAuNDUsMSwxLDFzMS0wLjQ1LDEtMXMtMC40NS0xLTEtMUgydi0yaDkuNSBDMTMuMTUsMTQsMTQuNSwxNS4zNSwxNC41LDE3eiBNMTksNi41QzE5LDQuNTcsMTcuNDMsMywxNS41LDNTMTIsNC41NywxMiw2LjVoMkMxNCw1LjY3LDE0LjY3LDUsMTUuNSw1UzE3LDUuNjcsMTcsNi41IFMxNi4zMyw4LDE1LjUsOEgydjJoMTMuNUMxNy40MywxMCwxOSw4LjQzLDE5LDYuNXogTTE4LjUsMTFIMnYyaDE2LjVjMC44MywwLDEuNSwwLjY3LDEuNSwxLjVTMTkuMzMsMTYsMTguNSwxNnYyIGMxLjkzLDAsMy41LTEuNTcsMy41LTMuNVMyMC40MywxMSwxOC41LDExeiIvPjwvZz48L2c+PC9zdmc+&labelColor=74FA4C&color=434343
[ccs811]: https://cdn-learn.adafruit.com/downloads/pdf/adafruit-ccs811-air-quality-sensor.pdf

[bmp280 badge]: https://img.shields.io/badge/Barometer-BMP280-informational?style=for-the-badge&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGhlaWdodD0iMjQiIHZpZXdCb3g9IjAgMCAyNCAyNCIgd2lkdGg9IjI0Ij48cGF0aCBkPSJNMCAwaDI0djI0SDB6IiBmaWxsPSJub25lIi8+PHBhdGggZD0iTTAgMGgyNHYyNEgwVjB6IiBmaWxsPSJub25lIi8+PHBhdGggZD0iTTggMTloM3YzaDJ2LTNoM2wtNC00LTQgNHptOC0xNWgtM1YxaC0ydjNIOGw0IDQgNC00ek00IDl2MmgxNlY5SDR6Ii8+PHBhdGggZD0iTTQgMTJoMTZ2Mkg0eiIvPjwvc3ZnPg==&labelColor=72F5F5&color=434343
[bmp280]: https://cdn-shop.adafruit.com/datasheets/BST-BMP280-DS001-11.pdf

[adxl345 badge]: https://img.shields.io/badge/Accelerometer-ADXL345-informational?style=for-the-badge&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGhlaWdodD0iMjQiIHZpZXdCb3g9IjAgMCAyNCAyNCIgd2lkdGg9IjI0Ij48cGF0aCBkPSJNMCAwaDI0djI0SDB6IiBmaWxsPSJub25lIi8+PHBhdGggZD0iTTIwLjM4IDguNTdsLTEuMjMgMS44NWE4IDggMCAwIDEtLjIyIDcuNThINS4wN0E4IDggMCAwIDEgMTUuNTggNi44NWwxLjg1LTEuMjNBMTAgMTAgMCAwIDAgMy4zNSAxOWEyIDIgMCAwIDAgMS43MiAxaDEzLjg1YTIgMiAwIDAgMCAxLjc0LTEgMTAgMTAgMCAwIDAtLjI3LTEwLjQ0em0tOS43OSA2Ljg0YTIgMiAwIDAgMCAyLjgzIDBsNS42Ni04LjQ5LTguNDkgNS42NmEyIDIgMCAwIDAgMCAyLjgzeiIvPjwvc3ZnPg==&labelColor=72F5F5&color=434343
[adxl345]: https://www.sparkfun.com/datasheets/Sensors/Accelerometer/ADXL345.pdf

[max33102 badge]: https://img.shields.io/badge/Pulse%20Oximeter-MAX33102-informational?style=for-the-badge&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGhlaWdodD0iMjQiIHZpZXdCb3g9IjAgMCAyNCAyNCIgd2lkdGg9IjI0Ij48cGF0aCBkPSJNMCAwaDI0djI0SDB6IiBmaWxsPSJub25lIi8+PHBhdGggZD0iTTEyIDIxLjM1bC0xLjQ1LTEuMzJDNS40IDE1LjM2IDIgMTIuMjggMiA4LjUgMiA1LjQyIDQuNDIgMyA3LjUgM2MxLjc0IDAgMy40MS44MSA0LjUgMi4wOUMxMy4wOSAzLjgxIDE0Ljc2IDMgMTYuNSAzIDE5LjU4IDMgMjIgNS40MiAyMiA4LjVjMCAzLjc4LTMuNCA2Ljg2LTguNTUgMTEuNTRMMTIgMjEuMzV6Ii8+PC9zdmc+&labelColor=FF9C91&color=434343
[max33102]: https://datasheets.maximintegrated.com/en/ds/MAX30102.pdf


[license file]: https://github.com/timoth-y/chainmetric-sensorsys/blob/main/LICENSE
