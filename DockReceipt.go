package main

import (
	
	"errors"
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)
type DR struct {
 po PO
}



type DRJSON struct {

UpdateTime string `json:"UpdateTime"`
DRSubmittedTime string  `json:"DRSubmittedTime"`

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



	func (t *DR) SubmitDoc(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		
		if len(args) != 3 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 3. Got: %d.", len(args))
		}

		ContractNo := args[0]
		UpdateTime := args[1]
		DRSubmittedTime := args[2]
		
		

		// Insert a row
	ok, err := stub.InsertRow("DockReceipt", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: "DR"}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractNo}},
			&shim.Column{Value: &shim.Column_String_{String_: UpdateTime}},
			&shim.Column{Value: &shim.Column_String_{String_: DRSubmittedTime}}},
		

	})

	if !ok && err == nil {
		return nil, errors.New("Document already exists.")
	}

		toSend := make ([]string, 2)
		toSend[0] = string(ContractNo)
		toSend[1] = "SubmitDR" 
			
		_, err = t.po.UpdatePO(stub, toSend)

				if err != nil {

					return nil, errors.New("Unable to Update PO Status with this ContractID ")
				}	


	return nil, err
}


func (t *DR) GetDR (stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1.")
		}


	ContractNo := args[0]

	// Get the row pertaining to this UID
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: "DR"}}
	columns = append(columns, col1)
	col2 := shim.Column{Value: &shim.Column_String_{String_: ContractNo}}
	columns = append(columns, col2)

	row, err := stub.GetRow("DockReceipt", columns)
	if err != nil {
		return nil, fmt.Errorf("Error: Failed retrieving document with ContractNo %s. Error %s", ContractNo, err.Error())
	}


	var drJSON DRJSON 
	
	// GetRows returns empty message if key does not exist
	if len(row.Columns) == 0 {
		
		drJSON.UpdateTime = ""
		drJSON.DRSubmittedTime = ""
	} else {

	drJSON.UpdateTime = row.Columns[2].GetString_()
	drJSON.DRSubmittedTime = row.Columns[3].GetString_()
	
	}


	jsonDR, err := json.Marshal(drJSON)

	if err != nil {

		return nil, err
	}

	fmt.Println(jsonDR)

 	return jsonDR, nil

	}
