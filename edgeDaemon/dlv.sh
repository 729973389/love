#!/bin/sh
dlv --listen=:9024 --headless=true --api-version=2 --accept-multiclient exec ./main
