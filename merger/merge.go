package merger

import (
	"context"
	"fmt"
	"io"
	"sort"
	"sync"

	// "github.com/pkg/errors"
	prom "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"golang.org/x/sync/errgroup"

	"math/rand"
	"net/http"
	"strings"
	"net"
	// "net/http/httptest"
	// "net/url"
	// "strconv"
	// "testing"
	"time"
	"gitee.com/g-devops/chisel-poll/chserver"
)

func (m *merger) merge(ctx context.Context, w io.Writer) error {

	mu := &sync.Mutex{}
	result := map[string]*prom.MetricFamily{}

	g, ctx := errgroup.WithContext(ctx)
	for _, source := range m.sources {
		source := source
		g.Go(func() error {
			//请求source, 取得body
			resp, err := m.client.Get(source.url) //if uds
			if err != nil {
				fmt.Println("[ERR] get url: %s", source.url)
				//return errors.Wrap(err, fmt.Sprintf("get url: %s", source.url))
				return nil //xiechen-func内? continue不可用
			}
			defer resp.Body.Close()

			//解析到out下 //map[string]*prom.MetricFamily
			tp := new(expfmt.TextParser)
			//prom2json.FetchMetricFamilies: https://github.com/LyridInc/prom2lyrid/blob/b2ee0cdb0b4d173bce7ffa6e0ac8993a16733c23/model/endpoint.go
			out, err := tp.TextToMetricFamilies(resp.Body)
			if err != nil {
				fmt.Println("[ERR] parse body: %s", source.url)
				// return errors.Wrap(err, fmt.Sprintf("parse url: %s", source.url))
				return nil
			}

			// 按out依次解析name
			mu.Lock()
			defer mu.Unlock()
			for name, metricFamily := range out {
				// fmt.Println("name: "+name)
				// append labels
				if len(source.labels) > 0 {
					for _, metric := range metricFamily.Metric {
						metric.Label = append(metric.Label, source.labels...)
					}
				}
				// append metrics
				if mfResult, ok := result[name]; ok { //key可获得value值时(*prom.MetricFamily)，value追加
					mfResult.Metric = append(mfResult.Metric, metricFamily.Metric...)
				} else {
					result[name] = metricFamily  //不存在时，按name加1条
				}
			}


			//+DO: 每个source下面增加一条 filterkey_$type=0的metric (用于filter标识)； ==> filter匹配: type|dash1,dash2,dashId3,..
			if ""==source.filter {
				source.filter= "filterkey_nontype"
			}
			metricFamily2, _:= filterMetric(source)
			// append metrics //同key 不同label的merge?
			if mfResult, ok := result[source.filter]; ok { //key可获得value值时(*prom.MetricFamily)，value追加
				mfResult.Metric = append(mfResult.Metric, metricFamily2.Metric...)
			} else {
				//result[name], name来自metricKeyName; 但这里的设定不会在http输出生效
				result[source.filter] = metricFamily2  //不存在时，按name加1条
			}
			return nil
		})
	}

	// TODO 解析tunService.cmap中status=CONNECT的指标;
	m.tunMerge(result)
	
	// wait to process all routines
	if err := g.Wait(); err != nil { //并发结束
		return err
	}

	// TODO filterkey_xxx 的label添加ALL,None 两个默认值
	//   addStaticFilterByType{加ALL,None} >>遍历type: filterkey_xxx
	/* for mtype in "es,kafka,redis,minio,..." {
		//source:= ...
		addStaticFilterByType(mtype, source, result) //All, None
	} */

	// sort names
	var names []string
	for n := range result {
		names = append(names, n)
	}
	sort.Strings(names)

	// write result
	enc := expfmt.NewEncoder(w, expfmt.FmtText)
	for _, n := range names {
		err := enc.Encode(result[n])
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *merger) tunMerge(result map[string]*prom.MetricFamily) error {
	mu := &sync.Mutex{}
	detailsMap := m.ReverseTunnelService.GetTunnelDetailsMap()
	for item := range detailsMap.IterBuffered() {
		epID:= item.Key
		tunnel := item.Val.(*chserver.TunnelDetails)
		fmt.Println("tun-Meta: ", tunnel.Meta)
		if tunnel.Status != "CONNECT" {
			fmt.Println("[ERR] tunnel.Status: %s", tunnel.Status)
			return nil //err
		}

		// resp, err := m.client.Get(source.url) //if uds
		requestURL := fmt.Sprintf("http://localhost/metrics")
		req, err := http.NewRequest(http.MethodGet, requestURL, nil)
		if err != nil {
			fmt.Println("[ERR] get url(req): %s", requestURL)
			return err
		}

		// Transport tr: unixMode
		rUds:= strings.ReplaceAll(tunnel.Meta.LocalUds, "/", "-") //chclient/poll.go: createTunnel
		sockpath := "/tmp/chserver-sock/"+string(epID)+rUds //"-tmp-node-exporter1.sock"
		tr := newSocketTransport(sockpath)

		httpClient := &http.Client{
			Transport: tr,
			Timeout:   10 * time.Second,
		}
		resp, err := httpClient.Do(req)
		if err != nil {
			fmt.Println("[ERR] get url(cli.Do): %s", sockpath, err)
			return err
		}
		defer resp.Body.Close()
		defer httpClient.CloseIdleConnections()
		// fmt.Printf("resp.Body: %v\n", resp.Body) //obj's addr


		var labels []*prom.LabelPair
		k:= "mtarget"
		v:= fmt.Sprintf("%s-%s", tunnel.Meta.Desc, tunnel.Meta.Target)
		labels = append(labels, &prom.LabelPair{Name: &k, Value: &v})
		fmt.Printf("[INFO] add url: tun-%s with labels: %v\n", tunnel.Meta.Target, labels)
		

		//====func (body, labels)================================
		//解析到out下 //map[string]*prom.MetricFamily
		tp := new(expfmt.TextParser)
		//prom2json.FetchMetricFamilies: https://github.com/LyridInc/prom2lyrid/blob/b2ee0cdb0b4d173bce7ffa6e0ac8993a16733c23/model/endpoint.go
		out, err := tp.TextToMetricFamilies(resp.Body)
		if err != nil {
			fmt.Println("[ERR] parse body: %s", sockpath)
			return nil
		}
		// fmt.Println("out:　", out)

		// 按out依次解析name
		mu.Lock()
		defer mu.Unlock()
		for name, metricFamily := range out {
			// append labels
			if len(labels) > 0 {
				for _, metric := range metricFamily.Metric {
					metric.Label = append(metric.Label, labels...)
				}
			}

			// append metrics
			if mfResult, ok := result[name]; ok { //key可获得value值时(*prom.MetricFamily)，value追加
				mfResult.Metric = append(mfResult.Metric, metricFamily.Metric...)
			} else {
				result[name] = metricFamily  //不存在时，按name加1条
			}
		}

		//
		//+DO: 每个source下面增加一条 filterkey_$type=0的metric (用于filter标识)； ==> filter匹配: type|dash1,dash2,dashId3,..
		source:= &source{}
		source.labels= labels
		if ""==source.filter {
			source.filter= "filterkey_node" //"filter_nontype"
		}
		metricFamily2, _:= filterMetric(source)
		// append metrics //同key 不同label的merge?
		if mfResult, ok := result[source.filter]; ok { //key可获得value值时(*prom.MetricFamily)，value追加
			mfResult.Metric = append(mfResult.Metric, metricFamily2.Metric...)
		} else {
			//result[name], name来自metricKeyName; 但这里的设定不会在http输出生效
			result[source.filter] = metricFamily2  //不存在时，按name加1条
		}
	}

	return nil
}

func newSocketTransport(socketPath string) *http.Transport {
	defaultTimeout := 10 * time.Second

	tr := new(http.Transport)
	tr.DisableCompression = true
	tr.DialContext = func(_ context.Context, _, _ string) (net.Conn, error) {
		return net.DialTimeout("unix", socketPath, defaultTimeout)
	}
	return tr
}

func addStaticFilterByType(mtype string, source *source, result map[string]*prom.MetricFamily) error {
	if ""==source.filter {
		source.filter= "filter_nontype"
	}
	metricFamily2, _:= filterMetric(source)
	//name来自metricKeyName; 但这里的设定不会在http输出生效
	// result[source.filter] = metricFamily2 //lazy 直接用最后一个  metricFamily
	
	// append metrics //同key 不同label的merge?
	if mfResult, ok := result[source.filter]; ok { //key可获得value值时(*prom.MetricFamily)，value追加
		mfResult.Metric = append(mfResult.Metric, metricFamily2.Metric...)
	} else {
		result[source.filter] = metricFamily2  //不存在时，按name加1条
	}

	return nil
}

// 所有指标: 追加mtarget标签;
func filterMetric(source *source) (*prom.MetricFamily, error){
	// ref: https://gitee.com/g-mids/fk-exporter_exporter/blob/master/http_test.go
	rand.Seed(time.Now().UnixNano())
	// helpMsg := "help msg"
	metricType := prom.MetricType_GAUGE
	// labelName := "label_sam123"

	m_num:= 1
	metrics := make([]*prom.Metric, m_num)
	for i, _ := range metrics {
		// labelValue := fmt.Sprint(rand.Int63())
		value := rand.Float64()
		value= 1
		// ts := time.Now().UnixNano()
		metrics[i] = &prom.Metric{
			/* Label: []*prom.LabelPair{
				&prom.LabelPair{
					Name:  &labelName,
					Value: &labelValue,
				},
			}, */
			Gauge: &prom.Gauge{
				Value: &value,
			},
			// TimestampMs: &ts,
		}
	}
	if len(source.labels) > 0 {
		for _, metric := range metrics {
			metric.Label = append(metric.Label, source.labels...)//arr

			// TODO: try +All, None (All: 避免非merger节点统计的展示指标 匹配不到value)

			//不是在此加label标签..??: filterMetric只filterkey_xx一个metric;
			//加mtarget标签： 得先加1新的filterkey_redixXX{target=All,None} >> 直接filterMetric 调两次(造source)
			// metric.Label = append(metric.Label, &prom.LabelPair{Name: &k, Value: &v}) //ref main.go
		}
	}

	metricName := source.filter //fmt.Sprintf("metric%d", rand.Int63())
	metricFamily2:= &prom.MetricFamily{
		Name:   &metricName,
		// Help:   &helpMsg,
		Type:   &metricType, //must?
		Metric: metrics,
	}
	return metricFamily2, nil
}
