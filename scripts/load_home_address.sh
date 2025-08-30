#!/bin/bash

db_file="./config/dbConfig.yaml"
rabbit_file="./config/rabbitConfig.yaml"

yq eval --inplace ".dbHost = \"127.0.0.1\"" $db_file
yq eval --inplace ".rabbit_host = \"127.0.0.1\"" $rabbit_file