#!/bin/bash

if [ "$1" == "skip" ]; then
    echo "Skipping Install & Instantiate"
    SKIP_INSTALL_INSTANTIATE="yes"
fi


SLEEP_TIME=3s

#1. Set the environmeent for dfarmadmin-peer1
. set-env.sh dfarmadmin
echo "ORGANIZATION_CONTEXT=$ORGANIZATION_CONTEXT"

#2. Set the chain code arguments
set-chain-env.sh   -l golang -p chaincode_example02 -n gocc -v 1 -C dfarmchannel \
                   -c '{"Args":["init","a","100","b","300"]}' -q '{"Args":["query","b"]}' -i  '{"Args":["invoke","a","b","5"]}'

#3. Install the chaincode
if [ "$SKIP_INSTALL_INSTANTIATE" != "yes" ]; then
    chain.sh   install
    #4. Instantiate the chaincode & query & invoke
    chain.sh   instantiate

    sleep  $SLEEP_TIME
fi

# Query
chain.sh query

#5. Set the environment for dfarmretail
. set-env.sh dfarmretail
echo "ORGANIZATION_CONTEXT=$ORGANIZATION_CONTEXT"

#6. Install the chaincode
if [ "$SKIP_INSTALL_INSTANTIATE" != "yes" ]; then
    chain.sh install
fi

#7. Query & Invoke
chain.sh query
chain.sh invoke
sleep  $SLEEP_TIME
chain.sh query

#8. Switch to dfarmadmin & query
. set-env.sh dfarmadmin
echo "ORGANIZATION_CONTEXT=$ORGANIZATION_CONTEXT"
chain.sh query
