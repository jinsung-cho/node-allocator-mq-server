#!/bin/sh

# RUN python flask server for argo workflow run
python3 argo_request_server.py &

# RUN node allocator
cd nodeAllocator/goVersion 
go run . &

# RUN go backend server
cd ../..
go run . &