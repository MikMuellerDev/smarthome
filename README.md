# Smarthome
# This repository has been transferred to the [Smarthome](https://github.com/smarthome-go/smarthome) organization.
**Version**: `0.0.26-beta-rc.4`


A completely self-built Smarthome-system written in Go.

[![Go Build](https://github.com/MikMuellerDev/smarthome/actions/workflows/go.yml/badge.svg)](https://github.com/MikMuellerDev/smarthome/actions/workflows/go.yml)
[![](https://tokei.rs/b1/github/MikMuellerDev/smarthome?category=code)](https://github.com/MikMuellerDev/smarthome).

## What is Smarthome?
Smarthome is a completely self-build home-automation system written in Go *(backend)* and Svelte *(frontend)*.
The system focuses on functionality and simplicity in order to guarantee a stable and reliable home-automation system which is actually helpful in automating common tasks.

### Concepts
- Completely self-hostable on your own infrastructure
- Simple setup: when version `1.0` is released, the entire configuration can be managed from the web interface
- Is able to operate without internet connection (except for the weather which relies on an API service)
- Privacy focused: Your data will stay on your system because Smarthome is not relying on cloud infrastructure
- An up-to-date docker-image is built and published to Docker Hub on every release 

## Hardware
As of April 27, 2022 the only way to make Smarthome interact with the real world is through the use of [smarthome-hw](https://github.com/MikMuellerDev/smarthome-hw), a Hardware interface which is required in order to interact with most generic 433mhz remote-sockets.
Naturally, the use of smarthome-hw requires physical hardware in order to communicate with remote sockets.

However, support for additional hardware, for example Zigbee devices is planned and would open additional possibilities for integration with other hardware.

## Getting Started
### The `setup.json`
Most of the configuration of the smarthome server can be achieved using the `setup.json` file.
This file is scanned and evaluated at startup.

```json
{
    "hardwareNodes": [
        {
            "name": "test raspberry pi",
            "url": "http://localhost:8070",
            "token": "smarthome"
        }
    ],
    "rooms": [
        {
            "data": {
                "id": "test_room",
                "name": "Test Room",
                "description": "This is a test room"
            },
            "switches": [
                {
                    "id": "s1",
                    "name": "Lamp1"
                },
                {
                    "id": "s2",
                    "name": "Lamp2"
                }
            ],
            "cameras": [
                {
                    "id": "test_camera",
                    "name": "Test Camera",
                    "url": "https://mik-mueller.de/assets/foo.png"
                }
            ]
        }
    ]
}
```
