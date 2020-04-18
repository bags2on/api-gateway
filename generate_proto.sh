#!/bin/bash
protoc ./protobuf/products.proto --go_out=plugins=grpc:.
