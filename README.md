# openwifiappliance

Software for wifi appliances built from the cheapest ali-express gray-market goods and the least amount of software.

## Hardware Platform

"Converse":

* raspberry pi 3
* usb powerbank
* 4 port usb hub
* 3x usb wifi adapters
* u-blox 7; gps
* GEMBIRD card-reader
* xGB sd card for data (nilfs2 suggested)
* 2GB microsd for OS image

Power draw: 5.8W (~1.3A). Potential for brownout if hub uses power over rpi.

Tested wifi:

* ar9271; 802.11n (best)
* rt3572; 802.11ac (ok)
* rt5370; (ok)
* mt7601u; 802.11n (ok; hangs if too many adapters)
* rt3070; 802.11n (poor)


## Building requirements

This should build wherever if you have docker on an otherwise reasonable linux system:

* make
* rsync
* docker

## Installing the OS image

```sh
OWA_DEV=/dev/sdx make install-image
```

## Data card

The system needs a flash card in addition to the system card. All scan data is written to the data card.

Each boot initializes a new data directory on the data card, `/owa/<boot-time>`, with the following directories:

* <mac>/ -- tcpdump logs for device, iwconfig data, lsusb data
* gps/ -- gps logs, lsusb data
* log/ -- system logs