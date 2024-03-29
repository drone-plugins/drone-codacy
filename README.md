# drone-codacy

[![Build Status](http://cloud.drone.io/api/badges/drone-plugins/drone-codacy/status.svg)](http://cloud.drone.io/drone-plugins/drone-codacy)
[![Gitter chat](https://badges.gitter.im/drone/drone.png)](https://gitter.im/drone/drone)
[![Join the discussion at https://discourse.drone.io](https://img.shields.io/badge/discourse-forum-orange.svg)](https://discourse.drone.io)
[![Drone questions at https://stackoverflow.com](https://img.shields.io/badge/drone-stackoverflow-orange.svg)](https://stackoverflow.com/questions/tagged/drone.io)
[![](https://images.microbadger.com/badges/image/plugins/codacy.svg)](https://microbadger.com/images/plugins/codacy "Get your own image badge on microbadger.com")
[![Go Doc](https://godoc.org/github.com/drone-plugins/drone-codacy?status.svg)](http://godoc.org/github.com/drone-plugins/drone-codacy)
[![Go Report](https://goreportcard.com/badge/github.com/drone-plugins/drone-codacy)](https://goreportcard.com/report/github.com/drone-plugins/drone-codacy)

Drone plugin to send coverage reports to [Codacy](https://www.codacy.com). For the usage information and a listing of the available options please take a look at [the docs](http://plugins.drone.io/drone-plugins/drone-codacy/).

## Build

Build the binary with the following command:

```bash
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
export GO111MODULE=on

go build -v -a -tags netgo -o release/linux/amd64/drone-codacy
```

## Docker

Build the Docker image with the following command:

```bash
docker build \
  --label org.label-schema.build-date=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --label org.label-schema.vcs-ref=$(git rev-parse --short HEAD) \
  --file docker/Dockerfile.linux.amd64 --tag plugins/codacy .
```

### Usage

```bash
docker run --rm \
  -e PLUGIN_TOKEN=xxx \
  -e PLUGIN_PATTERN="coverage.out" \
  -e DRONE_COMMIT_SHA=7fd1a60b01f91b314f59955a4e4d4e80d8edf11d \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  plugins/codacy
```
