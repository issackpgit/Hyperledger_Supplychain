package main

import (
	
	"errors"
	//"fmt"
	//"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)
type CISIPL struct {

}

func (t *CISIPL) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// Check if table already exists
	_, err := stub.GetTable("CISIPL")
	if err == nil {
		// Table already exists; do not recreate
		return nil, nil
	}


	// Create CISIPL Table
	err = stub.CreateTable("CISIPL", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Type", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "ContractNo", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "CaseMark", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "InvoiceNo", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "UpdateTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "CISIPLSubmittedTime", Type: shim.ColumnDefinition_STRING, Key: false},
		

	})
	if err != nil {
		return nil, errors.New("Failed creating CISIPLTable.")
	}


	return nil, nil

	}
