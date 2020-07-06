#!/bin/bash
protoc gRPCProto/gRPCDemo.proto --go_out=plugins=grpc:./gRPCProto