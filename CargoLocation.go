package main

import (
	
	"errors"
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)
type CL struct {

}


type CLJSON struct {

		CargoLocation  string `json:"CargoLocation"`
		ExpCargoLocationTime string `json:"ExpCargoLocationTime"`
		ExForwarderCargoLocationTime string `json:"ExForwarderCargoLocationTime"`
		ExShippingFirmLocationTime string `json:"ExShippingFirmLocationTime"` 
		ShippingCargoLocationTime string `json:"ShippingCargoLocationTime"`
		ImShipCargoLocationTime string `json:"ImShipCargoLocationTime"`
		ImForwarderCargoLocationTime string `json:"ImForwarderCargoLocationTime"`
		CargoReceivedLocationTime string `json:"CargoReceivedLocationTime"`


}



func (t *CL) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// Check if table already exists
	_, err := stub.GetTable("CargoLocation")
	if err == nil {
		// Table already exists; do not recreate
		return nil, nil
	}


	// Create CL Table
	err = stub.CreateTable("CargoLocation", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Type", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "ContractNo", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "CargoLocation", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ExpCargoLocationTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ExForwarderCargoLocationTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ExShippingFirmLocationTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ShippingCargoLocationTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ImShipCargoLocationTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ImForwarderCargoLocationTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "CargoReceivedLocationTime", Type: shim.ColumnDefinition_STRING, Key: false},

	})
	if err != nil {
		return nil, errors.New("Failed creating CargoLocationTable.")
	}

	return nil, nil

	}



func (t *CL) GetCargoLocation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
		if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1.")
		}

		ContractNo := args[0]

		// Get the row pertaining to this UID
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: "CL"}}
	columns = append(columns, col1)
	col2 := shim.Column{Value: &shim.Column_String_{String_: ContractNo}}
	columns = append(columns, col2)

	row, err := stub.GetRow("CargoLocation", columns)
	if err != nil {
		return nil, fmt.Errorf("Error: Failed retrieving document with ContractNo %s. Error %s", ContractNo, err.Error())
	}

	

	return []byte(row.Columns[2].GetString_()), nil

	}



func (t *CL) UpdateCargoLocation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){

		if len(args) != 3 {

		return nil, errors.New("Incorrect number of arguments. Expecting 3 args.")
		} 


		ExpCargoLocationTime := ""
		ExForwarderCargoLocationTime:= ""
		ExShippingFirmLocationTime := ""
		ShippingCargoLocationTime := ""
		ImShipCargoLocationTime := ""
		ImForwarderCargoLocationTime := ""
		CargoReceivedLocationTime := ""



		ContractNo := args[0]
		CargoLocation := args[1]

		if CargoLocation == "Exporter"{

			ExpCargoLocationTime = args[2]

		} 

		
		//var CurCargoLoc string
		//stateTransitionAllowed := false
		


// Insert a row
	ok, err := stub.InsertRow("CargoLocation", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: "CL"}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractNo}},
			&shim.Column{Value: &shim.Column_String_{String_: CargoLocation}},
			&shim.Column{Value: &shim.Column_String_{String_: ExpCargoLocationTime}},
			&shim.Column{Value: &shim.Column_String_{String_: ExForwarderCargoLocationTime}},
			&shim.Column{Value: &shim.Column_String_{String_: ExShippingFirmLocationTime}},
			&shim.Column{Value: &shim.Column_String_{String_: ShippingCargoLocationTime}},
			&shim.Column{Value: &shim.Column_String_{String_: ImShipCargoLocationTime}},
			&shim.Column{Value: &shim.Column_String_{String_: ImForwarderCargoLocationTime}},
			&shim.Column{Value: &shim.Column_String_{String_: CargoReceivedLocationTime}}},

	})

		if !ok && err == nil {

/*

			CNumb := make ([]string, 1)
			CNumb[0] = string(ContractNo)

			b, err := t.GetCargoLocation(stub, CNumb)

			if err != nil {

				return nil, errors.New("Unable to get Cargo Location.")
			}

		CurCargoLoc = string(b)

		if CurCargoLoc == "Exporter" && CargoLocation == "Ex FWD" {
		stateTransitionAllowed = true
        } else if CurCargoLoc == "Ex FWD" && CargoLocation == "Ex Ship"{
            stateTransitionAllowed = true
        } else if CurCargoLoc == "Ex Ship" && CargoLocation == "Shipping"{
            stateTransitionAllowed = true
        } else if CurCargoLoc == "Shipping" && CargoLocation == "Imp Ship"{
            stateTransitionAllowed = true
        } else if CurCargoLoc == "Imp Ship" && CargoLocation == "Imp FWD"{
            stateTransitionAllowed = true
        } else if CurCargoLoc == "Imp FWD" && CargoLocation == "Importer"{
            stateTransitionAllowed = true
        }


        if stateTransitionAllowed == false {
		return nil, errors.New("This state transition is not allowed.")
		}

*/

		// Get the row pertaining to this UID
		var columns []shim.Column
		col1 := shim.Column{Value: &shim.Column_String_{String_: "CL"}}
		columns = append(columns, col1)
		col2 := shim.Column{Value: &shim.Column_String_{String_: ContractNo}}
		columns = append(columns, col2)

		row, err := stub.GetRow("CargoLocation", columns)
		if err != nil {
		return nil, fmt.Errorf("Error: Failed retrieving document with ContractNo %s. Error %s", ContractNo, err.Error())
		}


		ExpCargoLocationTime = row.Columns[3].GetString_()
		ExForwarderCargoLocationTime = row.Columns[4].GetString_()
		ExShippingFirmLocationTime = row.Columns[5].GetString_()
		ShippingCargoLocationTime = row.Columns[6].GetString_()
		ImShipCargoLocationTime = row.Columns[7].GetString_()
		ImForwarderCargoLocationTime = row.Columns[8].GetString_()
		CargoReceivedLocationTime = row.Columns[9].GetString_()


		
		if CargoLocation == "Ex FWD"{
			ExForwarderCargoLocationTime = args[2]

		}else if CargoLocation == "Ex Ship"{
			ExShippingFirmLocationTime = args[2]

		}else if CargoLocation == "Shipping"{
			ShippingCargoLocationTime = args[2]

		}else if CargoLocation == "Imp Ship"{
			ImShipCargoLocationTime = args[2]

		}else if CargoLocation == "Imp FWD"{
			ImForwarderCargoLocationTime = args[2]

		} else if CargoLocation == "Importer"{
			CargoReceivedLocationTime = args[2]
		}



		
			//When row already exists, update the data
		
		ok, err := stub.ReplaceRow("CargoLocation", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: "CL"}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractNo}},
			&shim.Column{Value: &shim.Column_String_{String_: CargoLocation}},
			&shim.Column{Value: &shim.Column_String_{String_: ExpCargoLocationTime}},
			&shim.Column{Value: &shim.Column_String_{String_: ExForwarderCargoLocationTime}},
			&shim.Column{Value: &shim.Column_String_{String_: ExShippingFirmLocationTime}},
			&shim.Column{Value: &shim.Column_String_{String_: ShippingCargoLocationTime}},
			&shim.Column{Value: &shim.Column_String_{String_: ImShipCargoLocationTime}},
			&shim.Column{Value: &shim.Column_String_{String_: ImForwarderCargoLocationTime}},			
			&shim.Column{Value: &shim.Column_String_{String_: CargoReceivedLocationTime}},

		}})

			if !ok && err == nil {

				return nil, errors.New("Document unable to Update.")
			}

			
		}

		return nil, err
	}


func (t *CL) GetCL (stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1.")
		}


	ContractNo := args[0]

	// Get the row pertaining to this UID
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: "CL"}}
	columns = append(columns, col1)
	col2 := shim.Column{Value: &shim.Column_String_{String_: ContractNo}}
	columns = append(columns, col2)

	row, err := stub.GetRow("CargoLocation", columns)
	if err != nil {
		return nil, fmt.Errorf("Error: Failed retrieving document with ContractNo %s. Error %s", ContractNo, err.Error())
	}


	var clJSON CLJSON 

	// GetRows returns empty message if key does not exist
	if len(row.Columns) == 0 {
		return nil, errors.New("No record found")
	
	} 


		clJSON.CargoLocation  = row.Columns[2].GetString_()
		clJSON.ExpCargoLocationTime  = row.Columns[3].GetString_()
		clJSON.ExForwarderCargoLocationTime = row.Columns[4].GetString_()
		clJSON.ExShippingFirmLocationTime  = row.Columns[5].GetString_()
		clJSON.ShippingCargoLocationTime  = row.Columns[6].GetString_()
		clJSON.ImShipCargoLocationTime  = row.Columns[7].GetString_()
		clJSON.ImForwarderCargoLocationTime = row.Columns[8].GetString_()
		clJSON.CargoReceivedLocationTime = row.Columns[9].GetString_()



	jsonCL, err := json.Marshal(clJSON)

	if err != nil {

		return nil, err
	}


 	return jsonCL, nil

	}
