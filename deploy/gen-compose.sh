#!/bin/bash
nodecnt=${1:-"7"}
hacknodecnt=${2:-"2"}
testtps=${3:-"0"}

if [ $testtps -eq 0 ]; then
  echo "not test testtps"
fi
if [ $testtps -eq 1 ]; then
  echo "test testtps"
fi

composefile=docker-compose.yml

function addHeader() {
    echo 'version: "3.9"' > $composefile
    echo "" >> $composefile
    echo "services:" >> $composefile
}

function addBootNode() {
# add bootnode to compose file
echo "  bootnode:" >> $composefile
echo "    image: thor:latest" >> $composefile
echo "    container_name: thor-bootnode" >> $composefile
echo "    entrypoint: /usr/bin/thor --data-dir /root/node --network /root/genesis.json" >> $composefile
echo "    ports:" >> $composefile
echo "      - \"1000:8669\"" >> $composefile
echo "      - \"2000:11235\"" >> $composefile
echo "    volumes:" >> $composefile
echo "      - ./config/bootnode:/root/.org.vechain.thor" >> $composefile
echo "      - ./config/genesis.json:/root/genesis.json" >> $composefile
echo "      - ./data/bootnode:/root/node" >> $composefile
echo "    deploy:" >> $composefile
echo "      restart_policy:" >> $composefile
echo "        condition: on-failure" >> $composefile
echo "        delay: 15s" >> $composefile
echo "        max_attempts: 100" >> $composefile
echo "        window: 120s" >> $composefile
echo "    networks:" >> $composefile
echo "      thor-testnet:" >> $composefile
echo "        ipv4_address: 172.99.1.1" >> $composefile
}


function addHackCenter() {
# add vecenter to compose file
echo "  vecenter:" >> $composefile
echo "    image: vecenter:latest" >> $composefile
echo "    container_name: thor-vecenter" >> $composefile
echo "    entrypoint: /usr/bin/vecenter -c $hacknodecnt -begin 50 -port 9000" >> $composefile
echo "    ports:" >> $composefile
echo "      - \"9000:9000\"" >> $composefile
echo "    deploy:" >> $composefile
echo "      restart_policy:" >> $composefile
echo "        condition: on-failure" >> $composefile
echo "        delay: 15s" >> $composefile
echo "        max_attempts: 100" >> $composefile
echo "        window: 120s" >> $composefile
echo "    networks:" >> $composefile
echo "      thor-testnet:" >> $composefile
echo "        ipv4_address: 172.99.1.2" >> $composefile
}


function addNormalNode() {
    local i=$1
    echo "  node$i:" >> $composefile
    echo "    image: thor:latest" >> $composefile
    echo "    container_name: thor-node-$i" >> $composefile
    echo "    environment:" >> $composefile
    echo "      - BENEFICIARY=0x$(printf "%040d" $(($i+10)))" >> $composefile
    echo "      - BOOTNODE=enode://bc18b2d7dd0daf50073f53f5c8e7aecb41387275efb5fd0e41ec3b87ce2804353692c38a9774777ce39ba0de61648cd7adc70d3fc29692b46c5f520f542a7824@172.99.1.1:11235" >> $composefile
    echo "      - ACCOUNT_IDX=$i" >> $composefile
    echo "      - VE_P2P_SERVER_URL=vecenter:9000" >> $composefile
    echo "    ports:" >> $composefile
    echo "      - \"$(($i+10000)):8669\"" >> $composefile
    echo "      - \"$(($i+20000)):11235\"" >> $composefile
    echo "    volumes:" >> $composefile
    echo "      - ./config/keys/master.key.$i:/root/.org.vechain.thor/master.key" >> $composefile
    echo "      - ./config/genesis.json:/root/genesis.json" >> $composefile
    if [ $testtps -eq 0 ]; then
          echo "      - ./config/accounts.json:/root/account.json" >> $composefile
    fi
    echo "      - ./data/node$i:/root/node" >> $composefile
    echo "    depends_on:" >> $composefile
    echo "      - bootnode" >> $composefile
    echo "      - vecenter" >> $composefile
    echo "    deploy:" >> $composefile
    echo "      restart_policy:" >> $composefile
    echo "        condition: on-failure" >> $composefile
    echo "        delay: 15s" >> $composefile
    echo "        max_attempts: 100" >> $composefile
    echo "        window: 120s" >> $composefile
    echo "    networks:" >> $composefile
    echo "      thor-testnet:" >> $composefile
    echo "        ipv4_address: 172.99.1.$(($i+3))" >> $composefile
}


function addHackNode() {
    local i=$1
    local hackidx=$2
    echo "  node$i:" >> $composefile
    echo "    image: thor:latest" >> $composefile
    echo "    container_name: thor-node-hack-$i" >> $composefile
    echo "    environment:" >> $composefile
    echo "      - BENEFICIARY=0x$(printf "%040d" $(($i+10)))" >> $composefile
    echo "      - BOOTNODE=enode://bc18b2d7dd0daf50073f53f5c8e7aecb41387275efb5fd0e41ec3b87ce2804353692c38a9774777ce39ba0de61648cd7adc70d3fc29692b46c5f520f542a7824@172.99.1.1:11235" >> $composefile
    echo "      - ACCOUNT_IDX=$i" >> $composefile
    echo "      - VE_HACK_SERVER_URL=vecenter:9000" >> $composefile
    echo "      - VE_HACK_CLIENT_INDEX=$hackidx" >> $composefile
    echo "    ports:" >> $composefile
    echo "      - \"$(($i+10000)):8669\"" >> $composefile
    echo "      - \"$(($i+20000)):11235\"" >> $composefile
    echo "    volumes:" >> $composefile
    echo "      - ./config/keys/master.key.$i:/root/.org.vechain.thor/master.key" >> $composefile
    echo "      - ./config/genesis.json:/root/genesis.json" >> $composefile
    if [ $testtps -eq 0 ]; then
      echo "      - ./config/accounts.json:/root/account.json" >> $composefile
    fi
    echo "      - ./data/node$i:/root/node" >> $composefile
    echo "    depends_on:" >> $composefile
    echo "      - bootnode" >> $composefile
    echo "      - vecenter" >> $composefile
    echo "    deploy:" >> $composefile
    echo "      restart_policy:" >> $composefile
    echo "        condition: on-failure" >> $composefile
    echo "        delay: 15s" >> $composefile
    echo "        max_attempts: 100" >> $composefile
    echo "        window: 120s" >> $composefile
    echo "    networks:" >> $composefile
    echo "      thor-testnet:" >> $composefile
    echo "        ipv4_address: 172.99.1.$(($i+3))" >> $composefile
}

function addTxPress() {
echo "  txpress:" >> $composefile
echo "    image: tscel/txpress-vechain:0730" >> $composefile
echo "    container_name: thor-txpress" >> $composefile
echo "    entrypoint: /usr/bin/txpress --start --log /root/data/press.log" >> $composefile
echo "    volumes:" >> $composefile
echo "      - ./config/press-app.json:/root/app.json" >> $composefile
echo "      - ./config/accounts.json:/root/accounts.json" >> $composefile
echo "      - ./data/txpress:/root/data" >> $composefile
echo "    deploy:" >> $composefile
echo "      restart_policy:" >> $composefile
echo "        condition: on-failure" >> $composefile
echo "        delay: 15s" >> $composefile
echo "        max_attempts: 100" >> $composefile
echo "        window: 120s" >> $composefile
echo "    networks:" >> $composefile
echo "      - thor-testnet" >> $composefile
echo "    depends_on:" >> $composefile

for i in $(seq 0 $nodecnt)
do
    if [ $i -eq $nodecnt ]; then
      break
    fi
    echo "      - node$i" >> $composefile
done

}

function addTail() {
  echo "networks:" >> $composefile
    echo "  thor-testnet:" >> $composefile
    echo "    driver: bridge" >> $composefile
    echo "    ipam:" >> $composefile
    echo "      driver: default" >> $composefile
      echo "      config:" >> $composefile
      echo "        - subnet: 172.99.0.0/16" >> $composefile
}

addHeader
addBootNode
addHackCenter

let hackcnt=0
let nodeidx=0

for i in $(seq 0 $nodecnt)
do
  if [ $nodeidx -eq $nodecnt ]; then
    break
  fi
  # if i <= nodecnt/2 && hack node < hackcnt, add hack node
  if [ $nodeidx -ge $(($nodecnt/3)) ] && [ $hackcnt -lt $hacknodecnt ]; then
    addHackNode $nodeidx $hackcnt
    let hackcnt++
  else
    addNormalNode $nodeidx
  fi
  let nodeidx++
done

if [ $testtps -eq 1 ]; then
  addTxPress
fi

addTail

