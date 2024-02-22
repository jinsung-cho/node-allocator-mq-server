#!/bin/sh
export $(grep -v '^#' .env | xargs)

# Kill processes
fuser -k $GO_SERVER_PORT/tcp
fuser -k $PYTHON_SERVER_PORT/tcp

# RUN rabbitmq
docker-compose -f ./rabbitmq-docker/docker-compose.yaml up -d

# RUN python flask server for argo workflow run
cd ./argoIntegrationServer
pip3 install -r requirements.txt
python3 main.py &  
cd ..

pwd
# RUN node allocator(go)
cd nodeAllocator/golang 
go run . & 
cd ../..

# RUN node allocator(python)
cd nodeAllocator/python
pip3 install -r requirements.txt
python3 main.py &

# RUN go backend server
cd ../..
go run . 