all: go/bin/server static local

GOSRCS=go/src/server/main.go $(wildcard go/src/server/**/*.go)
go/bin/server: $(GOSRCS)
	GOPATH="$$GOPATH:$$(pwd)/go" go install -v server

.PHONY: static
static:
	mkdir -p static
	python compile_templates.py
	find static/ -regex ".*\.\(html\|css\|js\)" | xargs gzip -fk9

include local/Makefile
