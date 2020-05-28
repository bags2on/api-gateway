#!/bin/bash
cd ..
rm ./graph/*
protoc ./protobuf/products.proto --go_out=plugins=grpc:.
