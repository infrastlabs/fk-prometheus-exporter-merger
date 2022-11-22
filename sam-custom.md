# 

**TODO**

- merger ~~ep节点容错~~; ~~追加dash-id配置参数指标~~
- merger.json动态加载; 支持读取unixDomainSock
- mids_exporter绑定到本地uds

```bash
# metricFamily.Metric  
https://github.com/prometheus/client_model #v1.x 旧版本


```

## 

```bash
$ ./prometheus-exporter-merger  -config=example2.yaml
http://172.25.23.199:8080
http://172.25.23.199:8080/metrics

#=======================================
# HELP elasticsearch_breakers_estimated_size_bytes Estimated size in bytes of breaker
# TYPE elasticsearch_breakers_estimated_size_bytes gauge
elasticsearch_breakers_estimated_size_bytes{breaker="accounting",cluster="bigdata",es_client_node="true",es_data_node="true",es_ingest_node="true",es_master_node="true",host="172.25.20.217",name="host-172.25.20.217",appB="es7"} 2.4643461e+07
elasticsearch_breakers_estimated_size_bytes{breaker="accounting",cluster="bigdata",es_client_node="true",es_data_node="true",es_ingest_node="true",es_master_node="true",host="172.25.20.218",name="host-172.25.20.218",appB="es7"} 4.4000462e+07
.............
.............
.............
kminion_kafka_consumer_group_topic_partition_lag{group_id="myContainer0_batch02",partition_id="2",topic_name="dsg",appB="kminion"} 4
kminion_kafka_consumer_group_topic_partition_lag{group_id="myContainer0_batch02",partition_id="2",topic_name="topic1",appB="kminion"} 1
kminion_kafka_consumer_group_topic_partition_lag{group_id="myContainer1_parti02",partition_id="0",topic_name="topic2",appB="kminion"} 0
kminion_kafka_consumer_group_topic_partition_lag{group_id="myContainer1_parti02",partition_id="1",topic_name="topic2",appB="kminion"} 1
kminion_kafka_consumer_group_topic_partition_lag{group_id="myContainer1_parti02",partition_id="2",topic_name="topic2",appB="kminion"} 0
kminion_kafka_consumer_group_topic_partition_lag{group_id="myContainer2",partition_id="0",topic_name="dsg",appB="kminion"} 0
kminion_kafka_consumer_group_topic_partition_lag{group_id="myContainer2",partition_id="0",topic_name="test6",appB="kminion"} 0
kminion_kafka_consumer_group_topic_partition_lag{group_id="myContainer2",partition_id="0",topic_name="topic1",appB="kminion"} 0
kminion_kafka_consumer_group_topic_partition_lag{group_id="myContainer2",partition_id="0",topic_name="topic2",appB="kminion"} 0
kminion_kafka_consumer_group_topic_partition_lag{group_id="myContainer2",partition_id="1",topic_name="dsg",appB="kminion"} 0
kminion_kafka_consumer_group_topic_partition_lag{group_id="myContainer2",partition_id="1",topic_name="test6",appB="kminion"} 0
kminion_kafka_consumer_group_topic_partition_lag{group_id="myContainer2",partition_id="1",topic_name="topic1",appB="kminion"} 1

```