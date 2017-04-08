PLUGIN_TARGETS=$(wildcard local/plugins/*)
ifneq ($(PLUGIN_TARGETS),)
include local/plugins/*/Makefile
endif

# This rule will be executed whenever the global "make" is run.
local: $(PLUGIN_TARGETS)
