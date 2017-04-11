package main

import (
	
	"errors"
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	

)
type SN struct {


	po PO
}


type SNJSON struct {

		
		RefNo 		string `json:"RefNo"`
		ExporterName string `json:"ExporterName"`
		ImporterName string `json:"ImporterName"`
		Commodity string `json:"Commodity"`
		UnitPrice string `json:"UnitPrice"`
		Amount string `json:"Amount"`
		Currency string `json:"Currency"`
		Quantity string `json:"Quantity"`
		Weight 	string `json:"Weight"`
		TermsOfTrade string `json:"TermsOfTrade"`
		TermsOfInsurance string `json:"TermsOfInsurance"`
		TermsOfPayment string `json:"TermsOfPayment"`
		PackingMethod	string `json:"PackingMethod"`
		WayOfTransportation string `json:"WayOfTransportation"`
		TimeOfShipment string `json:"TimeOfShipment"`
		PortOfShipment string `json:"PortOfShipment"`
		PortOfDischarge string `json:"PortOfDischarge"`
		SNRejectReason string `json:"SNRejectReason"`
		PaymentDate string `json:"PaymentDate"`

}
	

func (t *SN) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// Check if table already exists
	_, err := stub.GetTable("SalesNote")
	if err == nil {
		// Table already exists; do not recreate
		return nil, nil
	}

	// Create SN Table
	err = stub.CreateTable("SalesNote", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Type", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "ContractNo", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "UpdateTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "SNSubmittedTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "SNConfirmedTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "SNRejectReason", Type: shim.ColumnDefinition_STRING, Key: false},


	})
	if err != nil {
		return nil, errors.New("Failed creating SalesNoteTable.")
	}

	return nil, nil

}


//SubmitDoc () inserts a new row in the table

func (t *SN) SubmitDoc(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		
		if len(args) != 3 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 3. Got: %d.", len(args))
		}

		ContractNo := args[0]
		UpdateTime := args[1]
		SNSubmittedTime := args[2]
		SNConfirmedTime := ""
		SNRejectReason := ""
		

		// Insert a row
	ok, err := stub.InsertRow("SalesNote", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: "SN"}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractNo}},
			&shim.Column{Value: &shim.Column_String_{String_: UpdateTime}},
			&shim.Column{Value: &shim.Column_String_{String_: SNSubmittedTime}},
			&shim.Column{Value: &shim.Column_String_{String_: SNConfirmedTime}},
			&shim.Column{Value: &shim.Column_String_{String_: SNRejectReason}}},
		

	})

	if !ok && err == nil {
		return nil, errors.New("Document already exists.")
	}

	toSend := make ([]string, 2)
		toSend[0] = string(ContractNo)
		toSend[1] = "SubmitSN"
			
		_, err = t.po.UpdatePO(stub, toSend)

				if err != nil {

					return nil, errors.New("Unable to Update PO Status with this ContractID ")
				}	


	return nil, err
}


	func (t *SN) UpdateSN(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){

		if len(args) != 4 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 4. Got: %d.", len(args))
		}

		ContractNo := args[0]
		rejectionComment := "" 
		snStatus := ""
		var UpdateTime string
		var SNConfirmedTime string


		if args[3] == "RejectSN" {
			
			rejectionComment = args[1]
			UpdateTime = args[2]
			snStatus = args[3]

		} else if args[3] == "AcceptSN"{
		
		UpdateTime = args[1]
		SNConfirmedTime = args[2]
		snStatus = args[3]

		} 


		// Get the row pertaining to this UID
		var columns []shim.Column
		col1 := shim.Column{Value: &shim.Column_String_{String_: "SN"}}
		columns = append(columns, col1)
		col2 := shim.Column{Value: &shim.Column_String_{String_: ContractNo}}
		columns = append(columns, col2)

		row, err := stub.GetRow("SalesNote", columns)
		if err != nil {
		return nil, fmt.Errorf("Error: Failed retrieving document with ContractNo %s. Error %s", ContractNo, err.Error())
		}

		// GetRows returns empty message if key does not exist
		if len(row.Columns) == 0 {
		return nil, err
		}



		SNSubmittedTime := row.Columns[3].GetString_()

		ok, err := stub.ReplaceRow("SalesNote", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: "SN"}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractNo}},
			&shim.Column{Value: &shim.Column_String_{String_: UpdateTime}},
			&shim.Column{Value: &shim.Column_String_{String_: SNSubmittedTime}},
			&shim.Column{Value: &shim.Column_String_{String_: SNConfirmedTime}},
			&shim.Column{Value: &shim.Column_String_{String_: rejectionComment}},
			
			
		}})

	if !ok && err == nil {

		return nil, errors.New("Document unable to Update.")
	}
		

		var newStatus string

		if snStatus == "AcceptSN"{

			newStatus = "S/N Confirmed"

		} else if snStatus == "RejectSN"{

			newStatus = "S/N Rejected"
		}

		if newStatus == "S/N Confirmed" {

		toSend := make ([]string, 2)
		toSend[0] = string(ContractNo)
		toSend[1] = "ConfirmSN"
			
		_, err := t.po.UpdatePO(stub, toSend)

				if err != nil {

					return nil, errors.New("Unable to Update PO Status with this ContractID ")
				}		
		} else if newStatus == "S/N Rejected" {

		toSend := make ([]string, 2)
		toSend[0] = string(ContractNo)
		toSend[1] = "RejectSN"
			
		_, err := t.po.UpdatePO(stub, toSend)

				if err != nil {

					return nil, errors.New("Unable to Update PO Status with this ContractID ")
				}
				
		}

	//End- Check that the currentStatus to newStatus transition is accurate
	
	
		return nil, nil
}



	func (t *SN) GetSN (stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1.")
		}

	ContractNo := args[0]


	// Get the row pertaining to this UID
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: "SN"}}
	columns = append(columns, col1)
	col2 := shim.Column{Value: &shim.Column_String_{String_: ContractNo}}
	columns = append(columns, col2)

	row, err := stub.GetRow("SalesNote", columns)
	if err != nil {
		return nil, fmt.Errorf("Error: Failed retrieving document with ContractNo %s. Error %s", ContractNo, err.Error())
	}

	var snJSON SNJSON 
	var poJSON POJSON
	var RReason string

	// GetRows returns empty message if key does not exist
	if len(row.Columns) == 0 {

		RReason = ""

	} else {

		RReason = row.Columns[5].GetString_()


	}

	

	b, err := t.po.GetPO(stub,args)

 	if err != nil {

 		return nil, fmt.Errorf("Error: Failed retrieving document with ContractNo %s. Error %s", ContractNo, err.Error())
 	}
  
  	err = json.Unmarshal(b, &poJSON)
		if err != nil{

			return nil, err

		}

	snJSON.RefNo = poJSON.RefNo
	snJSON.ExporterName = poJSON.ExporterName
	snJSON.ImporterName = poJSON.ImporterName
	snJSON.Commodity = poJSON.Commodity
	snJSON.UnitPrice = poJSON.UnitPrice
	snJSON.Amount = poJSON.Amount
	snJSON.Currency = poJSON.Currency
	snJSON.Quantity = poJSON.Quantity
	snJSON.Weight = poJSON.Weight
	snJSON.TermsOfTrade = poJSON.TermsOfTrade
	snJSON.TermsOfInsurance = poJSON.TermsOfInsurance
	snJSON.TermsOfPayment = poJSON.TermsOfPayment
	snJSON.PackingMethod = poJSON.PackingMethod
	snJSON.WayOfTransportation = poJSON.WayOfTransportation
	snJSON.TimeOfShipment = poJSON.TimeOfShipment
	snJSON.PortOfShipment = poJSON.PortOfShipment
	snJSON.PortOfDischarge = poJSON.PortOfDischarge
	snJSON.SNRejectReason = RReason
	snJSON.PaymentDate = poJSON.PaymentDate
	


	jsonSN, err := json.Marshal(snJSON)

	if err != nil {

		return nil, err
	}

	fmt.Println(jsonSN)

 	return jsonSN, nil

	}



	
