#!/bin/bash

docker build -t tpsi25_blog .

docker run -d -p 8081:8081 -p 8082:8082 --network bridge --name tpsi25_blog tpsi25_blog