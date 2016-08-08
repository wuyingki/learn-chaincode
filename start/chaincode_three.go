
package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	
    // Insurance demo variables
     var receiptId string
     var hkId string
     var companyId string
     var receiptAmount int
     var claimingAmount int
     var claimedAmount int //calcaulated var
     s:=[]string{} //slice of all receipts

    //var A, B string    // Entities
	//var Aval, Bval int // Asset holdings
	var err error

	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5")
	}

	// Initialize the chaincode
	/*
    A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}
	B = args[2]
	Bval, err = strconv.Atoi(args[3])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)
    */
    receiptId = args[0]
    hkId = args[1]
    companyId = args[2]
    receiptAmount = strconv.Aiot(args[3])
    claimingAmount = strconv.Atoi(args[4])
    claimedAmounr = 0
    
    fmt.Printf("Init - receiptId:"+ receiptId + " , hkId:" + hkId +" , companyId:" + companyId +" ,receiptAmount:" + receiptAmount + " , claimingAmount:" + claimingAmount )
	// Write the state to the ledger , which is receiptId-name, receiptAmount-value
	err = stub.PutState(receiptId, []byte(strconv.Itoa(receiptAmount)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState("AllReceipts", []byte(s)))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	}
/*
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var X int          // Transaction value
*/
    // Insurance demo variables
     var receiptId string
     var hkId string
     var companyId string
     var receiptAmount int
     var claimingAmount int
     var claimedAmount int //calcaulated var
     var numberOfReceipts int //calculated var
     s:=[]string{} //slice of all receipts
	 var err error

	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5")
	}
/*
	A = args[0]
	B = args[1]
*/
    
    receiptId = args[0]
    hkId = args[1]
    companyId = args[2]
    receiptAmount = strconv.Aiot(args[3])
    claimingAmount = strconv.Atoi(args[4])
    claimedAmounr = 0
    numberOfReceipts = 0
    fmt.Printf("Invoke - receiptId:"+ receiptId + " , hkId:" + hkId +" , companyId:" + companyId +" ,receiptAmount:" + receiptAmount + " , claimingAmount:" + claimingAmount )
    
    
	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
    
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	if Avalbytes == nil {
		return nil, errors.New("Entity not found")
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))

	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	if Bvalbytes == nil {
		return nil, errors.New("Entity not found")
	}
	Bval, _ = strconv.Atoi(string(Bvalbytes))

	// Perform the execution
	X, err = strconv.Atoi(args[2])
	Aval = Aval - X
	Bval = Bval + X
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state back to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function != "query" {
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}
	var A string // Entities
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return nil, errors.New(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return nil, errors.New(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return Avalbytes, nil
}

