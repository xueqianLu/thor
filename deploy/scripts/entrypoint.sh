#! /bin/sh

# exit script on any error
set -e

echo "runing thor"
/usr/bin/thor --data-dir /root/node --network /root/node/genesis.json >> /root/node/node.log 2>&1
