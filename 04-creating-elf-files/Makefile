SHELL := /usr/bin/env bash

elfheader.elf: elfheader.hex
	sed "s/#.*$$//g" $< | xxd -r -p > $@
