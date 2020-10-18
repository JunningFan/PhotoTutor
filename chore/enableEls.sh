#!/bin/bash 
# echo "http://${1}/picture" 

curl -X PUT "${1}/picture" \
-H 'Content-Type: application/json' \
-d '
{
  "mappings": {
    "properties": {
      "GeoHash": {
        "type": "geo_point"
      }
    }
  }
}'