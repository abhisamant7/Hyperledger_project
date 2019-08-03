#!/bin/bash

docker-compose down

# REMOVE the dev- container images also - TBD
docker rm $(docker ps -a -q)            &> /dev/null
docker rmi $(docker images dev-* -q)    &> /dev/null
sudo rm -rf $HOME/ledgers/ca &> /dev/null

docker-compose up -d

SLEEP_TIME=3s
echo    '========= Submitting txn for channel creation as DfarmadminAdmin ============'
export CHANNEL_TX_FILE=./config/dfarm-channel.tx
export ORDERER_ADDRESS=orderer.dfarmadmin.com:7050
# export FABRIC_LOGGING_SPEC=DEBUG
export CORE_PEER_LOCALMSPID=DfarmadminMSP
export CORE_PEER_MSPCONFIGPATH=$PWD/client/dfarmadmin/admin/msp
export CORE_PEER_ADDRESS=dfarmadmin-peer1.dfarmadmin.com:7051
peer channel create -o $ORDERER_ADDRESS -c dfarmchannel -f ./dfarmchannel.tx

echo    '========= Joining the dfarmadmin-peer1 to Dfarm channel ============'
DFARM_CHANNEL_BLOCK=./dfarmchannel.block
export CORE_PEER_ADDRESS=dfarmadmin-peer1.dfarmadmin.com:7051
peer channel join -o $ORDERER_ADDRESS -b $DFARM_CHANNEL_BLOCK
# Update anchor peer on channel for dfarmadmin
# sleep  3s
sleep $SLEEP_TIME
ANCHOR_UPDATE_TX=./config/dfarm-anchor-update-dfarmadmin.tx
peer channel update -o $ORDERER_ADDRESS -c dfarmchannel -f $ANCHOR_UPDATE_TX

echo    '========= Joining the dfarmretail-peer1 to Dfarm channel ============'
# peer channel fetch config $DFARM_CHANNEL_BLOCK -o $ORDERER_ADDRESS -c dfarmchannel
export CORE_PEER_LOCALMSPID=DfarmretailMSP
ORG_NAME=dfarmretail.com
export CORE_PEER_ADDRESS=dfarmretail-peer1.dfarmretail.com:8051
export CORE_PEER_MSPCONFIGPATH=$PWD/client/dfarmretail/admin/msp
peer channel join -o $ORDERER_ADDRESS -b $DFARM_CHANNEL_BLOCK
# Update anchor peer on channel for dfarmretail
sleep  $SLEEP_TIME
ANCHOR_UPDATE_TX=./config/dfarm-anchor-update-dfarmretail.tx
peer channel update -o $ORDERER_ADDRESS -c dfarmchannel -f $ANCHOR_UPDATE_TX

