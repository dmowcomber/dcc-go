# dcc-go
A command line DCC Throttle to control your model trains written in Golang. Inspired by [DCCute](https://github.com/deltaray/DCCute)

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

## Build docker image and run in docker
```
docker-compose up -d
```


# WARNING

dcc-go is alpha quality software at this time. Although it works, it has not
been tested with any commercial DCC system.
It has only been tested with Raspberry Pi, DCC++ on an Arduino Uno, and an ESU Loksound decoder.
The author of dcc-go claims no responsibility or liability for any damage to your system,
equipment or trains.

__USE AT YOUR OWN RISK.__
