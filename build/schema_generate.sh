#!/bin/bash
cd ..
rm ./graph/*
go run github.com/99designs/gqlgen generate
