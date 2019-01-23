#!/bin/bash

protoc -I backend/ backend/backend.proto --go_out=plugins=grpc:backend
