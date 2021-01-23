#/bin/sh
set -x

CI_COMMIT_TAG=$(git describe --always --tags)

docker build -t linclaus/grafana-operator:latest -f build/package/Dockerfile .