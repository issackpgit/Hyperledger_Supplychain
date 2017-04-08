package main

import (
	
	"errors"
	//"fmt"
	//"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)
type BL struct {

}

func (t *BL) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// Check if table already exists
	_, err := stub.GetTable("BillOfLading")
	if err == nil {
		// Table already exists; do not recreate
		return nil, nil
	}


	// Create BOL Table
	err = stub.CreateTable("BillOfLading", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Type", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "ContractNo", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "DateCargoReceived", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "UnknownClause", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "BLNo", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "UpdateTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "BLSubmittedTime", Type: shim.ColumnDefinition_STRING, Key: false},

	})
	if err != nil {
		return nil, errors.New("Failed creating BillOfLadingTable.")
	}


	return nil, nil

	}
