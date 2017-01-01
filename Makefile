COMPRESSABLE=$(wildcard static/*.html) $(wildcard static/*.css) $(wildcard static/*.js) $(wildcard static/**/*.html) $(wildcard static/**/*.css) $(wildcard static/**/*.js)
TO_GZ=$(addsuffix .gz,$(COMPRESSABLE))

all: go/bin/server static $(TO_GZ) local

GOSRCS=go/src/server/main.go $(wildcard go/src/server/**/*.go)
go/bin/server: $(GOSRCS)
	GOPATH="$$GOPATH:$$(pwd)/go" go install -v server

templates/blog/%.md.html: templates/blog/%.md
	cat $< | ~/Blogdown/Blogdown --footnote-prefix $(subst .md,,$(notdir $<)) --footnote-index-from 0 --footnote-backlinks --em-dashes > $@

.PHONY: static
static:
	mkdir -p static
	python compile_templates.py

static/%.gz: static/%
	gzip -fk9 $<

include local/Makefile
