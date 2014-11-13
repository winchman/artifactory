SHELL := /bin/bash

.PHONY: test
test: .fmtpolice
	go test ./...

fmtpolice:
	curl -sL https://raw.githubusercontent.com/rafecolton/fmtpolice/master/fmtpolice -o $@

.PHONY:
.fmtpolice: fmtpolice
	bash fmtpolice
	@find . -type f -name '*.test' -exec rm {} \;

.PHONY: get
	go get -d -t ./...
