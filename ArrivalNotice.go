package main

import (
	
	"errors"
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)
type AN struct {
 
 po PO
}

type ANJSON struct {

	EstimatedVesselArrivalDate string `json:"EstimatedVesselArrivalDate"`
	UpdateTime string `json:"UpdateTime"`
	ANSubmittedTime string  `json:"ANSubmittedTime"`
}


func (t *AN) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// Check if table already exists
	_, err := stub.GetTable("ArrivalNote")
	if err == nil {
		// Table already exists; do not recreate
		return nil, nil
	}


	// Create AN Table
	err = stub.CreateTable("ArrivalNote", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Type", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "ContractNo", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "EstimatedVesselArrivalDate", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "UpdateTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ANSubmittedTime", Type: shim.ColumnDefinition_STRING, Key: false},
		
	})
	if err != nil {
		return nil, errors.New("Failed creating ArrivalNoteTable.")
	}


	return nil, nil

	}

func (t *AN) SubmitDoc(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		
		if len(args) != 4 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 4. Got: %d.", len(args))
		}

		ContractNo := args[0]
		EstimatedVesselArrivalDate := args[1]
		UpdateTime := args[2]
		ANSubmittedTime := args[3]
		


		
		

		// Insert a row
	ok, err := stub.InsertRow("ArrivalNote", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: "BL"}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractNo}},
			&shim.Column{Value: &shim.Column_String_{String_: EstimatedVesselArrivalDate}},
			&shim.Column{Value: &shim.Column_String_{String_: UpdateTime}},
			&shim.Column{Value: &shim.Column_String_{String_: ANSubmittedTime}}},
		

	})

	if !ok && err == nil {
		return nil, errors.New("Document already exists.")
	}

		toSend := make ([]string, 2)
		toSend[0] = string(ContractNo)
		toSend[1] = "SubmitAN"
			
		_, err = t.po.UpdatePO(stub, toSend)

				if err != nil {

					return nil, errors.New("Unable to Update PO Status with this ContractID ")
				}	


	return nil, err
}


func (t *AN) GetAN (stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1.")
		}


	ContractNo := args[0]

	// Get the row pertaining to this UID
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: "AN"}}
	columns = append(columns, col1)
	col2 := shim.Column{Value: &shim.Column_String_{String_: ContractNo}}
	columns = append(columns, col2)

	row, err := stub.GetRow("ArrivalNote", columns)
	if err != nil {
		return nil, fmt.Errorf("Error: Failed retrieving document with ContractNo %s. Error %s", ContractNo, err.Error())
	}


	var anJSON ANJSON 

	// GetRows returns empty message if key does not exist
	if len(row.Columns) == 0 {
		
	anJSON.EstimatedVesselArrivalDate = ""
	anJSON.UpdateTime = ""
	anJSON.ANSubmittedTime = ""
	
	} else {
	anJSON.EstimatedVesselArrivalDate = row.Columns[2].GetString_()
	anJSON.UpdateTime = row.Columns[3].GetString_()
	anJSON.ANSubmittedTime = row.Columns[4].GetString_()


	}
	

	jsonAN, err := json.Marshal(anJSON)

	if err != nil {

		return nil, err
	}

	fmt.Println(jsonAN)

 	return jsonAN, nil

	}

