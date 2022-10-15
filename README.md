## Overview

This is a RADIUS server implementing Aruba's Multiple Pre-Shared Key (MPSK) authentication method.

PSKs are generated once a client is authenticating and are stored in Redis.

## Build

    $ make

## Usage

The server will listen to RADIUS requests on 0.0.0.0:1812 by default.

A local Redis server is required to run the server.

    $ ./bin/server -s <pre-shared key for NAS> -r <address to redis server>

To retrieve a PSK for a client, connect the device to your access point with any password. The server will then generate a PSK for the client and be available for use. Connect again with generated PSK.

## Roadmap

- [x] RADIUS server and Redis backend
- [ ] API endpoint to retrieve PSKs from web interface

## License

MIT
