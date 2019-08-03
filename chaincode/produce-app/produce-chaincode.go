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
	ProduceName    string `json:"ProduceName"`
	Health         string `json:"Health"`
	FarmID       string `json:"FarmID"`
	Owner          string `json:"Owner "`
	
	
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
	} else if function == "recordProduce" {
		return s.recordProduce(APIstub, args)
	} else if function == "queryAllProduces" {
		return s.queryAllProduces(APIstub)
	} else if function == "changeProduceOwner" {
		return s.changeProduceOwner(APIstub, args)
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
		Produce{ ProduceName: "Apple", Health: " Insect Pests",FarmID:"FARM101", Owner : "Tomoko"},
		Produce{ ProduceName: "Tomato", Health: "Good",FarmID:"FARM102", Owner : "Brad"},
		Produce{ ProduceName: " Corn", Health: " Insect Pests",FarmID:"FARM104", Owner : "Jin Soo"},
		Produce{ ProduceName: "Cotton", Health: "Very Good",FarmID:"FARM105", Owner : "Max"},
		Produce{ ProduceName: "Orange", Health: "Very Good",FarmID:"FARM103", Owner : "Adriana"},
		Produce{ ProduceName: "Corn", Health: "Good",FarmID:"FARM106", Owner : "Michel"},
		Produce{ ProduceName: "Rice", Health: "Very Good",FarmID:"FARM107", Owner : "Aarav"},
		Produce{ ProduceName: "PeaNuts", Health: "Very Good",FarmID:"FARM108", Owner : "Pari"},
		Produce{ ProduceName: "broccoli", Health: "Good",FarmID:"FARM101", Owner : "Valeria"},
		Produce{ ProduceName: "pumpkin", Health: "Very Good",FarmID:"FARM109", Owner : "Shotaro"},
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

func (s *SmartContract) recordProduce(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	var Produce = Produce{ProduceName: args[1], Health: args[2], Owner: args[3], FarmID: args[4] }

	ProduceAsBytes, _ := json.Marshal(Produce)
	APIstub.PutState(args[0], ProduceAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllProduces(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "PR1"
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

func (s *SmartContract) changeProduceOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	ProduceAsBytes, _ := APIstub.GetState(args[0])
	Produce := Produce{}

	json.Unmarshal(ProduceAsBytes, &Produce)
	Produce.Owner  = args[1]

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
