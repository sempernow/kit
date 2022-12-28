#!/usr/bin/env bash
#------------------------------------------------------------------------------
#  make : md2HTML.exe ...
# -----------------------------------------------------------------------------

find . -name '*.md'   |grep -v '/vendor' |grep -v '/modules' \
	|grep -v '/assets/src/content' |xargs -I{} md2HTML.exe "{}"

find . -name '*.html' |grep -v '/vendor' |grep -v '/modules' |xargs -I{} chmod 0600 "{}"

