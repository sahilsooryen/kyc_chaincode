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
    InfoElements []InfoElement `json:"infoElements"`;
}

// Document's Meta-Data structure
type InfoElement struct {
    Id string `json:"id"`;
		Title string `json:"title"`;
		ElementType string `json:"elementType"`;
		ElementValue string `json:"elementValue"`;
		ValidTill string `json:"validTill"`;
    Hash string `json:"hash"`;
		VerifiedOn string `json:"verifiedOn"`;
		VerificationProof string `json:"verificationProof"`;
    Status string `json:"status"`;
		Comments string `json:"comments"`;
}

// SimpleChaincode example simple Chaincode implementation
// type SimpleChaincode struct {
// }

func (kyc *KYCChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Init called, initializing chaincode")

	return nil, nil
}

func (kyc *KYCChaincode) createPerson(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("CHAINCODE: createPerson called")

	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	person := Person{}
	person.Id = args[0]

	infoElements := []InfoElement{}
	person.InfoElements = infoElements

	jsonAsBytes, _ := json.Marshal(person)
	err = stub.PutState(person.Id, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil

}

func (kyc *KYCChaincode) queryPerson(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("CHAINCODE: queryPerson called")

	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	Avalbytes, err := stub.GetState(args[0])
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	return Avalbytes, nil

}

func (kyc *KYCChaincode) updateInfoElement(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("CHAINCODE: updateInfoElement called")
	var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	personJSONAsBytes, err := stub.GetState(args[0])
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	person := Person{}
	json.Unmarshal(personJSONAsBytes, &person)
	fmt.Println("CHAINCODE: After Unmarshalling person")

	infoElement := InfoElement{}
	json.Unmarshal([]byte(args[1]), &infoElement)
	fmt.Println("CHAINCODE: After Unmarshalling infoElement")

	alteredInfoElements := []InfoElement{}

	if len(person.InfoElements) > 0 {
		for _, infoElement1 := range person.InfoElements {
			if infoElement1.Id == infoElement.Id {
				fmt.Println("CHAINCODE: Replacing the element found")
				alteredInfoElements = append(alteredInfoElements, infoElement)
			} else {
				fmt.Println("CHAINCODE: Keeping the old element")
				alteredInfoElements = append(alteredInfoElements, infoElement1)
			}
		}
		fmt.Println("CHAINCODE: Replacing the infoElementsList with the altered one")
		person.InfoElements = alteredInfoElements;
	} else {
			fmt.Println("CHAINCODE: Appending info element")
			person.InfoElements = append(person.InfoElements, infoElement)
	}

	fmt.Println("CHAINCODE: Writing person back to ledger")
	jsonAsBytes, _ := json.Marshal(person)
	err = stub.PutState(person.Id, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	fmt.Println("CHAINCODE: Returning from updateInfoElement")

	return nil, nil

}

func (kyc *KYCChaincode) deleteInfoElement(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("CHAINCODE: deleteInfoElement called")
	var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	personJSONAsBytes, err := stub.GetState(args[0])
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	person := Person{}
	json.Unmarshal(personJSONAsBytes, &person)
	fmt.Println("CHAINCODE: After Unmarshalling person")

	alteredInfoElements := []InfoElement{}

	if len(person.InfoElements) > 0 {
		for _, infoElement1 := range person.InfoElements {
			if infoElement1.Id != args[1] {
				fmt.Println("CHAINCODE: Keeping the old element")
				alteredInfoElements = append(alteredInfoElements, infoElement1)
			} else {
				fmt.Println("CHAINCODE: Removing info element with id" + infoElement1.Id)
			}
		}
		fmt.Println("CHAINCODE: Replacing the infoElementsList with the altered one")
		person.InfoElements = alteredInfoElements;
	} else {
		return nil, nil
	}

	fmt.Println("CHAINCODE: Writing person back to ledger")
	jsonAsBytes, _ := json.Marshal(person)
	err = stub.PutState(person.Id, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	fmt.Println("CHAINCODE: Returning from deleteInfoElement")

	return nil, nil

}

func (kyc *KYCChaincode) queryInfoElement(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("CHAINCODE: queryInfoElement called")

	var err error
	var infoElementExists bool = false
	var fetchedInfoElement InfoElement

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	personJSONAsBytes, err := stub.GetState(args[0])
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	person := Person{}
	json.Unmarshal(personJSONAsBytes, &person)
	fmt.Println("CHAINCODE: After Unmarshalling person")

	for _, infoElement := range person.InfoElements {
		if infoElement.Id == args[1]{
			fetchedInfoElement = infoElement
			infoElementExists = true
			break
		}
	}

	if infoElementExists == false {
		jsonResp := "{\"Error\":\"InfoElement with id " + args[1] + " does not exist \"}"
		return nil, errors.New(jsonResp)
	}

	infoElementAsJSONBytes, _ := json.Marshal(fetchedInfoElement)

	return infoElementAsJSONBytes, nil

}

// Deletes an entity from state
func (kyc *KYCChaincode) deletePerson(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("CHAINCODE: Running deletePerson")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	// Delete the key from the state in ledger
	err := stub.DelState(args[0])
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	return nil, nil
}

// Invoke callback representing the invocation of a chaincode
// This chaincode will manage two accounts A and B and will transfer X units from A to B upon invoke
func (kyc *KYCChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Invoke called, determining function")

	// Handle different functions
	if function == "createPerson" {
		fmt.Printf("Function is createPerson")
		return kyc.createPerson(stub, args)
	} else if function == "init" {
		fmt.Printf("Function is init")
		return kyc.Init(stub, function, args)
	} else if function == "updateInfoElement" {
		fmt.Printf("Function is updateInfoElement")
		return kyc.updateInfoElement(stub, args)
	} else if function == "deletePerson" {
		fmt.Printf("Function is deletePerson")
		return kyc.deletePerson(stub, args)
	} else if function == "deleteInfoElement" {
		fmt.Printf("Function is deleteInfoElement")
		return kyc.deleteInfoElement(stub, args)
	}

	return nil, errors.New("Received unknown function invocation")
}

// Query callback representing the query of a chaincode
func (kyc *KYCChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Query called, determining function")

	// Handle different functions
	if function == "queryPerson" {
		// Deletes an entity from its state
		fmt.Printf("Function is queryPerson")
		return kyc.queryPerson(stub, args)
	} else if function == "queryInfoElement" {
		fmt.Printf("Function is queryInfoElement")
		return kyc.queryInfoElement(stub, args)
	}

	return nil, errors.New("Received unknown function invocation")

}

func main() {
	err := shim.Start(new(KYCChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
