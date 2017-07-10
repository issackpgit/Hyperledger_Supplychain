package main

import (
	
	"errors"
	"fmt"
    "encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type DRR struct {
        
}

type DRRJSON struct {

		
		DeviceID 		string `json:"DeviceID"`
		ContractNo string `json:"ContractNo"`
		Email string `json:"Email"`
    
}

type ListDeviceDetails struct {
    deviceDetail    []ListDD `json:"listDD"`
}


type ListDD struct {
    
    DeviceID    string `json:"DeviceID"`
    ContractNo  string `json:"ContractNo"`
    Email       string `json:"Email"`
}


func (t *DRR) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	

_, err := stub.GetTable("DeviceRegister")
	if err == nil {
		// Table already exists; do not recreate
		return nil, nil
	}


	// Create DeviceRegister Table
	err = stub.CreateTable("DeviceRegister", []*shim.ColumnDefinition{
        &shim.ColumnDefinition{Name: "Type", Type: shim.ColumnDefinition_STRING, Key: true},
        &shim.ColumnDefinition{Name: "DeviceID", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "ContractNo", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Email", Type: shim.ColumnDefinition_STRING, Key: false},

	})
	if err != nil {
		return nil, errors.New("Failed creating DeviceRegister Table.")
	}

    return nil, nil
}

func (t* DRR) SubmitDoc(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
    
    if len(args) != 3 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 3. Got: %d.", len(args))
		}
    
        DeviceID := args[0]
        ContractNo := args[1]
		Email := args[2]
		
    // Insert a row
	ok, err := stub.InsertRow("DeviceRegister", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: "DRR"}},
			&shim.Column{Value: &shim.Column_String_{String_: DeviceID}},
            &shim.Column{Value: &shim.Column_String_{String_: ContractNo}},
			&shim.Column{Value: &shim.Column_String_{String_: Email}},
        }})
    
    if !ok && err == nil {
		return nil, errors.New("Document already exists.")
	}

    	return nil, err

    }

func (t *DRR) ReSubmitDoc(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		
		if len(args) != 3 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 3. Got: %d.", len(args))
		}

		DeviceID := args[0]
    
        var columns []shim.Column
		col1 := shim.Column{Value: &shim.Column_String_{String_: "DRR"}}
		columns = append(columns, col1)
		col2 := shim.Column{Value: &shim.Column_String_{String_: DeviceID}}
		columns = append(columns, col2)
    
        row, err := stub.GetRow("DeviceRegister", columns)
		if err != nil {
		return nil, fmt.Errorf("Error: Failed retrieving document with DeviceID %s. Error %s", DeviceID, err.Error())
		}

		// GetRows returns empty message if key does not exist
		if len(row.Columns) == 0 {
		return nil, err
		}
    
        ContractNo := args[1]
		Email := args[2]
		
    ok, err := stub.ReplaceRow("DeviceRegister", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: "DRR"}},
			&shim.Column{Value: &shim.Column_String_{String_: DeviceID}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractNo}},
			&shim.Column{Value: &shim.Column_String_{String_: Email}},
        }})

	if !ok && err == nil {

		return nil, errors.New("Document unable to Update.")
	}

	
	return nil, err
    
    
}

func (t* DRR) Getdrrdata (stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0.")
		}
    
    
    var listDeviceDetails ListDeviceDetails

		listDeviceDetails.deviceDetail = make([]ListDD, 0)
    
    
    
   /* DeviceID := args[0]*/

	// Get the row pertaining to this UID
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: "DRR"}}
	columns = append(columns, col1)
	/*col2 := shim.Column{Value: &shim.Column_String_{String_: DeviceID}}
	columns = append(columns, col2)*/

	rows, err := stub.GetRows("DeviceRegister", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve rows")
	}

    var listDD ListDD
    
    for row := range rows { 
    
	listDD.DeviceID = row.Columns[1].GetString_()
	listDD.ContractNo = row.Columns[2].GetString_()
	listDD.Email = row.Columns[3].GetString_()
                 
        listDeviceDetails.deviceDetail = append(listDeviceDetails.deviceDetail, listDD)
    }
    
    
    jsonDRR, err := json.Marshal(listDeviceDetails.deviceDetail)

	if err != nil {

		return nil, err
	}

	fmt.Println(jsonDRR)

 	return jsonDRR, nil

}

func (t *DRR) GetContractNo(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    
    if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1.")
	}
    DeviceID := args[0]
    
    var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: "DRR"}}
	columns = append(columns, col1)
	col2 := shim.Column{Value: &shim.Column_String_{String_: DeviceID}}
	columns = append(columns, col2)
    
    row, err := stub.GetRow("DeviceRegister", columns)
	if err != nil {
		return nil, fmt.Errorf("Error: Failed retrieving document with DeviceID %s. Error %s", DeviceID, err.Error())
	}

	// GetRows returns empty message if key does not exist
	if len(row.Columns) == 0 {
		return nil, nil
	}

	return []byte(row.Columns[2].GetString_()), nil
    
}