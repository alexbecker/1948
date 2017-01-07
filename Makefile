all: go/bin/server static local plugins/*

STATIC_DEPS=

COMPRESSIBLE=$(wildcard static/*.html) $(wildcard static/*.css) $(wildcard static/*.js) $(wildcard static/**/*.html) $(wildcard static/**/*.css) $(wildcard static/**/*.js)
.PHONY: gz
gz: $(addsuffix .gz,$(COMPRESSIBLE))

%.gz: %
	gzip -fk9 $<

include $(wildcard plugins/*/Makefile)
include local/Makefile

GOSRCS=go/src/server/main.go $(wildcard go/src/server/**/*.go) $(wildcard local/go/src/*.go) $(wildcard plugins/*/go/src/**/*.go)
space=$(eval) $(eval)
GOPATHS=$(subst $(space),:,$(abspath go local/go $(wildcard plugins/*/go)))
go/bin/server: $(GOSRCS)
	GOPATH="$$GOPATH:$(GOPATHS)" go install -v server

.PHONY: static
static: $(STATIC_DEPS)
	mkdir -p static
	python compile_templates.py
