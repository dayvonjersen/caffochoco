#\!/usr/bin/env bash
goimports -w *.go
go build && ./caffochoco |& pp
