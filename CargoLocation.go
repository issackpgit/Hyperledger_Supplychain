package main

import (
	
	"errors"
	"fmt"
	//"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)
type CL struct {

}

func (t *CL) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// Check if table already exists
	_, err := stub.GetTable("CargoLocation")
	if err == nil {
		// Table already exists; do not recreate
		return nil, nil
	}


	// Create CL Table
	err = stub.CreateTable("CargoLocation", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Type", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "ContractNo", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "CargoLocation", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ExpCargoLocationTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ExForwarderCargoLocationTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ExShippingFirmLocationTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ShippingCargoLocationTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ImShipCargoLocationTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "CargoReceivedLocationTime", Type: shim.ColumnDefinition_STRING, Key: false},

	})
	if err != nil {
		return nil, errors.New("Failed creating CargoLocationTable.")
	}

	return nil, nil

	}



func (t *CL) GetCargoLocation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
		if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1.")
		}

		ContractNo := args[0]

		// Get the row pertaining to this UID
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: "CL"}}
	columns = append(columns, col1)
	col2 := shim.Column{Value: &shim.Column_String_{String_: ContractNo}}
	columns = append(columns, col2)

	row, err := stub.GetRow("CargoLocation", columns)
	if err != nil {
		return nil, fmt.Errorf("Error: Failed retrieving document with ContractNo %s. Error %s", ContractNo, err.Error())
	}

	

	return []byte(row.Columns[2].GetString_()), nil

	}



	func (t *CL) UpdateCargoLocation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
		if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2.")
		}

		ContractNo := args[0]
		CargoLocation := args[1]
		
	
// Insert a row
	ok, err := stub.InsertRow("CargoLocation", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: "CL"}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractNo}},
			&shim.Column{Value: &shim.Column_String_{String_: CargoLocation}},
			&shim.Column{Value: &shim.Column_String_{String_: "DDMMYYYY"}},
			&shim.Column{Value: &shim.Column_String_{String_: "DDMMYYYY"}},
			&shim.Column{Value: &shim.Column_String_{String_: "DDMMYYYY"}},
			&shim.Column{Value: &shim.Column_String_{String_: "DDMMYYYY"}},
			&shim.Column{Value: &shim.Column_String_{String_: "DDMMYYYY"}},
			&shim.Column{Value: &shim.Column_String_{String_: "DDMMYYYY"}}},

	})

		if !ok && err == nil {
		
			//When row already exists, update the data
		ok, err := stub.ReplaceRow("CargoLocation", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: "CL"}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractNo}},
			&shim.Column{Value: &shim.Column_String_{String_: CargoLocation}},
			&shim.Column{Value: &shim.Column_String_{String_: "DDMMYYYY"}},
			&shim.Column{Value: &shim.Column_String_{String_: "DDMMYYYY"}},
			&shim.Column{Value: &shim.Column_String_{String_: "DDMMYYYY"}},
			&shim.Column{Value: &shim.Column_String_{String_: "DDMMYYYY"}},
			&shim.Column{Value: &shim.Column_String_{String_: "DDMMYYYY"}},
			&shim.Column{Value: &shim.Column_String_{String_: "DDMMYYYY"}},			
			
		}})

			if !ok && err == nil {

				return nil, errors.New("Document unable to Update.")
		}


		}

		return nil, err
	}


