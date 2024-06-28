#!/bin/bash
for ((i=0; i<=800; i+=10)); do
  curl -X POST -d "move $i $i" http://localhost:17000
  curl -X POST -d "update" http://localhost:17000
  sleep 1
done
