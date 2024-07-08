#!/bin/bash
nodecnt=${1:-"6"}
hacknodecnt=${2:-"2"}


let nodeidx=0

composefile=docker-compose-${nodecnt}.yml

function addHeader() {
    echo 'version: "3.9"' > $composefile
    echo "services:" >> $composefile
}

function addBootNode() {
# add bootnode to compose file
echo "bootnode:" >> $composefile
echo "  image: thor:base" >> $composefile
echo "  container_name: thor-bootnode" >> $composefile
echo "  entrypoint: /usr/bin/thor --data-dir /root/node --network /root/genesis.json" >> $composefile
echo "  ports:" >> $composefile
echo "    - \"1000:8669\"" >> $composefile
echo "    - \"2000:11235\"" >> $composefile
echo "  volumes:" >> $composefile
echo "    - ./config/bootnode:/root/.org.vechain.thor" >> $composefile
echo "    - ./config/genesis.json:/root/genesis.json" >> $composefile
echo "    - ./data/bootnode:/root/node" >> $composefile
echo "  deploy:" >> $composefile
echo "    restart_policy:" >> $composefile
echo "      condition: on-failure" >> $composefile
echo "      delay: 15s" >> $composefile
echo "      max_attempts: 100" >> $composefile
echo "      window: 120s" >> $composefile
echo "  networks:" >> $composefile
echo "    thor-testnet:" >> $composefile
echo "      ipv4_address: 172.99.1.1" >> $composefile

nodeidx+=1
}


function addHackCenter() {
# add vecenter to compose file
echo "vecenter:" >> $composefile
echo "  image: thor:base" >> $composefile
echo "  container_name: thor-vecenter" >> $composefile
echo "  entrypoint: /usr/bin/vecenter -c $hacknodecnt -begin 500 -port 9000" >> $composefile
echo "  ports:" >> $composefile
echo "    - \"9000:9000\"" >> $composefile
echo "  deploy:" >> $composefile
echo "    restart_policy:" >> $composefile
echo "      condition: on-failure" >> $composefile
echo "      delay: 15s" >> $composefile
echo "      max_attempts: 100" >> $composefile
echo "      window: 120s" >> $composefile
echo "  networks:" >> $composefile
echo "    thor-testnet:" >> $composefile
echo "      ipv4_address: 172.99.1.2" >> $composefile

nodeidx+=1
}


function addNormalNode() {
    echo "node$i:" >> $composefile
    echo "image: thor:base" >> $composefile
    echo "container_name: thor-node-$i" >> $composefile
    echo "environment:" >> $composefile
    echo "  - BENEFICIARY=0x$(printf "%040d" $i)" >> $composefile
    echo "  - BOOTNODE=enode://bc18b2d7dd0daf50073f53f5c8e7aecb41387275efb5fd0e41ec3b87ce2804353692c38a9774777ce39ba0de61648cd7adc70d3fc29692b46c5f520f542a7824@172.99.1.1:11235" >> $composefile
    echo "  - ACCOUNT_IDX=$i" >> $composefile
    echo "ports:" >> $composefile
    echo "  - \"$(($i+10000)):8669\"" >> $composefile
    echo "  - \"$(($i+20000)):11235\"" >> $composefile
    echo "volumes:" >> $composefile
    echo "  - ./config/keys/master.key.$i:/root/.org.vechain.thor" >> $composefile
    echo "  - ./config/genesis.json:/root/genesis.json" >> $composefile
    echo "  - ./config/account.json:/root/account.json" >> $composefile
    echo "  - ./data/node$i:/root/node" >> $composefile
    echo "depends_on:" >> $composefile
    echo "  - bootnode" >> $composefile
    echo "deploy:" >> $composefile
    echo "  restart_policy:" >> $composefile
    echo "    condition: on-failure" >> $composefile
    echo "    delay: 15s" >> $composefile
    echo "    max_attempts: 100" >> $composefile
    echo "    window: 120s" >> $composefile
    echo "networks:" >> $composefile
    echo "  thor-testnet:" >> $composefile
    echo "    ipv4_address: 172.99.1.$(($i+2))" >> $composefile
}




for i in $(seq 0 $nodecnt)
do
  echo "node$i:" >> $composefile
  echo "image: thor:base" >> $composefile
  echo "container_name: thor-node-$i" >> $composefile
  echo "environment:" >> $composefile
  echo "  - BENEFICIARY=0x$(printf "%040d" $i)" >> $composefile
  echo "  - BOOTNODE=enode://bc18b2d7dd0daf50073f53f5c8e7aecb41387275efb5fd0e41ec3b87ce2804353692c38a9774777ce39ba0de61648cd7adc70d3fc29692b46c5f520f542a7824@172.99.1.1:11235" >> $composefile
  echo "  - ACCOUNT_IDX=$i" >> $composefile
  echo "ports:" >> $composefile
  echo "  - \"$(($i+10000)):8669\"" >> $composefile
  echo "  - \"$(($i+20000)):11235\"" >> $composefile
  echo "volumes:" >> $composefile
  echo "  - ./config/keys/master.key.$i:/root/.org.vechain.thor" >> $composefile
  echo "  - ./config/genesis.json:/root/genesis.json" >> $composefile
  echo "  - ./config/account.json:/root/account.json" >> $composefile
  echo "  - ./data/node$i:/root/node" >> $composefile
  echo "depends_on:" >> $composefile
  echo "  - bootnode" >> $composefile
  echo "deploy:" >> $composefile
  echo "  restart_policy:" >> $composefile
  echo "    condition: on-failure" >> $composefile
  echo "    delay: 15s" >> $composefile
  echo "    max_attempts: 100" >> $composefile
  echo "    window: 120s" >> $composefile
  echo "networks:" >> $composefile
  echo "  thor-testnet:" >> $composefile
  echo "    ipv4_address: 172.99.1.$(($i+2))" >> $composefile
done
