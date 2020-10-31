#!/bin/bash 
docker-compose up -d --build
./chore/enableEls.sh http://localhost:8000/els
