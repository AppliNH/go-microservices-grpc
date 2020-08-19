#!/bin/bash

docker volume create mongodbdata
docker run -p 27017:27017 -v mongodbdata:/data/db --name mongo mongo