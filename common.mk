MAKEFLAGS += --warn-undefined-variables --no-builtin-rules

SHELL := /usr/bin/env bash
.SHELLFLAGS := --norc -uo pipefail -c
.DELETE_ON_ERROR:
.SUFFIXES:

RISCV_PREFIX ?= riscv64-elf-
