/*
So let's just play around a bit with chaincode!


*/

package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type TestChaincode struct {
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(TestChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *TestChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if len(args) < 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting at least 2")
	}
	aName := args[0]
	bName := args[1]
	
	// Write the state to the ledger
	err := stub.PutState(aName, []byte(strconv.Itoa(1000)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(bName, []byte(strconv.Itoa(1000)))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *TestChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "transfer" {
		return t.transfer(stub, args)
	} else if function == "earn" {
		return t.earn(stub, args)
	} else if function == "set" {
		return t.earn(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation")
}

// Query is our entry point for queries
func (t *TestChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "balance" {											//read a variable
		return t.getBalance(stub, args);
	}
	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query")
}

func (t *TestChaincode) getBalance(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	name := args[0]
	fmt.Println("reading value for " + name)
	val,err := stub.GetState(name)
	if err != nil {
		return nil,err
	}
	return val,nil
}

// from, to, amount
func (t *TestChaincode) transfer(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}
	from := args[0]
	to := args[1]
	amount,err := strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("3rd value must be integer")
	}
	
	fmt.Println("transferring " + args[2] + " from " + from + " to " + to)
	
	fromBal_,err := stub.GetState(from)
	fromBalance,err := strconv.Atoi(string(fromBal_))
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	if amount > fromBalance {
		return nil, errors.New("Not enough money in account to transfer")
	}
	toBal_,err := stub.GetState(to)
	toBalance,err := strconv.Atoi(string(toBal_))
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	fromBalance -= amount
	toBalance += amount
	fmt.Printf(from + " = %d, " + to + " = %d\n", fromBalance, toBalance)	
	
	// Write the state back to the ledger
	err = stub.PutState(from, []byte(strconv.Itoa(fromBalance)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(to, []byte(strconv.Itoa(toBalance)))
	if err != nil {
		return nil, err
	}
	
	return nil, nil
}

func (t *TestChaincode) earn(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	to := args[0]
	amount,err := strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("2nd value must be integer")
	}
	
	bal_,err := stub.GetState(to)
	balance,err := strconv.Atoi(string(bal_))
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	if amount + balance < 0 {
		return nil, errors.New("Not enough money in account")
	}
	
	balance += amount
	
	// Write the state back to the ledger
	err = stub.PutState(to, []byte(strconv.Itoa(balance)))
	if err != nil {
		return nil, err
	}
	return nil,nil
	
}

func (t *TestChaincode) set(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	to := args[0]
	amount,err := strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("2nd value must be integer")
	}
	
	// Write the state back to the ledger
	err = stub.PutState(to, []byte(strconv.Itoa(amount)))
	if err != nil {
		return nil, err
	}
	return nil,nil	
}


/*
func (t *TestChaincode) write(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")
	
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}
	
	key = args[0]
	value = args[1]
	err = stub.PutState(key, []byte(value))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *TestChaincode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
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
*/
