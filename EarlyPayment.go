package main

import (
	
	"errors"
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)
type EPO struct {
 
 
}

type EPOJSON struct {


	EarlyPaymentDate  string `json:"EarlyPaymentDate"`                                
	DaysAccelerated   string `json:"DaysAccelerated"`                                    
	DiscountRate       string `json:"DiscountRate"`                                      
	DiscountAmount     string `json:"DiscountAmount"`                                  
	TotalAmount       string `json:"TotalAmount"`         
	OfferStatus    string `json:"OfferStatus"`
	OfferedTime		string `json:"OfferedTime"`
	UpdateTime		string `json:"UpdateTime"`

	
}


func (t *EPO) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// Check if table already exists
	_, err := stub.GetTable("EarlyPayment")
	if err == nil {
		// Table already exists; do not recreate
		return nil, nil
	}


	// Create AN Table
	err = stub.CreateTable("EarlyPayment", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Type", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "ContractNo", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "EarlyPaymentDate", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "DaysAccelerated", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "DiscountRate", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "DiscountAmount", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "TotalAmount", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "OfferStatus", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "OfferedTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "UpdateTime", Type: shim.ColumnDefinition_STRING, Key: false},
		
	})
	if err != nil {
		return nil, errors.New("Failed creating EarlyPaymentTable.")
	}


	return nil, nil

	}



func (t *EPO) SubmitDoc(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		

		if len(args) != 1 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 1. Got: %d.", len(args))
		}

		ContractNo := args[0]
		
		// Insert a row
		ok, err := stub.InsertRow("EarlyPayment", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: "EPO"}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractNo}},
			&shim.Column{Value: &shim.Column_String_{String_: ""}},
			&shim.Column{Value: &shim.Column_String_{String_: ""}},
			&shim.Column{Value: &shim.Column_String_{String_: ""}},
			&shim.Column{Value: &shim.Column_String_{String_: ""}},
			&shim.Column{Value: &shim.Column_String_{String_: ""}},
			&shim.Column{Value: &shim.Column_String_{String_: "Blank"}},
			&shim.Column{Value: &shim.Column_String_{String_: ""}},
			&shim.Column{Value: &shim.Column_String_{String_: ""}}},

	})

	if !ok && err == nil {
		return nil, errors.New("Document already exists.")
	}

	return nil, err
}


func (t *EPO) InitEarlyPaymentOffer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		
		if len(args) != 1 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 1. Got: %d.", len(args))
		}


		ContractNo := args[0]

	
	ok, err := stub.ReplaceRow("EarlyPayment", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: "EPO"}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractNo}},
			&shim.Column{Value: &shim.Column_String_{String_: ""}},
			&shim.Column{Value: &shim.Column_String_{String_: ""}},
			&shim.Column{Value: &shim.Column_String_{String_: ""}},
			&shim.Column{Value: &shim.Column_String_{String_: ""}},
			&shim.Column{Value: &shim.Column_String_{String_: ""}},
			&shim.Column{Value: &shim.Column_String_{String_: "NoOffer"}},
			&shim.Column{Value: &shim.Column_String_{String_: ""}},
			&shim.Column{Value: &shim.Column_String_{String_: ""}},
			
			
		}})

	if !ok && err == nil {

		return nil, errors.New("Document unable to Update.")
	}

	return nil, err

	}



func (t *EPO) SubmitEarlyPaymentOffer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		
		if len(args) != 7 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 7. Got: %d.", len(args))
		}

		ContractNo 		:= args[0]
		EarlyPaymentDate  := args[1]                             
		DaysAccelerated    := args[2]                                   
		DiscountRate       := args[3]                                      
		DiscountAmount    := args[4]                              
		TotalAmount       := args[5]          
		OfferedTime    := args[6]  
		

		ok, err := stub.ReplaceRow("EarlyPayment", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: "EPO"}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractNo}},
			&shim.Column{Value: &shim.Column_String_{String_: EarlyPaymentDate}},
			&shim.Column{Value: &shim.Column_String_{String_: DaysAccelerated}},
			&shim.Column{Value: &shim.Column_String_{String_: DiscountRate}},
			&shim.Column{Value: &shim.Column_String_{String_: DiscountAmount}},
			&shim.Column{Value: &shim.Column_String_{String_: TotalAmount}},
			&shim.Column{Value: &shim.Column_String_{String_: "Offered"}},
			&shim.Column{Value: &shim.Column_String_{String_: OfferedTime}},
			&shim.Column{Value: &shim.Column_String_{String_: ""}},
			
			
		}})

	if !ok && err == nil {

		return nil, errors.New("Document unable to Update.")
	}

	return nil, err

	}


func (t *EPO) UpdateEarlyPaymentOffer (stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		
		if len(args) != 3 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 3. Got: %d.", len(args))
		}

		ContractNo 		:= args[0]        
		UpdateTime    := args[1] 
		status := args[2] 
		
		// Get the row pertaining to this UID
		var columns []shim.Column
		col1 := shim.Column{Value: &shim.Column_String_{String_: "EPO"}}
		columns = append(columns, col1)
		col2 := shim.Column{Value: &shim.Column_String_{String_: ContractNo}}
		columns = append(columns, col2)

		row, err := stub.GetRow("EarlyPayment", columns)
		if err != nil {
		return nil, fmt.Errorf("Error: Failed retrieving document with ContractNo %s. Error %s", ContractNo, err.Error())
		}

		// GetRows returns empty message if key does not exist
		if len(row.Columns) == 0 {
		return nil, err
		}


		ok, err := stub.ReplaceRow("EarlyPayment", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: "EPO"}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractNo}},
			&shim.Column{Value: &shim.Column_String_{String_: row.Columns[2].GetString_()}},
			&shim.Column{Value: &shim.Column_String_{String_: row.Columns[3].GetString_()}},
			&shim.Column{Value: &shim.Column_String_{String_: row.Columns[4].GetString_()}},
			&shim.Column{Value: &shim.Column_String_{String_: row.Columns[5].GetString_()}},
			&shim.Column{Value: &shim.Column_String_{String_: row.Columns[6].GetString_()}},
			&shim.Column{Value: &shim.Column_String_{String_: status}},
			&shim.Column{Value: &shim.Column_String_{String_: row.Columns[7].GetString_()}},
			&shim.Column{Value: &shim.Column_String_{String_: UpdateTime}},
			
			
		}})

	if !ok && err == nil {

		return nil, errors.New("Document unable to Update.")
	}


	return nil, err
	}





func (t *EPO) GetEPO (stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1.")
		}


	ContractNo := args[0]

	// Get the row pertaining to this UID
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: "EPO"}}
	columns = append(columns, col1)
	col2 := shim.Column{Value: &shim.Column_String_{String_: ContractNo}}
	columns = append(columns, col2)

	row, err := stub.GetRow("EarlyPayment", columns)
	if err != nil {
		return nil, fmt.Errorf("Error: Failed retrieving document with ContractNo %s. Error %s", ContractNo, err.Error())
	}


	var epoJSON EPOJSON 

	// GetRows returns empty message if key does not exist
	if len(row.Columns) == 0 {
		return nil, errors.New("No record found")
	
	} 


	epoJSON.EarlyPaymentDate = row.Columns[2].GetString_()
	epoJSON.DaysAccelerated	= row.Columns[3].GetString_()
	epoJSON.DiscountRate	= row.Columns[4].GetString_()
	epoJSON.DiscountAmount	= row.Columns[5].GetString_()
	epoJSON.TotalAmount		= row.Columns[6].GetString_()
	epoJSON.OfferStatus		= row.Columns[7].GetString_()
	epoJSON.OfferedTime		= row.Columns[8].GetString_()
	epoJSON.UpdateTime		= row.Columns[9].GetString_()

	

	jsonEPO, err := json.Marshal(epoJSON)

	if err != nil {

		return nil, err
	}


 	return jsonEPO, nil

	}



