Install & Instantiate sample chaincode
=========
Part-1
------
1. Initialize the dev environment
2. Set the context to dfarmadmin

Part-2
------
1. Set the name, version & path of the chaincode chaincode_example02
2. Install the chaincode
3. Confirm by listing

Part-3
------
1. Setup the constructor parameters

'{"Args":["init","a","100","b","300"]}'

2. Instantiate the chaincode
3. Confirm by listing



Solution
--------
Part-1
------
dev-init.sh
.  set-env.sh   dfarmadmin

Part-2
------
set-chain-env.sh   -n  gocc  -v  1.0   -p  chaincode_example02
chain.sh install
chain.sh list

Part-3
------
set-chain-env.sh   -c   '{"Args":["initLedger"]}'
chain.sh instantiate
chain.sh list


Exercise: Invoke & Query the sample chaincode
=========
This is in continuation of the first exercise

Part-1
------
1. Make sure the gocc chaincode is installed 
2. Check the chaincode env to see if the CC_NAME=gocc

Part-2
------
1. Setup the query argument for the chaincode

'{"Args":["queryProduce","PR001"]}'  <-- This will get us value of a-->

2. Execute the query

3. Change the query & check the value of b

Part-3
------
1. Setup the invoke argument for the chaincode

'{"Args":["invoke","a","b","5"]}'      <-- This will transfer 5 from a to b-->

2. Execute the invoke

<wait for few seconds>

3. Query the values for a & b - did the invoke work?


Exercise: Update the chaincode
=========
This is in continuation of the first exercise

Part-1
------
1. Install a new version of the chaincode by changing the version to 2.0
2. Upgrade the chaincode
3. List to confirm that chaincode version has changed

Part-2
------
1. Upgrade the chaincode by executing upgrade-auto 
2. List to check the version of upgraded chaincode


Exercise: Launch Peers in Dev Mode & Test
=========

Part-1
------
1. Launch the env in DEV mode
2. Install the sample chaincode
3. Execute the chaincode using 
cc-run.sh

<Open terminal#2>
4. Instantiate the chaincode 
chain.sh instantiate

5. Test by running - 
cc-run.sh 
=> Observe on terminal #1

6. Change the chaincode
7. 