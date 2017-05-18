package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"log"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"math"
	 "encoding/binary"

)
//var myLogger = logging.MustGetLogger("digital_im")
var (
    Trace   *log.Logger
    Info    *log.Logger
    Warning *log.Logger
    Error   *log.Logger
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

type Product struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
	Owner string  `json:"owner"`
	Productid string     `json:"productid"`
}

// SimpleChaincode2 example simple Chaincode implementation
type SimpleChaincode2 struct {
}

func main() {
	err := shim.Start(new(SimpleChaincode2))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode2) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
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
func (t *SimpleChaincode2) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if function == "addtransaction" {
		return t.addTransaction(stub, args)
	} else if function == "addproduct" {
		return t.addProduct(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode2) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	} else if function == "readtransaction" {
		return t.readTransaction(stub, args)
	} else if function == "readproduct" {
		return t.readProduct(stub, args)
	}  else if function == "readonetransaction" {
		return t.readOneTransaction(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode2) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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
func (t *SimpleChaincode2) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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

func (t *SimpleChaincode2) addTransaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("adding transaction information")
	errors.New("testing")
	if len(args) < 4 {
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
		InstitutionName: args[6],
	}

	bytes, err := json.Marshal(transaction)
	if err != nil {
		fmt.Println("Error marshaling transaction")
		return nil, errors.New("Error marshaling transaction")
	}
	fmt.Println("Error marshaling transaction"+transaction.PanNumber)
	
	err = stub.PutState(transaction.PanNumber+transaction.TransactionId, bytes)
	if err != nil {
		return nil, err
	}
return nil, nil
}



func (t *SimpleChaincode2) readTransaction1(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("read() is running")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. expecting 1")
	}

	myLogger2 := shim.NewLogger("Read Transaction Logger");
	infoLevel, _ := shim.LogLevel("INFO")
	myLogger2.SetLevel(infoLevel)
	myLogger2.Info("***********************************Read Transaction Logger************************");
	key := args[0] // name of Entity

//	myLogger.Debug("*********...")
	myLogger2.Info("Before get state " + key);
	object,err := stub.GetState(key)
	myLogger2.Info("After get state " + key);
	//var test *[] byte
	//test  = &object
//	defer iter.Close()


	if err != nil {
		fmt.Println("Error retrieving " + key)
		return nil, errors.New("Error retrieving " + key)
	}
	return object, nil
}


func (t *SimpleChaincode2) readTransaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("read() is running")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. expecting 1")
	}

	key := args[0] // name of Entity
	
	fmt.Println("key is ")
	fmt.Println(key)
	trans := Transaction {}
	tests  := []Transaction {trans}
	var p *[]byte
	for  i:=0; i>=0; i++ { 
		bytes, err := stub.GetState(key+strconv.Itoa(i))
		fmt.Println(bytes)
		if err != nil {
			fmt.Println("Error retrieving " + key)
			return nil, errors.New("Error retrieving " + key)
		}
		
		/*if(true){
			return nil,errors.New(string(bytes))
		}
		err = json.Unmarshal(bytes, trans)*/
		p = &bytes
		/*if err != nil {
			fmt.Println("Error unmarshelling" + trans.PanNumber)
			return nil, errors.New("Error Unmarshelling " + key+ ":"+err.Error())
		}*/
		tests = append(tests,trans)
	}
	sum := 0.0
	for i:=0;i<len(tests);i++ {
		sum = sum + tests[i].Amount
	}
	var tt []byte
	//tt=*p
	tt=Float64bytes(sum)
	p = &tt
	/*
	product := Product{}
	err = json.Unmarshal(bytes, &product)
	if err != nil {
		fmt.Println("Error Unmarshaling customerBytes")
		return nil, errors.New("Error Unmarshaling customerBytes")
	}
	
	bytes, err = json.Marshal(product)
	if err != nil {
		fmt.Println("Error marshaling customer")
		return nil, errors.New("Error marshaling customer")
	}
	fmt.Println(bytes)
	*/
	return *p, nil
}
func Float64bytes(float float64) []byte {
    bits := math.Float64bits(float)
    bytes := make([]byte, 8)
    binary.LittleEndian.PutUint64(bytes, bits)
    return bytes
}
func (t *SimpleChaincode2) readProduct(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("read() is running")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. expecting 1")
	}

	key := args[0] // name of Entity
	fmt.Println("key is ")
	fmt.Println(key)
	bytes, err := stub.GetState(args[0])
	fmt.Println(bytes)
	if err != nil {
		fmt.Println("Error retrieving " + key)
		return nil, errors.New("Error retrieving " + key)
	}
	/*
	product := Product{}
	err = json.Unmarshal(bytes, &product)
	if err != nil {
		fmt.Println("Error Unmarshaling customerBytes")
		return nil, errors.New("Error Unmarshaling customerBytes")
	}
	
	bytes, err = json.Marshal(product)
	if err != nil {
		fmt.Println("Error marshaling customer")
		return nil, errors.New("Error marshaling customer")
	}
	fmt.Println(bytes)
	*/
	return bytes, nil
}
func (t *SimpleChaincode2) addProduct(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("adding product information")
	if len(args) != 4 {
		return nil, errors.New("Incorrect Number of arguments.Expecting 4 for addProduct")
	}
	amt, err := strconv.ParseFloat(args[1], 64)
	

	product := Product{
		Name:   args[0],
		Amount: amt,
		Owner: args[2],
		Productid: args[3],
	}

	bytes, err := json.Marshal(product)
	if err != nil {
		fmt.Println("Error marshaling product")
		return nil, errors.New("Error marshaling product")
	}

	err = stub.PutState(product.Productid, bytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SimpleChaincode2) readOneTransaction (stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("read() is running")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. expecting 1")
	}

	key := args[0] // name of Entity
	fmt.Println("key is ")
	fmt.Println(key)
	bytes, err := stub.GetState(args[0]+"1")
	fmt.Println(bytes)
	if err != nil {
		fmt.Println("Error retrieving " + key)
		return nil, errors.New("Error retrieving " + key)
	}
	
	return bytes, nil
}