
 package main

 /* Imports
  * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
  * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
  */
 import (
	 "bytes"
	 "encoding/json"
	 "fmt"
	 "strconv"
 
	 "github.com/hyperledger/fabric/core/chaincode/shim"
	 sc "github.com/hyperledger/fabric/protos/peer"
 )
 
 // Define the Smart Contract structure
 type SmartContract struct {
 }
 
 // Define the Apple structure, with 4 properties.  Structure tags are used by encoding/json library
 type Apple struct {
	 PlantingMethod 	string `json:"PlantingMethod"`,
	 Type  string `json:"Type"`,
	 Health string `json:"Health"`,
	 Farmer  string `json:"Farmer"`
 }
 
 /*
  * The Init method is called when the Smart Contract "fabApple" is instantiated by the blockchain network
  * Best practice is to have any Ledger initialization in separate function -- see initLedger()
  */
 func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	 return shim.Success(nil)
 }
 
 /*
  * The Invoke method is called as a result of an application request to run the Smart Contract "Apple_register"
  * The calling application program has also specified the particular smart contract function to be called, with arguments
  */
 func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
 
	 // Retrieve the requested Smart Contract function and arguments
	 function, args := APIstub.GetFunctionAndParameters()
	 // Route to the appropriate handler function to interact with the ledger appropriately
	 if function == "queryApple" {
		 return s.queryApple(APIstub, args)
	 } else if function == "initLedger" {
		 return s.initLedger(APIstub)
	 } else if function == "createApple" {
		 return s.createApple(APIstub, args)
	 } else if function == "queryAllApples" {
		 return s.queryAllApples(APIstub)
	 } else if function == "changeAppleFarmer" {
		 return s.changeAppleFarmer(APIstub, args)
	 }
 
	 return shim.Error("Invalid Smart Contract function name.")
 }
 
 func (s *SmartContract) queryApple(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 if len(args) != 1 {
		 return shim.Error("Incorrect number of arguments. Expecting 1")
	 }
 
	 AppleAsBytes, _ := APIstub.GetState(args[0])
	 return shim.Success(AppleAsBytes)
 }
 
 func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	 Apples := []Apple{
		 Apple{PlantingMethod: "Irrigation", Type: "Jonagold Apple", Health: " Insect Pests", Farmer: "Tomoko"},
		 Apple{PlantingMethod: "Pruning", Type: "Cameo Apple", Health: "Good", Farmer: "Brad"},
		 Apple{PlantingMethod: "Pruning", Type: " Empire Apple", Health: " Insect Pests", Farmer: "Jin Soo"},
		 Apple{PlantingMethod: "Irrigation", Type: "McIntosh Apple", Health: "Very Good", Farmer: "Max"},
		 Apple{PlantingMethod: "Pruning", Type: "Golden Delicious Apple", Health: "Very Good", Farmer: "Adriana"},
		 Apple{PlantingMethod: "Irrigation", Type: "Braeburn Apple", Health: "Good", Farmer: "Michel"},
		 Apple{PlantingMethod: "Pruning", Type: "Cortland Apple", Health: "Very Good", Farmer: "Aarav"},
		 Apple{PlantingMethod: "Irrigation", Type: "Red Delicious Apple", Health: "Very Good", Farmer: "Pari"},
		 Apple{PlantingMethod: "Pruning", Type: "Gala Apple", Health: "Good", Farmer: "Valeria"},
		 Apple{PlantingMethod: "Irrigation", Type: "Granny Smith Apple", Health: "Very Good", Farmer: "Shotaro"},
	 }
 
	 i := 0
	 for i < len(Apples) {
		 fmt.Println("i is ", i)
		 AppleAsBytes, _ := json.Marshal(Apples[i])
		 APIstub.PutState("Apple"+strconv.Itoa(i), AppleAsBytes)
		 fmt.Println("Added", Apples[i])
		 i = i + 1
	 }
 
	 return shim.Success(nil)
 }
 
 func (s *SmartContract) createApple(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 if len(args) != 5 {
		 return shim.Error("Incorrect number of arguments. Expecting 5")
	 }
 
	 var Apple = Apple{PlantingMethod: args[1], Type: args[2], Health: args[3], Farmer: args[4]}
 
	 AppleAsBytes, _ := json.Marshal(Apple)
	 APIstub.PutState(args[0], AppleAsBytes)
 
	 return shim.Success(nil)
 }
 
 func (s *SmartContract) queryAllApples(APIstub shim.ChaincodeStubInterface) sc.Response {
 
	 startKey := "Apple0"
	 endKey := "Apple999"
 
	 resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	 if err != nil {
		 return shim.Error(err.Error())
	 }
	 defer resultsIterator.Close()
 
	 // buffer is a JSON array containing QueryResults
	 var buffer bytes.Buffer
	 buffer.WriteString("[")
 
	 bArrayMemberAlreadyWritten := false
	 for resultsIterator.HasNext() {
		 queryResponse, err := resultsIterator.Next()
		 if err != nil {
			 return shim.Error(err.Error())
		 }
		 // Add a comma before array members, suppress it for the first array member
		 if bArrayMemberAlreadyWritten == true {
			 buffer.WriteString(",")
		 }
		 buffer.WriteString("{\"Key\":")
		 buffer.WriteString("\"")
		 buffer.WriteString(queryResponse.Key)
		 buffer.WriteString("\"")
 
		 buffer.WriteString(", \"Record\":")
		 // Record is a JSON object, so we write as-is
		 buffer.WriteString(string(queryResponse.Value))
		 buffer.WriteString("}")
		 bArrayMemberAlreadyWritten = true
	 }
	 buffer.WriteString("]")
 
	 fmt.Printf("- queryAllApples:\n%s\n", buffer.String())
 
	 return shim.Success(buffer.Bytes())
 }
 
 func (s *SmartContract) changeAppleFarmer(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 if len(args) != 2 {
		 return shim.Error("Incorrect number of arguments. Expecting 2")
	 }
 
	 AppleAsBytes, _ := APIstub.GetState(args[0])
	 Apple := Apple{}
 
	 json.Unmarshal(AppleAsBytes, &Apple)
	 Apple.Farmer = args[1]
 
	 AppleAsBytes, _ = json.Marshal(Apple)
	 APIstub.PutState(args[0], AppleAsBytes)
 
	 return shim.Success(nil)
 }
 
 // The main function is only relevant in unit test mode. Only included here for completeness.
 func main() {
 
	 // Create a new Smart Contract
	 err := shim.Start(new(SmartContract))
	 if err != nil {
		 fmt.Printf("Error creating new Smart Contract: %s", err)
	 }
 }
 