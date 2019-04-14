SHELL := /usr/bin/env bash

SOURCES := $(wildcard *.hex)
TARGETS := $(patsubst %.hex,%.elf,$(SOURCES))

all: $(TARGETS)

%.elf: %.hex
	sed "s/#.*$$//g" $< | xxd -r -p > $@
