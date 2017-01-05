all: go/bin/server static local plugins/*

GOSRCS=go/src/server/main.go $(wildcard go/src/server/**/*.go) $(wildcard local/go/src/*.go) $(wildcard plugins/*/go/src/**/*.go)
space=$(eval) $(eval)
GOPATHS=$(subst $(space),:,$(abspath go local/go $(wildcard plugins/*/go)))
go/bin/server: $(GOSRCS)
	GOPATH="$$GOPATH:$(GOPATHS)" go install -v server

.PHONY: static
static:
	mkdir -p static
	python compile_templates.py
	find static/ -regex ".*\.\(html\|css\|js\)" | xargs gzip -fk9

include local/Makefile
include $(wildcard plugins/*/Makefile)
