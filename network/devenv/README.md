Shell Scripts
=============
All scripts under [bin] sub folder to be executed in VM/Shell

bin/init.sh         Initializes the dev environment
bin/launch.sh       Launches the dev enviornment
bin/stop.sh         Stops the dev environment

Scripts
=======
All scripts under [scripts] folder to be executed in tools/shell

Start & Validate
================
bin
bin/launch.sh restart




peer channel fetch config dfarmchannel.block -o $ORDERER_ADDRESS -c dfarmchannel
peer channel join -o $ORDERER_ADDRESS -b  dfarmchannel.block


Dev Setup
=========
Container#1   orderer.dfarmadmin.com
    - Type solo
Container#2   dfarmadmin-peer1.dfarmadmin.com
Container#3   dfarmretail-peer1.dfarmretail.com
Container#4   tools

Remove all images
=================
docker rmi  $(docker images -a -q)

export CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
export  CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=net_dfarm


Postgres Service
================
sudo service postgresql stop