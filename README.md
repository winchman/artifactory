artifactory
===========

[![Build Status](https://travis-ci.org/winchman/artifactory.svg?branch=master)](https://travis-ci.org/winchman/artifactory)
[![GoDoc](https://godoc.org/github.com/winchman/artifactory?status.png)](https://godoc.org/github.com/winchman/artifactory)
[![Coverage Status](https://img.shields.io/coveralls/winchman/artifactory.svg)](https://coveralls.io/r/winchman/artifactory?branch=master)

Pluck and store artifacts from Docker images

## Usage

See [\_example/](./_example) and [integration tests](./rw_artifactory_test.go)

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
