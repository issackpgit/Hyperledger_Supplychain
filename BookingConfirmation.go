package main

import (
	
	"errors"
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)
type BC struct {

    po PO
}

type BCJSON struct {

		
		ContainerNo 		string `json:"ContainerNo"`
		PlaceOfDelivery string `json:"PlaceOfDelivery"`
		NumberOfContainers string `json:"NumberOfContainers"`
		PlaceOfReceipt string `json:"PlaceOfReceipt"`
		ExpectedTimeOfDepature string `json:"ExpectedTimeOfDepature"`
		ExpectedTimeOfArrival string `json:"ExpectedTimeOfArrival"`
		CutOffDateTime string `json:"CutOffDateTime"`
		VesselName string `json:"VesselName"`
		VesselNo 	string `json:"VesselNo"`
		BookingNo string `json:"BookingNo"`
		Freight string `json:"Freight"`
		FreightPayment string `json:"FreightPayment"`
		UpdateTime	string `json:"UpdateTime"`
		BCSubmittedTime string `json:"BCSubmittedTime"`
		


}

func (t *BC) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// Check if table already exists
	_, err := stub.GetTable("BookingConfirmation")
	if err == nil {
		// Table already exists; do not recreate
		return nil, nil
	}


	// Create BC Table
	err = stub.CreateTable("BookingConfirmation", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Type", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "ContractNo", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "ContainerNo", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "PlaceOfDelivery", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "NumberOfContainers", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "PlaceOfReceipt", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ExpectedTimeOfDepature", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ExpectedTimeOfArrival", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "CutOffDateTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "VesselName", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "VesselNo", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "BookingNo", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Freight", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "FreightPayment", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "UpdateTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "BCSubmittedTime", Type: shim.ColumnDefinition_STRING, Key: false},


	})
	if err != nil {
		return nil, errors.New("Failed creating BookingConfirmationTable.")
	}


	return nil, nil

	}

//SubmitDoc () inserts a new row in the table

func (t *BC) SubmitDoc(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		
		if len(args) != 15 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 15. Got: %d.", len(args))
		}

		ContractNo := args[0]
		ContainerNo := args[1]
		PlaceOfDelivery := args[2]
		NumberOfContainers := args[3]
		PlaceOfReceipt := args[4]
		ExpectedTimeOfDepature := args[5]
		ExpectedTimeOfArrival := args[6]
		CutOffDateTime := args[7]
		VesselName:= args[8]
		VesselNo := args[9]
		BookingNo := args[10]
		Freight := args[11]
		FreightPayment := args[12]
		UpdateTime:= args[13]
		BCSubmittedTime := args[14]



		// Insert a row
	ok, err := stub.InsertRow("BookingConfirmation", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: "BC"}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractNo}},
			&shim.Column{Value: &shim.Column_String_{String_: ContainerNo}},
			&shim.Column{Value: &shim.Column_String_{String_: PlaceOfDelivery}},
			&shim.Column{Value: &shim.Column_String_{String_: NumberOfContainers}},
			&shim.Column{Value: &shim.Column_String_{String_: PlaceOfReceipt}},
			&shim.Column{Value: &shim.Column_String_{String_: ExpectedTimeOfDepature}},
			&shim.Column{Value: &shim.Column_String_{String_: ExpectedTimeOfArrival}},
			&shim.Column{Value: &shim.Column_String_{String_: CutOffDateTime}},
			&shim.Column{Value: &shim.Column_String_{String_: VesselName}},
			&shim.Column{Value: &shim.Column_String_{String_: VesselNo}},
			&shim.Column{Value: &shim.Column_String_{String_: BookingNo}},
			&shim.Column{Value: &shim.Column_String_{String_: Freight}},
			&shim.Column{Value: &shim.Column_String_{String_: FreightPayment}},
			&shim.Column{Value: &shim.Column_String_{String_: UpdateTime}},
			&shim.Column{Value: &shim.Column_String_{String_: BCSubmittedTime}},
			
        }})

	if !ok && err == nil {
		return nil, errors.New("Document already exists.")
	}

	
    		toSend := make ([]string, 2)
			toSend[0] = string(ContractNo)
			toSend[1] = "SubmitBC"
			
			_,poErr := t.po.UpdatePO(stub, toSend)
			if poErr != nil {
				return nil, poErr
			} 
    
    
	return nil, err
	}


func (t *BC) GetBC (stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1.")
		}


	ContractNo := args[0]

	// Get the row pertaining to this UID
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: "BC"}}
	columns = append(columns, col1)
	col2 := shim.Column{Value: &shim.Column_String_{String_: ContractNo}}
	columns = append(columns, col2)

	row, err := stub.GetRow("BookingConfirmation", columns)
	if err != nil {
		return nil, fmt.Errorf("Error: Failed retrieving document with ContractNo %s. Error %s", ContractNo, err.Error())
	}

	var bcJSON BCJSON 

	// GetRows returns empty message if key does not exist
	if len(row.Columns) == 0 {
		
	bcJSON.ContainerNo = ""
	bcJSON.PlaceOfDelivery = ""
	bcJSON.NumberOfContainers = ""
	bcJSON.PlaceOfReceipt = "" 
	bcJSON.ExpectedTimeOfDepature = ""
	bcJSON.ExpectedTimeOfArrival = ""
	bcJSON.CutOffDateTime = ""
	bcJSON.VesselName = ""
	bcJSON.VesselNo = ""
	bcJSON.BookingNo = ""
	bcJSON.Freight = ""
	bcJSON.FreightPayment = ""
	bcJSON.UpdateTime = ""
	bcJSON.BCSubmittedTime = ""
	} else {

	

	bcJSON.ContainerNo = row.Columns[2].GetString_()
	bcJSON.PlaceOfDelivery = row.Columns[3].GetString_()
	bcJSON.NumberOfContainers = row.Columns[4].GetString_()
	bcJSON.PlaceOfReceipt = row.Columns[5].GetString_() 
	bcJSON.ExpectedTimeOfDepature = row.Columns[6].GetString_()
	bcJSON.ExpectedTimeOfArrival = row.Columns[7].GetString_()
	bcJSON.CutOffDateTime = row.Columns[8].GetString_()
	bcJSON.VesselName = row.Columns[9].GetString_()
	bcJSON.VesselNo = row.Columns[10].GetString_()
	bcJSON.BookingNo = row.Columns[11].GetString_()
	bcJSON.Freight = row.Columns[12].GetString_()
	bcJSON.FreightPayment = row.Columns[13].GetString_()
	bcJSON.UpdateTime = row.Columns[14].GetString_()
	bcJSON.BCSubmittedTime = row.Columns[15].GetString_()
	
	}

	jsonBC, err := json.Marshal(bcJSON)

	if err != nil {

		return nil, err
	}

	fmt.Println(jsonBC)

 	return jsonBC, nil

	}

