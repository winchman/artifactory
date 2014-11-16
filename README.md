artifactory
===========

[![Build Status](https://drone.io/github.com/sylphon/artifactory/status.png)](https://drone.io/github.com/sylphon/artifactory/latest)
[![Build Status](https://travis-ci.org/sylphon/artifactory.svg?branch=master)](https://travis-ci.org/sylphon/artifactory)
[![GoDoc](https://godoc.org/github.com/sylphon/artifactory?status.png)](https://godoc.org/github.com/sylphon/artifactory)
[![Coverage Status](https://img.shields.io/coveralls/sylphon/artifactory.svg)](https://coveralls.io/r/sylphon/artifactory?branch=master)

Pluck and store artifacts from Docker images

## Usage

See [\_examples/](./_examples) and [integration tests](./rw_artifactory_test.go)

## Testing

```bash
# get dependencies
make get

# run tests
make test
```

## Integration Tests

Integration tests require a local docker daemon.  You may also need to
run `docker pull quay.io/rafecolton/docker-builder:latest` first.  To
run full test suite:

```bash
make integration
```
