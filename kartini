#!/bin/sh

KARTINI_ROOT="$(pwd)/container" \
KARTINI_PATH="$KARTINI_ROOT/var/db/kartini/repo" \
KARTINI_CACHE="$KARTINI_ROOT/var/cache/kartini" \
go run toml.go system.go helper.go repokeeper.go "$@"