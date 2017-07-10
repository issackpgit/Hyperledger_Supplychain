package main

import (
	
	"errors"
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"

	

)
type PO struct {

	
	cl	CL
	epo EPO
	crt CRT


}


type POJSON struct {

		
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
		ProcessStatus   string `json:"ProcessStatus"`
		POSubmittedTime string `json: "POSubmittedTime"`
		PORejectReason string `json:"PORejectReason"`
		PaymentDate string `json:"PaymentDate"`


}
	

	type ContractsList struct {
		ContractNo string `json:"ContractNo"`
	}

	type count struct {
		NumContracts int
	}

	type ListContracts struct{
		poDetail	[]ListPO `json:"listPO"`
	}

	type ListPO struct {

		ContractNo 		string `json:"ContractNo"`
		POSubmittedTime string `json: "POSubmittedTime"`
		ImporterName string `json:"ImporterName"`
		ExporterName string `json:"ExporterName"`
		Currency string `json:"Currency"`
		Amount string `json:"Amount"`
		Commodity string `json:"Commodity"`
		ProcessStatus   string `json:"ProcessStatus"`
		CargoLocation string `json:"CargoLocation"`	
		OfferStatus string `json:"OfferStatus"`   
		
			}

	type POList struct {

			contractNo  []ContractNo `json:"contractNo"`
		}

	type ContractNo struct {

		ContractNo  string `json:"ContractNo"`
	}


	type ClArgs struct {

		ContractNo 		string `json:"ContractNo"`
		CargoLocation 	string `json:"CargoLocation"`

	}  


	


func (t *PO) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// Check if table already exists
	_, err := stub.GetTable("PurchaseOrder")
	if err == nil {
		// Table already exists; do not recreate
		return nil, nil
	}

	// Create PO Table
	err = stub.CreateTable("PurchaseOrder", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Type", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "ContractNo", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "RefNo", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ExporterName", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ImporterName", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Commodity", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "UnitPrice", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Amount", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Currency", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Quantity", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Weight", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "TermsOfTrade", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "TermsOfInsurance", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "TermsOfPayment", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "PackingMethod", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "WayOfTransportation", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "TimeOfShipment", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "PortOfShipment", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "PortOfDischarge", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ProcessStatus", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "POInitialCreateTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "UpdateTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "POCreatedTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "POSubmittedTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "CompanyIdOfExporter", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "CompanyIdOfImporter", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "PORejectReason", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "PaymentDate", Type: shim.ColumnDefinition_STRING, Key: false},

	})
	if err != nil {
		return nil, errors.New("Failed creating PurchaseOrderTable.")
	}

	return nil, nil

}


//SubmitDoc () inserts a new row in the table

func (t *PO) SubmitDoc(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		

		if len(args) != 25 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 25. Got: %d.", len(args))
		}

		//Generate new ContractID
		toSend := make ([]string, 0)
		numbContract, err := t.NumbContracts(stub,toSend)
		if err != nil {
			return nil, err
		} 

		var count count
		err = json.Unmarshal(numbContract, &count)
		if err != nil{
			return nil, err
		}

		NewContractNumb := count.NumContracts + 1
		var str1 string
		
		if NewContractNumb  < 10 {
			str1 =  "CN0000" + strconv.Itoa(NewContractNumb) 
		} else if NewContractNumb < 100 {
			str1 =  "CN000" + strconv.Itoa(NewContractNumb) 
		} else if NewContractNumb < 1000 {
			str1 =  "CN00" + strconv.Itoa(NewContractNumb) 
		} else if NewContractNumb < 10000 {
			str1 =  "CN0" + strconv.Itoa(NewContractNumb) 
		} else if NewContractNumb < 100000 {
			str1 =  "CN" + strconv.Itoa(NewContractNumb) 
		} else {
			return nil, errors.New("Contract number count exceeded the limit.")
		}

		var contractNum ContractNo 
		contractNum.ContractNo = str1

		//End of ContractID generation

		ContractNo := str1
		RefNo := args[1]
		ExporterName := args[2]
		ImporterName := args[3]
		Commodity := args[4]
		UnitPrice := args[5]
		Amount := args[6]
		Currency := args[7]
		Quantity:= args[8]
		Weight := args[9]
		TermsOfTrade := args[10]
		TermsOfInsurance := args[11]
		TermsOfPayment := args[12]
		PackingMethod:= args[13]
		WayOfTransportation := args[14]
		TimeOfShipment := args[15]
		PortOfShipment := args[16]
		PortOfDischarge := args[17]
		//ProcessStatus := args[18]
		POInitialCreateTime := args[18]
		UpdateTime := args[19]
		POCreatedTime := args[20]
		POSubmittedTime := args[21]
		CompanyIdOfExporter := args[22]
		CompanyIdOfImporter := args[23]
		PORejectReason := ""
		PaymentDate := args[24]



		// Insert a row
	ok, err := stub.InsertRow("PurchaseOrder", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: "PO"}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractNo}},
			&shim.Column{Value: &shim.Column_String_{String_: RefNo}},
			&shim.Column{Value: &shim.Column_String_{String_: ExporterName}},
			&shim.Column{Value: &shim.Column_String_{String_: ImporterName}},
			&shim.Column{Value: &shim.Column_String_{String_: Commodity}},
			&shim.Column{Value: &shim.Column_String_{String_: UnitPrice}},
			&shim.Column{Value: &shim.Column_String_{String_: Amount}},
			&shim.Column{Value: &shim.Column_String_{String_: Currency}},
			&shim.Column{Value: &shim.Column_String_{String_: Quantity}},
			&shim.Column{Value: &shim.Column_String_{String_: Weight}},
			&shim.Column{Value: &shim.Column_String_{String_: TermsOfTrade}},
			&shim.Column{Value: &shim.Column_String_{String_: TermsOfInsurance}},
			&shim.Column{Value: &shim.Column_String_{String_: TermsOfPayment}},
			&shim.Column{Value: &shim.Column_String_{String_: PackingMethod}},
			&shim.Column{Value: &shim.Column_String_{String_: WayOfTransportation}},
			&shim.Column{Value: &shim.Column_String_{String_: TimeOfShipment}},
			&shim.Column{Value: &shim.Column_String_{String_: PortOfShipment}},
			&shim.Column{Value: &shim.Column_String_{String_: PortOfDischarge}},
			&shim.Column{Value: &shim.Column_String_{String_: "P/O Submitted"}},
			&shim.Column{Value: &shim.Column_String_{String_: POInitialCreateTime}},
			&shim.Column{Value: &shim.Column_String_{String_: UpdateTime}},
			&shim.Column{Value: &shim.Column_String_{String_: POCreatedTime}},
			&shim.Column{Value: &shim.Column_String_{String_: POSubmittedTime}},
			&shim.Column{Value: &shim.Column_String_{String_: CompanyIdOfExporter}},
			&shim.Column{Value: &shim.Column_String_{String_: CompanyIdOfImporter}},
			&shim.Column{Value: &shim.Column_String_{String_: PORejectReason}},
			&shim.Column{Value: &shim.Column_String_{String_: PaymentDate}}},

	})

	if !ok && err == nil {
		return nil, errors.New("Document already exists.")
	}

			

			toSend = make ([]string, 3)
			toSend[0] = string(ContractNo)
			toSend[1] = "Exporter"
			toSend[2] = POCreatedTime

			
			_,clErr := t.cl.UpdateCargoLocation(stub, toSend)
			if clErr != nil {
				
				return nil, clErr
			}


			

			//Insert data in EPO table

			Cont := make ([]string, 1)
			Cont[0] = string(ContractNo)

			_, Cerr := t.epo.SubmitDoc(stub, Cont) 

			if Cerr != nil {
				return nil, Cerr
			}

			_, CrtErr := t.crt.SubmitDoc(stub, Cont) 

			if CrtErr != nil {
				return nil, CrtErr
			}

			//return Contractnumber 
			jsonNewContract, err := json.Marshal(contractNum)

			if err != nil {
				return nil, err
			}

			fmt.Println(jsonNewContract)

 			return jsonNewContract, nil
			//return nil, err
		
	}


func (t *PO) ProcessChange(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		
		if len(args) != 3 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 3. Got: %d.", len(args))
		}

		ContractNo := args[0]
		Status := args[1]
		Time := args[2]
		

		SChange := make ([]string, 2)
		SChange[0] = string(ContractNo)
		SChange[1] = Status

		_, err := t.UpdatePO(stub,SChange)

		if err != nil {
			return nil, err
		}

		CRTChange := make ([]string, 3)
		CRTChange[0] = string(ContractNo)
		CRTChange[1] = Status
		CRTChange[2] = Time

		
		_, err = t.crt.UpdateDoc(stub,CRTChange)

		if err != nil {
			return nil, err
		}


 		return nil, err

	}





func (t *PO) ReSubmitDoc(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		
		if len(args) != 24 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 24. Got: %d.", len(args))
		}

		ContractNo := args[0]


// Get the row pertaining to this UID
		var columns []shim.Column
		col1 := shim.Column{Value: &shim.Column_String_{String_: "PO"}}
		columns = append(columns, col1)
		col2 := shim.Column{Value: &shim.Column_String_{String_: ContractNo}}
		columns = append(columns, col2)

		row, err := stub.GetRow("PurchaseOrder", columns)
		if err != nil {
		return nil, fmt.Errorf("Error: Failed retrieving document with ContractNo %s. Error %s", ContractNo, err.Error())
		}

		// GetRows returns empty message if key does not exist
		if len(row.Columns) == 0 {
		return nil, err
		}

		
		
		
		RefNo := args[1]
		ExporterName := args[2]
		ImporterName := args[3]
		Commodity := args[4]
		UnitPrice := args[5]
		Amount := args[6]
		Currency := args[7]
		Quantity:= args[8]
		Weight := args[9]
		TermsOfTrade := args[10]
		TermsOfInsurance := args[11]
		TermsOfPayment := args[12]
		PackingMethod:= args[13]
		WayOfTransportation := args[14]
		TimeOfShipment := args[15]
		PortOfShipment := args[16]
		PortOfDischarge := args[17]
		POInitialCreateTime := row.Columns[20].GetString_()
		UpdateTime := args[18]
		POCreatedTime := args[19]
		POSubmittedTime := args[20]
		CompanyIdOfExporter := args[21]
		CompanyIdOfImporter := args[22]
		PORejectReason := ""
		PaymentDate := args[23]




		ok, err := stub.ReplaceRow("PurchaseOrder", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: "PO"}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractNo}},
			&shim.Column{Value: &shim.Column_String_{String_: RefNo}},
			&shim.Column{Value: &shim.Column_String_{String_: ExporterName}},
			&shim.Column{Value: &shim.Column_String_{String_: ImporterName}},
			&shim.Column{Value: &shim.Column_String_{String_: Commodity}},
			&shim.Column{Value: &shim.Column_String_{String_: UnitPrice}},
			&shim.Column{Value: &shim.Column_String_{String_: Amount}},
			&shim.Column{Value: &shim.Column_String_{String_: Currency}},
			&shim.Column{Value: &shim.Column_String_{String_: Quantity}},
			&shim.Column{Value: &shim.Column_String_{String_: Weight}},
			&shim.Column{Value: &shim.Column_String_{String_: TermsOfTrade}},
			&shim.Column{Value: &shim.Column_String_{String_: TermsOfInsurance}},
			&shim.Column{Value: &shim.Column_String_{String_: TermsOfPayment}},
			&shim.Column{Value: &shim.Column_String_{String_: PackingMethod}},
			&shim.Column{Value: &shim.Column_String_{String_: WayOfTransportation}},
			&shim.Column{Value: &shim.Column_String_{String_: TimeOfShipment}},
			&shim.Column{Value: &shim.Column_String_{String_: PortOfShipment}},
			&shim.Column{Value: &shim.Column_String_{String_: PortOfDischarge}},
			&shim.Column{Value: &shim.Column_String_{String_: "P/O Submitted"}},
			&shim.Column{Value: &shim.Column_String_{String_: POInitialCreateTime}},
			&shim.Column{Value: &shim.Column_String_{String_: UpdateTime}},
			&shim.Column{Value: &shim.Column_String_{String_: POCreatedTime}},
			&shim.Column{Value: &shim.Column_String_{String_: POSubmittedTime}},
			&shim.Column{Value: &shim.Column_String_{String_: CompanyIdOfExporter}},
			&shim.Column{Value: &shim.Column_String_{String_: CompanyIdOfImporter}},
			&shim.Column{Value: &shim.Column_String_{String_: PORejectReason}},
			&shim.Column{Value: &shim.Column_String_{String_: PaymentDate}},
			
			
		}})

	if !ok && err == nil {

		return nil, errors.New("Document unable to Update.")
	}

	
	return nil, err
	}



	func (t *PO) UpdatePO(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){

		ContractNo := args[0]
		rejectionComment := "" 
		poStatus := ""

		if len(args) == 3 {
			poStatus = args[2]
			rejectionComment = args[1]

		} else if len(args) == 2 {

		poStatus = args[1]
		

		} else {

			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 2 or 3. Got: %d.", len(args))
		}

		// Get the row pertaining to this UID
		var columns []shim.Column
		col1 := shim.Column{Value: &shim.Column_String_{String_: "PO"}}
		columns = append(columns, col1)
		col2 := shim.Column{Value: &shim.Column_String_{String_: ContractNo}}
		columns = append(columns, col2)

		row, err := stub.GetRow("PurchaseOrder", columns)
		if err != nil {
		return nil, fmt.Errorf("Error: Failed retrieving document with ContractNo %s. Error %s", ContractNo, err.Error())
		}

		// GetRows returns empty message if key does not exist
		if len(row.Columns) == 0 {
		return nil, err
		}


		
		RefNo := row.Columns[2].GetString_()
		ExporterName := row.Columns[3].GetString_()
		ImporterName := row.Columns[4].GetString_()
		Commodity := row.Columns[5].GetString_()
		UnitPrice := row.Columns[6].GetString_()
		Amount := row.Columns[7].GetString_()
		Currency := row.Columns[8].GetString_()
		Quantity:= row.Columns[9].GetString_()
		Weight := row.Columns[10].GetString_()
		TermsOfTrade := row.Columns[11].GetString_()
		TermsOfInsurance := row.Columns[12].GetString_()
		TermsOfPayment := row.Columns[13].GetString_()
		PackingMethod:= row.Columns[14].GetString_()
		WayOfTransportation := row.Columns[15].GetString_()
		TimeOfShipment := row.Columns[16].GetString_()
		PortOfShipment := row.Columns[17].GetString_()
		PortOfDischarge := row.Columns[18].GetString_()
		ProcessStatus := row.Columns[19].GetString_()
		POInitialCreateTime := row.Columns[20].GetString_()
		UpdateTime := row.Columns[21].GetString_()
		POCreatedTime := row.Columns[22].GetString_()
		POSubmittedTime := row.Columns[23].GetString_()
		CompanyIdOfExporter := row.Columns[24].GetString_()
		CompanyIdOfImporter := row.Columns[25].GetString_()
		PORejectReason := rejectionComment
		PaymentDate := row.Columns[27].GetString_()




		var newStatus string

		/*if poStatus == "AcceptPO"{

			newStatus = "P/O Submitted"

		} else 
*/

		if poStatus == "RejectPO"{

			newStatus = "P/O Created"
            
        } else if poStatus == "SubmitBC"{
            
            newStatus = "B/C Submitted"

        }  else if poStatus == "SubmitSN"{
            
            newStatus = "S/N Submitted"

        }  else if poStatus == "ConfirmSN"{
            
            newStatus = "S/N Confirmed"

        } else if poStatus == "RejectSN"{
            
            newStatus = "P/O Submitted"

        }  else if poStatus == "RRCreated"{
            
            newStatus = "Reservation Requested"

        }else if poStatus == "SubmitCISIPL"{
            
            newStatus = "C/I S/I P/L Submitted"

        }else if poStatus == "SubmitDR"{
            
            newStatus = "D/R Submitted"

        }else if poStatus == "SubmitBL"{
            
            newStatus = "B/L Submitted"

        }else if poStatus == "SubmitAN"{
            
            newStatus = "A/N Submitted"

        }else if poStatus == "SubmitCRR"{
            
            newStatus = "Request FWD"

        } else if poStatus == "NotifyImShip"{
            
            newStatus = "Notify to Im Ship"

        } else if poStatus == "CargoReceived"{
            
            newStatus = "Cargo Received"
        }

       
       

		//Start- Check that the currentStatus to newStatus transition is accurate

		stateTransitionAllowed := false

		/*
		if ProcessStatus == "P/O Created" && newStatus == "P/O Submitted" {
		stateTransitionAllowed = true
		} else */


		if ProcessStatus == "P/O Submitted" && newStatus == "P/O Created" {
		stateTransitionAllowed = true
        } else if ProcessStatus == "P/O Submitted" && newStatus == "S/N Submitted"{
            stateTransitionAllowed = true
        } else if ProcessStatus == "S/N Submitted" && newStatus == "S/N Confirmed"{
            stateTransitionAllowed = true
        } else if ProcessStatus == "S/N Submitted" && newStatus == "P/O Submitted"{
            stateTransitionAllowed = true
        } else if ProcessStatus == "S/N Confirmed" && newStatus == "Reservation Requested"{
            stateTransitionAllowed = true
        } else if ProcessStatus == "Reservation Requested" && newStatus == "B/C Submitted"{
            stateTransitionAllowed = true
        } else if ProcessStatus == "B/C Submitted" && newStatus == "C/I S/I P/L Submitted"{
            stateTransitionAllowed = true
        }else if ProcessStatus == "C/I S/I P/L Submitted" && newStatus == "D/R Submitted"{
            stateTransitionAllowed = true
        }else if ProcessStatus == "D/R Submitted" && newStatus == "B/L Submitted"{
            stateTransitionAllowed = true
        }else if ProcessStatus == "B/L Submitted" && newStatus == "A/N Submitted"{
            stateTransitionAllowed = true
        }else if ProcessStatus == "A/N Submitted" && newStatus == "Request FWD"{
            stateTransitionAllowed = true
        }else if ProcessStatus == "Request FWD" && newStatus == "Notify to Im Ship"{
            stateTransitionAllowed = true
        }else if ProcessStatus == "Notify to Im Ship" && newStatus == "Cargo Received"{
            stateTransitionAllowed = true
        }



	if stateTransitionAllowed == false {
		return nil, errors.New("This state transition is not allowed.")
	}

	//End- Check that the currentStatus to newStatus transition is accurate

		ok, err := stub.ReplaceRow("PurchaseOrder", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: "PO"}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractNo}},
			&shim.Column{Value: &shim.Column_String_{String_: RefNo}},
			&shim.Column{Value: &shim.Column_String_{String_: ExporterName}},
			&shim.Column{Value: &shim.Column_String_{String_: ImporterName}},
			&shim.Column{Value: &shim.Column_String_{String_: Commodity}},
			&shim.Column{Value: &shim.Column_String_{String_: UnitPrice}},
			&shim.Column{Value: &shim.Column_String_{String_: Amount}},
			&shim.Column{Value: &shim.Column_String_{String_: Currency}},
			&shim.Column{Value: &shim.Column_String_{String_: Quantity}},
			&shim.Column{Value: &shim.Column_String_{String_: Weight}},
			&shim.Column{Value: &shim.Column_String_{String_: TermsOfTrade}},
			&shim.Column{Value: &shim.Column_String_{String_: TermsOfInsurance}},
			&shim.Column{Value: &shim.Column_String_{String_: TermsOfPayment}},
			&shim.Column{Value: &shim.Column_String_{String_: PackingMethod}},
			&shim.Column{Value: &shim.Column_String_{String_: WayOfTransportation}},
			&shim.Column{Value: &shim.Column_String_{String_: TimeOfShipment}},
			&shim.Column{Value: &shim.Column_String_{String_: PortOfShipment}},
			&shim.Column{Value: &shim.Column_String_{String_: PortOfDischarge}},
			&shim.Column{Value: &shim.Column_String_{String_: newStatus}},
			&shim.Column{Value: &shim.Column_String_{String_: POInitialCreateTime}},
			&shim.Column{Value: &shim.Column_String_{String_: UpdateTime}},
			&shim.Column{Value: &shim.Column_String_{String_: POCreatedTime}},
			&shim.Column{Value: &shim.Column_String_{String_: POSubmittedTime}},
			&shim.Column{Value: &shim.Column_String_{String_: CompanyIdOfExporter}},
			&shim.Column{Value: &shim.Column_String_{String_: CompanyIdOfImporter}},
			&shim.Column{Value: &shim.Column_String_{String_: PORejectReason}},
			&shim.Column{Value: &shim.Column_String_{String_: PaymentDate}},
			
			
		}})

	if !ok && err == nil {

		return nil, errors.New("Document unable to Update.")
	}



		/*

		if poStatus == "RRCreated"{
			toSend := make ([]string, 2)
			toSend[0] = string(ContractNo)
			toSend[1] = "Ex FWD"

			
			_,clErr := t.cl.UpdateCargoLocation(stub, toSend)
			if clErr != nil {
				return nil, clErr
			} 
		}  else if poStatus == "SubmitBC"{
			toSend := make ([]string, 2)
			toSend[0] = string(ContractNo)
			toSend[1] = "Ex Ship"

			
			_,clErr := t.cl.UpdateCargoLocation(stub, toSend)
			if clErr != nil {
				return nil, clErr
			} 
		} else if poStatus == "SubmitCISIPL"{
			toSend := make ([]string, 2)
			toSend[0] = string(ContractNo)
			toSend[1] = "Shipping"

			
			_,clErr := t.cl.UpdateCargoLocation(stub, toSend)
			if clErr != nil {
				return nil, clErr
			} 
		} else if poStatus == "SubmitAN"{
			toSend := make ([]string, 2)
			toSend[0] = string(ContractNo)
			toSend[1] = "Imp Ship"

			
			_,clErr := t.cl.UpdateCargoLocation(stub, toSend)
			if clErr != nil {
				return nil, clErr
			} 
		} else if poStatus == "SubmitCRR"{
			toSend := make ([]string, 2)
			toSend[0] = string(ContractNo)
			toSend[1] = "Imp FWD"

			
			_,clErr := t.cl.UpdateCargoLocation(stub, toSend)
			if clErr != nil {
				return nil, clErr
			} 
		}  else if poStatus == "CargoReceived"{
			toSend := make ([]string, 2)
			toSend[0] = string(ContractNo)
			toSend[1] = "Importer"

			
			_,clErr := t.cl.UpdateCargoLocation(stub, toSend)
			if clErr != nil {
				return nil, clErr
			} 

			
		}  */

        
		return nil, nil
}


	func (t *PO) GetContractStatus(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
		if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1.")
		}


		var poJSON POJSON

		podetails,_ := t.GetPO(stub, []string{args[0]})

		err := json.Unmarshal(podetails, &poJSON)
		if err != nil{

			return nil, err

		}

		
		return []byte(poJSON.ProcessStatus), nil

	}


	func (t *PO) GetPO (stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1.")
		}

	


	ContractNo := args[0]

	// Get the row pertaining to this UID
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: "PO"}}
	columns = append(columns, col1)
	col2 := shim.Column{Value: &shim.Column_String_{String_: ContractNo}}
	columns = append(columns, col2)

	row, err := stub.GetRow("PurchaseOrder", columns)
	if err != nil {
		return nil, fmt.Errorf("Error: Failed retrieving document with ContractNo %s. Error %s", ContractNo, err.Error())
	}

	// GetRows returns empty message if key does not exist
	if len(row.Columns) == 0 {
		return nil, nil
	}

	var poJSON POJSON 

	poJSON.RefNo = row.Columns[2].GetString_()
	poJSON.ExporterName = row.Columns[3].GetString_()
	poJSON.ImporterName = row.Columns[4].GetString_()
	poJSON.Commodity = row.Columns[5].GetString_() 
	poJSON.UnitPrice = row.Columns[6].GetString_()
	poJSON.Amount = row.Columns[7].GetString_()
	poJSON.Currency = row.Columns[8].GetString_()
	poJSON.Quantity = row.Columns[9].GetString_()
	poJSON.Weight = row.Columns[10].GetString_()
	poJSON.TermsOfTrade = row.Columns[11].GetString_()
	poJSON.TermsOfInsurance = row.Columns[12].GetString_()
	poJSON.TermsOfPayment = row.Columns[13].GetString_()
	poJSON.PackingMethod = row.Columns[14].GetString_()
	poJSON.WayOfTransportation = row.Columns[15].GetString_()
	poJSON.TimeOfShipment = row.Columns[16].GetString_()
	poJSON.PortOfShipment = row.Columns[17].GetString_()
	poJSON.PortOfDischarge = row.Columns[18].GetString_()
	poJSON.ProcessStatus = row.Columns[19].GetString_()
	poJSON.POSubmittedTime = row.Columns[23].GetString_()
	poJSON.PORejectReason = row.Columns[26].GetString_()
	poJSON.PaymentDate = row.Columns[27].GetString_()
	

	jsonPO, err := json.Marshal(poJSON)

	if err != nil {

		return nil, err
	}

	fmt.Println(jsonPO)

 	return jsonPO, nil

	}



func (t *PO) ListContractsByCompID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2.")
		}



		companyID := args[0]
		roleID := args[1]
	
		var listContracts ListContracts

		listContracts.poDetail = make([]ListPO, 0)
		

		if roleID == "1" || roleID == "4" {

			
			var columns []shim.Column
			col1 := shim.Column{Value: &shim.Column_String_{String_: "PO"}}
			columns = append(columns, col1)
			

			rows,err := stub.GetRows("PurchaseOrder", columns)
			
			if err != nil {
			return nil, fmt.Errorf("Error: Failed retrieving document Error %s", err.Error())
			}

			
			
			var contractIDOfUser ContractsList

			

				for row := range rows {
					
					contractIDOfUser.ContractNo = ""
				
					if len(row.Columns) == 0 { 

						break 
					
					} else if roleID == "1" {
						if row.Columns[24].GetString_() == companyID && row.Columns[19].GetString_() != "P/O Created"{

						CToSend := make ([]string, 1)
						CToSend[0] = row.Columns[1].GetString_()

						
						//get PO Table Contract status and if status is Cargo Received don't show that Contract to Exporter
						b, err := t.GetContractStatus(stub, CToSend)

						if err != nil {
							return nil, err
						}

						if string(b) != "Cargo Received"{

					 		contractIDOfUser.ContractNo = row.Columns[1].GetString_()
					 	}

						}

					} else if roleID == "4"{

						if row.Columns[25].GetString_() == companyID {

						CToSend := make ([]string, 1)
						CToSend[0] = row.Columns[1].GetString_()

						
						//get PO Table Contract status and if status is Cargo Received don't show that Contract to Importer
						b, err := t.GetContractStatus(stub, CToSend)

						if err != nil {
							return nil, err
						}

						if string(b) != "Cargo Received"{

					 	contractIDOfUser.ContractNo = row.Columns[1].GetString_()

					 	}

						}

					}  


					if contractIDOfUser.ContractNo != "" {

					b,_ := t.GetPO(stub, []string{contractIDOfUser.ContractNo})
					var listPO ListPO
					err = json.Unmarshal(b, &listPO)
					listPO.ContractNo = contractIDOfUser.ContractNo

					
					b,_ = t.cl.GetCargoLocation(stub, []string{contractIDOfUser.ContractNo})

					if err != nil {
        				return nil, err
        			}
					listPO.CargoLocation = string(b)


					var OfferStatusJSON EPOJSON
					b,_ = t.epo.GetEPO(stub, []string{contractIDOfUser.ContractNo})
					if err != nil {
        				return nil, err
        			}

        			err = json.Unmarshal(b, &OfferStatusJSON)

        			if err != nil{
						return nil, err
        			}

        			listPO.OfferStatus = OfferStatusJSON.OfferStatus

					listContracts.poDetail = append(listContracts.poDetail, listPO)

					}

				}

				

		} else if roleID == "2" || roleID == "3" || roleID == "6" {

			

			var columns []shim.Column
			col1 := shim.Column{Value: &shim.Column_String_{String_: "RR"}}
			columns = append(columns, col1)
			

			rows,err := stub.GetRows("RequestReservation", columns)
			
			if err != nil {
			return nil, fmt.Errorf("Error: Failed retrieving document Error %s", err.Error())
			
			}


			var contractIDOfUser ContractsList
			
				for row := range rows {
	
					contractIDOfUser.ContractNo = ""
					if len(row.Columns) == 0 {

						break 
					
					} else if roleID == "2" {
						if row.Columns[7].GetString_() == companyID{

						CToSend := make ([]string, 1)
						CToSend[0] = row.Columns[1].GetString_()

						
						//get PO Table Contract status and if status is Cargo Received don't show that Contract to Exp Cnt For
						b, err := t.GetContractStatus(stub, CToSend)

						if err != nil {
							return nil, err
						}

							if string(b) != "Cargo Received"{

					 		contractIDOfUser.ContractNo = row.Columns[1].GetString_()

					 		}
						}
					} else if roleID == "3" {


						if row.Columns[8].GetString_() == companyID{

						CToSend := make ([]string, 1)
						CToSend[0] = row.Columns[1].GetString_()

						
						//get PO Table Contract status and if status is Cargo Received don't show that Contract to Exp cty Shi
						b, err := t.GetContractStatus(stub, CToSend)

						if err != nil {
							return nil, err
						}

							if string(b) != "Cargo Received"{

					 			contractIDOfUser.ContractNo = row.Columns[1].GetString_()
					 		}
					 	
						}

					} else if roleID == "6"{



						if row.Columns[4].GetString_() == companyID{

						CToSend := make ([]string, 1)
						CToSend[0] = row.Columns[1].GetString_()

						
						//get PO Table Contract status and if status is Reservation Requested and Cargo Received don't show that Contract to Imp.Fwd
						b, err := t.GetContractStatus(stub, CToSend)

						if err != nil {
							return nil, err
						}

						if string(b) != "Reservation Requested" && string(b) != "Cargo Received" {

					 	contractIDOfUser.ContractNo = row.Columns[1].GetString_()

					 	}

					 	
						}

					} 

					if contractIDOfUser.ContractNo != "" {

					//var poJSON POJSON
					b,_ := t.GetPO(stub, []string{contractIDOfUser.ContractNo})
					var listPO ListPO
					err = json.Unmarshal(b, & listPO)
					listPO.ContractNo = contractIDOfUser.ContractNo


					b,_ = t.cl.GetCargoLocation(stub, []string{contractIDOfUser.ContractNo})
				
					listPO.CargoLocation = string(b)

					var OfferStatusJSON EPOJSON
					b,_ = t.epo.GetEPO(stub, []string{contractIDOfUser.ContractNo})
					if err != nil {
        				return nil, err
        			}

        			err = json.Unmarshal(b, &OfferStatusJSON)

        			if err != nil{
						return nil, err
        			}

        			listPO.OfferStatus = OfferStatusJSON.OfferStatus
					listContracts.poDetail = append(listContracts.poDetail, listPO)

				}

				}

		} else if roleID == "5" {

			var columns []shim.Column
			col1 := shim.Column{Value: &shim.Column_String_{String_: "CRR"}}
			columns = append(columns, col1)
			

			rows,err := stub.GetRows("CargoReceiveRequest", columns)
			
			if err != nil {
			return nil, fmt.Errorf("Error: Failed retrieving document Error %s", err.Error())
			}

			var contractIDOfUser ContractsList
			

				for row := range rows {

					contractIDOfUser.ContractNo = ""

					if len(row.Columns) == 0 {
					
						
						break 
					
					} else if roleID == "5" {
						if row.Columns[5].GetString_() == companyID{

						CToSend := make ([]string, 1)
						CToSend[0] = row.Columns[1].GetString_()

						
						//get PO Table Contract status and if status is Cargo Received don't show that Contract to Imp Cty Frw
						b, err := t.GetContractStatus(stub, CToSend)

						if err != nil {
							return nil, err
						}

							if string(b) != "Cargo Received"{

					 		contractIDOfUser.ContractNo = row.Columns[1].GetString_()

					 	}
					 		
						}
					}  

					if contractIDOfUser.ContractNo != "" {
					b,_ := t.GetPO(stub, []string{contractIDOfUser.ContractNo})
					var listPO ListPO
					err = json.Unmarshal(b, &listPO)
					listPO.ContractNo = contractIDOfUser.ContractNo



					b,_ = t.cl.GetCargoLocation(stub, []string{contractIDOfUser.ContractNo})

					listPO.CargoLocation = string(b)

					var OfferStatusJSON EPOJSON
					b,_ = t.epo.GetEPO(stub, []string{contractIDOfUser.ContractNo})
					if err != nil {
        				return nil, err
        			}

        			err = json.Unmarshal(b, &OfferStatusJSON)

        			if err != nil{
						return nil, err
        			}

        			listPO.OfferStatus = OfferStatusJSON.OfferStatus

					listContracts.poDetail = append(listContracts.poDetail, listPO)

					}

				}
			
		}


		return json.Marshal(listContracts.poDetail) 

	}


	func (t *PO) NumbContracts(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0.")
	}

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: "PO"}}
	columns = append(columns, col1)

	
	contractCounter := 0


	rows, err := stub.GetRows("PurchaseOrder", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve row")
	}

	for row := range rows {
		if len(row.Columns) != 0 {
			contractCounter++
		}
	}

	var c count
	c.NumContracts = contractCounter

	return json.Marshal(c)


	}


	func (t *PO) ListOfContracts(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
		if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0.")
		}

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: "PO"}}
	columns = append(columns, col1)
	

	rows, err := stub.GetRows("PurchaseOrder", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve row")
	}

 	var listOfContracts POList
 	var contractNo ContractNo 

 	listOfContracts.contractNo = make([]ContractNo, 0)


	for row := range rows {

		
		 contractNo.ContractNo = row.Columns[1].GetString_()

        listOfContracts.contractNo = append(listOfContracts.contractNo, contractNo)

	}

		return json.Marshal(listOfContracts.contractNo)

	}


	func (t *PO) GetContractReceivedList(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2.")
		}



		companyID := args[0]
		roleID := args[1]
	
		var listContracts ListContracts

		listContracts.poDetail = make([]ListPO, 0)
		

		if roleID == "1" || roleID == "4" {

			
			var columns []shim.Column
			col1 := shim.Column{Value: &shim.Column_String_{String_: "PO"}}
			columns = append(columns, col1)
			

			rows,err := stub.GetRows("PurchaseOrder", columns)
			
			if err != nil {
			return nil, fmt.Errorf("Error: Failed retrieving document Error %s", err.Error())
			}

			
			
			var contractIDOfUser ContractsList

			

				for row := range rows {
					
					contractIDOfUser.ContractNo = ""
				
					if len(row.Columns) == 0 { 

						break 
					
					} else if roleID == "1" {
						if row.Columns[24].GetString_() == companyID && row.Columns[19].GetString_() == "Cargo Received"{

					 		contractIDOfUser.ContractNo = row.Columns[1].GetString_()
					 
						}

					} else if roleID == "4"{

						if row.Columns[25].GetString_() == companyID && row.Columns[19].GetString_() == "Cargo Received"{
					 	contractIDOfUser.ContractNo = row.Columns[1].GetString_()

						}

					}  


					if contractIDOfUser.ContractNo != "" {

					b,_ := t.GetPO(stub, []string{contractIDOfUser.ContractNo})
					var listPO ListPO
					err = json.Unmarshal(b, &listPO)
					listPO.ContractNo = contractIDOfUser.ContractNo

					
					b,_ = t.cl.GetCargoLocation(stub, []string{contractIDOfUser.ContractNo})

					if err != nil {
        				return nil, err
        			}
					listPO.CargoLocation = string(b)


					var OfferStatusJSON EPOJSON
					b,_ = t.epo.GetEPO(stub, []string{contractIDOfUser.ContractNo})
					if err != nil {
        				return nil, err
        			}

        			err = json.Unmarshal(b, &OfferStatusJSON)

        			if err != nil{
						return nil, err
        			}

        			listPO.OfferStatus = OfferStatusJSON.OfferStatus

					listContracts.poDetail = append(listContracts.poDetail, listPO)

					}

				}

				

		} else if roleID == "2" || roleID == "3" || roleID == "6" {

			

			var columns []shim.Column
			col1 := shim.Column{Value: &shim.Column_String_{String_: "RR"}}
			columns = append(columns, col1)
			

			rows,err := stub.GetRows("RequestReservation", columns)
			
			if err != nil {
			return nil, fmt.Errorf("Error: Failed retrieving document Error %s", err.Error())
			
			}


			var contractIDOfUser ContractsList
			
				for row := range rows {
	
					contractIDOfUser.ContractNo = ""
					if len(row.Columns) == 0 {

						break 
					
					} else if roleID == "2" {
						if row.Columns[7].GetString_() == companyID{

						CToSend := make ([]string, 1)
						CToSend[0] = row.Columns[1].GetString_()

						
						//get PO Table Contract status and if status is Cargo Received show that Contract to Exp Cnt For
						b, err := t.GetContractStatus(stub, CToSend)

						if err != nil {
							return nil, err
						}

							if string(b) == "Cargo Received"{

					 		contractIDOfUser.ContractNo = row.Columns[1].GetString_()

					 		}
						}
					} else if roleID == "3" {


						if row.Columns[8].GetString_() == companyID{

						CToSend := make ([]string, 1)
						CToSend[0] = row.Columns[1].GetString_()

						
						//get PO Table Contract status and if status is Cargo Received show that Contract to Exp cty Shi
						b, err := t.GetContractStatus(stub, CToSend)

						if err != nil {
							return nil, err
						}

							if string(b) == "Cargo Received"{

					 			contractIDOfUser.ContractNo = row.Columns[1].GetString_()
					 		}
					 	
						}

					} else if roleID == "6"{



						if row.Columns[4].GetString_() == companyID{

						CToSend := make ([]string, 1)
						CToSend[0] = row.Columns[1].GetString_()

						
						//get PO Table Contract status and if status is Cargo Received show that Contract to Imp.Fwd
						b, err := t.GetContractStatus(stub, CToSend)

						if err != nil {
							return nil, err
						}

						if string(b) == "Cargo Received" {

					 	contractIDOfUser.ContractNo = row.Columns[1].GetString_()

					 	}

					 	
						}

					} 

					if contractIDOfUser.ContractNo != "" {

					//var poJSON POJSON
					b,_ := t.GetPO(stub, []string{contractIDOfUser.ContractNo})
					var listPO ListPO
					err = json.Unmarshal(b, & listPO)
					listPO.ContractNo = contractIDOfUser.ContractNo


					b,_ = t.cl.GetCargoLocation(stub, []string{contractIDOfUser.ContractNo})
				
					listPO.CargoLocation = string(b)

					var OfferStatusJSON EPOJSON
					b,_ = t.epo.GetEPO(stub, []string{contractIDOfUser.ContractNo})
					if err != nil {
        				return nil, err
        			}

        			err = json.Unmarshal(b, &OfferStatusJSON)

        			if err != nil{
						return nil, err
        			}

        			listPO.OfferStatus = OfferStatusJSON.OfferStatus
					listContracts.poDetail = append(listContracts.poDetail, listPO)

				}

				}

		} else if roleID == "5" {

			var columns []shim.Column
			col1 := shim.Column{Value: &shim.Column_String_{String_: "CRR"}}
			columns = append(columns, col1)
			

			rows,err := stub.GetRows("CargoReceiveRequest", columns)
			
			if err != nil {
			return nil, fmt.Errorf("Error: Failed retrieving document Error %s", err.Error())
			}

			var contractIDOfUser ContractsList
			

				for row := range rows {

					contractIDOfUser.ContractNo = ""

					if len(row.Columns) == 0 {
					
						
						break 
					
					} else if roleID == "5" {
						if row.Columns[5].GetString_() == companyID{

						CToSend := make ([]string, 1)
						CToSend[0] = row.Columns[1].GetString_()

						
						//get PO Table Contract status and if status is Cargo Received show that Contract to Imp Cty Frw
						b, err := t.GetContractStatus(stub, CToSend)

						if err != nil {
							return nil, err
						}

							if string(b) == "Cargo Received"{

					 		contractIDOfUser.ContractNo = row.Columns[1].GetString_()

					 	}
					 		
						}
					}  

					if contractIDOfUser.ContractNo != "" {
					b,_ := t.GetPO(stub, []string{contractIDOfUser.ContractNo})
					var listPO ListPO
					err = json.Unmarshal(b, &listPO)
					listPO.ContractNo = contractIDOfUser.ContractNo



					b,_ = t.cl.GetCargoLocation(stub, []string{contractIDOfUser.ContractNo})

					listPO.CargoLocation = string(b)

					var OfferStatusJSON EPOJSON
					b,_ = t.epo.GetEPO(stub, []string{contractIDOfUser.ContractNo})
					if err != nil {
        				return nil, err
        			}

        			err = json.Unmarshal(b, &OfferStatusJSON)

        			if err != nil{
						return nil, err
        			}

        			listPO.OfferStatus = OfferStatusJSON.OfferStatus

					listContracts.poDetail = append(listContracts.poDetail, listPO)

					}

				}
			
		}


		return json.Marshal(listContracts.poDetail) 

	}


func (t *PO) GetCargoLocationChangeList(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0.")
		}

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: "PO"}}
	columns = append(columns, col1)

	rows, err := stub.GetRows("PurchaseOrder", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve row")
	}

	var listContracts ListContracts
	listContracts.poDetail = make([]ListPO, 0)
	
				for row := range rows {

					if len(row.Columns) == 0 { 

						break 
					
					} else {

						CToSend := make ([]string, 1)
						CToSend[0] = row.Columns[1].GetString_()

						b,_ := t.GetPO(stub,CToSend)

						if err != nil {

							return nil, err
						}
						var listPO ListPO
						err = json.Unmarshal(b, &listPO)
						if err != nil {
							return nil, err
						}	

						b,_ = t.cl.GetCargoLocation(stub,CToSend)

						if err != nil {
        				return nil, err
        				}
						
						listPO.CargoLocation = string(b)
						listContracts.poDetail = append(listContracts.poDetail, listPO)

						}
					
					}

		return json.Marshal(listContracts.poDetail) 
}


func (t *PO) UpdatePOAmount(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		
		if len(args) != 2 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 2. Got: %d.", len(args))
		}

		ContractNo := args[0]
		amount := args[1]


// Get the row pertaining to this UID
		var columns []shim.Column
		col1 := shim.Column{Value: &shim.Column_String_{String_: "PO"}}
		columns = append(columns, col1)
		col2 := shim.Column{Value: &shim.Column_String_{String_: ContractNo}}
		columns = append(columns, col2)

		row, err := stub.GetRow("PurchaseOrder", columns)
		if err != nil {
		return nil, fmt.Errorf("Error: Failed retrieving document with ContractNo %s. Error %s", ContractNo, err.Error())
		}

		// GetRows returns empty message if key does not exist
		if len(row.Columns) == 0 {
		return nil, err
		}

		
		
		
		RefNo := row.Columns[2].GetString_()
		ExporterName := row.Columns[3].GetString_()
		ImporterName := row.Columns[4].GetString_()
		Commodity := row.Columns[5].GetString_()
		UnitPrice := row.Columns[6].GetString_()
		Amount := amount
		Currency := row.Columns[8].GetString_()
		Quantity:= row.Columns[9].GetString_()
		Weight := row.Columns[10].GetString_()
		TermsOfTrade := row.Columns[11].GetString_()
		TermsOfInsurance := row.Columns[12].GetString_()
		TermsOfPayment := row.Columns[13].GetString_()
		PackingMethod:= row.Columns[14].GetString_()
		WayOfTransportation := row.Columns[15].GetString_()
		TimeOfShipment := row.Columns[16].GetString_()
		PortOfShipment := row.Columns[17].GetString_()
		PortOfDischarge := row.Columns[18].GetString_()
		ProcessStatus := row.Columns[19].GetString_()
		POInitialCreateTime := row.Columns[20].GetString_()
		UpdateTime := row.Columns[21].GetString_()
		POCreatedTime := row.Columns[22].GetString_()
		POSubmittedTime := row.Columns[23].GetString_()
		CompanyIdOfExporter := row.Columns[24].GetString_()
		CompanyIdOfImporter := row.Columns[25].GetString_()
		PORejectReason := row.Columns[26].GetString_()
		PaymentDate := row.Columns[27].GetString_()




		ok, err := stub.ReplaceRow("PurchaseOrder", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: "PO"}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractNo}},
			&shim.Column{Value: &shim.Column_String_{String_: RefNo}},
			&shim.Column{Value: &shim.Column_String_{String_: ExporterName}},
			&shim.Column{Value: &shim.Column_String_{String_: ImporterName}},
			&shim.Column{Value: &shim.Column_String_{String_: Commodity}},
			&shim.Column{Value: &shim.Column_String_{String_: UnitPrice}},
			&shim.Column{Value: &shim.Column_String_{String_: Amount}},
			&shim.Column{Value: &shim.Column_String_{String_: Currency}},
			&shim.Column{Value: &shim.Column_String_{String_: Quantity}},
			&shim.Column{Value: &shim.Column_String_{String_: Weight}},
			&shim.Column{Value: &shim.Column_String_{String_: TermsOfTrade}},
			&shim.Column{Value: &shim.Column_String_{String_: TermsOfInsurance}},
			&shim.Column{Value: &shim.Column_String_{String_: TermsOfPayment}},
			&shim.Column{Value: &shim.Column_String_{String_: PackingMethod}},
			&shim.Column{Value: &shim.Column_String_{String_: WayOfTransportation}},
			&shim.Column{Value: &shim.Column_String_{String_: TimeOfShipment}},
			&shim.Column{Value: &shim.Column_String_{String_: PortOfShipment}},
			&shim.Column{Value: &shim.Column_String_{String_: PortOfDischarge}},
			&shim.Column{Value: &shim.Column_String_{String_: ProcessStatus}},
			&shim.Column{Value: &shim.Column_String_{String_: POInitialCreateTime}},
			&shim.Column{Value: &shim.Column_String_{String_: UpdateTime}},
			&shim.Column{Value: &shim.Column_String_{String_: POCreatedTime}},
			&shim.Column{Value: &shim.Column_String_{String_: POSubmittedTime}},
			&shim.Column{Value: &shim.Column_String_{String_: CompanyIdOfExporter}},
			&shim.Column{Value: &shim.Column_String_{String_: CompanyIdOfImporter}},
			&shim.Column{Value: &shim.Column_String_{String_: PORejectReason}},
			&shim.Column{Value: &shim.Column_String_{String_: PaymentDate}},
			
			
		}})

	if !ok && err == nil {

		return nil, errors.New("Document unable to Update.")
	}

	
	return nil, err
	}
