export CGO_ENABLED=0
go build -o prometheus-exporter-merger -v -ldflags "-s -w $flags" ./


tar -zcvf prometheus-exporter-merger.tar.gz prometheus-exporter-merger
ls -lh prometheus-exporter-merger*