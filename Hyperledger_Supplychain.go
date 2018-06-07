package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/crypto/primitives"
)



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
	epo  	EPO
	crt  	CRT
    iot     IOT
    drr     DRR
}

type DetailContractJSON struct{

		RefNo  string  `json:"RefNo"`
		ImporterName	string  `json:"ImporterName"` 
		ExporterName	string  `json:"ExporterName"` 
		ExpCtryForwardingCompany	string  `json:"ExpCtryForwardingCompany"` 
		ExpCtryShippingFirm	string `json:"ExpCtryShippingFirm"` 
		ImpCtryForwardingCompany	string `json:"ImpCtryForwardingCompany"` 
		ImpCtryShippingFirm	string  `json:"ImpCtryShippingFirm"`
		ProcessStatus  string  `json:"ProcessStatus"`
		CargoLocation string `json:"CargoLocation"`
		Commodity	string  `json:"Commodity"`
		UnitPrice	string `json:"UnitPrice"` 
		Amount	  string  `json:"Amount"`
		Currency	string `json:"Currency"`
		Quantity	string `json:"Quantity"`
		Weight		string `json:"Weight"`
		DateCargoReceived	string `json:"DateCargoReceived"`
		TermsOfPayment  string `json:"TermsOfPayment"` 
		TermsOfTrade string `json:"TermsOfTrade"` 
		TermsOfInsurance  string `json:"TermsOfInsurance"` 
		ContainerNo string `json:"ContainerNo"` 
		PlaceOfDelivery  string `json:"PlaceOfDelivery"` 
		NumberOfContainers  string `json:"NumberOfContainers"` 
		PackingMethod  string `json:"PackingMethod"` 
		WayOfTransportation string `json:"WayOfTransportation"` 
		TimeOfShipment  string `json:"TimeOfShipment"` 
		PortOfShipment   string `json:"PortOfShipment"` 
		PortOfDischarge  string `json:"PortOfDischarge"` 
		PlaceOfReceipt    string `json:"PlaceOfReceipt"` 
		ExpectedTimeOfDepature   string `json:"ExpectedTimeOfDepature"` 
 		ExpectedTimeOfArrival  string `json:"ExpectedTimeOfArrival"` 
		EstimatedVesselArrivalDate string `json:"EstimatedVesselArrivalDate"` 
		CutOffDateTime string `json:"CutOffDateTime"` 
		VesselName  string `json:"VesselName"` 
		VesselNo  string `json:"VesselNo"` 
		BookingNo  string `json:"BookingNo"` 
		Freight  string `json:"Freight"` 
		FreightPayment  string `json:"FreightPayment"`
		CaseMark   string `json:"CaseMark"` 
		InvoiceNo  string `json:"InvoiceNo"` 
		CISIPLUpdateTime string `json:"CISIPLUpdateTime"`
		CISIPLSubmittedTime string `json:"CISIPLSubmittedTime"`
		UnknownClause string `json:"UnknownClause"` 
		BLNo  string  `json:"BLNo"` 
		POSubmittedTime string  `json:"POSubmittedTime"`

		SNSubmittedTime string `json:"SNSubmittedTime"`
		SNConfirmedTime string `json:"SNConfirmedTime"`
		ReservationRequestedTime string `json:"ReservationRequestedTime"`
		BCSubmittedTime string `json:"BCSubmittedTime"`
		
		DRSubmittedTime string  `json:"DRSubmittedTime"`
		BLSubmittedTime string `json: "BLSubmittedTime"`
		ANSubmittedTime string  `json:"ANSubmittedTime"`
		ImpCtryForwardingTime string `json:"ImpCtryForwardingTime"`
		ImporterShippingTime string `json:"ImporterShippingTime"`
		ImporterReceivedTime string `json:"ImporterReceivedTime"`



		ExpCargoLocationTime string `json:"ExpCargoLocationTime"`
		ExForwarderCargoLocationTime string `json:"ExForwarderCargoLocationTime"`
		ExShippingFirmLocationTime string `json:"ExShippingFirmLocationTime"`
		ShippingCargoLocationTime string `json:"ShippingCargoLocationTime"`
		ImShipCargoLocationTime  string `json:"ImShipCargoLocationTime"`
		ImForwarderCargoLocationTime string `json:"ImForwarderCargoLocationTime"`
		CargoReceivedLocationTime string `json:"CargoReceivedLocationTime"`


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
	t.epo.Init(stub, function, args)
	t.crt.Init(stub, function, args)
    t.iot.Init(stub, function, args)
    t.drr.Init(stub, function, args)

	return nil, nil
}



// Invoke invokes the chaincode
func (t *SMBC) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "submitPO" {
		
		return t.po.SubmitDoc(stub, args)

	} else if function == "rejectPO"{

		args := append(args, "RejectPO")
		return t.po.UpdatePO(stub, args)

	}  else if function == "submitSN"{

		return t.sn.SubmitDoc(stub, args)

	} else if function == "acceptSN"{

			args := append(args, "AcceptSN")
			return t.sn.UpdateSN(stub, args)


	}else if function == "rejectSN"{

			args := append(args, "RejectSN")
			return t.sn.UpdateSN(stub, args)
		
	}else if function == "submitBC"{

		return t.bc.SubmitDoc(stub, args)
	} else if function == "submitRR"{

		return t.rr.SubmitDoc(stub, args)

	}  else if function == "submitDR"{

		return t.dr.SubmitDoc(stub, args)
	} else if function == "submitBL"{

		return t.bl.SubmitDoc(stub, args)
	} else if function == "submitAN"{

		return t.an.SubmitDoc(stub, args)
	} else if function == "submitCRR"{

		return t.crr.SubmitDoc(stub, args)

	} else if function == "reSubmitPO"{

		return t.po.ReSubmitDoc(stub, args)
		
	} else if function == "submitCISIPL"{

		return t.cisipl.SubmitDoc(stub, args)

	}  else if function == "processChange"{

		return t.po.ProcessChange(stub, args)

	} else if function == "initEarlyPaymentOffer"{

		return t.epo.InitEarlyPaymentOffer(stub, args)

	} else if function == "submitEarlyPaymentOffer"{

		return t.epo.SubmitEarlyPaymentOffer(stub, args)

	} else if function == "acceptEarlyPaymentOffer"{

		args := append(args, "Accepted")


			toSend := make ([]string, 1)
			toSend[0] = args[0]
			
			_,err := t.UpdatePO(stub, toSend)
			if err != nil {
				return nil, err
			}

		return t.epo.UpdateEarlyPaymentOffer(stub, args)

	} else if function == "rejectEarlyPaymentOffer"{

		args := append(args, "NoOffer")
		return t.epo.UpdateEarlyPaymentOffer(stub, args)
        
	} else if function == "submitIOTdata"{

		return t.iot.SubmitDoc(stub, args)
        
	} else if function == "deviceRegister"{

		return t.drr.SubmitDoc(stub, args)
        
	}  else if function == "deviceReRegister"{

		return t.drr.ReSubmitDoc(stub, args)
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
    } else if function == "getSN"{
        
        return t.sn.GetSN(stub,args)

    } else if function == "getRR"{
        
        return t.rr.GetRR(stub,args)

    } else if function == "getBC"{
        
        return t.bc.GetBC(stub,args)
    } else if function == "detailContractByContractID"{

    	return t.DetailContractByContractID(stub,args)

    }else if function == "getDR"{

    	return t.dr.GetDR(stub,args)

    }else if function == "getBL"{

    	return t.dr.GetDR(stub,args)

    }else if function == "getAN"{

    	return t.an.GetAN(stub,args)

    }else if function == "getCRR"{

    	return t.crr.GetCRR(stub,args)

    }else if function == "listOfContracts"{

    	return t.po.ListOfContracts(stub,args)

    }else if function == "getCISIPL"{
        
        return t.cisipl.GetCISIPL(stub,args)

    } else if function == "getEPO"{

    	return t.epo.GetEPO(stub,args)
    	
    } else if function == "getCL"{

    	return t.cl.GetCL(stub,args)

    }  else if function == "getContractReceivedList"{

    	return t.po.GetContractReceivedList(stub,args)

    } else if function == "getCargoLocationChangeList"{

    	return t.po.GetCargoLocationChangeList(stub,args)

    } else if function == "getIOTdata"{

    	return t.iot.GetIOTdata(stub,args)

    } else if function == "getdrrdata"{

    	return t.drr.Getdrrdata(stub,args)

    } else if function == "getContractNo"{

    	return t.drr.GetContractNo(stub,args)

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


func (t *SMBC) DetailContractByContractID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
		if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1.")
		}

		

		var detailContractJSON DetailContractJSON
		var poJSON POJSON
		var rrJSON RRJSON
		var snJSON SNJSON 
		var bcJSON BCJSON
		var cisiplJSON CISIPLJSON 
		var crrJSON CRRJSON
		var anJSON  ANJSON
		var blJSON BLJSON
		var drJSON DRJSON
		var clJSON CLJSON
		var crtJSON CRTJSON

		var CLocation string


		b, err := t.po.GetPO (stub, args)
		err = json.Unmarshal(b, &poJSON)
		 if err != nil{

		 	return nil, fmt.Errorf("Unable to retrieve PO data: %s.", err.Error())
		 }

		b, err = t.rr.GetRR (stub, args)

		err = json.Unmarshal(b, &rrJSON)
		 if err != nil{

		 	return nil, fmt.Errorf("Unable to retrieve RR data: %s.", err.Error())
		 }

		 b, err = t.sn.GetSN (stub, args)
		err = json.Unmarshal(b, &snJSON)
		 if err != nil{

		 	return nil, fmt.Errorf("Unable to retrieve SN data: %s.", err.Error())
		 }

		 b, err = t.bc.GetBC (stub, args)
		err = json.Unmarshal(b, &bcJSON)
		 if err != nil{

		 	return nil, fmt.Errorf("Unable to retrieve BC data: %s.", err.Error())
		 }

		 b, err = t.cisipl.GetCISIPL (stub, args)
		err = json.Unmarshal(b, &cisiplJSON)
		 if err != nil{

		 	return nil, fmt.Errorf("Unable to retrieve CISIPL data: %s.", err.Error())
		 }

		b, err = t.crr.GetCRR (stub, args)
		err = json.Unmarshal(b, &crrJSON)
		 if err != nil{

		 	return nil, fmt.Errorf("Unable to retrieve CRR data: %s.", err.Error())
		 }

		 b, err = t.an.GetAN (stub, args)
		err = json.Unmarshal(b, &anJSON)
		 if err != nil{

		 	return nil, fmt.Errorf("Unable to retrieve AN data: %s.", err.Error())
		 }

		 b, err = t.bl.GetBL (stub, args)
		err = json.Unmarshal(b, &blJSON)
		 if err != nil{

		 	return nil, fmt.Errorf("Unable to retrieve BL data: %s.", err.Error())
		 }

		 b, err = t.cl.GetCargoLocation (stub, args)

		 if err != nil {

		 	return nil, fmt.Errorf("Unable to retrieve CL status: %s.", err.Error())
		 }

		 CLocation = string(b)

		 if err != nil{

		 	return nil, fmt.Errorf("Unable to retrieve CL data: %s.", err.Error())
		 }


		b, err = t.dr.GetDR(stub, args)

		 err = json.Unmarshal(b, &drJSON)

		 if err != nil{

		 	return nil, fmt.Errorf("Unable to retrieve DR data: %s.", err.Error())
		 }

		 b, err = t.cl.GetCL(stub, args)

		 err = json.Unmarshal(b, &clJSON)

		 if err != nil{

		 	return nil, fmt.Errorf("Unable to retrieve CL data: %s.", err.Error())
		 }

		 b, err = t.crt.GetCRT(stub, args)

		 err = json.Unmarshal(b, &crtJSON)

		 if err != nil{

		 	return nil, fmt.Errorf("Unable to retrieve CRT data: %s.", err.Error())
		 }

		
		detailContractJSON.RefNo  = poJSON.RefNo
		detailContractJSON.ImporterName	 = poJSON.ImporterName	
		detailContractJSON.ExporterName	 = poJSON.ExporterName
		detailContractJSON.ExpCtryForwardingCompany	= rrJSON.ExpCtryForwardingCompany
		detailContractJSON.ExpCtryShippingFirm	= rrJSON.ExpCtryShippingFirm
		detailContractJSON.ImpCtryForwardingCompany	= crrJSON.ImpCtryForwardingCompany
		detailContractJSON.ImpCtryShippingFirm	= rrJSON.ImpCtryShippingFirm
		detailContractJSON.ProcessStatus  = poJSON.ProcessStatus
		detailContractJSON.CargoLocation  = CLocation
		detailContractJSON.Commodity	 = poJSON.Commodity
		detailContractJSON.UnitPrice	 = poJSON.UnitPrice
		detailContractJSON.Amount	  = poJSON.Amount	
		detailContractJSON.Currency	 = poJSON.Currency
		detailContractJSON.Quantity	 = poJSON.Quantity	 
		detailContractJSON.Weight	 = poJSON.Weight
		detailContractJSON.DateCargoReceived = blJSON.DateCargoReceived
		detailContractJSON.TermsOfPayment   = poJSON.TermsOfPayment
		detailContractJSON.TermsOfTrade  = poJSON.TermsOfTrade 
		detailContractJSON.TermsOfInsurance   = poJSON.TermsOfInsurance
		detailContractJSON.ContainerNo = bcJSON.ContainerNo
		detailContractJSON.PlaceOfDelivery  = bcJSON.PlaceOfDelivery  
		detailContractJSON.NumberOfContainers  = bcJSON.NumberOfContainers 
		detailContractJSON.PackingMethod  = poJSON.PackingMethod
		detailContractJSON.WayOfTransportation  = poJSON.WayOfTransportation 
		detailContractJSON.TimeOfShipment   = poJSON.TimeOfShipment 
		detailContractJSON.PortOfShipment    = poJSON.PortOfShipment
		detailContractJSON.PortOfDischarge   = poJSON.PortOfDischarge 
		detailContractJSON.PlaceOfReceipt   = bcJSON.PlaceOfReceipt 
		detailContractJSON.ExpectedTimeOfDepature   = bcJSON.ExpectedTimeOfDepature  
 		detailContractJSON.ExpectedTimeOfArrival = bcJSON.ExpectedTimeOfArrival  
		detailContractJSON.EstimatedVesselArrivalDate = anJSON.EstimatedVesselArrivalDate 
		detailContractJSON.CutOffDateTime = bcJSON.CutOffDateTime  
		detailContractJSON.VesselName  = bcJSON.VesselName 
		detailContractJSON.VesselNo  = bcJSON.VesselNo  
		detailContractJSON.BookingNo  = bcJSON.BookingNo 
		detailContractJSON.Freight  = bcJSON.Freight  
		detailContractJSON.FreightPayment  = bcJSON.FreightPayment 
		detailContractJSON.CaseMark  = cisiplJSON.CaseMark
		detailContractJSON.InvoiceNo = cisiplJSON.InvoiceNo 
		detailContractJSON.CISIPLUpdateTime = cisiplJSON.UpdateTime
		detailContractJSON.CISIPLSubmittedTime =cisiplJSON.CISIPLSubmittedTime
		detailContractJSON.UnknownClause = blJSON.UnknownClause 
		detailContractJSON.BLNo  = blJSON.BLNo
		detailContractJSON.POSubmittedTime = poJSON.POSubmittedTime
		detailContractJSON.SNSubmittedTime = snJSON.SNSubmittedTime
		detailContractJSON.SNConfirmedTime = snJSON.SNConfirmedTime
		detailContractJSON.ReservationRequestedTime = rrJSON.ReservationRequestedTime
		detailContractJSON.BCSubmittedTime = bcJSON.BCSubmittedTime
		
		detailContractJSON.DRSubmittedTime = drJSON.DRSubmittedTime
		detailContractJSON.BLSubmittedTime = blJSON.BLSubmittedTime
		detailContractJSON.ANSubmittedTime = anJSON.ANSubmittedTime
		detailContractJSON.ImpCtryForwardingTime = crrJSON.ImpCtryForwardingTime
		detailContractJSON.ImporterShippingTime = crtJSON.ImporterShippingTime 
		detailContractJSON.ImporterReceivedTime = crtJSON.ImporterReceivedTime 


		detailContractJSON.ExpCargoLocationTime = clJSON.ExpCargoLocationTime
		detailContractJSON.ExForwarderCargoLocationTime = clJSON.ExForwarderCargoLocationTime
		detailContractJSON.ExShippingFirmLocationTime = clJSON.ExShippingFirmLocationTime
		detailContractJSON.ShippingCargoLocationTime = clJSON.ShippingCargoLocationTime
		detailContractJSON.ImShipCargoLocationTime = clJSON.ImShipCargoLocationTime
		detailContractJSON.ImForwarderCargoLocationTime = clJSON.ImForwarderCargoLocationTime
		detailContractJSON.CargoReceivedLocationTime = clJSON.CargoReceivedLocationTime




		jsonDetailOfContract, err := json.Marshal(detailContractJSON)

		if err != nil {

		return nil, err
		}

	

 	return jsonDetailOfContract, nil



}


func (t *SMBC) UpdatePO(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
		if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1.")
		}

		var epoJSON EPOJSON

		b, err := t.epo.GetEPO(stub, args)

		err = json.Unmarshal(b, &epoJSON)

		 if err != nil{

		 	return nil, fmt.Errorf("Unable to retrieve EPO data: %s.", err.Error())
		 }


		 TotalAmount := epoJSON.TotalAmount
		
		 args = append (args, TotalAmount)

		 _, err = t.po.UpdatePOAmount(stub,args)

		 if err != nil {

		 	return nil, err
		 }

		return nil, err 

	}
