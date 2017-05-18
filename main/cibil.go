package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"log"
	"math"
	"strconv"
	"strings"
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
	PanNumber       string  `json:"PanNumber"`
	TransactionType string  `json:"TransactionType"`
	LoanId          string  `json:"LoanId"`
	TransactionId   string  `json:"TransactionId"`
	Amount          float64 `json:"Amount"`
	TransactionDate string  `json:"TransactionDate"`
	InstitutionName string  `json:"InstitutionName"`
}

type Product struct {
	Name      string  `json:"name"`
	Amount    float64 `json:"amount"`
	Owner     string  `json:"owner"`
	Productid string  `json:"productid"`
}

type Response2 struct {
	Page   int      `json:"Page"`
	Fruits []string `json:"Fruits"`
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
	} else if function == "readonetransaction" {
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
		PanNumber:       args[0],
		TransactionType: args[1],
		LoanId:          args[2],
		TransactionId:   args[3],
		Amount:          amt,
		TransactionDate: args[5],
		InstitutionName: args[6],
	}

	bytes, err := json.Marshal(transaction)
	if err != nil {
		fmt.Println("Error marshaling transaction")
		return nil, errors.New("Error marshaling transaction")
	}
	fmt.Println("Error marshaling transaction" + transaction.PanNumber)

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

	myLogger2 := shim.NewLogger("Read Transaction Logger")
	infoLevel, _ := shim.LogLevel("INFO")
	myLogger2.SetLevel(infoLevel)
	myLogger2.Info("***********************************Read Transaction Logger************************")
	key := args[0]

	myLogger2.Info("Before get state " + key)
	object, err := stub.GetState(key)
	myLogger2.Info("After get state " + key)

	if err != nil {
		fmt.Println("Error retrieving " + key)
		return nil, errors.New("Error retrieving " + key)
	}
	return object, nil
}

func (t *SimpleChaincode2) readTransaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	log := "Start of the read process"
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. expecting 1")
	}
	key := args[0]

	log = log + " key is " + key + " "
	inputs := [100]string{}
	for count := 0; count < 100; count++ {
		bytes, err := stub.GetState(args[0] + strconv.Itoa(count))
		if err != nil {
			break
		}
		inputs[count] = string(bytes)

	}
	sumOfEMIs := 0.0
	sumOfOs := 0.0
	loanAmt := 0.0
	score := 600.0

	for j := 0; j < len(inputs); j++ {
		str := inputs[j]
		fmt.Printf("\n")
		fmt.Printf("\n")
		fmt.Printf(str)
		str = strings.Replace(str, "{", "", -1)
		str = strings.Replace(str, "}", "", -1)
		values := strings.Split(str, ",")
		transaction := Transaction{}
		for i := 0; i < len(values); i++ {
			replacedStr := strings.Replace(values[i], "\"", "", -1)
			tokens := strings.Split(replacedStr, ":")
			//fmt.Printf("\nKey " + tokens[0] + " " + " Value " + tokens[1])

			if "TransactionType" == tokens[0] {
				transaction.TransactionType = tokens[1]
				fmt.Printf("\n Transaction type !!!!!!!" + transaction.TransactionType)
			}

			if "PanNumber" == tokens[0] {
				transaction.PanNumber = tokens[1]
				fmt.Printf("\n PAN Number !!!!!!!" + transaction.PanNumber)
			}

			if "Amount" == tokens[0] {
				amt, err := strconv.ParseFloat(tokens[1], 64)
				transaction.Amount = amt
				if err != nil {

				}
			}

			if "TransactionId" == tokens[0] {
				transaction.TransactionId = tokens[1]
				fmt.Printf("\n TransactionId !!!!!!!" + transaction.TransactionId)
			}
		}
		if transaction.TransactionType == "EMI" {
			sumOfEMIs = sumOfEMIs + transaction.Amount
		} else if transaction.TransactionType == "OUTSTANDING" {
			sumOfOs = sumOfOs + transaction.Amount
		} else if transaction.TransactionType == "LOAN" {
			loanAmt = transaction.Amount
		}

	}
	fmt.Printf("sum of Emi %f", sumOfEMIs)
	fmt.Printf(" sum of os %f", sumOfOs)
	fmt.Printf(" loan amt %f", loanAmt)

	diff := sumOfEMIs - sumOfOs
	percentage := (diff / loanAmt) * 100

	spread := 300 * (percentage / 100)
	score = score + spread

	if score > 900 {
		score = 900
	} else if score < 300 {
		score = 300
	}
	fmt.Printf("score: %f", score)

	scoreBytes := []byte(strconv.FormatFloat(score, 'f', -1, 64))
	//scoreBytes := Float64bytes(score)


	return scoreBytes, nil
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
	return bytes, nil
}
func (t *SimpleChaincode2) addProduct(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("adding product information")
	if len(args) != 4 {
		return nil, errors.New("Incorrect Number of arguments.Expecting 4 for addProduct")
	}
	amt, err := strconv.ParseFloat(args[1], 64)

	product := Product{
		Name:      args[0],
		Amount:    amt,
		Owner:     args[2],
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

func (t *SimpleChaincode2) readOneTransaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("read() is running")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. expecting 1")
	}
	key := args[0]
	bytes, err := stub.GetState(args[0] + args[1])
	fmt.Println(bytes)
	if err != nil {
		fmt.Println("Error retrieving " + key)
		return nil, errors.New("Error retrieving " + key+args[1])
	}

	return bytes, nil
}
