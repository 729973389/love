#!/bin/sh
s=${1}
t=${2}
sed -i 's/\"serialNumber\":\s*\"\w*\"/\"serialNumber\": \"'${s}'\"/g' server.json
sed -i 's/\"token\":\s*\"\w*\"/\"token\": \"'${2}'\"/g' server.json
