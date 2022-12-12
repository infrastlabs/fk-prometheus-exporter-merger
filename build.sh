#!/bin/bash
cur=$(cd "$(dirname "$0")"; pwd)

APPNAME=prometheus-exporter-merger

export CGO_ENABLED=0
go build -o $APPNAME -v -ldflags "-s -w $flags" ./
GOARM=7 GOARCH=arm64 GOOS=linux \
go build -o $APPNAME-arm64 -v -ldflags "-s -w $flags" ./
./$APPNAME -h |grep "\-\-ext"

tar --exclude-from=.tarignore -zcvf $APPNAME.tar.gz $APPNAME merger.json run.sh logs scripts
tar --exclude-from=.tarignore -zcvf $APPNAME-arm64.tar.gz $APPNAME-arm64 merger.json run.sh logs scripts
ls -lh $APPNAME* |grep -Ev "*.go$"
