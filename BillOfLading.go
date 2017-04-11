package main

import (
	
	"errors"
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)
type BL struct {

po PO
}


type BLJSON struct {


DateCargoReceived string `json: "DateCargoReceived"` 
UnknownClause string `json: "UnknownClause"`
BLNo string `json: "BLNo" `
UpdateTime string `json: "UpdateTime"`
BLSubmittedTime string `json: "BLSubmittedTime"`

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


	func (t *BL) SubmitDoc(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		
		if len(args) != 6 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 6. Got: %d.", len(args))
		}

		ContractNo := args[0]
		DateCargoReceived := args[1]
		UnknownClause := args[2]
		BLNo := args[3]
		UpdateTime := args[4]
		BLSubmittedTime:= args[5]


		
		

		// Insert a row
	ok, err := stub.InsertRow("BillOfLading", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: "BL"}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractNo}},
			&shim.Column{Value: &shim.Column_String_{String_: DateCargoReceived}},
			&shim.Column{Value: &shim.Column_String_{String_: UnknownClause}},
			&shim.Column{Value: &shim.Column_String_{String_: BLNo}},
			&shim.Column{Value: &shim.Column_String_{String_: UpdateTime}},
			&shim.Column{Value: &shim.Column_String_{String_: BLSubmittedTime}}},
		

	})

	if !ok && err == nil {
		return nil, errors.New("Document already exists.")
	}

		toSend := make ([]string, 2)
		toSend[0] = string(ContractNo)
		toSend[1] = "SubmitBL" 
			
		_, err = t.po.UpdatePO(stub, toSend)

				if err != nil {

					return nil, errors.New("Unable to Update PO Status with this ContractID ")
				}	


	return nil, err
}


func (t *BL) GetBL (stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1.")
		}


	ContractNo := args[0]

	// Get the row pertaining to this UID
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: "BL"}}
	columns = append(columns, col1)
	col2 := shim.Column{Value: &shim.Column_String_{String_: ContractNo}}
	columns = append(columns, col2)

	row, err := stub.GetRow("BillOfLading", columns)
	if err != nil {
		return nil, fmt.Errorf("Error: Failed retrieving document with ContractNo %s. Error %s", ContractNo, err.Error())
	}

	var blJSON BLJSON

	// GetRows returns empty message if key does not exist
	if len(row.Columns) == 0 {
		
	blJSON.DateCargoReceived = ""
	blJSON.UnknownClause = ""
	blJSON.BLNo = ""
	blJSON.UpdateTime = ""
	blJSON.BLSubmittedTime = ""

	} else {


	blJSON.DateCargoReceived = row.Columns[2].GetString_()
	blJSON.UnknownClause = row.Columns[3].GetString_()
	blJSON.BLNo = row.Columns[4].GetString_()
	
	blJSON.UpdateTime = row.Columns[5].GetString_()
	blJSON.BLSubmittedTime = row.Columns[6].GetString_()

	}

	jsonBL, err := json.Marshal(blJSON)

	if err != nil {

		return nil, err
	}

	fmt.Println(jsonBL)

 	return jsonBL, nil

	}

