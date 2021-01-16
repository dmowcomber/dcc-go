# dcc-go
A DCC Throttle to control model trains. Inspired by [DCCute](https://github.com/deltaray/DCCute)

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

### docker
```
docker build . -t github.com/dmowcomber/dcc-go
docker run -p 8080:8080 --device /dev/ttyACM0 github.com/dmowcomber/dcc-go
```

# Screenshot
![GitHub Logo](/screenshot.png)

# WARNING

dcc-go is alpha quality software at this time. Although it works, it has not
been tested with any commercial DCC system.
It has only been tested with Raspberry Pi, DCC++ on an Arduino Uno, and an ESU Loksound decoder.
The author of dcc-go claims no responsibility or liability for any damage to your system,
equipment or trains.

__USE AT YOUR OWN RISK.__
