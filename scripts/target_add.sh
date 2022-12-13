#!/bin/bash
function errExit(){
    echo "err: $1"
    exit 0
}

target=172.25.23.199:8089
test -z "$target" && errExit "target 为空"
echo "Query parameter: target=$target"
tg=$(echo $target |sed "s/:/_/g")

dest="/tmp/target-$tg.yml"
cat > $dest<<EOF
- targets: ['$target']
  labels:
    app:    'Merger-$target'
EOF
ls -lh /tmp |grep target
cat $dest