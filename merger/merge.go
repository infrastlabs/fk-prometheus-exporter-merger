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
	// "net/http"
	// "net/http/httptest"
	// "net/url"
	// "strconv"
	// "testing"
	"time"

)

func (m *merger) merge(ctx context.Context, w io.Writer) error {

	mu := &sync.Mutex{}
	result := map[string]*prom.MetricFamily{}

	g, ctx := errgroup.WithContext(ctx)
	for _, source := range m.sources {
		source := source
		g.Go(func() error {
			//请求source, 取得body
			resp, err := m.client.Get(source.url)
			if err != nil {
				fmt.Println("[ERR] cli.Get")
				// fmt.Println("[ERR] get url: %s", source.url)
				//return errors.Wrap(err, fmt.Sprintf("get url: %s", source.url))
				// xiechen-func内? continue不可用
				return nil
			}
			defer resp.Body.Close()

			//解析到out下 //map[string]*prom.MetricFamily
			tp := new(expfmt.TextParser)
			//prom2json.FetchMetricFamilies: https://github.com/LyridInc/prom2lyrid/blob/b2ee0cdb0b4d173bce7ffa6e0ac8993a16733c23/model/endpoint.go
			out, err := tp.TextToMetricFamilies(resp.Body)
			if err != nil {
				fmt.Println("[ERR] TextToMetricFamilies")
				// fmt.Println("[ERR] parse url: %s", source.url)
				// return errors.Wrap(err, fmt.Sprintf("parse url: %s", source.url))
				return nil
			}

			// 按out依次解析name
			mu.Lock()
			defer mu.Unlock()
			// var metricFamily2 *prom.MetricFamily
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
				// metricFamily2= metricFamily
			}

			// TODO result: 每个source下面增加一条 filter_$type=0的metric (用于filter标识)；
			//  ===> filter匹配: type|dash1,dash2,dashId3,..
			/* _= &prom.MetricFamily{
				"",
				"",
				nil,
				nil,

			} */

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
		})
	}

	// wait to process all routines
	if err := g.Wait(); err != nil {
		return err
	}

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
			metric.Label = append(metric.Label, source.labels...)
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
