#!/bin/bash
cd ..
protoc ./protobuf/products.proto --go_out=plugins=grpc:.
