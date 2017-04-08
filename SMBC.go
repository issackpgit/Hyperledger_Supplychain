package main

import (
	//"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/crypto/primitives"
)

const accessControlFlag bool = false

// SMBC is a high level smart contract that SMBCs together business artifact based smart contracts
type SMBC struct {
	po      PO
	sn 		SN
	rr 		RR
	cl 		CL
	bc		BC
	bl 		BL
	an 		AN 
	crr 	CRR 
	cisipl 	CISIPL
	dr 		DR
}


func (t *SMBC) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	
	t.po.Init(stub, function, args)
	t.sn.Init(stub, function, args)
	t.rr.Init(stub, function, args)
	t.cl.Init(stub, function, args)
	t.bc.Init(stub, function, args)
	t.bl.Init(stub, function, args)
	t.an.Init(stub, function, args)
	t.crr.Init(stub, function, args)
	t.cisipl.Init(stub, function, args)
	t.dr.Init(stub, function, args)

	

	return nil, nil
}





// Invoke invokes the chaincode
func (t *SMBC) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "submitPO" {
		
		return t.po.SubmitDoc(stub, args)

	}  else if function == "acceptPO"{


		args := append(args, "AcceptPO")
		return t.po.UpdatePO(stub,args)


	}else if function == "rejectPO"{

		args := append(args, "RejectPO")
		return t.po.UpdatePO(stub, args)
	} else if function == "submitBC"{

		return t.bc.SubmitDoc(stub, args)
	} 

return nil, errors.New("Invalid invoke function name.")
}




// Query callback representing the query of a chaincode
func (t *SMBC) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
 	
 	if function == "getPO"{

 		return t.po.GetPO(stub,args)

 	} else if function == "contractStatus" {

		return t.po.GetContractStatus(stub, args)

	} else if function == "listContractsByCompID"{

		return t.po.ListContractsByCompID(stub, args)
	}else if function == "numbContracts"{

		return t.po.NumbContracts(stub, args)
	} else if function == "getCargoLocation"{

		return t.cl.GetCargoLocation(stub, args)
    } else if function == "getBC"{
        
        return t.bc.GetBC(stub,args)
    }

return nil, errors.New("Invalid query function name.")
}


func main() {
	primitives.SetSecurityLevel("SHA3", 256)
	err := shim.Start(new(SMBC))
	if err != nil {
		fmt.Printf("Error starting SMBC: %s", err)
	}
}