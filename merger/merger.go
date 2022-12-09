package merger

import (
	"context"
	"io"
	"net/http"
	"sync"
	"time"

	prom "github.com/prometheus/client_model/go"
	"gitee.com/g-devops/chisel-poll/chserver"
	"gitee.com/g-devops/chisel-poll/chserver/chisel"
)

type Merger interface {
	Merge(w io.Writer) error
	AddSource(url string, filter string, labels []*prom.LabelPair)
	AddChiselService(chiselService *chisel.Service)
}

type merger struct {
	mu            sync.Mutex
	scrapeTimeout time.Duration
	client        *http.Client
	sources       []*source
	ReverseTunnelService chserver.ReverseTunnelService //+
}

type source struct {
	url    string
	filter string
	labels []*prom.LabelPair
}

func New(scrapeTimeout time.Duration) Merger {
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   false,
			DisableCompression:  false,
			MaxIdleConns:        1,
			MaxIdleConnsPerHost: 1,
			MaxConnsPerHost:     10,
			IdleConnTimeout:     5 * time.Minute,
		},
		Timeout: scrapeTimeout,
	}
	return &merger{
		scrapeTimeout: scrapeTimeout,
		client:        client,
	}
}

// AddSource new source
func (m *merger) AddSource(url string, filter string, labels []*prom.LabelPair) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sources = append(m.sources, &source{url: url, filter: filter, labels: labels})
}


func (m *merger) AddChiselService(chiselService *chisel.Service) {
	// m.mu.Lock()
	// defer m.mu.Unlock()
	m.ReverseTunnelService = chiselService
}

// Merge sources
func (m *merger) Merge(w io.Writer) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), m.scrapeTimeout)
	defer cancel()
	return m.merge(ctx, w)
}
