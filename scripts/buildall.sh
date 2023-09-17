#!/bin/bash

# Build all builtin plugins
for pluginFolder in plugins/builtin/*; do
  if [ -d "$pluginFolder" ]; then
    echo "Going to $pluginFolder"
    pushd "$pluginFolder"
      go get -d -v .
    popd
    go build -o $(basename "$pluginFolder") "$pluginFolder"/main.go
  fi
done

# Build the main application
go get -d -v .
go build -o my-helper-bot main.go
