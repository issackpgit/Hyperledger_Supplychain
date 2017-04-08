package main

import (
	
	"errors"
	//"fmt"
	//"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)
type CRR struct {

}

func (t *CRR) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// Check if table already exists
	_, err := stub.GetTable("CargoReceiveRequest")
	if err == nil {
		// Table already exists; do not recreate
		return nil, nil
	}


	// Create AN Table
	err = stub.CreateTable("CargoReceiveRequest", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Type", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "ContractNo", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "ImpCtryForwardingCompany", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "UpdateTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "CargoRecievedTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "CompanyIdOfImpCtryForwardingCompany", Type: shim.ColumnDefinition_STRING, Key: false},
		
	})
	if err != nil {
		return nil, errors.New("Failed creating CargoReceiveRequestTable.")
	}


	return nil, nil

	}
