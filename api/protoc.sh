#!/bin/bash

shopt -s globstar
protoc -I ./ --go_opt=paths=source_relative --go_out=plugins=grpc:. **/*.proto
