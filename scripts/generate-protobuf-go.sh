#!/bin/bash

protoc --go_out=. --go-grpc_out=. --proto_path=plugins/ plugins/proto/*.proto
