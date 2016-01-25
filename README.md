# pzsvc-lasinfo

## Overview

`pzsvc-lasinfo` was initially developed as a demo application to test the VeniceGeo DevOps infrastructure, and modeled after https://github.com/venicegeo/refapp-devops.

`pzsvc-lasinfo` is written entirely in Go with no runtime dependencies.

All commits to master will be pushed to Cloud Foundry.

Calls to http://pzsvc-lasinfo.cf.piazzageo.io should return simply "Hi!".

Calls to http://pzsvc-lasinfo.cf.piazzageo.io/info will return a JSON message summarizing the LAS header info.
