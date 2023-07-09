#!/bin/sh
set -eux
gofmt -s -w .
~/prog/go/bin/goimports -w .
~/prog/go/bin/easyjson -all ./api/types.go
go run ./cmd/futbin