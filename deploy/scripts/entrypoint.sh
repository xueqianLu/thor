#! /bin/sh

# exit script on any error
set -e

echo "start thor with beneficiary $BENEFICIARY and bootnode $BOOTNODE"
/usr/bin/thor --data-dir /root/node --network /root/genesis.json --beneficiary ${BENEFICIARY} --api-addr "0.0.0.0:8669" --bootnode ${BOOTNODE} >> /root/node/node.log 2>&1
