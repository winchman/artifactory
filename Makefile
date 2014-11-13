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

coverage:
	curl -sL https://raw.githubusercontent.com/rafecolton/fmtpolice/master/coverage -o $@

.PHONY: .coverage
.coverage: coverage
	go get -u code.google.com/p/go.tools/cmd/cover || go get -u golang.org/x/tools/cmd/cover
	go get -u github.com/axw/gocov/gocov
	bash coverage

.PHONY: goveralls
goveralls: .coverage
	go get -u github.com/mattn/goveralls
	@echo "goveralls -coverprofile=gover.coverprofile -repotoken <redacted>"
	@goveralls -coverprofile=gover.coverprofile -repotoken $(GOVERALLS_TOKEN)
