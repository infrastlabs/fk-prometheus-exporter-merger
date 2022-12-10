package cmd

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	prom "github.com/prometheus/client_model/go"
	"github.com/vadv/prometheus-exporter-merger/merger"
	"github.com/vadv/prometheus-exporter-merger/hook"

	sercmd "gitee.com/g-devops/chisel-poll/chserver/cmd"
	// clicmd "gitee.com/g-devops/chisel-poll/chclient/cmd"
	"github.com/gorilla/mux"
)

func Execute() {

	var (
		configPath = flag.String("config", "/config/prometheus-exporter-merger.yaml", "Path to config")
	)
	flag.Parse()

	c, err := parseConfig(*configPath)
	if err != nil {
		panic(err)
	}

	m := merger.New(c.ScrapeTimeout)
	for _, s := range c.Sources {
		var labels []*prom.LabelPair
		for k, v := range s.Labels {
			k, v := k, v
			labels = append(labels, &prom.LabelPair{Name: &k, Value: &v})
		}
		log.Printf("[INFO] add url: %s with labels: %v\n", s.Url, s.Labels)
		m.AddSource(s.Url, s.Filter, labels)
	}


	r := mux.NewRouter()
	prefix:= "/api/endpoints"
	reverseTunnelService:= sercmd.MuxHandle(r, prefix)

	prefix= "/api/hook"
	hook.SetVars(prefix)
	hook.MuxHandle(r, prefix)

	r.PathPrefix("/").Handler(&handler{m: m})
	//m.设置ReverseTunnelService
	m.AddChiselService(reverseTunnelService)

	// localUds := "/tmp/chclient-001.sock"
	// clicmd.LocalFileServer(conf.ListenUds)
	// clicmd.StartPollService(conf.PollServerAddr, prefix, conf.PollAgentId, conf.ListenUds)


	// srv := &http.Server{Addr: c.Listen, Handler: &handler{m: m}} //TODO mux 集成chisel
	srv := &http.Server{Addr: c.Listen, Handler: r}
	log.Printf("[INFO] starting listen %s\n", c.Listen)
	go srv.ListenAndServe()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}
}
