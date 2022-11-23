# source /etc/profile #source err with github's flow env
export |grep DOCKER_REG
repo=registry.cn-shenzhen.aliyuncs.com
echo "${DOCKER_REGISTRY_PW_infrastSubUser2}" |docker login --username=${DOCKER_REGISTRY_USER_infrastSubUser2} --password-stdin $repo

# repoHub=docker.io
# echo "${DOCKER_REGISTRY_PW_dockerhub}" |docker login --username=${DOCKER_REGISTRY_USER_dockerhub} --password-stdin $repoHub
        
ns=infrastlabs
ver=v5

img=$repo/$ns/prom-exporters-down:latest
docker build -t $img .
docker push $img
