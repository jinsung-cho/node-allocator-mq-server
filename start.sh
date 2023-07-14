#!/bin/sh
export $(grep -v '^#' .env | xargs)

# Kill processes
fuser -k $GO_SERVER_PORT/tcp
fuser -k $PYTHON_SERVER_PORT/tcp

# RUN rabbitmq
docker-compose -f ./rabbitmq/docker-compose.yaml up -d

# RUN python flask server for argo workflow run
pip3 install -r requirements.txt
python3 argo_request_server.py &

# RUN node allocator
cd nodeAllocator/goVersion 
go run . &

# RUN go backend server
cd ../..
go run . &