module github.com/vadv/prometheus-exporter-merger

go 1.14

require (
	gitee.com/g-devops/chisel-poll v0.0.0-20221202080939-ed5d76c30d93
	github.com/gorilla/mux v1.7.3
	github.com/ncarlier/webhookd v1.15.1
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_model v0.2.0 //v1
	github.com/prometheus/common v0.10.0 //old?
	golang.org/x/sync v0.0.0-20181221193216-37e7f081c4d4
	gopkg.in/yaml.v2 v2.2.4
)

replace (
	gitee.com/g-devops/chisel => ./_links/chisel
	gitee.com/g-devops/chisel-poll => ./_links/chisel-poll
// github.com/ncarlier/webhookd => ./_links/webhookd
)
