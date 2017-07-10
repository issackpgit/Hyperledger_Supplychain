package main

import (
	
	"errors"
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)
type CRT struct {

}


type CRTJSON struct {

	ImporterShippingTime  string `json:"ImporterShippingTime"`
	ImporterReceivedTime  string `json:"ImporterReceivedTime"`

}



func (t *CRT) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// Check if table already exists
	_, err := stub.GetTable("CargoReceiveTime")
	if err == nil {
		// Table already exists; do not recreate
		return nil, nil
	}


	// Create CL Table
	err = stub.CreateTable("CargoReceiveTime", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Type", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "ContractNo", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "ImporterShippingTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ImporterReceivedTime", Type: shim.ColumnDefinition_STRING, Key: false},
		

	})
	if err != nil {
		return nil, errors.New("Failed creating CargoReceiveTimeTable.")
	}

	return nil, nil

	}



	func (t *CRT) SubmitDoc(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		
		if len(args) != 1 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 1. Got: %d.", len(args))
		}

		ContractNo := args[0]
		ImporterShippingTime := ""
		ImporterReceivedTime:= ""

		// Insert a row
	ok, err := stub.InsertRow("CargoReceiveTime", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: "CRT"}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractNo}},
			&shim.Column{Value: &shim.Column_String_{String_: ImporterShippingTime}},
			&shim.Column{Value: &shim.Column_String_{String_: ImporterReceivedTime}}},
		

	})

	if !ok && err == nil {
		return nil, errors.New("Document already exists.")

	}

		return nil, err
}


func (t *CRT) UpdateDoc(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		
		if len(args) != 3 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 3. Got: %d.", len(args))
		}

		ContractNo := args[0]
		status := args[1]
		
		

		// Get the row pertaining to this UID
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: "CRT"}}
	columns = append(columns, col1)
	col2 := shim.Column{Value: &shim.Column_String_{String_: ContractNo}}
	columns = append(columns, col2)

	row, err := stub.GetRow("CargoReceiveTime", columns)
	if err != nil {
		return nil, fmt.Errorf("Error: Failed retrieving document with ContractNo %s. Error %s", ContractNo, err.Error())
	}


		ImporterShippingTime := row.Columns[2].GetString_()
		ImporterReceivedTime:= row.Columns[3].GetString_()


		if status == "NotifyImShip" {

			ImporterShippingTime = args[2]

		} else if status == "CargoReceived" {

			ImporterReceivedTime = args[2]
		}


		ok, err := stub.ReplaceRow("CargoReceiveTime", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: "CRT"}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractNo}},
			&shim.Column{Value: &shim.Column_String_{String_: ImporterShippingTime}},			
			&shim.Column{Value: &shim.Column_String_{String_: ImporterReceivedTime}},

		}})

			if !ok && err == nil {

				return nil, errors.New("Document unable to Update.")
			}		


	return nil ,err

	}

	func (t *CRT) GetCRT (stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1.")
		}


	ContractNo := args[0]

	// Get the row pertaining to this UID
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: "CRT"}}
	columns = append(columns, col1)
	col2 := shim.Column{Value: &shim.Column_String_{String_: ContractNo}}
	columns = append(columns, col2)

	row, err := stub.GetRow("CargoReceiveTime", columns)
	if err != nil {
		return nil, fmt.Errorf("Error: Failed retrieving document with ContractNo %s. Error %s", ContractNo, err.Error())
	}


	var crtJSON CRTJSON 

	// GetRows returns empty message if key does not exist
	if len(row.Columns) == 0 {
		return nil, errors.New("No record found")
	
	} 


		crtJSON.ImporterShippingTime  = row.Columns[2].GetString_()
		crtJSON.ImporterReceivedTime = row.Columns[3].GetString_()
		
		


	jsonCRT, err := json.Marshal(crtJSON)

	if err != nil {

		return nil, err
	}


 	return jsonCRT, nil

	}