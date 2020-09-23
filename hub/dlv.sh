#!/bin/sh
dlv --listen=:9026 --headless=true --api-version=2 --accept-multiclient exec ./main
