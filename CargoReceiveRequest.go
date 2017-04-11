package main

import (
	
	"errors"
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)
type CRR struct {

po PO
}

type CRRJSON struct {

ImpCtryForwardingCompany  string `json:"ImpCtryForwardingCompany"`
UpdateTime string `json:"UpdateTime"`
CargoRecievedTime string `json:"CargoRecievedTime"`
CompanyIdOfImpCtryForwardingCompany string `json:"CompanyIdOfImpCtryForwardingCompany"`

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


func (t *CRR) SubmitDoc(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		
		if len(args) != 5 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 5. Got: %d.", len(args))
		}

		ContractNo := args[0]
		ImpCtryForwardingCompany := args[1]
		UpdateTime := args[2]
		CargoRecievedTime := args[3]
		CompanyIdOfImpCtryForwardingCompany := args[4]
		
		

		// Insert a row
	ok, err := stub.InsertRow("CargoReceiveRequest", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: "CRR"}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractNo}},
			&shim.Column{Value: &shim.Column_String_{String_: ImpCtryForwardingCompany}},
			&shim.Column{Value: &shim.Column_String_{String_: UpdateTime}},
			&shim.Column{Value: &shim.Column_String_{String_: CargoRecievedTime}},
			&shim.Column{Value: &shim.Column_String_{String_: CompanyIdOfImpCtryForwardingCompany}}},
		

	})

	if !ok && err == nil {
		return nil, errors.New("Document already exists.")
	}

		toSend := make ([]string, 2)
		toSend[0] = string(ContractNo)
		toSend[1] = "SubmitCRR" 
			
		_, err = t.po.UpdatePO(stub, toSend)

				if err != nil {

					return nil, errors.New("Unable to Update PO Status with this ContractID ")
				}	


	return nil, err
}




func (t *CRR) GetCRR (stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1.")
		}


	ContractNo := args[0]

	// Get the row pertaining to this UID
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: "CRR"}}
	columns = append(columns, col1)
	col2 := shim.Column{Value: &shim.Column_String_{String_: ContractNo}}
	columns = append(columns, col2)

	row, err := stub.GetRow("CargoReceiveRequest", columns)
	if err != nil {
		return nil, fmt.Errorf("Error: Failed retrieving document with ContractNo %s. Error %s", ContractNo, err.Error())
	}


	var crrJSON CRRJSON
	// GetRows returns empty message if key does not exist
	if len(row.Columns) == 0 {
		
	crrJSON.ImpCtryForwardingCompany = ""
	crrJSON.UpdateTime = ""
	crrJSON.CargoRecievedTime = ""
	crrJSON.CargoRecievedTime = ""


	} else {

	crrJSON.ImpCtryForwardingCompany = row.Columns[2].GetString_()
	crrJSON.UpdateTime = row.Columns[3].GetString_()
	crrJSON.CargoRecievedTime = row.Columns[4].GetString_()
	crrJSON.CargoRecievedTime = row.Columns[5].GetString_()
	
	}
	

	jsonCRR, err := json.Marshal(crrJSON)

	if err != nil {

		return nil, err
	}

	fmt.Println(jsonCRR)

 	return jsonCRR, nil

	}
