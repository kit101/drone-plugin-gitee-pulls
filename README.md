# drone-plugin-gitee-pulls

[![Build Status](http://cloud.drone.io/api/badges/kit101/drone-plugin-gitee-pulls/status.svg)](http://cloud.drone.io/kit101/drone-plugin-gitee-pulls)
[![Join the discussion at https://discourse.drone.io](https://img.shields.io/badge/discourse-forum-orange.svg)](https://discourse.drone.io)
[![Go Doc](https://godoc.org/github.com/kit101/drone-plugin-gitee-pulls?status.svg)](https://godoc.org/github.com/kit101/drone-plugin-gitee-pulls)
[![Go Report](https://goreportcard.com/badge/github.com/kit101/drone-plugin-gitee-pulls)](https://goreportcard.com/report/github.com/kit101/drone-plugin-gitee-pulls)

[中文文档](./README.zh_CN.md)

<div style="display: none;">
Drone plugin to create comment and label in PR to Gitee.

For the usage information and a listing of the available options please take a look
at [the docs](http://plugins.drone.io/kit101/drone-plugin-gitee-pulls/).
</div>

## Build

Build the binary with the following command:

```console
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
export GO111MODULE=on

go build -v -a -tags netgo -o release/linux/amd64/drone-plugin-gitee-pulls
```

## Docker

Build the Docker image with the following command:

```console
docker build \
  --label org.label-schema.build-date=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --label org.label-schema.vcs-ref=$(git rev-parse --short HEAD) \
  --file docker/Dockerfile.linux.amd64 --tag kit101z/drone-plugin-gitee-pulls .
```

## Using in docker

```console
docker run --rm \
  -e PLUGIN_ACCESS_TOKEN=your-access-token \
  -e PLUGIN_IS_RUNNING=true \
  -e DRONE_SYSTEM_HOST=your.drone.host \
  -e DRONE_SYSTEM_PROTO=https \
  -e DRONE_REPO=kit101/demo1 \
  -e DRONE_PULL_REQUEST=11 \
  -e DRONE_BUILD_LINK=https://your.drone.host/api/badges/kit101/demo1/status.svg?ref\=refs/pull/11/head \
  -e DRONE_STAGE_STATUS=success \
  -e DRONE_COMMIT_REF=refs/pull/11/head \
  kit101z/drone-plugin-gitee-pulls
```

## Using in drone

### Parameter Reference

**debug**:
enable debug mode, default: `false`.

**api_server**:
the gitee api server url, default: `https://gitee.com/api/v5`.

**access_token**: 
gitee access token, you can generate personal access token.

**is_running**:
is the build running , default: `false`.

**comment_disabled**:
disable automatic updating of the comment with build status, default: `false`.

**label_disabled**:
disable automatic updating of the label with build status, default: `false`.

**running_label**:
set the name and color of the running label, default: `drone-build/running,E6A23C`.

**success_label**:
set the name and color of the success label, default: `drone-build/success,67C23A`.

**failure_label**:
set the name and color of the failure label, default: `drone-build/failure,DB2828`.

**test_disabled**:
disable automatic updating of the test status, default: `false`

### Example
```yaml
---
name: default
kind: pipeline
type: docker

#label has default values
#environment:
#  PLUGIN_GITEE_RUNNING_LABEL: drone-build/running,E6A23C
#  PLUGIN_GITEE_SUCCESS_LABEL: drone-build/success,67C23A
#  PLUGIN_GITEE_FAILURE_LABEL: drone-build/failure,DB2828

steps:
  - name: pr-enhance/start
    pull: always
    image: kit101z/drone-plugin-gitee-pulls
    settings:
      # should set `is_running: true` in the first step
      is_running: true
      access_token:
        from_secret: GITEE_ACCESS_TOKEN
    when:
      event:
        - pull_request

  - name: env
    image: alpine
    commands:
      - env

  - name: pr-enhance/end
    image: kit101z/drone-plugin-gitee-pulls
    settings:
      access_token:
        from_secret: GITEE_ACCESS_TOKEN
    when:
      event:
        - pull_request
      status:
        - failure
        - success
```

## Screen Shot

![pulls page](./docs/img/pulls.jpg)

![drone is https](./docs/img/https.jpg)

![pulls is http](./docs/img/http.jpg)

## FAQ

_Create an issue and ask questions_