#!/bin/bash
protoc-linux/bin/protoc --go_out=paths=source_relative:. -I. cache_old.proto
