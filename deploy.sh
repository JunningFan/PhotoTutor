#!/bin/bash 
mkdir els_data 
chmod 777 els_data
docker-compose up -d --build
./chore/enableEls.sh http://localhost:8000/els
cp avatar.jpg img/small/avatar.jpg