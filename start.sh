#!/bin/bash

sudo docker build -t go-rest-api docker/
sudo docker run -p 27017:27017 go-rest-api 

sudo docker run -d -it -p 27017:27017 --name=go-rest-api-1 go-rest-api /bin/sh
sudo docker exec -i -t go-rest-api-1 /bin/bash


docker-compose up -d