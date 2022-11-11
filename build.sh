export CGO_ENABLED=0
go build -o prometheus-exporter-merger -v -ldflags "-s -w $flags" ./
