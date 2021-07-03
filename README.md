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

## Supported IO

### Digital sensors

| 📷                | Sensor               | Interface   | Metrics                                                             | Driver                                   |
| :---------------- | :------------------- | :---------- |:------------------------------------------------------------------- | :--------------------------------------- |
| ![max44009 image] | [MAX44009][max44009] | `I²C`       | ![luminosity badge][]                                               | [Custom implementation][max44009 driver] |
| ![si1145 image]   | [SI1145][si1145]     | `I²C`       | ![uv badge][] ![ir badge][] ![visible badge][] ![proximity badge][] | [Custom implementation][si1145 driver]   |
| ![hdc1080 image]  | [HDC1080][hdc1080]   | `I²C`       | ![temp badge][] ![humidity badge][]                                 | [Custom implementation][hdc1080 driver]  |
| ![dht11 image]    | [DHT11/22][dht22]    | `1-Wire`    | ![temp badge][] ![humidity badge][]                                 | Library [d2r2/go-dht](https://github.com/d2r2/go-dht) with [custom wrapper][dht22 driver] |
| ![ccs811 image]   | [CCS811][ccs811]     | `I²C`       | ![c02 badge][] ![tvoc badge][]                                      | [Custom implementation][ccs811 driver]   |
| ![bmp280 image]   | [BMP280][bmp280]     | `I²C`       | ![pressure badge][] ![altitude badge][] ![temp badge][]             | Library [google/periph](https://github.com/google/periph/tree/main/devices/bmxx80) with [custom wrapper][bmp280 driver] |
| ![adxl345 image]  | [ADXL345][adxl345]   | `I²C`       | ![acceleration badge][]                                             | [Custom implementation][adxl345 driver]  |
| ![lsm303c image]  | [LSM303C][lsm303c]   | `I²C`       | ![acceleration badge][] ![magnetism badge][] ![temp badge][]        | Fork [bskari/go-lsm303](https://github.com/bskari/go-lsm303) with [custom wrapper][lsm303c driver] |
| ![max30102 image] | [MAX30102][max30102] | `I²C`       | ![heart rate badge][] ![blood oxidation badge][]                    | Library [cgxeiji/max3010x](https://github.com/cgxeiji/max3010x) with [custom wrapper][max30102 driver] |

Digital sensors natively supports hotswap, so that it is possible to add, replace, or remove such sensors on fly,
without device restart or reconfiguration. This is possible due to combination of the address assigned for each `I²C` chip
and `CHIP_ID` register, which together should be unique. The exception is of course sensors based on `1-Wire` communication
interface, they must instead be registered as static sensors.

### Analog sensors

| 📷                    | Sensor                             | Interface                        | Metrics                | Driver                                       | ADC Driver                              |
| :-------------------- | :--------------------------------- | :------------------------------- | :--------------------- | :------------------------------------------- | :-------------------------------------- |
| ![analog hall image]  | [Hall Effect][analog hall]         | Analog with `I²C` [ADC][ads1115] | ![magnetism badge][]   | [Custom implementation][analog hall driver]  | Library [MichaelS11/go-ads][go-ads lib] |
| ![analog mic image]   | [Microphone][analog mic]           | Analog with `I²C` [ADC][ads1115] | ![noise level badge][] | [Custom implementation][analog mic driver]   | Library [MichaelS11/go-ads][go-ads lib] |
| ![analog piezo image] | [Piezoelectric film][analog piezo] | Analog with `I²C` [ADC][ads1115] | ![vibration badge][]   | [Custom implementation][analog piezo driver] | Library [MichaelS11/go-ads][go-ads lib] |
| ![analog mq9 image]   | [Gas (MQ-9)][analog mq9]           | Analog with `I²C` [ADC][ads1115] | ![lpg badge][]         | [Custom implementation][analog mq9 driver]   | Library [MichaelS11/go-ads][go-ads lib] |
| ![analog flame image] | [Flame detector][analog flame]     | Analog with `I²C` [ADC][ads1115] | ![flame badge][]       | [Custom implementation][analog flame driver] | Library [MichaelS11/go-ads][go-ads lib] |

Hotswap capabilities is also supported for analog sensors, and although these do not have any unique identifier
to be detectable by, the ADC chip does. So, the solution here is to attach ADC chip to each analog sensor
and setup different address for each used sensor. There is a limitation in this method, since ADC available addresses is finite.
For [ADS1115][ads1115] used this project we are bounded to 4 addresses (0x48, 0x49, 0x4A, 0x4B).

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
[lsm303c]: https://www.st.com/resource/en/datasheet/lsm303c.pdf
[max30102]: https://datasheets.maximintegrated.com/en/ds/MAX30102.pdf
[analog hall]: https://arduinomodules.info/ky-003-hall-magnetic-sensor-module
[analog mic]: https://datasheets.maximintegrated.com/en/ds/MAX9814.pdf
[analog piezo]: https://estheryudina.blogspot.com/2019/08/tzt-5v-piezoelectric-film-vibration.html
[analog mq9]: https://www.pololu.com/file/0J314/MQ9.pdf
[analog flame]: https://rogerbit.com/wprb/wp-content/uploads/2018/01/Flame-sensor-arduino.pdf

[luminosity badge]: https://img.shields.io/badge/Luminosity-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGVuYWJsZS1iYWNrZ3JvdW5kPSJuZXcgMCAwIDI0IDI0IiBoZWlnaHQ9IjI0IiB2aWV3Qm94PSIwIDAgMjQgMjQiIHdpZHRoPSIyNCI+PGc+PHBhdGggZD0iTTAsMGgyNHYyNEgwVjB6IiBmaWxsPSJub25lIi8+PC9nPjxnPjxnPjxyZWN0IGhlaWdodD0iMyIgd2lkdGg9IjIiIHg9IjExIiB5PSIxOSIvPjxyZWN0IGhlaWdodD0iMiIgd2lkdGg9IjMiIHg9IjIiIHk9IjExIi8+PHJlY3QgaGVpZ2h0PSIyIiB3aWR0aD0iMyIgeD0iMTkiIHk9IjExIi8+PHJlY3QgaGVpZ2h0PSIzIiB0cmFuc2Zvcm09Im1hdHJpeCgwLjcwNzEgLTAuNzA3MSAwLjcwNzEgMC43MDcxIC03LjY2NjUgMTcuODAxNCkiIHdpZHRoPSIxLjk5IiB4PSIxNi42NiIgeT0iMTYuNjYiLz48cmVjdCBoZWlnaHQ9IjEuOTkiIHRyYW5zZm9ybT0ibWF0cml4KDAuNzA3MSAtMC43MDcxIDAuNzA3MSAwLjcwNzEgLTEwLjk3OTEgOS44MDQxKSIgd2lkdGg9IjMiIHg9IjQuODUiIHk9IjE3LjE2Ii8+PHBhdGggZD0iTTE1LDguMDJWM0g5djUuMDJDNy43OSw4Ljk0LDcsMTAuMzcsNywxMmMwLDIuNzYsMi4yNCw1LDUsNXM1LTIuMjQsNS01QzE3LDEwLjM3LDE2LjIxLDguOTQsMTUsOC4wMnogTTExLDVoMnYyLjEgQzEyLjY4LDcuMDQsMTIuMzQsNywxMiw3cy0wLjY4LDAuMDQtMSwwLjFWNXoiLz48L2c+PC9nPjwvc3ZnPg==&labelColor=FBE967&color=434343
[uv badge]: https://img.shields.io/badge/UV%20Index-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0iYmxhY2siIHdpZHRoPSIxNHB4IiBoZWlnaHQ9IjE0cHgiPjxwYXRoIGQ9Ik0wIDBoMjR2MjRIMHoiIGZpbGw9Im5vbmUiLz48cGF0aCBkPSJNMjAgMTUuMzFMMjMuMzEgMTIgMjAgOC42OVY0aC00LjY5TDEyIC42OSA4LjY5IDRINHY0LjY5TC42OSAxMiA0IDE1LjMxVjIwaDQuNjlMMTIgMjMuMzEgMTUuMzEgMjBIMjB2LTQuNjl6TTEyIDE4Yy0zLjMxIDAtNi0yLjY5LTYtNnMyLjY5LTYgNi02IDYgMi42OSA2IDYtMi42OSA2LTYgNnoiLz48L3N2Zz4=&labelColor=FBE967&color=434343
[ir badge]: https://img.shields.io/badge/IR%20Light-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0iYmxhY2siIHdpZHRoPSIxNHB4IiBoZWlnaHQ9IjE0cHgiPjxwYXRoIGQ9Ik0wIDBoMjR2MjRIMHoiIGZpbGw9Im5vbmUiLz48cGF0aCBkPSJNMjAgMTUuMzFMMjMuMzEgMTIgMjAgOC42OVY0aC00LjY5TDEyIC42OSA4LjY5IDRINHY0LjY5TC42OSAxMiA0IDE1LjMxVjIwaDQuNjlMMTIgMjMuMzEgMTUuMzEgMjBIMjB2LTQuNjl6TTEyIDE4Yy0zLjMxIDAtNi0yLjY5LTYtNnMyLjY5LTYgNi02IDYgMi42OSA2IDYtMi42OSA2LTYgNnoiLz48L3N2Zz4=&labelColor=FBE967&color=434343
[visible badge]: https://img.shields.io/badge/Visible%20Light-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0iYmxhY2siIHdpZHRoPSIxNHB4IiBoZWlnaHQ9IjE0cHgiPjxwYXRoIGQ9Ik0wIDBoMjR2MjRIMHoiIGZpbGw9Im5vbmUiLz48cGF0aCBkPSJNMjAgMTUuMzFMMjMuMzEgMTIgMjAgOC42OVY0aC00LjY5TDEyIC42OSA4LjY5IDRINHY0LjY5TC42OSAxMiA0IDE1LjMxVjIwaDQuNjlMMTIgMjMuMzEgMTUuMzEgMjBIMjB2LTQuNjl6TTEyIDE4Yy0zLjMxIDAtNi0yLjY5LTYtNnMyLjY5LTYgNi02IDYgMi42OSA2IDYtMi42OSA2LTYgNnoiLz48L3N2Zz4=&labelColor=FBE967&color=434343
[proximity badge]: https://img.shields.io/badge/Proximity-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGhlaWdodD0iMjRweCIgdmlld0JveD0iMCAwIDI0IDI0IiB3aWR0aD0iMjRweCIgZmlsbD0iIzAwMDAwMCI+PHBhdGggZD0iTTAgMGgyNHYyNEgweiIgZmlsbD0ibm9uZSIvPjxwYXRoIGQ9Ik0yMSA2SDNjLTEuMSAwLTIgLjktMiAydjhjMCAxLjEuOSAyIDIgMmgxOGMxLjEgMCAyLS45IDItMlY4YzAtMS4xLS45LTItMi0yem0wIDEwSDNWOGgydjRoMlY4aDJ2NGgyVjhoMnY0aDJWOGgydjRoMlY4aDJ2OHoiLz48L3N2Zz4=&labelColor=50B1AA&color=434343
[temp badge]: https://img.shields.io/badge/Temperature-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGNvbG9yPSJ3aGl0ZSIgaGVpZ2h0PSIxNCIgdmlld0JveD0iMCAwIDI0IDI0IiB3aWR0aD0iMTQiPjxwYXRoIGQ9Ik0wIDBoMjR2MjRIMHoiIGZpbGw9Im5vbmUiLz48cGF0aCBkPSJNMTUgMTNWNWMwLTEuNjYtMS4zNC0zLTMtM1M5IDMuMzQgOSA1djhjLTEuMjEuOTEtMiAyLjM3LTIgNCAwIDIuNzYgMi4yNCA1IDUgNXM1LTIuMjQgNS01YzAtMS42My0uNzktMy4wOS0yLTR6bS00LThjMC0uNTUuNDUtMSAxLTFzMSAuNDUgMSAxaC0xdjFoMXYyaC0xdjFoMXYyaC0yVjV6Ii8+PC9zdmc+&labelColor=72F5F5&color=434343
[humidity badge]: https://img.shields.io/badge/Humidity-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGhlaWdodD0iMjRweCIgdmlld0JveD0iMCAwIDI0IDI0IiB3aWR0aD0iMjRweCIgZmlsbD0iIzAwMDAwMCI+PHBhdGggZD0iTTAgMGgyNHYyNEgweiIgZmlsbD0ibm9uZSIvPjxwYXRoIGQ9Ik0xOS4zNSAxMC4wNEMxOC42NyA2LjU5IDE1LjY0IDQgMTIgNCA5LjExIDQgNi42MSA1LjY0IDUuMzYgOC4wNCAyLjM1IDguMzYgMCAxMC45IDAgMTRjMCAzLjMxIDIuNjkgNiA2IDZoMTNjMi43NiAwIDUtMi4yNCA1LTUgMC0yLjY0LTIuMDUtNC43OC00LjY1LTQuOTZ6TTE5IDE4SDZjLTIuMjEgMC00LTEuNzktNC00czEuNzktNCA0LTQgNCAxLjc5IDQgNGgyYzAtMi43Ni0xLjg2LTUuMDgtNC40LTUuNzhDOC42MSA2Ljg4IDEwLjIgNiAxMiA2YzMuMDMgMCA1LjUgMi40NyA1LjUgNS41di41SDE5YzEuNjUgMCAzIDEuMzUgMyAzcy0xLjM1IDMtMyAzeiIvPjwvc3ZnPg==&labelColor=72F5F5&color=434343
[c02 badge]: https://img.shields.io/badge/Air%20Quality%20(CO2)-informational?style=flat&logo=data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0idXRmLTgiPz4NCjwhLS0gU3ZnIFZlY3RvciBJY29ucyA6IGh0dHA6Ly93d3cub25saW5ld2ViZm9udHMuY29tL2ljb24gLS0+DQo8IURPQ1RZUEUgc3ZnIFBVQkxJQyAiLS8vVzNDLy9EVEQgU1ZHIDEuMS8vRU4iICJodHRwOi8vd3d3LnczLm9yZy9HcmFwaGljcy9TVkcvMS4xL0RURC9zdmcxMS5kdGQiPg0KPHN2ZyB2ZXJzaW9uPSIxLjEiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIgeG1sbnM6eGxpbms9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkveGxpbmsiIHg9IjBweCIgeT0iMHB4IiB2aWV3Qm94PSIwIDAgMTAwMCAxMDAwIiBlbmFibGUtYmFja2dyb3VuZD0ibmV3IDAgMCAxMDAwIDEwMDAiIHhtbDpzcGFjZT0icHJlc2VydmUiPg0KPG1ldGFkYXRhPiBTdmcgVmVjdG9yIEljb25zIDogaHR0cDovL3d3dy5vbmxpbmV3ZWJmb250cy5jb20vaWNvbiA8L21ldGFkYXRhPg0KPGc+PGcgdHJhbnNmb3JtPSJ0cmFuc2xhdGUoMC4wMDAwMDAsNTExLjAwMDAwMCkgc2NhbGUoMC4xMDAwMDAsLTAuMTAwMDAwKSI+PHBhdGggZD0iTTI0MzIuNSw0Nzk2LjFjLTYzOC4xLTcyLjgtMTA3OS4yLTY5OC4xLTkzMS41LTEzMjEuMmM4OS45LTM4MS4yLDM3NC43LTY4Ny40LDc0My4xLTc5OC43YzExOS45LTM2LjQsMTg4LjQtNDUsMzcwLjUtMzYuNGwyMjIuNyw2LjRsNDI4LjMtNzc3LjNsNDI4LjMtNzc5LjVsLTEzNy4xLTEzN2MtNDgxLjgtNDczLjItNjQyLjQtMTE1NC4yLTQyNi4xLTE3OTIuM2MyNy44LTgzLjUsNzkuMi0yMDEuMywxMTUuNi0yNjMuNGw2NC4yLTExMy41bC03NjAuMi02ODcuNGwtNzU4LTY4NS4ybC03MC43LDQyLjhjLTIwOS44LDEzMC42LTU5MSwxNjctODU4LjcsODUuN2MtMzU1LjUtMTA5LjItNjQ0LjUtNDIxLjktNzMyLjMtNzk0LjRjLTQwLjctMTczLjUtNDAuNy0zMTQuOCwwLTQ4OC4zYzE2NC45LTcwMC4yLDkzMS41LTEwNDUsMTU4NC42LTcxMy4xYzI0Ni4zLDEyNC4yLDQ3Ny41LDQzNC43LDU0Niw3MzAuMmMzNi40LDE1NC4yLDE3LjEsNDI4LjMtNDIuOCw1OTkuNmwtNDcuMSwxMzQuOWw2MCw2MGMzNC4zLDMyLjEsMzcyLjYsMzQwLjUsNzUzLjcsNjgzLjFsNjg5LjUsNjIzLjJsMTQ1LjYtOTQuMmMzMDguNC0xOTkuMSw1NzguMi0yNzguNCw5NTcuMi0yNzYuMmMyMTQuMSwwLDI4OS4xLDEwLjcsNDQ1LjQsNTcuOGMxMDQuOSwzMCwyMzcuNyw4MS40LDI5NS41LDExMS40YzU5LjksMzAsMTE1LjYsNTEuNCwxMjYuMyw0Ny4xYzE3LjEtNi40LDU5OS42LTcwMC4yLDU5OS42LTcxNS4yYzAtNi40LTMwLTY2LjQtNjQuMi0xMzdjLTM3NC43LTc0MC45LDEzOS4yLTE1ODguOSw5NjMuNi0xNTg4LjljODAwLjksMCwxMzI1LjUsODMwLjgsOTc4LjYsMTU1Mi41Yy0xMzAuNiwyNzQuMS0zNzAuNCw0ODMuOS02NTEsNTcxLjdjLTE5Mi43LDYwLTQ3MS4xLDU3LjgtNjU3LjQtMi4xbC0xMzkuMi00Ny4xbC01MS40LDYyLjFjLTMxNi45LDM3Ni45LTU1MC4zLDY2MS43LTU1MC4zLDY3NC41YzAsNi40LDMyLjEsNTUuNyw3MC43LDEwNy4xYzE1Mi4xLDIwMS4zLDI0Ni4zLDQxMy4zLDMwNi4yLDY3NC41YzM0LjMsMTU0LjIsMzQuMywzNTMuMywyLjEsNzQ3LjNsLTQuMyw2Mi4xTDcyMDEuMyw1MjJsNzUzLjgsMzQyLjZMODA4MS40LDc0OWM0MzAuNC0zOTguMywxMDc3LjEtMzgzLjMsMTQ5Ni44LDM2LjRjMjQ2LjMsMjQ0LjEsMzU5LjcsNTc4LjIsMzEwLjUsOTEyLjJjLTc5LjIsNTM5LjYtNTM3LjUsOTQwLjEtMTA3Mi44LDkzNy45Yy0zMTYuOS0yLjEtNTUyLjUtOTguNS03NzAuOS0zMTkuMWMtMjI0LjgtMjI0LjgtMzI5LjgtNDg0LTMyMS4yLTc5Mi4zbDQuMy0xNjAuNmwtNzQ3LjMtMzM4LjNjLTQxMS4yLTE4Ni4zLTc1OC4xLTMzNi4yLTc3MC45LTMzNGMtMTIuOSw0LjMtNjYuNCw2MC0xMTkuOSwxMjQuMmMtMjE4LjQsMjcyLTU1OC45LDQ4OC4yLTkwNS44LDU4MC4zYy0yMzUuNiw2Mi4xLTU5Ny40LDYyLjEtODMzLDBjLTk0LjItMjUuNy0xNzEuMy00Mi44LTE3NS42LTQwLjdjLTIuMiw0LjMtMTk3LDM1Ny42LTQzMi42LDc4NS45bC00MzAuNCw3NzkuNWw4My41LDk4LjVDNDAyNS43LDM3NjguMywzNDA2LjgsNDkwOS42LDI0MzIuNSw0Nzk2LjF6IE00NzYyLjMsODg2bC02LjQtMTY3bC0xMzkuMi0xNWMtNDA2LjktNDUtNzY2LjYtNDAwLjQtODM5LjQtODI2LjZsLTI1LjctMTU0LjJoLTE1Ni4zaC0xNTYuM3Y4MS40YzAsMTI2LjMsNzcuMSw0MDAuNCwxNTIuMSw1NDEuOGM4NS43LDE1OC41LDI4NC44LDM3Ni45LDQzMC40LDQ3My4yYzE5NywxMjguNSw0ODEuOCwyMjQuOCw2NzYuNywyMjkuMWw3MC43LDIuMUw0NzYyLjMsODg2eiIvPjwvZz48L2c+DQo8L3N2Zz4=&labelColor=74FA4C&color=434343
[tvoc badge]: https://img.shields.io/badge/Air%20Quality%20(TVOC)-informational?style=flat&logo=data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0idXRmLTgiPz4NCjwhLS0gU3ZnIFZlY3RvciBJY29ucyA6IGh0dHA6Ly93d3cub25saW5ld2ViZm9udHMuY29tL2ljb24gLS0+DQo8IURPQ1RZUEUgc3ZnIFBVQkxJQyAiLS8vVzNDLy9EVEQgU1ZHIDEuMS8vRU4iICJodHRwOi8vd3d3LnczLm9yZy9HcmFwaGljcy9TVkcvMS4xL0RURC9zdmcxMS5kdGQiPg0KPHN2ZyB2ZXJzaW9uPSIxLjEiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIgeG1sbnM6eGxpbms9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkveGxpbmsiIHg9IjBweCIgeT0iMHB4IiB2aWV3Qm94PSIwIDAgMTAwMCAxMDAwIiBlbmFibGUtYmFja2dyb3VuZD0ibmV3IDAgMCAxMDAwIDEwMDAiIHhtbDpzcGFjZT0icHJlc2VydmUiPg0KPG1ldGFkYXRhPiBTdmcgVmVjdG9yIEljb25zIDogaHR0cDovL3d3dy5vbmxpbmV3ZWJmb250cy5jb20vaWNvbiA8L21ldGFkYXRhPg0KPGc+PGcgdHJhbnNmb3JtPSJ0cmFuc2xhdGUoMC4wMDAwMDAsNTExLjAwMDAwMCkgc2NhbGUoMC4xMDAwMDAsLTAuMTAwMDAwKSI+PHBhdGggZD0iTTI0MzIuNSw0Nzk2LjFjLTYzOC4xLTcyLjgtMTA3OS4yLTY5OC4xLTkzMS41LTEzMjEuMmM4OS45LTM4MS4yLDM3NC43LTY4Ny40LDc0My4xLTc5OC43YzExOS45LTM2LjQsMTg4LjQtNDUsMzcwLjUtMzYuNGwyMjIuNyw2LjRsNDI4LjMtNzc3LjNsNDI4LjMtNzc5LjVsLTEzNy4xLTEzN2MtNDgxLjgtNDczLjItNjQyLjQtMTE1NC4yLTQyNi4xLTE3OTIuM2MyNy44LTgzLjUsNzkuMi0yMDEuMywxMTUuNi0yNjMuNGw2NC4yLTExMy41bC03NjAuMi02ODcuNGwtNzU4LTY4NS4ybC03MC43LDQyLjhjLTIwOS44LDEzMC42LTU5MSwxNjctODU4LjcsODUuN2MtMzU1LjUtMTA5LjItNjQ0LjUtNDIxLjktNzMyLjMtNzk0LjRjLTQwLjctMTczLjUtNDAuNy0zMTQuOCwwLTQ4OC4zYzE2NC45LTcwMC4yLDkzMS41LTEwNDUsMTU4NC42LTcxMy4xYzI0Ni4zLDEyNC4yLDQ3Ny41LDQzNC43LDU0Niw3MzAuMmMzNi40LDE1NC4yLDE3LjEsNDI4LjMtNDIuOCw1OTkuNmwtNDcuMSwxMzQuOWw2MCw2MGMzNC4zLDMyLjEsMzcyLjYsMzQwLjUsNzUzLjcsNjgzLjFsNjg5LjUsNjIzLjJsMTQ1LjYtOTQuMmMzMDguNC0xOTkuMSw1NzguMi0yNzguNCw5NTcuMi0yNzYuMmMyMTQuMSwwLDI4OS4xLDEwLjcsNDQ1LjQsNTcuOGMxMDQuOSwzMCwyMzcuNyw4MS40LDI5NS41LDExMS40YzU5LjksMzAsMTE1LjYsNTEuNCwxMjYuMyw0Ny4xYzE3LjEtNi40LDU5OS42LTcwMC4yLDU5OS42LTcxNS4yYzAtNi40LTMwLTY2LjQtNjQuMi0xMzdjLTM3NC43LTc0MC45LDEzOS4yLTE1ODguOSw5NjMuNi0xNTg4LjljODAwLjksMCwxMzI1LjUsODMwLjgsOTc4LjYsMTU1Mi41Yy0xMzAuNiwyNzQuMS0zNzAuNCw0ODMuOS02NTEsNTcxLjdjLTE5Mi43LDYwLTQ3MS4xLDU3LjgtNjU3LjQtMi4xbC0xMzkuMi00Ny4xbC01MS40LDYyLjFjLTMxNi45LDM3Ni45LTU1MC4zLDY2MS43LTU1MC4zLDY3NC41YzAsNi40LDMyLjEsNTUuNyw3MC43LDEwNy4xYzE1Mi4xLDIwMS4zLDI0Ni4zLDQxMy4zLDMwNi4yLDY3NC41YzM0LjMsMTU0LjIsMzQuMywzNTMuMywyLjEsNzQ3LjNsLTQuMyw2Mi4xTDcyMDEuMyw1MjJsNzUzLjgsMzQyLjZMODA4MS40LDc0OWM0MzAuNC0zOTguMywxMDc3LjEtMzgzLjMsMTQ5Ni44LDM2LjRjMjQ2LjMsMjQ0LjEsMzU5LjcsNTc4LjIsMzEwLjUsOTEyLjJjLTc5LjIsNTM5LjYtNTM3LjUsOTQwLjEtMTA3Mi44LDkzNy45Yy0zMTYuOS0yLjEtNTUyLjUtOTguNS03NzAuOS0zMTkuMWMtMjI0LjgtMjI0LjgtMzI5LjgtNDg0LTMyMS4yLTc5Mi4zbDQuMy0xNjAuNmwtNzQ3LjMtMzM4LjNjLTQxMS4yLTE4Ni4zLTc1OC4xLTMzNi4yLTc3MC45LTMzNGMtMTIuOSw0LjMtNjYuNCw2MC0xMTkuOSwxMjQuMmMtMjE4LjQsMjcyLTU1OC45LDQ4OC4yLTkwNS44LDU4MC4zYy0yMzUuNiw2Mi4xLTU5Ny40LDYyLjEtODMzLDBjLTk0LjItMjUuNy0xNzEuMy00Mi44LTE3NS42LTQwLjdjLTIuMiw0LjMtMTk3LDM1Ny42LTQzMi42LDc4NS45bC00MzAuNCw3NzkuNWw4My41LDk4LjVDNDAyNS43LDM3NjguMywzNDA2LjgsNDkwOS42LDI0MzIuNSw0Nzk2LjF6IE00NzYyLjMsODg2bC02LjQtMTY3bC0xMzkuMi0xNWMtNDA2LjktNDUtNzY2LjYtNDAwLjQtODM5LjQtODI2LjZsLTI1LjctMTU0LjJoLTE1Ni4zaC0xNTYuM3Y4MS40YzAsMTI2LjMsNzcuMSw0MDAuNCwxNTIuMSw1NDEuOGM4NS43LDE1OC41LDI4NC44LDM3Ni45LDQzMC40LDQ3My4yYzE5NywxMjguNSw0ODEuOCwyMjQuOCw2NzYuNywyMjkuMWw3MC43LDIuMUw0NzYyLjMsODg2eiIvPjwvZz48L2c+DQo8L3N2Zz4=&labelColor=74FA4C&color=434343
[lpg badge]: https://img.shields.io/badge/Gas%20(LPG)-informational?style=flat&logo=data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0idXRmLTgiPz4NCjwhLS0gU3ZnIFZlY3RvciBJY29ucyA6IGh0dHA6Ly93d3cub25saW5ld2ViZm9udHMuY29tL2ljb24gLS0+DQo8IURPQ1RZUEUgc3ZnIFBVQkxJQyAiLS8vVzNDLy9EVEQgU1ZHIDEuMS8vRU4iICJodHRwOi8vd3d3LnczLm9yZy9HcmFwaGljcy9TVkcvMS4xL0RURC9zdmcxMS5kdGQiPg0KPHN2ZyB2ZXJzaW9uPSIxLjEiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIgeG1sbnM6eGxpbms9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkveGxpbmsiIHg9IjBweCIgeT0iMHB4IiB2aWV3Qm94PSIwIDAgMTAwMCAxMDAwIiBlbmFibGUtYmFja2dyb3VuZD0ibmV3IDAgMCAxMDAwIDEwMDAiIHhtbDpzcGFjZT0icHJlc2VydmUiPg0KPG1ldGFkYXRhPiBTdmcgVmVjdG9yIEljb25zIDogaHR0cDovL3d3dy5vbmxpbmV3ZWJmb250cy5jb20vaWNvbiA8L21ldGFkYXRhPg0KPGc+PGcgdHJhbnNmb3JtPSJ0cmFuc2xhdGUoMC4wMDAwMDAsNTExLjAwMDAwMCkgc2NhbGUoMC4xMDAwMDAsLTAuMTAwMDAwKSI+PHBhdGggZD0iTTI0MzIuNSw0Nzk2LjFjLTYzOC4xLTcyLjgtMTA3OS4yLTY5OC4xLTkzMS41LTEzMjEuMmM4OS45LTM4MS4yLDM3NC43LTY4Ny40LDc0My4xLTc5OC43YzExOS45LTM2LjQsMTg4LjQtNDUsMzcwLjUtMzYuNGwyMjIuNyw2LjRsNDI4LjMtNzc3LjNsNDI4LjMtNzc5LjVsLTEzNy4xLTEzN2MtNDgxLjgtNDczLjItNjQyLjQtMTE1NC4yLTQyNi4xLTE3OTIuM2MyNy44LTgzLjUsNzkuMi0yMDEuMywxMTUuNi0yNjMuNGw2NC4yLTExMy41bC03NjAuMi02ODcuNGwtNzU4LTY4NS4ybC03MC43LDQyLjhjLTIwOS44LDEzMC42LTU5MSwxNjctODU4LjcsODUuN2MtMzU1LjUtMTA5LjItNjQ0LjUtNDIxLjktNzMyLjMtNzk0LjRjLTQwLjctMTczLjUtNDAuNy0zMTQuOCwwLTQ4OC4zYzE2NC45LTcwMC4yLDkzMS41LTEwNDUsMTU4NC42LTcxMy4xYzI0Ni4zLDEyNC4yLDQ3Ny41LDQzNC43LDU0Niw3MzAuMmMzNi40LDE1NC4yLDE3LjEsNDI4LjMtNDIuOCw1OTkuNmwtNDcuMSwxMzQuOWw2MCw2MGMzNC4zLDMyLjEsMzcyLjYsMzQwLjUsNzUzLjcsNjgzLjFsNjg5LjUsNjIzLjJsMTQ1LjYtOTQuMmMzMDguNC0xOTkuMSw1NzguMi0yNzguNCw5NTcuMi0yNzYuMmMyMTQuMSwwLDI4OS4xLDEwLjcsNDQ1LjQsNTcuOGMxMDQuOSwzMCwyMzcuNyw4MS40LDI5NS41LDExMS40YzU5LjksMzAsMTE1LjYsNTEuNCwxMjYuMyw0Ny4xYzE3LjEtNi40LDU5OS42LTcwMC4yLDU5OS42LTcxNS4yYzAtNi40LTMwLTY2LjQtNjQuMi0xMzdjLTM3NC43LTc0MC45LDEzOS4yLTE1ODguOSw5NjMuNi0xNTg4LjljODAwLjksMCwxMzI1LjUsODMwLjgsOTc4LjYsMTU1Mi41Yy0xMzAuNiwyNzQuMS0zNzAuNCw0ODMuOS02NTEsNTcxLjdjLTE5Mi43LDYwLTQ3MS4xLDU3LjgtNjU3LjQtMi4xbC0xMzkuMi00Ny4xbC01MS40LDYyLjFjLTMxNi45LDM3Ni45LTU1MC4zLDY2MS43LTU1MC4zLDY3NC41YzAsNi40LDMyLjEsNTUuNyw3MC43LDEwNy4xYzE1Mi4xLDIwMS4zLDI0Ni4zLDQxMy4zLDMwNi4yLDY3NC41YzM0LjMsMTU0LjIsMzQuMywzNTMuMywyLjEsNzQ3LjNsLTQuMyw2Mi4xTDcyMDEuMyw1MjJsNzUzLjgsMzQyLjZMODA4MS40LDc0OWM0MzAuNC0zOTguMywxMDc3LjEtMzgzLjMsMTQ5Ni44LDM2LjRjMjQ2LjMsMjQ0LjEsMzU5LjcsNTc4LjIsMzEwLjUsOTEyLjJjLTc5LjIsNTM5LjYtNTM3LjUsOTQwLjEtMTA3Mi44LDkzNy45Yy0zMTYuOS0yLjEtNTUyLjUtOTguNS03NzAuOS0zMTkuMWMtMjI0LjgtMjI0LjgtMzI5LjgtNDg0LTMyMS4yLTc5Mi4zbDQuMy0xNjAuNmwtNzQ3LjMtMzM4LjNjLTQxMS4yLTE4Ni4zLTc1OC4xLTMzNi4yLTc3MC45LTMzNGMtMTIuOSw0LjMtNjYuNCw2MC0xMTkuOSwxMjQuMmMtMjE4LjQsMjcyLTU1OC45LDQ4OC4yLTkwNS44LDU4MC4zYy0yMzUuNiw2Mi4xLTU5Ny40LDYyLjEtODMzLDBjLTk0LjItMjUuNy0xNzEuMy00Mi44LTE3NS42LTQwLjdjLTIuMiw0LjMtMTk3LDM1Ny42LTQzMi42LDc4NS45bC00MzAuNCw3NzkuNWw4My41LDk4LjVDNDAyNS43LDM3NjguMywzNDA2LjgsNDkwOS42LDI0MzIuNSw0Nzk2LjF6IE00NzYyLjMsODg2bC02LjQtMTY3bC0xMzkuMi0xNWMtNDA2LjktNDUtNzY2LjYtNDAwLjQtODM5LjQtODI2LjZsLTI1LjctMTU0LjJoLTE1Ni4zaC0xNTYuM3Y4MS40YzAsMTI2LjMsNzcuMSw0MDAuNCwxNTIuMSw1NDEuOGM4NS43LDE1OC41LDI4NC44LDM3Ni45LDQzMC40LDQ3My4yYzE5NywxMjguNSw0ODEuOCwyMjQuOCw2NzYuNywyMjkuMWw3MC43LDIuMUw0NzYyLjMsODg2eiIvPjwvZz48L2c+DQo8L3N2Zz4=&labelColor=74FA4C&color=434343
[pressure badge]: https://img.shields.io/badge/Pressure-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGhlaWdodD0iMjQiIHZpZXdCb3g9IjAgMCAyNCAyNCIgd2lkdGg9IjI0Ij48cGF0aCBkPSJNMCAwaDI0djI0SDB6IiBmaWxsPSJub25lIi8+PHBhdGggZD0iTTAgMGgyNHYyNEgwVjB6IiBmaWxsPSJub25lIi8+PHBhdGggZD0iTTggMTloM3YzaDJ2LTNoM2wtNC00LTQgNHptOC0xNWgtM1YxaC0ydjNIOGw0IDQgNC00ek00IDl2MmgxNlY5SDR6Ii8+PHBhdGggZD0iTTQgMTJoMTZ2Mkg0eiIvPjwvc3ZnPg==&labelColor=72F5F5&color=434343
[altitude badge]: https://img.shields.io/badge/Altitude-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGhlaWdodD0iMjRweCIgdmlld0JveD0iMCAwIDI0IDI0IiB3aWR0aD0iMjRweCIgZmlsbD0iIzAwMDAwMCI+PHBhdGggZD0iTTAgMGgyNHYyNEgweiIgZmlsbD0ibm9uZSIvPjxwYXRoIGQ9Ik0xNCA2bC0zLjc1IDUgMi44NSAzLjgtMS42IDEuMkM5LjgxIDEzLjc1IDcgMTAgNyAxMGwtNiA4aDIyTDE0IDZ6Ii8+PC9zdmc+&labelColor=72F5F5&color=434343
[acceleration badge]: https://img.shields.io/badge/Acceleration-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGhlaWdodD0iMjQiIHZpZXdCb3g9IjAgMCAyNCAyNCIgd2lkdGg9IjI0Ij48cGF0aCBkPSJNMCAwaDI0djI0SDB6IiBmaWxsPSJub25lIi8+PHBhdGggZD0iTTIwLjM4IDguNTdsLTEuMjMgMS44NWE4IDggMCAwIDEtLjIyIDcuNThINS4wN0E4IDggMCAwIDEgMTUuNTggNi44NWwxLjg1LTEuMjNBMTAgMTAgMCAwIDAgMy4zNSAxOWEyIDIgMCAwIDAgMS43MiAxaDEzLjg1YTIgMiAwIDAgMCAxLjc0LTEgMTAgMTAgMCAwIDAtLjI3LTEwLjQ0em0tOS43OSA2Ljg0YTIgMiAwIDAgMCAyLjgzIDBsNS42Ni04LjQ5LTguNDkgNS42NmEyIDIgMCAwIDAgMCAyLjgzeiIvPjwvc3ZnPg==&labelColor=50B1AA&color=434343
[magnetism badge]: https://img.shields.io/badge/Magnetism-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyBoZWlnaHQ9JzMwMHB4JyB3aWR0aD0nMzAwcHgnICBmaWxsPSIjMDAwMDAwIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHhtbG5zOnhsaW5rPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5L3hsaW5rIiB2ZXJzaW9uPSIxLjEiIHg9IjBweCIgeT0iMHB4IiB2aWV3Qm94PSIwIDAgMTAwIDEwMCIgZW5hYmxlLWJhY2tncm91bmQ9Im5ldyAwIDAgMTAwIDEwMCIgeG1sOnNwYWNlPSJwcmVzZXJ2ZSI+PGc+PHBvbHlnb24gcG9pbnRzPSI2OS4xOTQsNTMuMDUgODAuNDYyLDcwLjAyNiA5MS43Niw2MS45NzcgODAuMTExLDQ1LjE5NCAgIj48L3BvbHlnb24+PHBhdGggZD0iTTM4Ljk5NCw0MS41MDdsMTUuNTg2LTcuMDI1TDQ1LjUyLDE2LjI1bC0xNi43MDYsNy42MjVjLTE2LjMwOSw5LjQxNi0yMS44OTYsMzAuMjctMTIuNDgxLDQ2LjU3OCAgIHMzMC4xMDMsMjEuNjA4LDQ2LjQxMiwxMi4xOTNsMTQuOTU3LTEwLjY1NUw2Ni40NDMsNTUuMDNsLTEzLjg3Nyw5Ljk4NmMtNi41MjQsMy43NjYtMTQuNzQ5LDEuNzMzLTE4LjUxNS00Ljc5MSAgIFMzMi40Nyw0NS4yNzMsMzguOTk0LDQxLjUwN3oiPjwvcGF0aD48cG9seWdvbiBwb2ludHM9IjYxLjIyMiw5LjA4MyA0OC42MDMsMTQuODQzIDU3LjY3LDMzLjA5IDY5LjkzMiwyNy41NjMgICI+PC9wb2x5Z29uPjwvZz48L3N2Zz4=&labelColor=50B1AA&color=434343
[heart rate badge]: https://img.shields.io/badge/Heart%20Rate-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGhlaWdodD0iMjQiIHZpZXdCb3g9IjAgMCAyNCAyNCIgd2lkdGg9IjI0Ij48cGF0aCBkPSJNMCAwaDI0djI0SDB6IiBmaWxsPSJub25lIi8+PHBhdGggZD0iTTEyIDIxLjM1bC0xLjQ1LTEuMzJDNS40IDE1LjM2IDIgMTIuMjggMiA4LjUgMiA1LjQyIDQuNDIgMyA3LjUgM2MxLjc0IDAgMy40MS44MSA0LjUgMi4wOUMxMy4wOSAzLjgxIDE0Ljc2IDMgMTYuNSAzIDE5LjU4IDMgMjIgNS40MiAyMiA4LjVjMCAzLjc4LTMuNCA2Ljg2LTguNTUgMTEuNTRMMTIgMjEuMzV6Ii8+PC9zdmc+&labelColor=FF9C91&color=434343
[blood oxidation badge]: https://img.shields.io/badge/Blood%20Oxidation-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyBoZWlnaHQ9JzMwMHB4JyB3aWR0aD0nMzAwcHgnICBmaWxsPSIjMDAwMDAwIiB4bWxuczp4PSJodHRwOi8vbnMuYWRvYmUuY29tL0V4dGVuc2liaWxpdHkvMS4wLyIgeG1sbnM6aT0iaHR0cDovL25zLmFkb2JlLmNvbS9BZG9iZUlsbHVzdHJhdG9yLzEwLjAvIiB4bWxuczpncmFwaD0iaHR0cDovL25zLmFkb2JlLmNvbS9HcmFwaHMvMS4wLyIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIiB4bWxuczp4bGluaz0iaHR0cDovL3d3dy53My5vcmcvMTk5OS94bGluayIgdmVyc2lvbj0iMS4xIiB4PSIwcHgiIHk9IjBweCIgdmlld0JveD0iMCAwIDEwMCAxMDAiIHN0eWxlPSJlbmFibGUtYmFja2dyb3VuZDpuZXcgMCAwIDEwMCAxMDA7IiB4bWw6c3BhY2U9InByZXNlcnZlIj48c3dpdGNoPjxmb3JlaWduT2JqZWN0IHJlcXVpcmVkRXh0ZW5zaW9ucz0iaHR0cDovL25zLmFkb2JlLmNvbS9BZG9iZUlsbHVzdHJhdG9yLzEwLjAvIiB4PSIwIiB5PSIwIiB3aWR0aD0iMSIgaGVpZ2h0PSIxIj48L2ZvcmVpZ25PYmplY3Q+PGcgaTpleHRyYW5lb3VzPSJzZWxmIj48cGF0aCBkPSJNNTAsMi41Yy0zLjIsMC02LjEsMS43LTcuNyw0LjRDMzYuNCwxNi43LDE3LDQ5LjksMTcsNjQuNWMwLDE4LjIsMTQuOCwzMywzMywzM2MxOC4yLDAsMzMtMTQuOCwzMy0zMyAgICBjMC0xNC42LTE5LjMtNDcuOC0yNS4zLTU3LjdDNTYuMSw0LjIsNTMuMiwyLjUsNTAsMi41eiBNNDMuNSw0OC42Yy0wLjEsMC4xLTkuMSwxNC43LTIuNSwyNC42YzEuNCwyLjEsMC44LDQuOS0xLjIsNi4yICAgIGMtMC44LDAuNS0xLjYsMC44LTIuNSwwLjhjLTEuNCwwLTIuOS0wLjctMy43LTJjLTkuOS0xNC43LDEuOC0zMy42LDIuMy0zNC40YzEuMy0yLjEsNC4xLTIuNyw2LjItMS40ICAgIEM0NC4yLDQzLjcsNDQuOCw0Ni41LDQzLjUsNDguNnoiPjwvcGF0aD48L2c+PC9zd2l0Y2g+PC9zdmc+&labelColor=FF9C91&color=434343
[noise level badge]: https://img.shields.io/badge/Noise%20Level-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGhlaWdodD0iMjRweCIgdmlld0JveD0iMCAwIDI0IDI0IiB3aWR0aD0iMjRweCIgZmlsbD0iIzAwMDAwMCI+PHBhdGggZD0iTTAgMGgyNHYyNEgweiIgZmlsbD0ibm9uZSIvPjxwYXRoIGQ9Ik0zIDl2Nmg0bDUgNVY0TDcgOUgzem0xMy41IDNjMC0xLjc3LTEuMDItMy4yOS0yLjUtNC4wM3Y4LjA1YzEuNDgtLjczIDIuNS0yLjI1IDIuNS00LjAyek0xNCAzLjIzdjIuMDZjMi44OS44NiA1IDMuNTQgNSA2Ljcxcy0yLjExIDUuODUtNSA2LjcxdjIuMDZjNC4wMS0uOTEgNy00LjQ5IDctOC43N3MtMi45OS03Ljg2LTctOC43N3oiLz48L3N2Zz4=&labelColor=50B1AA&color=434343
[vibration badge]: https://img.shields.io/badge/Vibration-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGhlaWdodD0iMjRweCIgdmlld0JveD0iMCAwIDI0IDI0IiB3aWR0aD0iMjRweCIgZmlsbD0iIzAwMDAwMCI+PHBhdGggZD0iTTAgMGgyNHYyNEgweiIgZmlsbD0ibm9uZSIvPjxwYXRoIGQ9Ik0wIDE1aDJWOUgwdjZ6bTMgMmgyVjdIM3YxMHptMTktOHY2aDJWOWgtMnptLTMgOGgyVjdoLTJ2MTB6TTE2LjUgM2gtOUM2LjY3IDMgNiAzLjY3IDYgNC41djE1YzAgLjgzLjY3IDEuNSAxLjUgMS41aDljLjgzIDAgMS41LS42NyAxLjUtMS41di0xNWMwLS44My0uNjctMS41LTEuNS0xLjV6TTE2IDE5SDhWNWg4djE0eiIvPjwvc3ZnPg==&labelColor=50B1AA&color=434343
[flame badge]: https://img.shields.io/badge/Flame-informational?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGhlaWdodD0iMjRweCIgdmlld0JveD0iMCAwIDI0IDI0IiB3aWR0aD0iMjRweCIgZmlsbD0iIzAwMDAwMCI+PHBhdGggZD0iTTAgMGgyNHYyNEgweiIgZmlsbD0ibm9uZSIvPjxwYXRoIGQ9Ik0xMy41LjY3cy43NCAyLjY1Ljc0IDQuOGMwIDIuMDYtMS4zNSAzLjczLTMuNDEgMy43My0yLjA3IDAtMy42My0xLjY3LTMuNjMtMy43M2wuMDMtLjM2QzUuMjEgNy41MSA0IDEwLjYyIDQgMTRjMCA0LjQyIDMuNTggOCA4IDhzOC0zLjU4IDgtOEMyMCA4LjYxIDE3LjQxIDMuOCAxMy41LjY3ek0xMS43MSAxOWMtMS43OCAwLTMuMjItMS40LTMuMjItMy4xNCAwLTEuNjIgMS4wNS0yLjc2IDIuODEtMy4xMiAxLjc3LS4zNiAzLjYtMS4yMSA0LjYyLTIuNTguMzkgMS4yOS41OSAyLjY1LjU5IDQuMDQgMCAyLjY1LTIuMTUgNC44LTQuOCA0Ljh6Ii8+PC9zdmc+&labelColor=FBE967&color=434343

[max44009 driver]: https://github.com/timoth-y/chainmetric-iot/blob/main/drivers/sensors/max44009.go
[si1145 driver]: https://github.com/timoth-y/chainmetric-iot/blob/main/drivers/sensors/si1145.go
[hdc1080 driver]: https://github.com/timoth-y/chainmetric-iot/blob/main/drivers/sensors/hdc1080.go
[dht22 driver]: https://github.com/timoth-y/chainmetric-iot/blob/main/drivers/sensors/dht22.go
[ccs811 driver]: https://github.com/timoth-y/chainmetric-iot/blob/main/drivers/sensors/ccs811.go
[bmp280 driver]: https://github.com/timoth-y/chainmetric-iot/blob/main/drivers/sensors/bmp280.go
[adxl345 driver]: https://github.com/timoth-y/chainmetric-iot/blob/main/drivers/sensors/adxl345.go
[lsm303c driver]: https://github.com/timoth-y/chainmetric-iot/blob/main/drivers/sensors/lsm303c.go
[max30102 driver]: https://github.com/timoth-y/chainmetric-iot/blob/main/drivers/sensors/max30102.go
[analog hall driver]: https://github.com/timoth-y/chainmetric-iot/blob/main/drivers/sensors/adc_hall.go
[analog mic driver]: https://github.com/timoth-y/chainmetric-iot/blob/main/drivers/sensors/adc_microphone.go
[analog piezo driver]: https://github.com/timoth-y/chainmetric-iot/blob/main/drivers/sensors/adc_piezo.go
[analog mq9 driver]: https://github.com/timoth-y/chainmetric-iot/blob/main/drivers/sensors/adc_mq9.go
[analog flame driver]: https://github.com/timoth-y/chainmetric-iot/blob/main/drivers/sensors/adc_flame.go

[go-ads lib]: https://github.com/MichaelS11/go-ads
[adc driver]: https://github.com/timoth-y/chainmetric-iot/blob/main/drivers/periphery/adc.go

### Power

| 📷                | Chip                 | Interface | Hardware options         | Driver                                   |
| :---------------- | :------------------- | :-------- | :---------------------   | :--------------------------------------- |
| ![ups-lite image] | [MAX17040][max17040] | `I²C`     | [UPS-Lite][ups-lite]     | [Custom implementation][max17040 driver] |

[ups-lite image]: https://github.com/timoth-y/chainmetric-iot/blob/github/update_readme/docs/ups-lite.png?raw=true
[ups-lite]: https://hackaday.io/project/173847-ups-lite
[max17040]: https://cdn.hackaday.io/files/1738477437870048/MAX17040.pdf
[max17040 driver]: https://github.com/timoth-y/chainmetric-iot/blob/main/drivers/power/ups.go

### Displays

| 📷                              | Chip                   | Interface | Hardware options                   | Driver                                   |
| :------------------------------ | :--------------------- | :-------- | :--------------------------------- | :--------------------------------------- |
| [![e-ink image]][e-ink display] | [e-Ink (EPD)][e-ink]   | `SPI`     | [2.13' e-Paper HAT][e-ink display] | [Custom implementation][e-ink driver]    |
| [![st7789 image]][lcd display]  | [ST7789][st7789]       | `SPI`     | [2' IPS LCD][lcd display]          | [Custom implementation][st7789 driver]   |

[e-ink image]: https://github.com/timoth-y/chainmetric-iot/blob/github/update_readme/docs/e-ink.png?raw=true
[st7789 image]: https://github.com/timoth-y/chainmetric-iot/blob/github/update_readme/docs/st7789.png?raw=true

[e-ink]: https://www.waveshare.com/w/upload/e/e6/2.13inch_e-Paper_Datasheet.pdf
[st7789]: https://www.buydisplay.com/download/ic/ST7789.pdf

[e-ink display]: https://www.waveshare.com/wiki/2.13inch_e-Paper_HAT
[lcd display]: https://www.waveshare.com/wiki/2inch_LCD_Module

[e-ink driver]: https://github.com/timoth-y/chainmetric-iot/blob/main/drivers/display/eink.go
[st7789 driver]: https://github.com/timoth-y/chainmetric-iot/blob/return_of_st7789_display_driver/drivers/display/lcd.go

### Bluetooth

| Protocol | Service          | UUID                                   | Description                                        | Driver                                          |
| :------- | :--------------- | :------------------------------------- | :------------------------------------------------- | :---------------------------------------------- |
| BLE      | Location service | `F8AE4978-5AAB-46C3-A8CB-127F347EAA01` | Enables location tethering with mobile application | Library [go-ble](https://github.com/go-ble/ble) |


## Firmware architecture

The design decisions firmware development was mostly based on the idea of enabling wide range of use cases for IoT device.
With that in mind the architecture itself is based on the concept of modularity, thus allowing feature set to be
easily extendable and adaptable for new hardware, areas of application, and deployment environments.

### Drivers

On the lower level we got drivers, which of course implementing direct communication and control of the hardware:

- [`drivers/periphery`][drivers/periphery] - communication protocols implementation and wrappers (`I²C`, `SPI`, `GPIO`)
- [`drivers/sensors`][drivers/sensors] - each sensor custom drivers or library wrappers reduced to a single interface structure
- [`drivers/display`][drivers/display] - possible visual output drivers 
- [`drivers/power`][drivers/power] - UPS's and power management drivers

[drivers/periphery]: https://github.com/timoth-y/chainmetric-iot/blob/main/drivers/periphery
[drivers/sensors]: https://github.com/timoth-y/chainmetric-iot/blob/main/drivers/sensors
[drivers/display]: https://github.com/timoth-y/chainmetric-iot/blob/main/drivers/sensors
[drivers/power]: https://github.com/timoth-y/chainmetric-iot/blob/main/drivers/power

### Network

Since being a firmware oriented on IoT device in the Blockchain infrastructure,
it is of course cannot come without network layer:

- [`network/blockchain`][network/blockchain] - the Hyperledger Fabric related clients as well as Smart Contract RPCs.
- [`network/localnet`][network/localnet] - package proving interface for low range radius communication with other devices

[network/blockchain]: https://github.com/timoth-y/chainmetric-iot/blob/main/network/blockchain
[network/localnet]: https://github.com/timoth-y/chainmetric-iot/blob/main/network/localnet

### Controllers

Lastly, on the higher level we got business logic driven controllers, which by taking favor of the previous layers
define required feature set of the IoT device.

Besides [`controllers/gui`][controllers/gui] and [`controllers/storage`][controllers/storage],
which by definition do exactly what their names stand for, this layer holds [`controllers/device`][controllers/device],
which is the controller for device itself where the most of the domain logic are.
This is the place where the mentioned above modularity really starts to shine.

[controllers/device]: https://github.com/timoth-y/chainmetric-iot/blob/main/controllers/device
[controllers/gui]: https://github.com/timoth-y/chainmetric-iot/blob/main/controllers/gui
[controllers/storage]: https://github.com/timoth-y/chainmetric-iot/blob/main/controllers/storage
[controllers/engine]: https://github.com/timoth-y/chainmetric-iot/blob/main/controllers/engine

#### Logical modules

While the `Device` type holds the current state of the device, along with cached operational data,
up-to-date sensors registry and the main context, it actually does not hold any feature-related functionality.
Instead, it delegates that to the *logical modules*, each containing its own atomic portion of the business logic, responsibilities,
and is capable of mutating state of the `Device`.

| Module              | Description                                                                                                                     | Implementation                                            |
| :------------------ | :------------------------------------------------------------------------------------------------------------------------------ | :-------------------------------------------------------- |
| `LIFECYCLE_MANAGER` | Manages device initialization, registration on network, updates device state on startup and shutdown                            | [`modules/lifecycle_manager`][modules/lifecycle_manager]  |
| `ENGINE_OPERATOR`   | Operates `SensorsReader` engine, handles readings requests and posts results on chain                                           | [`modules/engine_operator`][modules/engine_operator]      |
| `EVENTS_OBSERVER`   | Listens to changes in assets, requirements or device state on the Blockchain, handles them accordingly                          | [`modules/events_observer`][modules/events_observer]      |
| `CACHE_MANAGER`     | Caches operational data on device startup, updates or flushes cache when needed                                                 | [`modules/cache_manager`][modules/cache_manager]          |
| `FAILOVER_HANDLER`  | Handles network issues which lead to inability of posting readings by storing them in the persistent storage (embedded LevelDB) | [`modules/failover_handler`][modules/failover_handler]    |
| `HOTSWAP_DETECTOR`  | Monitors and detects changes in device's periphery, updates available sensors pool                                              | [`modules/hotswap_detector`][modules/hotswap_detector]    |
| `REMOTE_CONTROLLER` | Listens to remote commands directed to the current device, performs command execution against the device                        | [`modules/remote_controller`][modules/remote_controller]  |
| `POWER_MANAGER`     | Monitors device power consumption and battery level, updates device state on chain                                              | [`modules/power_manager`][modules/power_manager]          |
| `LOCATION_MANAGER`  | Manages device physical location, updates device state on chain                                                                 | [`modules/location_manager`][modules/location_manager]    |
| `GUI_RENDERER`      | Displays device specs, requests throughput, and other useful data on the display if such is available                           | [`modules/gui_renderer`][modules/gui_renderer]            |

Logical modules are implemented in such a way, so they cannot directly communicate with each other
and foremost not being aware of other modules' existence. Yet, some functionality requires exactly that,
e.g. requires intermediate input or must be triggered by some occurred event.

For such purposes modules can utilize shared state of the device or more often local operational events.
Such idea is borrowed from the Event Driven Architecture (EDA) and is achieved with help of [`timoth-y/go-eventdriver`](https://github.com/timoth-y/go-eventdriver) package.

Although such approach definitely can increase complexity of the codebase it still has some major benefits,
which have been considered to be worth-taking trade off. Some of them are stronger abstraction, low logic blocks coupling,
higher extensibility with lower risc of broking something, and more.

#### Sensors reading engine

One other essential component of the device functionality is the [`controllers/engine`][controllers/engine] package,
whose responsibility includes proving interface to receive sensor reading requests, subscribe handlers for outcome results,
control, initialize and deallocate physical sensor modules if such are on stand-by for certain amount of time.

Engine leverages the convenience of the Go concurrency model, allowing to harvest data from multiply sensors
and for multiply subscribers at the same time, while not blocking other components and modules execution. 

[modules/lifecycle_manager]: https://github.com/timoth-y/chainmetric-iot/blob/main/controllers/device/modules/lifecycle_manager.go
[modules/engine_operator]: https://github.com/timoth-y/chainmetric-iot/blob/main/controllers/device/modules/engine_operator.go
[modules/events_observer]: https://github.com/timoth-y/chainmetric-iot/blob/main/controllers/device/modules/events_observer.go
[modules/cache_manager]: https://github.com/timoth-y/chainmetric-iot/blob/main/controllers/device/modules/cache_manager.go
[modules/failover_handler]: https://github.com/timoth-y/chainmetric-iot/blob/main/controllers/device/modules/failover_handler.go
[modules/hotswap_detector]: https://github.com/timoth-y/chainmetric-iot/blob/main/controllers/device/modules/hotswap_detector.go
[modules/remote_controller]: https://github.com/timoth-y/chainmetric-iot/blob/main/controllers/device/modules/remote_controller.go
[modules/power_manager]: https://github.com/timoth-y/chainmetric-iot/blob/main/controllers/device/modules/power_manager.go
[modules/location_manager]: https://github.com/timoth-y/chainmetric-iot/blob/main/controllers/device/modules/location_manager.go
[modules/gui_renderer]: https://github.com/timoth-y/chainmetric-iot/blob/main/controllers/device/modules/gui_renderer.go

## Requirements
- [Raspberry Pi 3/4/Zero][raspberry pi] or other microcomputer board with `GPIO`, `I²C`, and `SPI` available, as well as Internet connection capabilities, preferably with Wi-Fi module. Based on considerations of portability and relative cheapness this project intends to use [RPi Zero W][rpi zero w]
- Sensors modules mentioned in the [above section](#supported-io)
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
