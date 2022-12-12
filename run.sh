#!/bin/bash
cur=$(cd "$(dirname "$0")"; pwd)
# ./prometheus-exporter-merger  -config=example2.yaml
# go run ./ -config=merger.json #example2.yaml


APPNAME=prometheus-exporter-merger
matchgo=$(which go)
if [ -f "$cur/main.go" ] && [ ! -z "$matchgo" ]; then
  echo "go-run"
  go run ./
else
  echo "exec-bin"
  match1=$(uname  -a |grep aarch64)
  test  -z "$match1" && exec ./$APPNAME -config=merger.json || exec ./$APPNAME-arm64 -config=merger.json
fi
