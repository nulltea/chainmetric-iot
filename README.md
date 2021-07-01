# ChainMetric: IoT

[![golang badge]][golang]&nbsp;
[![commit activity badge]][repo commit activity]&nbsp;
[![hardware badge]][raspberry pi]&nbsp;
[![license badge]][license url]

## Overview

_**Chainmetric Sensor System**_, being an embedded IoT sensors-equipped device, designed to be compatible with a permissioned blockchain network based on Hyperledger Fabric stack. 

By leveraging highly concurrent engine driver such implementation is ideal for harvesting environmental surrounding conditions and further publishing them onto the distributed, immutable ledger, where such data can be validated by on-chain Smart Contacts against previously assigned requirements. 

Both hardware specification and firmware architecture designed to modular and extendable to support vast variety of use cases, deployment environments, and business needs.  

The device itself is intended for deployment in the areas where assets requiring monitoring are stored or being delivered on. Thus, providing a general-purpose supply-chain monitoring solution.

![device photo]
*Chainmetric edge IoT device (development stage beta build)*

## Supports

### Digital sensors

| 📷                  | Sensor               | Metrics   | Description               |
| :------------------ | :------------------- | :-------- | :------------------------- |
| ![max44009 image][] | [MAX44009][max44009] | ![luminosity badge][]     | |
| ![si1145 image][]   | [SI1145][si1145]     | ![uv badge][] ![ir badge][] ![visible badge][] | |  |
| ![hdc1080 image][]  | [HDC1080][hdc1080]   | ![temperature badge][] ![humidity badge][] | |  |
| ![dht11 image][]    | [DHT11/22][dht22]    | ![temperature badge][] ![humidity badge][] | |
| ![ccs811 image][]   | [CCS811][ccs811]     | ![c02 badge][] ![tvoc badge][] | |
| ![bmp280 image][]   | [BMP280][bmp280]     | ![pressure badge][] ![altitude badge][] | |
| ![adxl345 image][]  | [ADXL345][adxl345]   | ![acceleration badge][] | |
| ![lsm303c image][]  | [LSM303C][lsm303c]   | ![acceleration badge][] ![magnetism badge][] ![temperature badge][] | |
| ![max30102 image][] | [MAX30102][max30102] | ![heart rate badge][] ![blood oxidation badge][] | |

### Analog sensors

| 📷                  | Sensor               | Metrics   | Description               |
| :------------------ | :------------------- | :-------- | :------------------------- |
| ![analog hall image][] | [Hall Effect][max44009] | ![magnetism badge][]     | |
| ![analog mic image][] | [Microphone][max44009] | ![noise level badge][]     | |
| ![analog piezo image][] | [Piezoelectric film][max44009] | ![vibration badge][]     | |
| ![analog mq9 image][] | [Gas (MQ-9)][max44009] | ![lpg badge][]     | |
| ![analog flame image][] | [Flame detector][max44009] | ![flame badge][]     | |



[max44009 image]: https://github.com/timoth-y/chainmetric-iot/blob/github/update_readme/docs/max44009.png?raw=true
[si1145 image]: https://github.com/timoth-y/chainmetric-iot/blob/github/update_readme/docs/si1145.png?raw=true
[hdc1080 image]: https://github.com/timoth-y/chainmetric-iot/blob/github/update_readme/docs/hdc1080.png?raw=true
[dht11 image]: https://github.com/timoth-y/chainmetric-iot/blob/github/update_readme/docs/dht11.png?raw=true
[ccs811 image]: https://github.com/timoth-y/chainmetric-iot/blob/github/update_readme/docs/ccs811.png?raw=true
[bmp280 image]: https://github.com/timoth-y/chainmetric-iot/blob/github/update_readme/docs/bmp280.png?raw=true
[adxl345 image]: https://github.com/timoth-y/chainmetric-iot/blob/github/update_readme/docs/adxl345.png?raw=true
[lsm303c image]: https://github.com/timoth-y/chainmetric-iot/blob/github/update_readme/docs/lsm303c.png?raw=true
[max30102 image]: https://github.com/timoth-y/chainmetric-iot/blob/github/update_readme/docs/max30102.png?raw=true
[analog hall image]: https://github.com/timoth-y/chainmetric-iot/blob/github/update_readme/docs/analog-hall.png?raw=true
[analog mic image]: https://github.com/timoth-y/chainmetric-iot/blob/github/update_readme/docs/analog-microphone.png?raw=true
[analog piezo image]: https://github.com/timoth-y/chainmetric-iot/blob/github/update_readme/docs/analog-piezo.png?raw=true
[analog mq9 image]: https://github.com/timoth-y/chainmetric-iot/blob/github/update_readme/docs/analog-mq-9.png?raw=true
[analog flame image]: https://github.com/timoth-y/chainmetric-iot/blob/github/update_readme/docs/analog-flame.png?raw=true

[max44009]: https://datasheets.maximintegrated.com/en/ds/MAX44009.pdf
[si1145]: https://cdn-shop.adafruit.com/datasheets/Si1145-46-47.pdf
[hdc1080]: https://www.ti.com/lit/ds/symlink/hdc1080.pdf
[dht22]: https://www.sparkfun.com/datasheets/Sensors/Temperature/DHT22.pdf
[ccs811]: https://cdn-learn.adafruit.com/downloads/pdf/adafruit-ccs811-air-quality-sensor.pdf
[bmp280]: https://cdn-shop.adafruit.com/datasheets/BST-BMP280-DS001-11.pdf
[adxl345]: https://www.sparkfun.com/datasheets/Sensors/Accelerometer/ADXL345.pdf
[adxl345]: https://www.sparkfun.com/datasheets/Sensors/Accelerometer/ADXL345.pdf
[lsm303c]: https://www.st.com/resource/en/datasheet/lsm303c.pdf
[max30102]: https://datasheets.maximintegrated.com/en/ds/MAX30102.pdf
[analog hall]: https://arduinomodules.info/ky-003-hall-magnetic-sensor-module
[analog mic]: https://datasheets.maximintegrated.com/en/ds/MAX9814.pdf

[luminosity badge]: https://img.shields.io/badge/Luminosity-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGVuYWJsZS1iYWNrZ3JvdW5kPSJuZXcgMCAwIDI0IDI0IiBoZWlnaHQ9IjI0IiB2aWV3Qm94PSIwIDAgMjQgMjQiIHdpZHRoPSIyNCI+PGc+PHBhdGggZD0iTTAsMGgyNHYyNEgwVjB6IiBmaWxsPSJub25lIi8+PC9nPjxnPjxnPjxyZWN0IGhlaWdodD0iMyIgd2lkdGg9IjIiIHg9IjExIiB5PSIxOSIvPjxyZWN0IGhlaWdodD0iMiIgd2lkdGg9IjMiIHg9IjIiIHk9IjExIi8+PHJlY3QgaGVpZ2h0PSIyIiB3aWR0aD0iMyIgeD0iMTkiIHk9IjExIi8+PHJlY3QgaGVpZ2h0PSIzIiB0cmFuc2Zvcm09Im1hdHJpeCgwLjcwNzEgLTAuNzA3MSAwLjcwNzEgMC43MDcxIC03LjY2NjUgMTcuODAxNCkiIHdpZHRoPSIxLjk5IiB4PSIxNi42NiIgeT0iMTYuNjYiLz48cmVjdCBoZWlnaHQ9IjEuOTkiIHRyYW5zZm9ybT0ibWF0cml4KDAuNzA3MSAtMC43MDcxIDAuNzA3MSAwLjcwNzEgLTEwLjk3OTEgOS44MDQxKSIgd2lkdGg9IjMiIHg9IjQuODUiIHk9IjE3LjE2Ii8+PHBhdGggZD0iTTE1LDguMDJWM0g5djUuMDJDNy43OSw4Ljk0LDcsMTAuMzcsNywxMmMwLDIuNzYsMi4yNCw1LDUsNXM1LTIuMjQsNS01QzE3LDEwLjM3LDE2LjIxLDguOTQsMTUsOC4wMnogTTExLDVoMnYyLjEgQzEyLjY4LDcuMDQsMTIuMzQsNywxMiw3cy0wLjY4LDAuMDQtMSwwLjFWNXoiLz48L2c+PC9nPjwvc3ZnPg==&labelColor=FBE967&color=434343
[uv badge]: https://img.shields.io/badge/UV%20Index-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0iYmxhY2siIHdpZHRoPSIxNHB4IiBoZWlnaHQ9IjE0cHgiPjxwYXRoIGQ9Ik0wIDBoMjR2MjRIMHoiIGZpbGw9Im5vbmUiLz48cGF0aCBkPSJNMjAgMTUuMzFMMjMuMzEgMTIgMjAgOC42OVY0aC00LjY5TDEyIC42OSA4LjY5IDRINHY0LjY5TC42OSAxMiA0IDE1LjMxVjIwaDQuNjlMMTIgMjMuMzEgMTUuMzEgMjBIMjB2LTQuNjl6TTEyIDE4Yy0zLjMxIDAtNi0yLjY5LTYtNnMyLjY5LTYgNi02IDYgMi42OSA2IDYtMi42OSA2LTYgNnoiLz48L3N2Zz4=&labelColor=FBE967&color=434343
[ir badge]: https://img.shields.io/badge/IR%20Light-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0iYmxhY2siIHdpZHRoPSIxNHB4IiBoZWlnaHQ9IjE0cHgiPjxwYXRoIGQ9Ik0wIDBoMjR2MjRIMHoiIGZpbGw9Im5vbmUiLz48cGF0aCBkPSJNMjAgMTUuMzFMMjMuMzEgMTIgMjAgOC42OVY0aC00LjY5TDEyIC42OSA4LjY5IDRINHY0LjY5TC42OSAxMiA0IDE1LjMxVjIwaDQuNjlMMTIgMjMuMzEgMTUuMzEgMjBIMjB2LTQuNjl6TTEyIDE4Yy0zLjMxIDAtNi0yLjY5LTYtNnMyLjY5LTYgNi02IDYgMi42OSA2IDYtMi42OSA2LTYgNnoiLz48L3N2Zz4=&labelColor=FBE967&color=434343
[visible badge]: https://img.shields.io/badge/Visible%20Light-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0iYmxhY2siIHdpZHRoPSIxNHB4IiBoZWlnaHQ9IjE0cHgiPjxwYXRoIGQ9Ik0wIDBoMjR2MjRIMHoiIGZpbGw9Im5vbmUiLz48cGF0aCBkPSJNMjAgMTUuMzFMMjMuMzEgMTIgMjAgOC42OVY0aC00LjY5TDEyIC42OSA4LjY5IDRINHY0LjY5TC42OSAxMiA0IDE1LjMxVjIwaDQuNjlMMTIgMjMuMzEgMTUuMzEgMjBIMjB2LTQuNjl6TTEyIDE4Yy0zLjMxIDAtNi0yLjY5LTYtNnMyLjY5LTYgNi02IDYgMi42OSA2IDYtMi42OSA2LTYgNnoiLz48L3N2Zz4=&labelColor=FBE967&color=434343
[temperature badge]: https://img.shields.io/badge/Temperature-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGNvbG9yPSJ3aGl0ZSIgaGVpZ2h0PSIxNCIgdmlld0JveD0iMCAwIDI0IDI0IiB3aWR0aD0iMTQiPjxwYXRoIGQ9Ik0wIDBoMjR2MjRIMHoiIGZpbGw9Im5vbmUiLz48cGF0aCBkPSJNMTUgMTNWNWMwLTEuNjYtMS4zNC0zLTMtM1M5IDMuMzQgOSA1djhjLTEuMjEuOTEtMiAyLjM3LTIgNCAwIDIuNzYgMi4yNCA1IDUgNXM1LTIuMjQgNS01YzAtMS42My0uNzktMy4wOS0yLTR6bS00LThjMC0uNTUuNDUtMSAxLTFzMSAuNDUgMSAxaC0xdjFoMXYyaC0xdjFoMXYyaC0yVjV6Ii8+PC9zdmc+&labelColor=72F5F5&color=434343
[humidity badge]: https://img.shields.io/badge/Humidity-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGhlaWdodD0iMjRweCIgdmlld0JveD0iMCAwIDI0IDI0IiB3aWR0aD0iMjRweCIgZmlsbD0iIzAwMDAwMCI+PHBhdGggZD0iTTAgMGgyNHYyNEgweiIgZmlsbD0ibm9uZSIvPjxwYXRoIGQ9Ik0xOS4zNSAxMC4wNEMxOC42NyA2LjU5IDE1LjY0IDQgMTIgNCA5LjExIDQgNi42MSA1LjY0IDUuMzYgOC4wNCAyLjM1IDguMzYgMCAxMC45IDAgMTRjMCAzLjMxIDIuNjkgNiA2IDZoMTNjMi43NiAwIDUtMi4yNCA1LTUgMC0yLjY0LTIuMDUtNC43OC00LjY1LTQuOTZ6TTE5IDE4SDZjLTIuMjEgMC00LTEuNzktNC00czEuNzktNCA0LTQgNCAxLjc5IDQgNGgyYzAtMi43Ni0xLjg2LTUuMDgtNC40LTUuNzhDOC42MSA2Ljg4IDEwLjIgNiAxMiA2YzMuMDMgMCA1LjUgMi40NyA1LjUgNS41di41SDE5YzEuNjUgMCAzIDEuMzUgMyAzcy0xLjM1IDMtMyAzeiIvPjwvc3ZnPg==&labelColor=72F5F5&color=434343
[c02 badge]: https://img.shields.io/badge/Air%20Quality%20(CO2)-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGVuYWJsZS1iYWNrZ3JvdW5kPSJuZXcgMCAwIDI0IDI0IiBoZWlnaHQ9IjI0IiB2aWV3Qm94PSIwIDAgMjQgMjQiIHdpZHRoPSIyNCI+PGc+PHBhdGggZD0iTTAsMGgyNHYyNEgwVjB6IiBmaWxsPSJub25lIi8+PC9nPjxnPjxnPjxwYXRoIGQ9Ik0xNC41LDE3YzAsMS42NS0xLjM1LDMtMywzcy0zLTEuMzUtMy0zaDJjMCwwLjU1LDAuNDUsMSwxLDFzMS0wLjQ1LDEtMXMtMC40NS0xLTEtMUgydi0yaDkuNSBDMTMuMTUsMTQsMTQuNSwxNS4zNSwxNC41LDE3eiBNMTksNi41QzE5LDQuNTcsMTcuNDMsMywxNS41LDNTMTIsNC41NywxMiw2LjVoMkMxNCw1LjY3LDE0LjY3LDUsMTUuNSw1UzE3LDUuNjcsMTcsNi41IFMxNi4zMyw4LDE1LjUsOEgydjJoMTMuNUMxNy40MywxMCwxOSw4LjQzLDE5LDYuNXogTTE4LjUsMTFIMnYyaDE2LjVjMC44MywwLDEuNSwwLjY3LDEuNSwxLjVTMTkuMzMsMTYsMTguNSwxNnYyIGMxLjkzLDAsMy41LTEuNTcsMy41LTMuNVMyMC40MywxMSwxOC41LDExeiIvPjwvZz48L2c+PC9zdmc+&labelColor=74FA4C&color=434343
[tvoc badge]: https://img.shields.io/badge/Air%20Quality%20(TVOC)-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGVuYWJsZS1iYWNrZ3JvdW5kPSJuZXcgMCAwIDI0IDI0IiBoZWlnaHQ9IjI0IiB2aWV3Qm94PSIwIDAgMjQgMjQiIHdpZHRoPSIyNCI+PGc+PHBhdGggZD0iTTAsMGgyNHYyNEgwVjB6IiBmaWxsPSJub25lIi8+PC9nPjxnPjxnPjxwYXRoIGQ9Ik0xNC41LDE3YzAsMS42NS0xLjM1LDMtMywzcy0zLTEuMzUtMy0zaDJjMCwwLjU1LDAuNDUsMSwxLDFzMS0wLjQ1LDEtMXMtMC40NS0xLTEtMUgydi0yaDkuNSBDMTMuMTUsMTQsMTQuNSwxNS4zNSwxNC41LDE3eiBNMTksNi41QzE5LDQuNTcsMTcuNDMsMywxNS41LDNTMTIsNC41NywxMiw2LjVoMkMxNCw1LjY3LDE0LjY3LDUsMTUuNSw1UzE3LDUuNjcsMTcsNi41IFMxNi4zMyw4LDE1LjUsOEgydjJoMTMuNUMxNy40MywxMCwxOSw4LjQzLDE5LDYuNXogTTE4LjUsMTFIMnYyaDE2LjVjMC44MywwLDEuNSwwLjY3LDEuNSwxLjVTMTkuMzMsMTYsMTguNSwxNnYyIGMxLjkzLDAsMy41LTEuNTcsMy41LTMuNVMyMC40MywxMSwxOC41LDExeiIvPjwvZz48L2c+PC9zdmc+&labelColor=74FA4C&color=434343
[lpg badge]: https://img.shields.io/badge/Air%20Quality%20(LPG)-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGVuYWJsZS1iYWNrZ3JvdW5kPSJuZXcgMCAwIDI0IDI0IiBoZWlnaHQ9IjI0IiB2aWV3Qm94PSIwIDAgMjQgMjQiIHdpZHRoPSIyNCI+PGc+PHBhdGggZD0iTTAsMGgyNHYyNEgwVjB6IiBmaWxsPSJub25lIi8+PC9nPjxnPjxnPjxwYXRoIGQ9Ik0xNC41LDE3YzAsMS42NS0xLjM1LDMtMywzcy0zLTEuMzUtMy0zaDJjMCwwLjU1LDAuNDUsMSwxLDFzMS0wLjQ1LDEtMXMtMC40NS0xLTEtMUgydi0yaDkuNSBDMTMuMTUsMTQsMTQuNSwxNS4zNSwxNC41LDE3eiBNMTksNi41QzE5LDQuNTcsMTcuNDMsMywxNS41LDNTMTIsNC41NywxMiw2LjVoMkMxNCw1LjY3LDE0LjY3LDUsMTUuNSw1UzE3LDUuNjcsMTcsNi41IFMxNi4zMyw4LDE1LjUsOEgydjJoMTMuNUMxNy40MywxMCwxOSw4LjQzLDE5LDYuNXogTTE4LjUsMTFIMnYyaDE2LjVjMC44MywwLDEuNSwwLjY3LDEuNSwxLjVTMTkuMzMsMTYsMTguNSwxNnYyIGMxLjkzLDAsMy41LTEuNTcsMy41LTMuNVMyMC40MywxMSwxOC41LDExeiIvPjwvZz48L2c+PC9zdmc+&labelColor=74FA4C&color=434343
[pressure badge]: https://img.shields.io/badge/Pressure-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGhlaWdodD0iMjQiIHZpZXdCb3g9IjAgMCAyNCAyNCIgd2lkdGg9IjI0Ij48cGF0aCBkPSJNMCAwaDI0djI0SDB6IiBmaWxsPSJub25lIi8+PHBhdGggZD0iTTAgMGgyNHYyNEgwVjB6IiBmaWxsPSJub25lIi8+PHBhdGggZD0iTTggMTloM3YzaDJ2LTNoM2wtNC00LTQgNHptOC0xNWgtM1YxaC0ydjNIOGw0IDQgNC00ek00IDl2MmgxNlY5SDR6Ii8+PHBhdGggZD0iTTQgMTJoMTZ2Mkg0eiIvPjwvc3ZnPg==&labelColor=72F5F5&color=434343
[altitude badge]: https://img.shields.io/badge/Altitude-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGhlaWdodD0iMjRweCIgdmlld0JveD0iMCAwIDI0IDI0IiB3aWR0aD0iMjRweCIgZmlsbD0iIzAwMDAwMCI+PHBhdGggZD0iTTAgMGgyNHYyNEgweiIgZmlsbD0ibm9uZSIvPjxwYXRoIGQ9Ik0xNCA2bC0zLjc1IDUgMi44NSAzLjgtMS42IDEuMkM5LjgxIDEzLjc1IDcgMTAgNyAxMGwtNiA4aDIyTDE0IDZ6Ii8+PC9zdmc+&labelColor=72F5F5&color=434343
[acceleration badge]: https://img.shields.io/badge/Acceleration-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGhlaWdodD0iMjQiIHZpZXdCb3g9IjAgMCAyNCAyNCIgd2lkdGg9IjI0Ij48cGF0aCBkPSJNMCAwaDI0djI0SDB6IiBmaWxsPSJub25lIi8+PHBhdGggZD0iTTIwLjM4IDguNTdsLTEuMjMgMS44NWE4IDggMCAwIDEtLjIyIDcuNThINS4wN0E4IDggMCAwIDEgMTUuNTggNi44NWwxLjg1LTEuMjNBMTAgMTAgMCAwIDAgMy4zNSAxOWEyIDIgMCAwIDAgMS43MiAxaDEzLjg1YTIgMiAwIDAgMCAxLjc0LTEgMTAgMTAgMCAwIDAtLjI3LTEwLjQ0em0tOS43OSA2Ljg0YTIgMiAwIDAgMCAyLjgzIDBsNS42Ni04LjQ5LTguNDkgNS42NmEyIDIgMCAwIDAgMCAyLjgzeiIvPjwvc3ZnPg==&labelColor=72F5F5&color=434343
[magnetism badge]: https://img.shields.io/badge/Magnetism-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyBoZWlnaHQ9JzMwMHB4JyB3aWR0aD0nMzAwcHgnICBmaWxsPSIjMDAwMDAwIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHhtbG5zOnhsaW5rPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5L3hsaW5rIiB2ZXJzaW9uPSIxLjEiIHg9IjBweCIgeT0iMHB4IiB2aWV3Qm94PSIwIDAgMTAwIDEwMCIgZW5hYmxlLWJhY2tncm91bmQ9Im5ldyAwIDAgMTAwIDEwMCIgeG1sOnNwYWNlPSJwcmVzZXJ2ZSI+PGc+PHBvbHlnb24gcG9pbnRzPSI2OS4xOTQsNTMuMDUgODAuNDYyLDcwLjAyNiA5MS43Niw2MS45NzcgODAuMTExLDQ1LjE5NCAgIj48L3BvbHlnb24+PHBhdGggZD0iTTM4Ljk5NCw0MS41MDdsMTUuNTg2LTcuMDI1TDQ1LjUyLDE2LjI1bC0xNi43MDYsNy42MjVjLTE2LjMwOSw5LjQxNi0yMS44OTYsMzAuMjctMTIuNDgxLDQ2LjU3OCAgIHMzMC4xMDMsMjEuNjA4LDQ2LjQxMiwxMi4xOTNsMTQuOTU3LTEwLjY1NUw2Ni40NDMsNTUuMDNsLTEzLjg3Nyw5Ljk4NmMtNi41MjQsMy43NjYtMTQuNzQ5LDEuNzMzLTE4LjUxNS00Ljc5MSAgIFMzMi40Nyw0NS4yNzMsMzguOTk0LDQxLjUwN3oiPjwvcGF0aD48cG9seWdvbiBwb2ludHM9IjYxLjIyMiw5LjA4MyA0OC42MDMsMTQuODQzIDU3LjY3LDMzLjA5IDY5LjkzMiwyNy41NjMgICI+PC9wb2x5Z29uPjwvZz48L3N2Zz4=&labelColor=72F5F5&color=434343
[heart rate badge]: https://img.shields.io/badge/Heart%20Rate-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGhlaWdodD0iMjQiIHZpZXdCb3g9IjAgMCAyNCAyNCIgd2lkdGg9IjI0Ij48cGF0aCBkPSJNMCAwaDI0djI0SDB6IiBmaWxsPSJub25lIi8+PHBhdGggZD0iTTEyIDIxLjM1bC0xLjQ1LTEuMzJDNS40IDE1LjM2IDIgMTIuMjggMiA4LjUgMiA1LjQyIDQuNDIgMyA3LjUgM2MxLjc0IDAgMy40MS44MSA0LjUgMi4wOUMxMy4wOSAzLjgxIDE0Ljc2IDMgMTYuNSAzIDE5LjU4IDMgMjIgNS40MiAyMiA4LjVjMCAzLjc4LTMuNCA2Ljg2LTguNTUgMTEuNTRMMTIgMjEuMzV6Ii8+PC9zdmc+&labelColor=FF9C91&color=434343
[blood oxidation badge]: https://img.shields.io/badge/Blood%20Oxidation-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGhlaWdodD0iMjQiIHZpZXdCb3g9IjAgMCAyNCAyNCIgd2lkdGg9IjI0Ij48cGF0aCBkPSJNMCAwaDI0djI0SDB6IiBmaWxsPSJub25lIi8+PHBhdGggZD0iTTEyIDIxLjM1bC0xLjQ1LTEuMzJDNS40IDE1LjM2IDIgMTIuMjggMiA4LjUgMiA1LjQyIDQuNDIgMyA3LjUgM2MxLjc0IDAgMy40MS44MSA0LjUgMi4wOUMxMy4wOSAzLjgxIDE0Ljc2IDMgMTYuNSAzIDE5LjU4IDMgMjIgNS40MiAyMiA4LjVjMCAzLjc4LTMuNCA2Ljg2LTguNTUgMTEuNTRMMTIgMjEuMzV6Ii8+PC9zdmc+&labelColor=FF9C91&color=434343
[noise level badge]: https://img.shields.io/badge/Noise%20Level-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGhlaWdodD0iMjRweCIgdmlld0JveD0iMCAwIDI0IDI0IiB3aWR0aD0iMjRweCIgZmlsbD0iIzAwMDAwMCI+PHBhdGggZD0iTTAgMGgyNHYyNEgweiIgZmlsbD0ibm9uZSIvPjxwYXRoIGQ9Ik0zIDl2Nmg0bDUgNVY0TDcgOUgzem0xMy41IDNjMC0xLjc3LTEuMDItMy4yOS0yLjUtNC4wM3Y4LjA1YzEuNDgtLjczIDIuNS0yLjI1IDIuNS00LjAyek0xNCAzLjIzdjIuMDZjMi44OS44NiA1IDMuNTQgNSA2Ljcxcy0yLjExIDUuODUtNSA2LjcxdjIuMDZjNC4wMS0uOTEgNy00LjQ5IDctOC43N3MtMi45OS03Ljg2LTctOC43N3oiLz48L3N2Zz4=&labelColor=72F5F5&color=434343
[vibration badge]: https://img.shields.io/badge/Vibration-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGhlaWdodD0iMjRweCIgdmlld0JveD0iMCAwIDI0IDI0IiB3aWR0aD0iMjRweCIgZmlsbD0iIzAwMDAwMCI+PHBhdGggZD0iTTAgMGgyNHYyNEgweiIgZmlsbD0ibm9uZSIvPjxwYXRoIGQ9Ik0wIDE1aDJWOUgwdjZ6bTMgMmgyVjdIM3YxMHptMTktOHY2aDJWOWgtMnptLTMgOGgyVjdoLTJ2MTB6TTE2LjUgM2gtOUM2LjY3IDMgNiAzLjY3IDYgNC41djE1YzAgLjgzLjY3IDEuNSAxLjUgMS41aDljLjgzIDAgMS41LS42NyAxLjUtMS41di0xNWMwLS44My0uNjctMS41LTEuNS0xLjV6TTE2IDE5SDhWNWg4djE0eiIvPjwvc3ZnPg==&labelColor=72F5F5&color=434343
[flame badge]: https://img.shields.io/badge/Flame-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGhlaWdodD0iMjRweCIgdmlld0JveD0iMCAwIDI0IDI0IiB3aWR0aD0iMjRweCIgZmlsbD0iIzAwMDAwMCI+PHBhdGggZD0iTTAgMGgyNHYyNEgweiIgZmlsbD0ibm9uZSIvPjxwYXRoIGQ9Ik0xMy41LjY3cy43NCAyLjY1Ljc0IDQuOGMwIDIuMDYtMS4zNSAzLjczLTMuNDEgMy43My0yLjA3IDAtMy42My0xLjY3LTMuNjMtMy43M2wuMDMtLjM2QzUuMjEgNy41MSA0IDEwLjYyIDQgMTRjMCA0LjQyIDMuNTggOCA4IDhzOC0zLjU4IDgtOEMyMCA4LjYxIDE3LjQxIDMuOCAxMy41LjY3ek0xMS43MSAxOWMtMS43OCAwLTMuMjItMS40LTMuMjItMy4xNCAwLTEuNjIgMS4wNS0yLjc2IDIuODEtMy4xMiAxLjc3LS4zNiAzLjYtMS4yMSA0LjYyLTIuNTguMzkgMS4yOS41OSAyLjY1LjU5IDQuMDQgMCAyLjY1LTIuMTUgNC44LTQuOCA0Ljh6Ii8+PC9zdmc+&labelColor=FBE967&color=434343


## Requirements
- [Raspberry Pi 3/4/Zero][raspberry pi] or other microcomputer board with GPIO, I²C, and SPI available, as well as Internet connection capabilities, preferably with Wi-Fi module. Based on considerations of portability and relative cheapness this project intends to use [RPi Zero W][rpi zero w]
- Sensors modules mentioned in the [above section](#supports)
- Assigned additional I²C buses [by utilizing spare GPIO pins][multiple i2c buses]
- Deployed and available [Blockchain network][chainmetric network repo] with its [specification configured][network spec] in `connection.yaml` file in the root directory

## Deployment

The Makefile in the root directory contains rule `sync` for syncing local project codebase with a remote device via IP address, which can be set as environmental variables:
```
$ export REMOTE_IP '192.168.31.180'
$ export REMOTE_DIR '/home/pi/sensorsys'

$ make sync
```
For the building step, there are two options either to build directly from the device on native ARM architecture via `build` make-rule or use `build-remote` to build on x86 processors with `GOARCH=arm64` option.

```
$ make build
```
To send cryptographic materials in the device for it to be able to connect to permissioned blockchain network use `crypto-sync` make-rule:
```
$ export CRYPTO_DIR '../network/crypto-config/...'
$ make crypto-sync
```

## Usage

- The device should be deployed in the same area with controlled assets (warehouse, delivery truck, etc)
- In case the device is being used for the first time it must be registered via [dedicated mobile application][chainmetric app repo] via QR code which will be automatically displayed on the embedded screen (currently [ST7789][st7789] is the only supported driver). The generated QR code will contain the device's specification: network info, supported metrics, etc.
- It is allowed to use any I²C bus for any sensor modules, the device will perform a scan to detect the location of sensors on startup.
- As soon as the device will be registered on the network it will detect surrounding assets and requirements assigned to them and will start posting sensor reading to the blockchain
- Further device management can be performed from [dedicated mobile application][chainmetric app repo]
- The registered device will automatically post its status on the startup and shutdown

## Roadmap

- [x] Caching on network connection absence [(#3)](https://github.com/timoth-y/chainmetric-iot/pull/3)
- [x] Sensor modules hot-swap support [(#1)](https://github.com/timoth-y/chainmetric-iot/pull/1)
- [x] Analog sensors ([Hall-effect sensor][hall-effect], microphone) support via ~~[MCP3008][mcp3008]~~ [ADS1115][ads1115] [(#4)]( https://github.com/timoth-y/chainmetric-iot/pull/4)
- [x] [E-Ink display][e-ink display] support [(#5)](https://github.com/timoth-y/chainmetric-iot/pull/5)
- [x] GUI for displaying statistics data and emergency warnings [(#9)](https://github.com/timoth-y/chainmetric-iot/pull/9)
- [x] Location tracking via bluetooth pairing with [mobile app][chainmetric app repo] [(#10)](https://github.com/timoth-y/chainmetric-iot/pull/10)
- [ ] A device as a blockchain node
- [ ] Video-camera driver

## Wrap up

Chainmetric device is designed for providing a real-time continuous stream of sensor-sourced environmental metrics readings to the [distributed secure ledger][chainmetric network repo] for further validation by on-chain [Smart Contracts][chainmetric contracts repo]. Such a core-concept in combination with a dedicated cross-platform [mobile application][chainmetric app repo] makes Chainmetric project an ambitious general-purpose requirements control solution for supply-chains.

## License

Licensed under the [Apache 2.0][license file].



[golang badge]: https://img.shields.io/badge/Code-Golang-informational?style=flat&logo=go&logoColor=white&color=6AD7E5
[lines counter]: https://img.shields.io/tokei/lines/github/timoth-y/chainmetric-sensorsys?color=teal&label=Lines
[commit activity badge]: https://img.shields.io/github/commit-activity/m/timoth-y/chainmetric-sensorsys?label=Commit%20activity&color=teal
[hardware badge]: https://img.shields.io/badge/Hardware-Raspberry%20Pi-informational?style=flat&logo=Raspberry%20Pi&color=953347
[license badge]: https://img.shields.io/badge/License-Apache%202.0-informational?style=flat&color=blue

[device photo]: https://github.com/timoth-y/chainmetric-iot/blob/main/docs/edge-device-betav1.png?raw=true

[this repo]: https://github.com/timoth-y/chainmetric-iot
[golang]: https://golang.org
[repo commit activity]: https://github.com/timoth-y/kicksware-api/graphs/commit-activity
[raspberry pi]: https://www.raspberrypi.org
[license url]: https://www.apache.org/licenses/LICENSE-2.0

[rpi zero w]: https://www.raspberrypi.org/products/raspberry-pi-zero-w/
[multiple i2c buses]: https://www.instructables.com/Raspberry-PI-Multiple-I2c-Devices
[network spec]: https://hyperledger-fabric.readthedocs.io/en/release-2.2/developapps/connectionprofile.html

[e-ink display]: https://www.waveshare.com/wiki/2.13inch_e-Paper_HAT
[mcp3008]: https://learn.adafruit.com/raspberry-pi-analog-to-digital-converters/mcp3008
[ads1115]: https://www.ti.com/lit/ds/symlink/ads1115.pdf
[hall-effect]: https://www.ti.com/lit/ds/symlink/drv5053.pdf

[chainmetric network repo]: https://github.com/timoth-y/chainmetric-network
[chainmetric contracts repo]: https://github.com/timoth-y/chainmetric-contracts
[chainmetric app repo]: https://github.com/timoth-y/chainmetric-app


[st7789]: https://www.newhavendisplay.com/appnotes/datasheets/LCDs/ST7789V.pdf


[license file]: https://github.com/timoth-y/chainmetric-iot/blob/main/LICENSE
