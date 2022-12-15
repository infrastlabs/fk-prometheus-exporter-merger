FROM infrastlabs/alpine-ext:weak

# download
WORKDIR /download
# first x3: kaf, es, redis
RUN \
wget https://github.com/danielqsj/kafka_exporter/releases/download/v1.6.0/kafka_exporter-1.6.0.linux-amd64.tar.gz; \
wget https://github.com/redpanda-data/kminion/releases/download/v2.2.0/kminion_2.2.0_linux_amd64.tar.gz; \
wget https://github.com/prometheus-community/elasticsearch_exporter/releases/download/v1.5.0/elasticsearch_exporter-1.5.0.linux-amd64.tar.gz; \
wget https://github.com/oliver006/redis_exporter/releases/download/v1.45.0/redis_exporter-v1.45.0.linux-amd64.tar.gz; \
ls -lh;

# merger
RUN \
# url=https://gitee.com/g-mids/fk-prometheus-exporter-merger/releases/download/v22.11.17/prometheus-exporter-merger.tar.gz; \
url=https://gitee.com/g-mids/fk-prometheus-exporter-merger/releases/download/v22.11.23/prometheus-exporter-merger.tar.gz; \
wget $url; \
ls -lh; 


# base: node, blackbox
WORKDIR /download/base
RUN \
curl -fSL -O https://github.com/prometheus/node_exporter/releases/download/v1.4.0/node_exporter-1.4.0.linux-amd64.tar.gz; \
curl -fSL -O https://github.com/prometheus/blackbox_exporter/releases/download/v0.22.0/blackbox_exporter-0.22.0.linux-amd64.tar.gz; \
ls -lh; 

# OTH
# - https://github.com/ClickHouse/clickhouse_exporter #无发版
WORKDIR /download/oth
RUN \
curl -fSL -O https://github.com/prometheus/mysqld_exporter/releases/download/v0.14.0/mysqld_exporter-0.14.0.linux-amd64.tar.gz; \
curl -fSL -O https://github.com/prometheus-community/postgres_exporter/releases/download/v0.11.1/postgres_exporter-0.11.1.linux-amd64.tar.gz; \
curl -fSL -O https://github.com/percona/mongodb_exporter/releases/download/v0.35.0/mongodb_exporter-0.35.0.linux-amd64.tar.gz; \
\
curl -fSL -O https://github.com/nginxinc/nginx-prometheus-exporter/releases/download/v0.11.0/nginx-prometheus-exporter_0.11.0_linux_amd64.tar.gz; \
#  #无arm
curl -fSL -O https://github.com/ribbybibby/s3_exporter/releases/download/v0.5.0/s3_exporter_0.5.0_linux_amd64.tar.gz; \
\
curl -fSL -O https://github.com/prometheus/snmp_exporter/releases/download/v0.20.0/snmp_exporter-0.20.0.linux-amd64.tar.gz; \
curl -fSL -O https://github.com/prometheus/statsd_exporter/releases/download/v0.22.8/statsd_exporter-0.22.8.linux-amd64.tar.gz; \
curl -fSL -O https://github.com/ncabatoff/process-exporter/releases/download/v0.7.10/process-exporter-0.7.10.linux-amd64.tar.gz; \
curl -fSL -O https://github.com/prometheus-community/systemd_exporter/releases/download/v0.5.0/systemd_exporter-0.5.0.linux-amd64.tar.gz; \
ls -lh; 

# ext1: +ch
WORKDIR /download/oth
RUN \
curl -fSL -O https://gitee.com/g-mids/fk-clickhouse_exporter/releases/download/v22.12.13/clickhouse_exporter.tar.gz; \
ls -lh; 


# TOOLS
WORKDIR /download/tool
RUN \
curl -fSL -O https://github.com/itchyny/gojq/releases/download/v0.12.9/gojq_v0.12.9_linux_amd64.tar.gz; \
curl -fSL -O https://github.com/mikefarah/yq/releases/download/v4.30.4/yq_linux_amd64.tar.gz; \
ls -lh; 


