#!/bin/sh
#export IP=$(ifconfig en0 | grep inet | awk '$1=="inet" {print $2}')
export IP=$(ifconfig bridge100 | grep inet | awk '$1=="inet" {print $2}')

docker-compose up
