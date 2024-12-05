#!/bin/bash
num=${1:-"1"}

curl -s "http://127.0.0.1:10001/blocks/${num}" | jq .id
sudo grep -nr "new block packed" ./data/node* | grep "#$num"
