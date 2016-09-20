package main

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	"strings"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation


type SimpleChaincode struct {
}

var claimIndexStr = "_claimindex"				//name for the key/value that will store a list of all known insurance claims
var openTradesStr = "_opentrades"				//name for the key/value that will store all open trades
var allstr string

type InsuranceClaim struct{
	/*
    Name string `json:"name"`					//the fieldtags are needed to keep case from bouncing around
	Color string `json:"color"`
	Size int `json:"size"`
	User string `json:"user"`
    */
    
    ReceiptId string `json:"receiptid"`					
	Hkid string `json:"hkid"`
	Amount int `json:"amount"`
    ClaimAmount int `json:"claimamount"`
	Company string `json:"company"`
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	var Aval int
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	// Initialize the chaincode
	Aval, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}

	// Write the state to the ledger
	err = stub.PutState("Testing", []byte(strconv.Itoa(Aval)))				//making a test var "abc", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}
	
	var empty []string
	jsonAsBytes, _ := json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	err = stub.PutState(claimIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ============================================================================================================================
// Run - Our entry point for Invocations - [LEGACY] obc-peer 4/25/2016
// ============================================================================================================================
func (t *SimpleChaincode) Run(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)
	return t.Invoke(stub, function, args)
}

// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "delete" {										//deletes an entity from its state
		return t.Delete(stub, args)
	} else if function == "write" {											//writes a value to the chaincode state
		return t.Write(stub, args)
	} else if function == "init_claim" {									//create a new insurance claim
		return t.init_claim(stub, args)
	} else if function == "set_user" {										//change owner of a insurance claim
		return t.set_user(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation")
}

// ============================================================================================================================
// Query - Our entry point for Queries
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" {													//read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query")
}

// ============================================================================================================================
// Read - read a variable from chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var receiptId, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}
    // check whether receiptid exists
	receiptId = args[0]
	valAsbytes, err := stub.GetState(receiptId )									//get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + receiptId  + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil													//send it onward
}

// ============================================================================================================================
// Delete - remove a key/value pair from state
// ============================================================================================================================
func (t *SimpleChaincode) Delete(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	
	name := args[0]
	err := stub.DelState(name)													//remove the key from chaincode state
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	//get the claim index
	claimAsBytes, err := stub.GetState(claimIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get claim index")
	}
	var claimIndex []string
	json.Unmarshal(claimAsBytes, &claimIndex)								//un stringify it aka JSON.parse()
	
	//remove claim from index
	for i,val := range claimIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + name)
		if val == name{															//find the correct claim
			fmt.Println("found claim")
			claimIndex = append(claimIndex[:i], claimIndex[i+1:]...)			//remove it
			for x:= range claimIndex{											//debug prints...
				fmt.Println(string(x) + " - " + claimIndex[x])
			}
			break
		}
	}
	jsonAsBytes, _ := json.Marshal(claimIndex)									//save new index
	err = stub.PutState(claimIndexStr, jsonAsBytes)
	return nil, nil
}

// ============================================================================================================================
// Write - write variable into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) Write(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var name, value string // Entities
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the variable and value to set")
	}

	name = args[0]															//rename for funsies
	value = args[1]
	err = stub.PutState(name, []byte(value))								//write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// ============================================================================================================================
// Init insurance claim - create a new insurance claim, store into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) init_claim(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var err error
    //var str string
	// args[]  0          1        2           3             4
	//   "receiptid", "hkid", "amount", "claimamount",  "Company"
    //   "IC2222" , "A1234567", "1000",     "500",      "Company A"
    
    
	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5")
	}

	fmt.Println("- start init claim -")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument -receiptid- must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument -hkid- must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("3rd argument -amount- must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("4th argument -claimamount- must be a non-empty string")
	}
    if len(args[4]) <= 0 {
		return nil, errors.New("5th argument -company- must be a non-empty string")
	}
	
    //amount, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("3rd -amount- argument must be a numeric ")
	}
    //claimamount, err := strconv.Atoi(args[3])
	if err != nil {
		return nil, errors.New("4th -claimamount- argument must be a numeric ")
	}
	
	// color := strings.ToLower(args[1])
	// user := strings.ToLower(args[3])

	str := `{"receiptid":"` + args[0] + `","hkid":"` + args[1] + `","amount":` + args[2] + `,"claimamount":` + args[3] +  `,"company":"` + args[4] + `"}`
    //str2 := `{"ReceiptId":"` + args[0] + `","Hkid":"` + args[1] + `","Smount":` + args[2] + `,"ClaimAmount":` + args[3] +  `,"company":"` + args[4] + `"}`
    
    fmt.Println("- debug str: " + str)
    
	err = stub.PutState(args[0], []byte(str))								//store insurance receiptid as key
	if err != nil {
		return nil, err
	}

	//get the claim index
	claimAsBytes, err := stub.GetState(claimIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get claim index")
	}
    
    fmt.Println("-----------------------------------")
    fmt.Println("- 20sep2016 claimIndexStr: ", claimIndexStr)
    fmt.Println("-----------------------------------")
    fmt.Println("- 20sep2016 claimAsBytes: ", claimAsBytes)
    fmt.Println("-----------------------------------")
   
    
	var claimIndex []string
	json.Unmarshal(claimAsBytes, &claimIndex)							//un stringify it aka JSON.parse()
	fmt.Println("- json.Unmarshal(claimAsBytes, &claimIndex): ", claimIndex)
    fmt.Println("-----------------------------------")
    //check if duplicated
   
    
    allstr = CToGoString(claimAsBytes[:])
    
    if strings.Contains(allstr,args[0]){
        fmt.Printf("Found receiptId in claimAsBytes " + args[0] + " \n") 
        //  t.read(stub, args)
        claimAsBytes, _ := stub.GetState(args[0])
        fmt.Printf(" claimAsBytes " + string(claimAsBytes)  + " \n") 
    } else {
        fmt.Printf("receiptId is not in claimAsBytes " + args[0] + " \n")
    }

    
	claimIndex = append(claimIndex, args[0])								//add claim to index list
	fmt.Println("! claim index: ", claimIndex)
	jsonAsBytes, _ := json.Marshal(claimIndex)
	err = stub.PutState(claimIndexStr, jsonAsBytes)						//store name of claim

	//fmt.Println("- end init claim")
	return nil, nil
}

// ============================================================================================================================
// Set User Permission on InsuranceClaim
// ============================================================================================================================
func (t *SimpleChaincode) set_user(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var err error
	
	//   0       1
	// "name", "bob"
	if len(args) < 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	
	fmt.Println("- start set user")
	fmt.Println(args[0] + " - " + args[1])
	claimAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get thing")
	}
	res := InsuranceClaim{}
	json.Unmarshal(claimAsBytes, &res)										//un stringify it aka JSON.parse()
    res.Amount,err = strconv.Atoi(args[2])										//change the user
	if err != nil {
		return nil, err
	}
    
	jsonAsBytes, _ := json.Marshal(res)
	err = stub.PutState(args[0], jsonAsBytes)								//rewrite the claim with id as key
	if err != nil {
		return nil, err
	}
	
	fmt.Println("- end set user")
	return nil, nil
}

func CToGoString(c []byte) string {
    n := -1
    for i, b := range c {
        if b == 0 {
            break
        }
        n = i
    }
    return string(c[:n+1])
}