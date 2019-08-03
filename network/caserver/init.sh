docker-compose -f docker-compose-ca.yaml down
rm -rf ./server/*
rm -rf ./client/*
cp fabric-ca-server-config.yaml ./server
docker-compose -f docker-compose-ca.yaml up -d

sleep 3s

# Bootstrap enrollment
export FABRIC_CA_CLIENT_HOME=$PWD/client/caserver/admin
fabric-ca-client enroll -u http://admin:adminpw@localhost:7054


######################
# Admin registration #
######################
echo "Registering: dfarmadmin-admin"
ATTRIBUTES='"hf.Registrar.Roles=peer,user,client","hf.AffiliationMgr=true","hf.Revoker=true","hf.Registrar.Attributes=*"'
fabric-ca-client register --id.type client --id.name dfarmadmin-admin --id.secret adminpw --id.affiliation dfarmadmin --id.attrs $ATTRIBUTES

# 3. Register dfarmretail-admin
echo "Registering: dfarmretail-admin"
ATTRIBUTES='"hf.Registrar.Roles=peer,user,client","hf.AffiliationMgr=true","hf.Revoker=true","hf.Registrar.Attributes=*"'
fabric-ca-client register --id.type client --id.name dfarmretail-admin --id.secret adminpw --id.affiliation dfarmretail --id.attrs $ATTRIBUTES

# 4. Register orderer-admin
echo "Registering: orderer-admin"
ATTRIBUTES='"hf.Registrar.Roles=orderer"'
fabric-ca-client register --id.type client --id.name orderer-admin --id.secret adminpw --id.affiliation orderer --id.attrs $ATTRIBUTES


####################
# Admin Enrollment #
####################
export FABRIC_CA_CLIENT_HOME=$PWD/client/dfarmadmin/admin
fabric-ca-client enroll -u http://dfarmadmin-admin:adminpw@localhost:7054
mkdir -p $FABRIC_CA_CLIENT_HOME/msp/admincerts
cp $FABRIC_CA_CLIENT_HOME/../../caserver/admin/msp/signcerts/*  $FABRIC_CA_CLIENT_HOME/msp/admincerts

export FABRIC_CA_CLIENT_HOME=$PWD/client/dfarmretail/admin
fabric-ca-client enroll -u http://dfarmretail-admin:adminpw@localhost:7054
mkdir -p $FABRIC_CA_CLIENT_HOME/msp/admincerts
cp $FABRIC_CA_CLIENT_HOME/../../caserver/admin/msp/signcerts/*  $FABRIC_CA_CLIENT_HOME/msp/admincerts

export FABRIC_CA_CLIENT_HOME=$PWD/client/orderer/admin
fabric-ca-client enroll -u http://orderer-admin:adminpw@localhost:7054
mkdir -p $FABRIC_CA_CLIENT_HOME/msp/admincerts
cp $FABRIC_CA_CLIENT_HOME/../../caserver/admin/msp/signcerts/*  $FABRIC_CA_CLIENT_HOME/msp/admincerts

#################
# Org MSP Setup #
#################
# Path to the CA certificate
ROOT_CA_CERTIFICATE=./server/ca-cert.pem
mkdir -p ./client/orderer/msp/admincerts
mkdir ./client/orderer/msp/cacerts
mkdir ./client/orderer/msp/keystore
cp $ROOT_CA_CERTIFICATE ./client/orderer/msp/cacerts
cp ./client/orderer/admin/msp/signcerts/* ./client/orderer/msp/admincerts   

mkdir -p ./client/dfarmadmin/msp/admincerts
mkdir ./client/dfarmadmin/msp/cacerts
mkdir ./client/dfarmadmin/msp/keystore
cp $ROOT_CA_CERTIFICATE ./client/dfarmadmin/msp/cacerts
cp ./client/dfarmadmin/admin/msp/signcerts/* ./client/dfarmadmin/msp/admincerts   

mkdir -p ./client/dfarmretail/msp/admincerts
mkdir ./client/dfarmretail/msp/cacerts
mkdir ./client/dfarmretail/msp/keystore
cp $ROOT_CA_CERTIFICATE ./client/dfarmretail/msp/cacerts
cp ./client/dfarmretail/admin/msp/signcerts/* ./client/dfarmretail/msp/admincerts   

######################
# Orderer Enrollment #
######################
export FABRIC_CA_CLIENT_HOME=$PWD/client/orderer/admin
fabric-ca-client register --id.type orderer --id.name orderer --id.secret adminpw --id.affiliation orderer 
export FABRIC_CA_CLIENT_HOME=$PWD/client/orderer/orderer
fabric-ca-client enroll -u http://orderer:adminpw@localhost:7054
cp -a $PWD/client/orderer/admin/msp/signcerts  $FABRIC_CA_CLIENT_HOME/msp/admincerts

####################
# Peer Enrollments #
####################
export FABRIC_CA_CLIENT_HOME=$PWD/client/dfarmadmin/admin
fabric-ca-client register --id.type peer --id.name dfarmadmin-peer1 --id.secret adminpw --id.affiliation dfarmadmin 
export FABRIC_CA_CLIENT_HOME=$PWD/client/dfarmadmin/peer1
fabric-ca-client enroll -u http://dfarmadmin-peer1:adminpw@localhost:7054
cp -a $PWD/client/dfarmadmin/admin/msp/signcerts  $FABRIC_CA_CLIENT_HOME/msp/admincerts

export FABRIC_CA_CLIENT_HOME=$PWD/client/dfarmretail/admin
fabric-ca-client register --id.type peer --id.name dfarmretail-peer1 --id.secret adminpw --id.affiliation dfarmretail
export FABRIC_CA_CLIENT_HOME=$PWD/client/dfarmretail/peer1
fabric-ca-client enroll -u http://dfarmretail-peer1:adminpw@localhost:7054
cp -a $PWD/client/dfarmretail/admin/msp/signcerts  $FABRIC_CA_CLIENT_HOME/msp/admincerts


##############################
# User Enrollments Dfarmadmin only #
##############################
export FABRIC_CA_CLIENT_HOME=$PWD/client/dfarmadmin/admin
ATTRIBUTES='"hf.AffiliationMgr=false:ecert","hf.Revoker=false:ecert","app.accounting.role=manager:ecert","department=accounting:ecert"'
fabric-ca-client register --id.type user --id.name mary --id.secret pw --id.affiliation dfarmadmin --id.attrs $ATTRIBUTES
export FABRIC_CA_CLIENT_HOME=$PWD/client/dfarmadmin/mary
fabric-ca-client enroll -u http://mary:pw@localhost:7054
cp -a $PWD/client/dfarmadmin/admin/msp/signcerts  $FABRIC_CA_CLIENT_HOME/msp/admincerts

export FABRIC_CA_CLIENT_HOME=$PWD/client/dfarmadmin/admin
ATTRIBUTES='"hf.AffiliationMgr=false:ecert","hf.Revoker=false:ecert","app.accounting.role=accountant:ecert","department=accounting:ecert"'
fabric-ca-client register --id.type user --id.name john --id.secret pw --id.affiliation dfarmadmin --id.attrs $ATTRIBUTES
export FABRIC_CA_CLIENT_HOME=$PWD/client/dfarmadmin/john
fabric-ca-client enroll -u http://john:pw@localhost:7054
cp -a $PWD/client/dfarmadmin/admin/msp/signcerts  $FABRIC_CA_CLIENT_HOME/msp/admincerts

export FABRIC_CA_CLIENT_HOME=$PWD/client/dfarmadmin/admin
ATTRIBUTES='"hf.AffiliationMgr=false:ecert","hf.Revoker=false:ecert","department=logistics:ecert","app.logistics.role=specialis:ecert"'
fabric-ca-client register --id.type user --id.name anil --id.secret pw --id.affiliation dfarmadmin --id.attrs $ATTRIBUTES
export FABRIC_CA_CLIENT_HOME=$PWD/client/dfarmadmin/anil
fabric-ca-client enroll -u http://anil:pw@localhost:7054
cp -a $PWD/client/dfarmadmin/admin/msp/signcerts  $FABRIC_CA_CLIENT_HOME/msp/admincerts

# Shutdown CA
docker-compose -f docker-compose-ca.yaml down

# Setup network config
export FABRIC_CFG_PATH=$PWD/config
configtxgen -outputBlock  ./config/orderer/dfarm-genesis.block -channelID ordererchannel  -profile DfarmOrdererGenesis
configtxgen -outputCreateChannelTx  ./config/dfarmchannel.tx -channelID dfarmchannel  -profile DfarmChannel

ANCHOR_UPDATE_TX=./config/dfarm-anchor-update-dfarmadmin.tx
configtxgen -profile DfarmChannel -outputAnchorPeersUpdate $ANCHOR_UPDATE_TX -channelID dfarmchannel -asOrg DfarmadminMSP

ANCHOR_UPDATE_TX=./config/dfarm-anchor-update-dfarmretail.tx
configtxgen -profile DfarmChannel -outputAnchorPeersUpdate $ANCHOR_UPDATE_TX -channelID dfarmchannel -asOrg DfarmretailMSP
