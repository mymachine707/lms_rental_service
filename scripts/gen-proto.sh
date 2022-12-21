#!/bin/bash
CURRENT_DIR=$(pwd)
for x in $(find ${CURRENT_DIR}/lms_proto/* -type d); do
  sudo protoc --plugin="protoc-gen-go=${GOPATH}/bin/protoc-gen-go" --plugin="protoc-gen-go-grpc=${GOPATH}/bin/protoc-gen-go-grpc" -I=${x} -I=${CURRENT_DIR}/lms_proto -I /usr/local/include --go_out=${CURRENT_DIR}/protogen \
   --go-grpc_out=${CURRENT_DIR}/protogen ${x}/*.proto
done