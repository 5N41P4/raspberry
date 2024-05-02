# RaspberryPineapple

School project aiming to provide a lightweight interface for WiFi-reconnaissance.

The backend and the API are written in Go, because it is lightweight, and performant when compiled. The idea is to give the user a Tool which can be easily installed on a small device such as a Raspberry Pi and enabling them to easily control and do reconnaissance on WiFi networks nearby.

The interface is located in a seperate repository called pineapple.

It is not necessary to clone the pineapple repo if you don't want to make any changes in the interface, as I've included the built version in this repo aswell for ease of use and installation.

## Installation

First we need an updated OS, I run UbuntuServer on my Raspberry Pi because it has some nice networking tools already baked in which you can use to set up the possibility to connect the devices via USB-C.

To update and install the dependencies run the following on Ubuntu:

```sh
sudo apt update && sudo apt upgrade -y && sudo apt autoremove -y
sudo apt install make aircrack-ng
sudo snap install go --classic
```

Once you're done updating, upgrading and installing dependencies you should be ready to clone the repo:

```sh
git clone git@github:5N41P4/raspberry
or
git clone https://github.com/5N41P4/raspberry.git
```

After successfully cloning the repo, change into the directory and run the install script with sudo permissions:

```sh
cd raspberry
sudo ./install.sh
```

If there are dependencies not installed ond the machine, the script will tell you which ones, please install them through your package manager.

## Usage

Once everything ran successfully you can run the program with `raspberry`.

The first time you run it, the configuration prompt will pop up and you need to define the following:

- inet interface, which is the interface with an internet connection, preferrably use one which can't do monitoring or injection.
- IP address, this is the IP on which the interface will be available, make sure it's reachable through the connection you want.
- Port, on which the interface will be reachable.
- interfaces, a whitespace seperated list of all the interfaces which can do monitoring and or injection connected to the device.

Once the configuration is done, re-run the program with `raspberry` if you want to change the configuration run `raspberry -config` or change it by hand in the directory `/usr/local/raspberry/`.

If everything is configured and you have a connection to the device then you should now be able to get to the IP:port in a browser and see the interface. The rest of the functionality is reachable through the interface.

## Technologies Used

- Go
- Vue

## Features

### Reconnaissance

Gather information about all the accesspoints and client devices in your area.

### Capture

Collect the traffic that is sent between AP's and clients in the area.

## Optional Setup

There are some configurations and setups that make the application more usefull and easier to use. They however need some additional effort to be configured and are a little more in depth, I want to however link some amazing guides here:

### Raspberry Pi

If you haven't for whatever reason ever used a Raspberry Pi, they have an amazing library of documentation and guides you can follow to make whatever you think of or anyone else has ever thought of.

But all that can be overwhelming and one can get lost not knowing where to start, in that case i would advise you to look at this [Getting Started](https://www.raspberrypi.com/documentation/computers/getting-started.html) from the official Website.

Note: If you have some but not much experience and are interested i would go for the guide in the DHCP through USB-C caption.

### DHCP through USB-C

[Here](https://docs.beamnetworks.dev/en/linux/networking/pi-usb-ipad) is a nice guide by Beam Networks on how to get a DHCP server running on your Raspberry Pi, for it to be able to connect devices through the USB-C interface, so it can be used with one cable between your laptop and the Raspberry Pi.

## License

The project is licensed under the [GLWTPL](./LICENSE).

## Authors

- Aurel Corti
