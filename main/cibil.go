package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//Product - Structure for transactions used in buy goods
type Transaction struct {
	PanNumber   string  `json:"pan"` 
	TransactionType string `json:"transaction_type"`
	LoanId string `json:"loanId"`
	TransactionId string `json:"transactionid"`
	Amount float64 `json:"amount"`
	Date string  `json:"date"`
	InstitutionName string     `json:"bank"`
}

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("Genesis", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if function == "addtransaction" {
		return t.addTransaction(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	} else if function == "readtransaction" {
		return t.readTransaction(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

func (t *SimpleChaincode) addTransaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("adding transaction information")
	if len(args) != 4 {
		return nil, errors.New("Incorrect Number of arguments.Expecting 4 for addTransaction")
	}
	amt, err := strconv.ParseFloat(args[4], 64)
	

	transaction := Transaction{
		PanNumber:   args[0],
		TransactionType : args[1],
		LoanId : 	args[2],
		TransactionId : args[3],
		Amount: amt,
		Date: 		args[5],
		InstitutionName: args[6]
	}

	bytes, err := json.Marshal(transaction)
	if err != nil {
		fmt.Println("Error marshaling transaction")
		return nil, errors.New("Error marshaling transaction")
	}
	//b := [5]string{transaction.PanNumber, transaction.TransactionType, transaction.LoanId,transaction.TransactionId}
	//compKey,err = stub.CreateCompositeKey(objectType string, b) (string, error)
	fmt.Println("Error marshaling transaction"+transaction.PanNumber)
	
	err := stub.PutState(transaction.PanNumber, bytes)
	if err != nil {
		return nil, err
	}
return nil, nil
}

func (t *SimpleChaincode) readTransaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("read() is running")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. expecting 1")
	}

	key := args[0] // name of Entity
	fmt.Println("key is ")
	fmt.Println(key)
	b := [4]string{transaction.PanNumber, transaction.TransactionType, transaction.LoanId,transaction.TransactionId}
	coloredMarbleResultsIterator, err := GetStateByPartialCompositeKey(args[0], b)
	if err != nil {
		return shim.Error(err.Error())
	}
	
	for i = 0; coloredMarbleResultsIterator.HasNext(); i++ {
		responseRange, err := coloredMarbleResultsIterator.Next()
		objectType, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		
		if err != nil {
			return shim.Error(err.Error())
		}
		panNumber := compositeKeyParts[0]
		fmt.Println(responseRange.Value+" PAN NUMBER" +panNumber);
		
//		TransactionType:= compositeKeyParts[0]
//		LoanId:= compositeKeyParts[0]
//		TransactionId:= compositeKeyParts[0]
//		Amountamt:= compositeKeyParts[0]
//		Date:= compositeKeyParts[0]
//		InstitutionName:= compositeKeyParts[0]
		
//		returnedMarbleName := compositeKeyParts[1]
//		fmt.Printf("- found a marble from index:%s color:%s name:%s\n", objectType, returnedColor, returnedMarbleName)

		// Now call the transfer function for the found marble.
		// Re-use the same function that is used to transfer individual marbles
//		response := t.transferMarble(stub, []string{returnedMarbleName, newOwner})
		// if the transfer failed break out of loop and return error
//		if response.Status != shim.OK {
//			return shim.Error("Transfer failed: " + response.Message)
//		}
	}
//	bytes, err := stub.GetState(args[0])
//	fmt.Println(bytes)
	if err != nil {
		fmt.Println("Error retrieving " + key)
		return nil, errors.New("Error retrieving " + key)
	}
	/*
	transaction := Product{}
	err = json.Unmarshal(bytes, &transaction)
	if err != nil {
		fmt.Println("Error Unmarshaling customerBytes")
		return nil, errors.New("Error Unmarshaling customerBytes")
	}
	
	bytes, err = json.Marshal(transaction)
	if err != nil {
		fmt.Println("Error marshaling customer")
		return nil, errors.New("Error marshaling customer")
	}
	fmt.Println(bytes)
	*/
	return bytes, nil
}


