export FABRIC_CFG_PATH=$PWD/config
configtxgen -outputBlock  ./config/dfarmgenesis.block -channelID ordererchannel  -profile DfarmOrdererGenesis
configtxgen -outputCreateChannelTx  ./config/dfarmchannel.tx -channelID dfarmchannel  -profile DfarmChannel