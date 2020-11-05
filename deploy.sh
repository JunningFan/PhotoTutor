#!/bin/bash 
mkdir els_data 
chmod 777 els_data
docker-compose up -d --build
cp avatar.jpg img/small/avatar.jpg
sleep 5 && ./chore/enableEls.sh http://localhost:8000/els