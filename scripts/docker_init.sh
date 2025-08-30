#!/bin/bash

docker build -t lenic .

docker run -d -p 8081:8081 -p 8082:8082 --network bridge --name lenic lenic