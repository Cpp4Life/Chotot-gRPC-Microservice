#! /bin/bash

protoc --go_out=. --go_opt=module=Chotot-Microservice --go-grpc_out=. --go-grpc_opt=module=Chotot-Microservice cmd/proto/*.proto