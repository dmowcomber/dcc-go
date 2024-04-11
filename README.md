# dcc-go
A DCC (Digital Command Control) Throttle to control model trains. DCC-GO can be installed on a computer like the raspberry pi to talk to an Arduino running [DCC-EX CommandStation](https://github.com/DCC-EX/CommandStation-EX) or [DCCPlusPlus BaseStation](https://github.com/DccPlusPlus/BaseStation)

DCC-GO is written in Golang and was inspired by [DCCute](https://github.com/deltaray/DCCute).

This repo contains:
* A Golang DCC Throttle library [github.com/dmowcomber/dcc-go/throttle](http://github.com/dmowcomber/dcc-go/blob/master/throttle/throttle.go)
* A DCC Throttle HTTP API
* A DCC Throttle Web page

# Find Arduino USB device
```
ls /dev |egrep 'ttyACM|ttyUSB'
ttyACM0
# found /dev/ttyACM0
```

# Installation
## Install binary and run
```
# install
go get github.com/dmowcomber/dcc-go

# run with default address 3 and default device /dev/ttyACM0
dcc-go

# run with overridden address and usb device
dcc-go -address=3 -device=/dev/ttyACM0
```

## Build docker image and run using docker-compose or docker
### docker-compose
```
docker-compose up -d
```
After making changes or pulling latest you'll need to rebuild and restart the service:
```
docker-compose build
docker-compose kill
docker-compose rm -f
docker-compose up -d
```

### docker
```
docker build . -t github.com/dmowcomber/dcc-go
docker run -p 8080:8080 --device /dev/ttyACM0 github.com/dmowcomber/dcc-go
```

### local UI/API testing with a virtual device
You can test out the UI and API against a virtual device which can be useful if you don't have an actually a DCC++ or DCC-EX Arduino connected but you'd still like to test the integration between the UI and the API.

Create a virtual device using a local file: ./testdev
```
socat -d -d -4 PTY,link="./testdev",raw,echo=0 STDIO
```
In a separate terminal, run dcc-go
```
go run . -device $(readlink ./testdev)
```

# Screenshot
![GitHub Logo](/screenshot.png)

# WARNING

dcc-go is alpha quality software at this time. Although it works, it has not
been tested with any commercial DCC system.
It has only been tested with Raspberry Pi, DCC++ or DCC-EX on an Arduino Uno, and an ESU Loksound decoder.
The author of dcc-go claims no responsibility or liability for any damage to your system,
equipment or trains.

__USE AT YOUR OWN RISK.__
