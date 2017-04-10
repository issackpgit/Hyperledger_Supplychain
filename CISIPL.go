package main

import (
	
	"errors"
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)
type CISIPL struct {

    po PO
}

type CISIPLJSON struct {

		
		CaseMark 		string `json:"CaseMark"`
		InvoiceNo string `json:"InvoiceNo"`
		UpdateTime string `json:"UpdateTime"`
		CISIPLSubmittedTime string `json:"CISIPLSubmittedTime"`
				


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

func (t *CISIPL) SubmitDoc(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		
		if len(args) != 15 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 5. Got: %d.", len(args))
		}

		ContractNo := args[0]
		CaseMark := args[1]
		InvoiceNo := args[2]
		UpdateTime := args[3]
		CISIPLSubmittedTime := args[4]
		

		// Insert a row
	ok, err := stub.InsertRow("CISIPL", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: "CISIPL"}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractNo}},
			&shim.Column{Value: &shim.Column_String_{String_: CaseMark}},
			&shim.Column{Value: &shim.Column_String_{String_: InvoiceNo}},
			&shim.Column{Value: &shim.Column_String_{String_: UpdateTime}},
			&shim.Column{Value: &shim.Column_String_{String_: CISIPLSubmittedTime}},
						
        }})

	if !ok && err == nil {
		return nil, errors.New("Document already exists.")
	}

	
    /*toSend := make ([]string, 2)
			toSend[0] = string(ContractNo)
			toSend[1] = "SubmitCISIPL"
			
			_,poErr := t.po.UpdatePO(stub, toSend)
			if poErr != nil {
				return nil, poErr
			} */
    
    
	return nil, err
	}

func (t *CISIPL) GetCISIPL (stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1.")
		}


	ContractNo := args[0]

	// Get the row pertaining to this UID
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: "CISIPL"}}
	columns = append(columns, col1)
	col2 := shim.Column{Value: &shim.Column_String_{String_: ContractNo}}
	columns = append(columns, col2)

	row, err := stub.GetRow("CISIPL", columns)
	if err != nil {
		return nil, fmt.Errorf("Error: Failed retrieving document with ContractNo %s. Error %s", ContractNo, err.Error())
	}

	// GetRows returns empty message if key does not exist
	if len(row.Columns) == 0 {
		return nil, nil
	}

	var cisiplJSON CISIPLJSON 

	cisiplJSON.CaseMark = row.Columns[2].GetString_()
	cisiplJSON.InvoiceNo = row.Columns[3].GetString_()
	cisiplJSON.UpdateTime = row.Columns[4].GetString_()
	cisiplJSON.CISIPLSubmittedTime = row.Columns[5].GetString_() 
	
	

	jsonCISIPL, err := json.Marshal(cisiplJSON)

	if err != nil {

		return nil, err
	}

	fmt.Println(jsonCISIPL)

 	return jsonCISIPL, nil

	}



