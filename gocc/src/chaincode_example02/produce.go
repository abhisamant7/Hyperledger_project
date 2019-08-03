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

type SmartContract struct {
}

type Produce struct {
	PlantingMethod string `json:"PlantingMethod"`
	ProduceName    string `json:"ProduceName"`
	Health         string `json:"Health"`
	Farmer         string `json:"Farmer"`
	FarmID       string `json:"FarmID"`
	
}

/*
 * The Init method is called when the Smart Contract "fabProduce" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "Produce_register"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryProduce" {
		return s.queryProduce(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createProduce" {
		return s.createProduce(APIstub, args)
	} else if function == "queryAllProduces" {
		return s.queryAllProduces(APIstub)
	} else if function == "changeProduceStatus" {
		return s.changeProduceStatus(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryProduce(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ProduceAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(ProduceAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	Produces := []Produce{
		Produce{PlantingMethod: "Irrigation", ProduceName: "Jonagold Apple", Health: " Insect Pests", Farmer: "Tomoko",FarmID:"FARM101"},
		Produce{PlantingMethod: "Pruning", ProduceName: "Tomato", Health: "Good", Farmer: "Brad",FarmID:"FARM102"},
		Produce{PlantingMethod: "Pruning", ProduceName: " Corn", Health: " Insect Pests", Farmer: "Jin Soo",FarmID:"FARM104"},
		Produce{PlantingMethod: "Irrigation", ProduceName: "Cotton", Health: "Very Good", Farmer: "Max",FarmID:"FARM105"},
		Produce{PlantingMethod: "Pruning", ProduceName: "Orange", Health: "Very Good", Farmer: "Adriana",FarmID:"FARM103"},
		Produce{PlantingMethod: "Irrigation", ProduceName: "Corn", Health: "Good", Farmer: "Michel",FarmID:"FARM106"},
		Produce{PlantingMethod: "Pruning", ProduceName: "Rice", Health: "Very Good", Farmer: "Aarav",FarmID:"FARM107"},
		Produce{PlantingMethod: "Irrigation", ProduceName: "PeaNuts", Health: "Very Good", Farmer: "Pari",FarmID:"FARM108"},
		Produce{PlantingMethod: "Pruning", ProduceName: "broccoli", Health: "Good", Farmer: "Valeria",FarmID:""},
		Produce{PlantingMethod: "Irrigation", ProduceName: "pumpkin", Health: "Very Good", Farmer: "Shotaro",FarmID:"FARM109"},
	}

	i := 0
	for i < len(Produces) {
		fmt.Println("i is ", i)
		ProduceAsBytes, _ := json.Marshal(Produces[i])
		APIstub.PutState("PR"+strconv.Itoa(i), ProduceAsBytes)
		fmt.Println("Added", Produces[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createProduce(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	var Produce = Produce{PlantingMethod: args[1], ProduceName: args[2], Health: args[3], Farmer: args[4],FarmID: args[5]}

	ProduceAsBytes, _ := json.Marshal(Produce)
	APIstub.PutState(args[0], ProduceAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllProduces(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "PR001"
	endKey := "PR999"

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

	fmt.Printf("- queryAllProduces:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) changeProduceStatus(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	ProduceAsBytes, _ := APIstub.GetState(args[0])
	Produce := Produce{}

	json.Unmarshal(ProduceAsBytes, &Produce)
	Produce.Farmer = args[1]

	ProduceAsBytes, _ = json.Marshal(Produce)
	APIstub.PutState(args[0], ProduceAsBytes)

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
