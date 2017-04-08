package main

import (
	
	"errors"
	//"fmt"
	//"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)
type DR struct {

}

func (t *DR) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// Check if table already exists
	_, err := stub.GetTable("DockReceipt")
	if err == nil {
		// Table already exists; do not recreate
		return nil, nil
	}


	// Create DR Table
	err = stub.CreateTable("DockReceipt", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Type", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "ContractNo", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "UpdateTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "DRSubmittedTime", Type: shim.ColumnDefinition_STRING, Key: false},
		
		
	})
	if err != nil {
		return nil, errors.New("Failed creating DockReceiptTable.")
	}


	return nil, nil

	}
