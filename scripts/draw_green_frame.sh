#!/bin/bash
curl -X POST -d "white" http://localhost:17000
curl -X POST -d "bgrect 0 0 800 800" http://localhost:17000
curl -X POST -d "green" http://localhost:17000
curl -X POST -d "update" http://localhost:17000
