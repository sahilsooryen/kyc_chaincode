/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at
  http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"errors"
	"fmt"
	// "strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// KYCChaincode structure.
type KYCChaincode struct {
}

// Person structure
type Person struct {
    Id string `json:"id"`;
    DocsMetaData []DocMetaData `json:"docsMetaData"`;
}

// Document's Meta-Data structure
type DocMetaData struct {
    Id uint `json:"id"`;
    Hash string `json:"hash"`;
    Status string `json:"status"`;
}

// SimpleChaincode example simple Chaincode implementation
// type SimpleChaincode struct {
// }

func (kyc *KYCChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Init called, initializing chaincode")

	// if len(args) != 2 {
	// 	return nil, errors.New("Incorrect number of arguments. Expecting 2")
	// }
	//
	// var err error
	//
	// err = stub.PutState(args[0], []byte(args[1]))
	// if err != nil {
	// 	return nil, err
	// }

	// if len(args) != 1 {
	// 	return nil, errors.New("Incorrect number of arguments. Expecting 1")
	// }
	//
	// personAsJSON := args[0];
	// personAsBytes := []byte(personAsJSON);
	//
	// person := Person{};
	// unmarshalingError := json.Unmarshal(personAsBytes, &person);
	// if unmarshalingError != nil {
	// 		return nil, unmarshalingError;
	// }
	//
	// creatingErr := stub.PutState(person.Id, personAsBytes);
	// if creatingErr != nil {
	// 		return nil, creatingErr;
	// }

	return nil, nil
}

// Transaction makes payment of X units from A to B
func (t *KYCChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Printf("Running invoke")

	// if len(args) != 2 {
	// 	return nil, errors.New("Incorrect number of arguments. Expecting 2")
	// }
	//
	// var err error
	//
	// err = stub.PutState(args[0], []byte(args[1]))
	// if err != nil {
	// 	return nil, err
	// }

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	personAsJSON := args[0];
	personAsBytes := []byte(personAsJSON);

	person := Person{};
	unmarshalingError := json.Unmarshal(personAsBytes, &person);
	if unmarshalingError != nil {
			return nil, unmarshalingError;
	}

	creatingErr := stub.PutState(person.Id, personAsBytes);
	if creatingErr != nil {
			return nil, creatingErr;
	}

	return nil, nil
}

// Deletes an entity from state
func (kyc *KYCChaincode) delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Printf("Running delete")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	return nil, nil
}

// Invoke callback representing the invocation of a chaincode
// This chaincode will manage two accounts A and B and will transfer X units from A to B upon invoke
func (kyc *KYCChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Invoke called, determining function")

	// Handle different functions
	if function == "invoke" {
		// Transaction makes payment of X units from A to B
		fmt.Printf("Function is invoke")
		return kyc.invoke(stub, args)
	} else if function == "init" {
		fmt.Printf("Function is init")
		return kyc.Init(stub, function, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		fmt.Printf("Function is delete")
		return kyc.delete(stub, args)
	}

	return nil, errors.New("Received unknown function invocation")
}

func (kyc* KYCChaincode) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Run called, passing through to Invoke (same function)")

	// Handle different functions
	if function == "invoke" {
		// Transaction makes payment of X units from A to B
		fmt.Printf("Function is invoke")
		return kyc.invoke(stub, args)
	} else if function == "init" {
		fmt.Printf("Function is init")
		return kyc.Init(stub, function, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		fmt.Printf("Function is delete")
		return kyc.delete(stub, args)
	}

	return nil, errors.New("Received unknown function invocation")
}

// Query callback representing the query of a chaincode
func (kyc *KYCChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Query called, determining function")

	if function != "query" {
		fmt.Printf("Function is query")
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	personId := args[0];
	personAsBytes, queryErr := stub.GetState(personId);
	if queryErr != nil {
        jsonResp := "{\"Error\":\"Failed to get state for person with (" + personId + ") GUID\"}";
		return nil, errors.New(jsonResp);
	}

	return personAsBytes, nil

	// valueBytes, err := stub.GetState(args[0])
	// if err != nil {
	// 	jsonResp := "{\"Error\":\"Failed to get state for \"}"
	// 	return nil, errors.New(jsonResp)
	// }

	// return valueBytes, nil

}

func main() {
	err := shim.Start(new(KYCChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
