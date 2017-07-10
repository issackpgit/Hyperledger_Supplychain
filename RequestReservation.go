package main

import (
	
	"errors"
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)
type RR struct {
po PO
}

type RRJSON struct {

		ExpCtryForwardingCompany string `json:"ExpCtryForwardingCompany"`
		ExpCtryShippingFirm string `json:"ExpCtryShippingFirm"`
		ImpCtryShippingFirm string `json:"ImpCtryShippingFirm"`
		UpdateTime string `json:"UpdateTime"`
		ReservationRequestedTime string `json:"ReservationRequestedTime"`
		CompanyIdOfExpCtryForwardingCompany string `json:"CompanyIdOfExpCtryForwardingCompany"`
		CompanyIdOfExpCtryShippingFirm string `json:"CompanyIdOfExpCtryShippingFirm"`
		

}

func (t *RR) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// Check if table already exists
	_, err := stub.GetTable("RequestReservation")
	if err == nil {
		// Table already exists; do not recreate
		return nil, nil
	}


	// Create RR Table
	err = stub.CreateTable("RequestReservation", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Type", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "ContractNo", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "ExpCtryForwardingCompany", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ExpCtryShippingFirm", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ImpCtryShippingFirm", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "UpdateTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ReservationRequestedTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "CompanyIdOfExpCtryForwardingCompany", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "CompanyIdOfExpCtryShippingFirm", Type: shim.ColumnDefinition_STRING, Key: false},

	})
	if err != nil {
		return nil, errors.New("Failed creating RequestReservationTable.")
	}


	return nil, nil

	}


//SubmitDoc () inserts a new row in the table

func (t *RR) SubmitDoc(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		
		if len(args) != 8 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 8. Got: %d.", len(args))
		}

		ContractNo := args[0]
		ExpCtryForwardingCompany := args[1]
		ExpCtryShippingFirm := args[2]
		ImpCtryShippingFirm := args[3]
		UpdateTime := args[4]
		ReservationRequestedTime := args[5]
		CompanyIdOfExpCtryForwardingCompany := args[6]
		CompanyIdOfExpCtryShippingFirm := args[7]
		


		// Insert a row
	ok, err := stub.InsertRow("RequestReservation", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: "RR"}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractNo}},
			&shim.Column{Value: &shim.Column_String_{String_: ExpCtryForwardingCompany}},
			&shim.Column{Value: &shim.Column_String_{String_: ExpCtryShippingFirm}},
			&shim.Column{Value: &shim.Column_String_{String_: ImpCtryShippingFirm}},
			&shim.Column{Value: &shim.Column_String_{String_: UpdateTime}},
			&shim.Column{Value: &shim.Column_String_{String_: ReservationRequestedTime}},
			&shim.Column{Value: &shim.Column_String_{String_: CompanyIdOfExpCtryForwardingCompany}},
			&shim.Column{Value: &shim.Column_String_{String_: CompanyIdOfExpCtryShippingFirm}},
			
        }})

	if !ok && err == nil {
		return nil, errors.New("Document already exists.")
	}

	
    		toSend := make ([]string, 2)
			toSend[0] = string(ContractNo)
			toSend[1] = "RRCreated"
			
			_,poErr := t.po.UpdatePO(stub, toSend)
			if poErr != nil {
				return nil, poErr
			} 
    
    
	return nil, err
	}



func (t *RR) GetRR (stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1.")
		}


	ContractNo := args[0]

	// Get the row pertaining to this UID
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: "RR"}}
	columns = append(columns, col1)
	col2 := shim.Column{Value: &shim.Column_String_{String_: ContractNo}}
	columns = append(columns, col2)

	row, err := stub.GetRow("RequestReservation", columns)
	if err != nil {
		return nil, fmt.Errorf("Error: Failed retrieving document with ContractNo %s. Error %s", ContractNo, err.Error())
	}


	var rrJSON RRJSON 


	// GetRows returns empty message if key does not exist
	if len(row.Columns) == 0 {

	rrJSON.ExpCtryForwardingCompany = ""
	rrJSON.ExpCtryShippingFirm = ""
	rrJSON.ImpCtryShippingFirm = ""
	rrJSON.UpdateTime = ""
	rrJSON.ReservationRequestedTime = ""
	rrJSON.CompanyIdOfExpCtryForwardingCompany = ""
	rrJSON.CompanyIdOfExpCtryShippingFirm = ""

		
	} else {

	rrJSON.ExpCtryForwardingCompany = row.Columns[2].GetString_()
	rrJSON.ExpCtryShippingFirm = row.Columns[3].GetString_()
	rrJSON.ImpCtryShippingFirm = row.Columns[4].GetString_()
	rrJSON.UpdateTime = row.Columns[5].GetString_() 
	rrJSON.ReservationRequestedTime = row.Columns[6].GetString_()
	rrJSON.CompanyIdOfExpCtryForwardingCompany = row.Columns[7].GetString_()
	rrJSON.CompanyIdOfExpCtryShippingFirm = row.Columns[8].GetString_()

	}


	jsonRR, err := json.Marshal(rrJSON)

	if err != nil {
		
		return nil, err
	}

	fmt.Println(jsonRR)

 	return jsonRR, nil

	}