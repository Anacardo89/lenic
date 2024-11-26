#!/bin/bash

db_name=mysql
rabbit_name=rabbitmq
network_name=tpsi25_blog_default
db_file="./config/dbConfig.yaml"
rabbit_file="./config/rabbitConfig.yaml"

db_ip=$(docker inspect $db_name | jq -r ".[0].NetworkSettings.Networks.$network_name.IPAddress")
rabbit_ip=$(docker inspect $rabbit_name | jq -r ".[0].NetworkSettings.Networks.$network_name.IPAddress")

yq eval --inplace ".dbHost = \"$db_ip\"" $db_file
yq eval --inplace ".rabbit_host = \"$rabbit_ip\"" $rabbit_file