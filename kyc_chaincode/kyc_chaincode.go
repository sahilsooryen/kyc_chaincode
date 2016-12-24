package main;

import (
	"errors"
    "encoding/json"
	"fmt"
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

// ============================
// Init - reset all the things
// ============================
func (kyc *KYCChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {


    return nil, nil;
}

// ======================================
// Invoke - entry point for Invocations
// ======================================
func (kyc *KYCChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("invoke is running " + function);

    if(function == "updatePerson"){
        return kyc.updatePerson(stub, args);
    } else if(function == "createPerson"){
        return kyc.createPerson(stub, args);
    }

    fmt.Println("invoke did not find func: " + function);
    return nil, errors.New("Received unknown function invocation");
}

// =================================
// Query - entry point for Queries
// =================================
func (kyc *KYCChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function);

    if(function == "queryPerson"){
        return kyc.queryPerson(stub, args);
    }

    fmt.Println("query did not find func: " + function);
    return nil, errors.New("Received unknown function query");
}

// ======================================================================
// CreatePerson - this method writes a new Person object into the ledger
// ======================================================================
func (kyc *KYCChaincode) createPerson(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    err := kyc.checkArguments(args);
    if err != nil {
        return nil, err;
    }

    personAsJSON := args[0];
    personAsBytes := []byte(personAsJSON);
    person := Person{};
    unmarshalingError := json.Unmarshal(personAsBytes, &person);
    if unmarshalingError != nil {
        return nil, unmarshalingError;
    }

		fmt.Println("SAHIL: Person ID: " + person.Id);
    // creatingErr := stub.PutState(person.Id, personAsBytes);
		creatingErr := stub.PutState("testId", []byte("testValue"));
    if creatingErr != nil {
        return nil, creatingErr;
    }

    return nil, nil;
}

// ================================================================
// QueryPerson - this method reads a Person object from the ledger
// ================================================================
func (kyc *KYCChaincode) queryPerson(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    err := kyc.checkArguments(args);
    if err != nil {
        return nil, err;
    }

    personGUID := args[0];
	personAsBytes, queryErr := stub.GetState(personGUID);
	if queryErr != nil {
        jsonResp := "{\"Error\":\"Failed to get state for person with (" + personGUID + ") GUID\"}";
		return nil, errors.New(jsonResp);
	}

    return personAsBytes, nil;
}

// =================================================================
// UpdatePerson - this method updates a Person object in the ledger
// =================================================================
func (kyc *KYCChaincode) updatePerson(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    err := kyc.checkArguments(args);
    if err != nil {
        return nil, err;
    }

    person := Person{};
    personAsJSON := args[0];
    personAsBytes := []byte(personAsJSON);
    unmarshalingError := json.Unmarshal(personAsBytes, &person);
    if unmarshalingError != nil {
        return nil, unmarshalingError;
    }

    updateErr := stub.PutState(person.Id, personAsBytes);
    if updateErr != nil {
        return nil, updateErr;
    }

    return nil, nil;
}

// ===================================================================
// CheckArguments - this method checks if the incoming Args are valid
// ===================================================================
func (kyc *KYCChaincode) checkArguments(args []string) (error) {
    if len(args) == 1 {
        return nil;
    }

    jsonResp := "{\"Error\":\"Wrong arguments, this API this API takes only 1 argument.\"}";
    return errors.New(jsonResp);
}

// =====
// Main
// =====
func main() {
	err := shim.Start(new(KYCChaincode))
	if err != nil {
		fmt.Printf("Error starting KYC chaincode: %s", err)
	}
}
