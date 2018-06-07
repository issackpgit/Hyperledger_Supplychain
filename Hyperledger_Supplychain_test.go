/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (

	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"
	"os"
	"path/filepath"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/container"
	"github.com/hyperledger/fabric/core/crypto"
	"github.com/hyperledger/fabric/core/db"
	"github.com/hyperledger/fabric/core/ledger"
	"github.com/hyperledger/fabric/core/util"
	"github.com/hyperledger/fabric/membersrvc/ca"
	pb "github.com/hyperledger/fabric/protos"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

const (
	chaincodeStartupTimeoutDefault int = 5000
)

var (
	testLogger = logging.MustGetLogger("test")

	lis net.Listener

	administrator crypto.Client
	alice         crypto.Client
	bob           crypto.Client

	server *grpc.Server
	aca    *ca.ACA
	eca    *ca.ECA
	tca    *ca.TCA
	tlsca  *ca.TLSCA
)
	

func TestMain(m *testing.M) {
	removeFolders()
	setup()
	go initMembershipSrvc()

	fmt.Println("Wait for some secs for OBCCA")
	time.Sleep(2 * time.Second)

	go initVP()

	fmt.Println("Wait for some secs for VP")
	time.Sleep(2 * time.Second)

	go initTFChaincode()

	fmt.Println("Wait for some secs for Chaincode")
	time.Sleep(2 * time.Second)

	if err := initClients(); err != nil {
		panic(err)
	}

	fmt.Println("Wait for 5 secs for chaincode to be started")
	time.Sleep(5 * time.Second)

	ret := m.Run()

	closeListenerAndSleep(lis)

	defer removeFolders()
	os.Exit(ret)
}


func deploy(admCert crypto.CertificateHandler) error {
	// Prepare the spec. The metadata includes the role of the users allowed to assign assets
	spec := &pb.ChaincodeSpec{
		Type:                 1,
		ChaincodeID:          &pb.ChaincodeID{Name: "mycc"},
		CtorMsg:              &pb.ChaincodeInput{Args: util.ToChaincodeArgs("init")},
		Metadata:             []byte("issuer"),
		ConfidentialityLevel: pb.ConfidentialityLevel_PUBLIC,
	}

	// First build and get the deployment spec
	var ctx = context.Background()
	chaincodeDeploymentSpec, err := getDeploymentSpec(ctx, spec)
	if err != nil {
		return err
	}

	tid := chaincodeDeploymentSpec.ChaincodeSpec.ChaincodeID.Name

	// Now create the Transactions message and send to Peer.
	transaction, err := administrator.NewChaincodeDeployTransaction(chaincodeDeploymentSpec, tid)
	if err != nil {
		return fmt.Errorf("Error deploying chaincode: %s ", err)
	}

	ledger, err := ledger.GetLedger()
	ledger.BeginTxBatch("1")
	_, _, err = chaincode.Execute(ctx, chaincode.GetChain(chaincode.DefaultChain), transaction)
	if err != nil {
		return fmt.Errorf("Error deploying chaincode: %s", err)
	}
	ledger.CommitTxBatch("1", []*pb.Transaction{transaction}, nil, nil)

	return err
}


func setup() {
	// Conf
	viper.SetConfigName("SMBC") // name of config file (without extension)
	viper.AddConfigPath(".")     // path to look for the config file in
	err := viper.ReadInConfig()  // Find and read the config file
	if err != nil {              // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file [%s] \n", err))
	}

	// Logging
	var formatter = logging.MustStringFormatter(
		`%{color}[%{module}] %{shortfunc} [%{shortfile}] -> %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	logging.SetFormatter(formatter)

	logging.SetLevel(logging.DEBUG, "peer")
	logging.SetLevel(logging.DEBUG, "chaincode")
	logging.SetLevel(logging.DEBUG, "cryptochain")

	// Init the crypto layer
	if err := crypto.Init(); err != nil {
		panic(fmt.Errorf("Failed initializing the crypto layer [%s]", err))
	}

	removeFolders()

	// Start db
	db.Start()
}

func initMembershipSrvc() {
	// ca.LogInit seems to have been removed
	//ca.LogInit(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr, os.Stdout)
	ca.CacheConfiguration() // Cache configuration
	aca = ca.NewACA()
	eca = ca.NewECA(aca)
	tca = ca.NewTCA(eca)
	tlsca = ca.NewTLSCA(eca)

	var opts []grpc.ServerOption
	if viper.GetBool("peer.pki.tls.enabled") {
		// TLS configuration
		creds, err := credentials.NewServerTLSFromFile(
			filepath.Join(viper.GetString("server.rootpath"), "tlsca.cert"),
			filepath.Join(viper.GetString("server.rootpath"), "tlsca.priv"),
		)
		if err != nil {
			panic("Failed creating credentials for Membersrvc: " + err.Error())
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	fmt.Printf("open socket...\n")
	sockp, err := net.Listen("tcp", viper.GetString("server.port"))
	if err != nil {
		panic("Cannot open port: " + err.Error())
	}
	fmt.Printf("open socket...done\n")

	server = grpc.NewServer(opts...)

	aca.Start(server)
	eca.Start(server)
	tca.Start(server)
	tlsca.Start(server)

	fmt.Printf("start serving...\n")
	server.Serve(sockp)
}

func initVP() {
	var opts []grpc.ServerOption
	if viper.GetBool("peer.tls.enabled") {
		creds, err := credentials.NewServerTLSFromFile(viper.GetString("peer.tls.cert.file"), viper.GetString("peer.tls.key.file"))
		if err != nil {
			grpclog.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)

	//lis, err := net.Listen("tcp", viper.GetString("peer.address"))

	//use a different address than what we usually use for "peer"
	//we override the peerAddress set in chaincode_support.go
	peerAddress := "0.0.0.0:40404"
	var err error
	lis, err = net.Listen("tcp", peerAddress)
	if err != nil {
		return
	}

	getPeerEndpoint := func() (*pb.PeerEndpoint, error) {
		return &pb.PeerEndpoint{ID: &pb.PeerID{Name: "testpeer"}, Address: peerAddress}, nil
	}

	ccStartupTimeout := time.Duration(chaincodeStartupTimeoutDefault) * time.Millisecond
	userRunsCC := true

	// Install security object for peer
	var secHelper crypto.Peer
	if viper.GetBool("security.enabled") {
		enrollID := viper.GetString("security.enrollID")
		enrollSecret := viper.GetString("security.enrollSecret")
		var err error

		if viper.GetBool("peer.validator.enabled") {
			testLogger.Debugf("Registering validator with enroll ID: %s", enrollID)
			if err = crypto.RegisterValidator(enrollID, nil, enrollID, enrollSecret); nil != err {
				panic(err)
			}
			testLogger.Debugf("Initializing validator with enroll ID: %s", enrollID)
			secHelper, err = crypto.InitValidator(enrollID, nil)
			if nil != err {
				panic(err)
			}
		} else {
			testLogger.Debugf("Registering non-validator with enroll ID: %s", enrollID)
			if err = crypto.RegisterPeer(enrollID, nil, enrollID, enrollSecret); nil != err {
				panic(err)
			}
			testLogger.Debugf("Initializing non-validator with enroll ID: %s", enrollID)
			secHelper, err = crypto.InitPeer(enrollID, nil)
			if nil != err {
				panic(err)
			}
		}
	}

	pb.RegisterChaincodeSupportServer(grpcServer,
		chaincode.NewChaincodeSupport(chaincode.DefaultChain, getPeerEndpoint, userRunsCC,
			ccStartupTimeout, secHelper))

	grpcServer.Serve(lis)
}

func initTFChaincode() {
	err := shim.Start(new(SMBC))
	if err != nil {
		panic(err)
	}
}

func initClients() error {
	// Administrator
	if err := crypto.RegisterClient("admin", nil, "admin", "6avZQLwcUe9b"); err != nil {
		return err
	}
	var err error
	administrator, err = crypto.InitClient("admin", nil)
	if err != nil {
		return err
	}

	// Alice
	if err := crypto.RegisterClient("alice", nil, "alice", "NPKYL39uKbkj"); err != nil {
		return err
	}
	alice, err = crypto.InitClient("alice", nil)
	if err != nil {
		return err
	}

	// Bob
	if err := crypto.RegisterClient("bob", nil, "bob", "DRJ23pEQl16a"); err != nil {
		return err
	}
	bob, err = crypto.InitClient("bob", nil)
	if err != nil {
		return err
	}

	return nil
}

func closeListenerAndSleep(l net.Listener) {
	l.Close()
	time.Sleep(2 * time.Second)
}

func getDeploymentSpec(context context.Context, spec *pb.ChaincodeSpec) (*pb.ChaincodeDeploymentSpec, error) {
	fmt.Printf("getting deployment spec for chaincode spec: %v\n", spec)
	var codePackageBytes []byte
	//if we have a name, we don't need to deploy (we are in userRunsCC mode)
	if spec.ChaincodeID.Name == "" {
		var err error
		codePackageBytes, err = container.GetChaincodePackageBytes(spec)
		if err != nil {
			return nil, err
		}
	}
	chaincodeDeploymentSpec := &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec, CodePackage: codePackageBytes}
	return chaincodeDeploymentSpec, nil
}

func removeFolders() {
	fmt.Println("-------------------------")
	if err := os.RemoveAll(viper.GetString("peer.fileSystemPath")); err != nil {
		fmt.Printf("Failed removing [%s] [%s]\n", "hyperledger", err)
	}
}


func TestCC(t *testing.T) {

       

        // Administrator deploy the chaicode
        adminCert, err := administrator.GetTCertificateHandlerNext("role")
        if err != nil {
                t.Fatal(err)
        }

        if err := deploy(adminCert); err != nil {
                t.Fatal(err)
        }

        type count struct {
		NumContracts int `json:"NumContracts"`
		}

		type ContractNo struct {

			ContractNo  string `json:"ContractNo"`
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

		
        /* WORKFLOW 1: Start */
        // Happy path: submit po, accept po, po is positive
        // This must succeed
        //if
         b, err := submitPO(adminCert, "1000", "AA-123","XYZ","ABC","Air-Compressor","10000","20000","USD","20 set","5 t","T/T","EXW","ICC(A)",
        	"On wooden skid","By sea","DDMMYYYY","PortA","PortB","DDMMYYYY","DDMMYYYY","DDMMYYYY","DDYYMMMM","XYZ","ABC", "DDMMYYYY")
        	if  err != nil {
                t.Fatal(err)
        }
        
        var con ContractNo

        err = json.Unmarshal(b,&con)
        if err != nil{
        	t.Fatal(err)
        }
        fmt.Println("Con0")
        fmt.Println(con.ContractNo)

        

        ContractNumb1 := con.ContractNo

        // This must succeed
        b, err = contractStatus(ContractNumb1)
		
       if err != nil || string(b) != "P/O Submitted"{

        t.Fatal(err)
        } 

        // This must succeed
        b, err = getPO(ContractNumb1)
        if err != nil || b == nil {
                t.Fatal(err)
        }
        var poJSON POJSON
        err = json.Unmarshal(b, &poJSON)
		if err != nil {
		t.Fatal(err)
		} else {

		fmt.Printf("Pocheck")
        fmt.Println(poJSON)
		}



		

/*
		err = acceptPO(adminCert, "1000")
		if err != nil  {
			t.Fatal(err)
		}


		
		// This must succeed
        b, err = contractStatus("1000")
		
       if err != nil || string(b) != "P/O Submitted"{

        t.Fatal(err)
        } 
	

*/
        b, err = detailContractByContractID(ContractNumb1)
       if err != nil {
      
			t.Fatal(err)
		}

		var dContract DetailContractJSON

		err = json.Unmarshal(b, &dContract)
		if err != nil{
			
			t.Fatal(err)	
		} 

		fmt.Println("find")
        fmt.Println(dContract)


         // This must succeed
        b, err = getCargoLocation(ContractNumb1)

        if err != nil {
        	t.Fatal(err)
        } else {

        	fmt.Println("CargoLocation1")
        	fmt.Println(string(b))
        } 


        err = submitSN (adminCert, ContractNumb1, "DDMMYYYY","DDMMYYYY")

        if err != nil {

        	t.Fatal(err)

        }
        

        var snJs SNJSON

        b, err = getSN (ContractNumb1)

        err = json.Unmarshal(b, &snJs)

        if err != nil {

        	t.Fatal(err)

        } else {

        	fmt.Println ("SNis")
        	fmt.Println (snJs)
        }


        err = rejectSN(adminCert, ContractNumb1, "Rejection Comment","DDMMYYYY")

        if err != nil {

        	t.Fatal(err)

        }

		b, err = getPO(ContractNumb1)
        if err != nil || b == nil {
                t.Fatal(err)
        }
       

        var pJson POJSON
        err = json.Unmarshal(b, &pJson)
		if err != nil {
		t.Fatal(err)


		fmt.Println("POIS")
		fmt.Println(pJson)


		} else if pJson.ProcessStatus != "P/O Submitted"{

			t.Fatal(err)
		
		} 

		 err = submitSN (adminCert, ContractNumb1, "DDMMYYYY","DDMMYYYY")

        if err != nil {

        	t.Fatal(err)

        }
        
        //var snJs SNJSON

        b, err = getSN (ContractNumb1)

        err = json.Unmarshal(b, &snJs)

        if err != nil {

        	t.Fatal(err)

        } else {

        	fmt.Println ("SNafterresubmission")
        	fmt.Println (snJs)
        }


		b, err = getPO(ContractNumb1)
        if err != nil || b == nil {
                t.Fatal(err)
        }
       

       // var pJson POJSON
        err = json.Unmarshal(b, &pJson)
		if err != nil {
		t.Fatal(err)


		} else if pJson.ProcessStatus != "S/N Submitted"{

			t.Fatal(err)
		
		} 


		var epoJSON EPOJSON
		b, err = getEPO(ContractNumb1)

		if err != nil {
			t.Fatal(err)
			} else {
				err = json.Unmarshal(b,&epoJSON)

			}

			fmt.Println("EPOis")
			fmt.Println(epoJSON)


			err = submitEarlyPaymentOffer(adminCert,ContractNumb1,"DDMMYYYY","DDMMYYYY","10per", "1000","30030","DDMMYYYY")
			if err != nil {
			t.Fatal(err)
			} 


			b, err = getEPO(ContractNumb1)

			if err != nil {
			t.Fatal(err)
			} else {
				err = json.Unmarshal(b,&epoJSON)

			}

			fmt.Println("EPOafterSubmitearly")
			fmt.Println(epoJSON)


			err = AcceptEarlyPaymentOffer(adminCert,ContractNumb1,"OfferAcceptedTime")
			if err != nil {
			t.Fatal(err)
			} 

			b, err = getEPO(ContractNumb1)

			if err != nil {
			t.Fatal(err)
			} else {
				err = json.Unmarshal(b,&epoJSON)

			}

			fmt.Println("EPOafteraccept")
			fmt.Println(epoJSON)

			b, err = getPO(ContractNumb1)
        if err != nil || b == nil {
                t.Fatal(err)
        }
       // var poJSON POJSON
        err = json.Unmarshal(b, &poJSON)
		if err != nil {
		t.Fatal(err)
		} else {

		fmt.Printf("POafterEPOchange")
        fmt.Println(poJSON)
		}

        /* WORKFLOW 1: End */

    	/* WORKFLOW 2: Start */

		// Reject condition

        //if 
        b, err = submitPO(adminCert, "1001", "AA-123","ABC0 Comp","XYZ1 Comp","Air-Compressor","10000","20000","USD","20 set","5 t","T/T","EXW","ICC(A)",
        	"On wooden skid","By sea","DDMMYYYY","PortA","PortB","DDMMYYYY","DDMMYYYY","DDMMYYYY","DDYYMMMM","XYZ0","ABC1", "DDMMYYYY")
        	 if err != nil {
                t.Fatal(err)
        } 

        //var con ContractNo
        err = json.Unmarshal(b,&con)
        if err != nil{
        	t.Fatal(err)
        }

        ContractNumb2 := con.ContractNo
        fmt.Println("Con1")
        fmt.Println(con.ContractNo)
        // This must succeed
        b, err = contractStatus(ContractNumb2)
		
       if err != nil || string(b) != "P/O Submitted"{
        t.Fatal(err)

        } 

       	err = rejectPO(adminCert, ContractNumb2, "Rejection comment")
		if err != nil  {
			t.Fatal(err)
		}

		// This must succeed
        b, err = getPO(ContractNumb2)
        if err != nil || b == nil {
                t.Fatal(err)
        }
       // var poJSON POJSON
        err = json.Unmarshal(b, &poJSON)
		if err != nil {
		t.Fatal(err)
		} else {

		fmt.Printf("PODetails")
        fmt.Println(poJSON)
		}


		listCont, err := listContractsByCompID("XYZ0", "1")
        if err != nil  {
			t.Fatal(err)
		}
		var listContracts ListContracts
		err = json.Unmarshal(listCont, &listContracts.poDetail)
		if err != nil {
		t.Fatal(err)
		} else {

		fmt.Println("l3")
        fmt.Println(listContracts)
		}


		err = reSubmitPO (adminCert, ContractNumb2, "AA-123","ABC0 Comp","XYZ1 Comp","Air-Compressor","10000","20000","USD","13 set","5 t","T/T","EXW","ICC(A)",
        	"On wooden skid","By sea","DDMMYYYY","PortA","PortB","DDMMYYYY","DDMMYYYY","DDYYMMMM","XYZ0","ABC1", "DDMMYYYY") 
		if err != nil{

		
		t.Fatal(err)	
		} 

		// This must succeed
        b, err = getPO(ContractNumb2)
        if err != nil || b == nil {
                t.Fatal(err)
        }
       // var poJSON POJSON
        err = json.Unmarshal(b, &poJSON)
		if err != nil {
		t.Fatal(err)
		} else {

		fmt.Printf("POResubmittedDetails ")
        fmt.Println(poJSON)
        fmt.Println(poJSON.ProcessStatus)
		}




		// This must succeed
        b, err = contractStatus(ContractNumb2)
        if err != nil || string(b) != "P/O Submitted"{
        	t.Fatal(err)
        }


        listCont, err = listContractsByCompID("ABC", "4")
        if err != nil  {
			t.Fatal(err)
		}
		//var listContracts ListContracts
		err = json.Unmarshal(listCont, &listContracts.poDetail)
		if err != nil {
		t.Fatal(err)
		} else {

		fmt.Println("l2 ")
        fmt.Println(listContracts)
		}




		b, err = numbContracts()
		var Count count

		err = json.Unmarshal(b, &Count)
       
       if err != nil {

        t.Fatal(err)
        } else {
        	fmt.Println("Contracts numb")
        	fmt.Println(Count.NumContracts)
        }


         // This must succeed
        b, err = getCargoLocation(ContractNumb2)

        if err != nil {
        	t.Fatal(err)
        } else {

        	fmt.Println("CargoLocation2")
        	fmt.Println(string(b))
        } 
        

        b, err = listOfContracts()

        if err != nil {
        	t.Fatal(err)
        }

        var listOfContracts POList
        err = json.Unmarshal(b, &listOfContracts.contractNo)
        if err != nil {
		t.Fatal(err)
		} else {

		fmt.Println("listofcontractstocheck ")
        fmt.Println(listOfContracts)
		}

        
       b, err = detailContractByContractID(ContractNumb1)
       if err != nil {
			t.Fatal(err)
		}

		//var dContract DetailContractJSON

		err = json.Unmarshal(b, &dContract)
		if err != nil{
			t.Fatal(err)	
		} 

		fmt.Println("DetailContract2")
        fmt.Println(dContract)





    }




    /* WORKFLOW 2: END */


	/* WORKFLOW final: Start Testing Auditing Services APIs */
	
	

func submitPO(admCert crypto.CertificateHandler, ContractNo string, 
		RefNo string, ExporterName string, ImporterName string, Commodity string, UnitPrice string,
		 Amount string, Currency string, Quantity string, Weight string,TermsOfTrade string,
		 TermsOfInsurance string, TermsOfPayment string, PackingMethod string, WayOfTransportation string, TimeOfShipment string, 
		 PortOfShipment string, PortOfDischarge string, POInitialCreateTime string, UpdateTime string, POCreatedTime string, POSubmittedTime string, 
		 CompanyIdOfExporter string, CompanyIdOfImporter string, PaymentDate string ) ([]byte, error) {
        // Get a transaction handler to be used to submit the execute transaction
        // and bind the chaincode access control logic using the binding

        submittingCertHandler, err := administrator.GetTCertificateHandlerNext()
        if err != nil {
                return nil,err
        }
        txHandler, err := submittingCertHandler.GetTransactionHandler()
        if err != nil {
                return nil,err
        }
        binding, err := txHandler.GetBinding()
        if err != nil {
                return nil,err
        }

       
        chaincodeInput := &pb.ChaincodeInput{Args: [][]byte{[]byte("submitPO"), []byte(ContractNo),[]byte(RefNo),
        []byte(ExporterName),[]byte(ImporterName),[]byte(Commodity),[]byte(UnitPrice),[]byte(Amount),
        []byte(Currency),[]byte(Quantity),[]byte(Weight),[]byte(TermsOfTrade),[]byte(TermsOfInsurance),
        []byte(TermsOfPayment),[]byte(PackingMethod),[]byte(WayOfTransportation),[]byte(TimeOfShipment),[]byte(PortOfShipment),
        []byte(PortOfDischarge),[]byte(POInitialCreateTime),[]byte(UpdateTime),[]byte(POCreatedTime),
        []byte(POSubmittedTime),[]byte(CompanyIdOfExporter),[]byte(CompanyIdOfImporter),[]byte(PaymentDate)}}

        chaincodeInputRaw, err := proto.Marshal(chaincodeInput)
        if err != nil {
                return nil,err
        }

        // Access control. Administrator signs chaincodeInputRaw || binding to confirm his identity
        sigma, err := admCert.Sign(append(chaincodeInputRaw, binding...))
        if err != nil {
                return nil,err
        }

        // Prepare spec and submit
        spec := &pb.ChaincodeSpec{
                Type:                 1,
                ChaincodeID:          &pb.ChaincodeID{Name: "mycc"},
                CtorMsg:              chaincodeInput,
                Metadata:             sigma, // Proof of identity
                ConfidentialityLevel: pb.ConfidentialityLevel_PUBLIC,
        }

        var ctx = context.Background()
        chaincodeInvocationSpec := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}

        tid := chaincodeInvocationSpec.ChaincodeSpec.ChaincodeID.Name

        // Now create the Transactions message and send to Peer.
        transaction, err := txHandler.NewChaincodeExecute(chaincodeInvocationSpec, tid)
        if err != nil {
                return nil,fmt.Errorf("Error deploying chaincode: %s ", err)
        }

        ledger, err := ledger.GetLedger()
        ledger.BeginTxBatch("1")
        result, _, err := chaincode.Execute(ctx, chaincode.GetChain(chaincode.DefaultChain), transaction)
        if err != nil {
                return nil, fmt.Errorf("Error deploying chaincode: %s", err)
        }
        ledger.CommitTxBatch("1", []*pb.Transaction{transaction}, nil, nil)

        return result, err
}



//getED
func getPO(ContractNo string) ([]byte, error) {
        
        chaincodeInput := &pb.ChaincodeInput{Args: [][]byte{[]byte("getPO"), []byte(ContractNo)}}

        // Prepare spec and submit
        spec := &pb.ChaincodeSpec{
                Type:                 1,
                ChaincodeID:          &pb.ChaincodeID{Name: "mycc"},
                CtorMsg:              chaincodeInput,
                ConfidentialityLevel: pb.ConfidentialityLevel_PUBLIC,
        }

        var ctx = context.Background()
        chaincodeInvocationSpec := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}

        tid := chaincodeInvocationSpec.ChaincodeSpec.ChaincodeID.Name

        // Now create the Transactions message and send to Peer.
        transaction, err := administrator.NewChaincodeQuery(chaincodeInvocationSpec, tid)
        if err != nil {
                return nil, fmt.Errorf("Error deploying chaincode: %s ", err)
        }

        ledger, err := ledger.GetLedger()
        ledger.BeginTxBatch("1")
        result, _, err := chaincode.Execute(ctx, chaincode.GetChain(chaincode.DefaultChain), transaction)
        if err != nil {
                return nil, fmt.Errorf("Error deploying chaincode: %s", err)
        }
        ledger.CommitTxBatch("1", []*pb.Transaction{transaction}, nil, nil)

        return result, err
}

func acceptPO(admCert crypto.CertificateHandler, ContractNo string) error {
        // Get a transaction handler to be used to submit the execute transaction
        // and bind the chaincode access control logic using the binding

        submittingCertHandler, err := administrator.GetTCertificateHandlerNext()
        if err != nil {
                return err
        }
        txHandler, err := submittingCertHandler.GetTransactionHandler()
        if err != nil {
                return err
        }
        binding, err := txHandler.GetBinding()
        if err != nil {
                return err
        }

      
        chaincodeInput := &pb.ChaincodeInput{Args: [][]byte{[]byte("acceptPO"), []byte(ContractNo)}}
        chaincodeInputRaw, err := proto.Marshal(chaincodeInput)
        if err != nil {
                return err
        }

        // Access control. Administrator signs chaincodeInputRaw || binding to confirm his identity
        sigma, err := admCert.Sign(append(chaincodeInputRaw, binding...))
        if err != nil {
                return err
        }

        // Prepare spec and submit
        spec := &pb.ChaincodeSpec{
                Type:                 1,
                ChaincodeID:          &pb.ChaincodeID{Name: "mycc"},
                CtorMsg:              chaincodeInput,
                Metadata:             sigma, // Proof of identity
                ConfidentialityLevel: pb.ConfidentialityLevel_PUBLIC,
        }

        var ctx = context.Background()
        chaincodeInvocationSpec := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}

        tid := chaincodeInvocationSpec.ChaincodeSpec.ChaincodeID.Name

        // Now create the Transactions message and send to Peer.
        transaction, err := txHandler.NewChaincodeExecute(chaincodeInvocationSpec, tid)
        if err != nil {
                return fmt.Errorf("Error deploying chaincode: %s ", err)
        }

        ledger, err := ledger.GetLedger()
        ledger.BeginTxBatch("1")
        _, _, err = chaincode.Execute(ctx, chaincode.GetChain(chaincode.DefaultChain), transaction)
        if err != nil {
                return fmt.Errorf("Error deploying chaincode: %s", err)
        }
        ledger.CommitTxBatch("1", []*pb.Transaction{transaction}, nil, nil)

        return err
}


func rejectPO(admCert crypto.CertificateHandler, ContractNo string , Comment string) error {
        // Get a transaction handler to be used to submit the execute transaction
        // and bind the chaincode access control logic using the binding

        submittingCertHandler, err := administrator.GetTCertificateHandlerNext()


        if err != nil {
                return err
        }
        txHandler, err := submittingCertHandler.GetTransactionHandler()
        if err != nil {
                return err
        }
        binding, err := txHandler.GetBinding()
        if err != nil {
                return err
        }

        
        chaincodeInput := &pb.ChaincodeInput{Args: [][]byte{[]byte("rejectPO"), []byte(ContractNo),[]byte(Comment)}}
        chaincodeInputRaw, err := proto.Marshal(chaincodeInput)
        if err != nil {
                return err
        }

        // Access control. Administrator signs chaincodeInputRaw || binding to confirm his identity
        sigma, err := admCert.Sign(append(chaincodeInputRaw, binding...))
        if err != nil {
                return err
        }

        // Prepare spec and submit
        spec := &pb.ChaincodeSpec{
                Type:                 1,
                ChaincodeID:          &pb.ChaincodeID{Name: "mycc"},
                CtorMsg:              chaincodeInput,
                Metadata:             sigma, // Proof of identity
                ConfidentialityLevel: pb.ConfidentialityLevel_PUBLIC,
        }

        var ctx = context.Background()
        chaincodeInvocationSpec := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}

        tid := chaincodeInvocationSpec.ChaincodeSpec.ChaincodeID.Name

        // Now create the Transactions message and send to Peer.
        transaction, err := txHandler.NewChaincodeExecute(chaincodeInvocationSpec, tid)
        if err != nil {
                return fmt.Errorf("Error deploying chaincode: %s ", err)
        }

        ledger, err := ledger.GetLedger()
        ledger.BeginTxBatch("1")
        _, _, err = chaincode.Execute(ctx, chaincode.GetChain(chaincode.DefaultChain), transaction)
        if err != nil {
                return fmt.Errorf("Error deploying chaincode: %s", err)
        }
        ledger.CommitTxBatch("1", []*pb.Transaction{transaction}, nil, nil)

        return err
}  

func contractStatus(ContractNo string) ([]byte, error) {
       
        chaincodeInput := &pb.ChaincodeInput{Args: [][]byte{[]byte("contractStatus"), []byte(ContractNo)}}

        // Prepare spec and submit
        spec := &pb.ChaincodeSpec{
                Type:                 1,
                ChaincodeID:          &pb.ChaincodeID{Name: "mycc"},
                CtorMsg:              chaincodeInput,
                ConfidentialityLevel: pb.ConfidentialityLevel_PUBLIC,
        }

        var ctx = context.Background()
        chaincodeInvocationSpec := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}

        tid := chaincodeInvocationSpec.ChaincodeSpec.ChaincodeID.Name

        // Now create the Transactions message and send to Peer.
        transaction, err := administrator.NewChaincodeQuery(chaincodeInvocationSpec, tid)
        if err != nil {
                return nil, fmt.Errorf("Error deploying chaincode: %s ", err)
        }

        ledger, err := ledger.GetLedger()
        ledger.BeginTxBatch("1")
        result, _, err := chaincode.Execute(ctx, chaincode.GetChain(chaincode.DefaultChain), transaction)
        if err != nil {
                return nil, fmt.Errorf("Error deploying chaincode: %s", err)
        }
        ledger.CommitTxBatch("1", []*pb.Transaction{transaction}, nil, nil)

        return result, err
}


func listContractsByCompID(CompanyName string, CompanyID string) ([]byte, error) {
        
        chaincodeInput := &pb.ChaincodeInput{Args: [][]byte{[]byte("listContractsByCompID"), []byte(CompanyName),[]byte(CompanyID)}}

        // Prepare spec and submit
        spec := &pb.ChaincodeSpec{
                Type:                 1,
                ChaincodeID:          &pb.ChaincodeID{Name: "mycc"},
                CtorMsg:              chaincodeInput,
                ConfidentialityLevel: pb.ConfidentialityLevel_PUBLIC,
        }

        var ctx = context.Background()
        chaincodeInvocationSpec := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}

        tid := chaincodeInvocationSpec.ChaincodeSpec.ChaincodeID.Name

        // Now create the Transactions message and send to Peer.
        transaction, err := administrator.NewChaincodeQuery(chaincodeInvocationSpec, tid)
        if err != nil {
                return nil, fmt.Errorf("Error deploying chaincode: %s ", err)
        }

        ledger, err := ledger.GetLedger()
        ledger.BeginTxBatch("1")
        result, _, err := chaincode.Execute(ctx, chaincode.GetChain(chaincode.DefaultChain), transaction)
        if err != nil {
                return nil, fmt.Errorf("Error deploying chaincode: %s", err)
        }
        ledger.CommitTxBatch("1", []*pb.Transaction{transaction}, nil, nil)

        return result, err
}


func numbContracts() ([]byte, error) {
       
        chaincodeInput := &pb.ChaincodeInput{Args: [][]byte{[]byte("numbContracts")}}

        // Prepare spec and submit
        spec := &pb.ChaincodeSpec{
                Type:                 1,
                ChaincodeID:          &pb.ChaincodeID{Name: "mycc"},
                CtorMsg:              chaincodeInput,
                ConfidentialityLevel: pb.ConfidentialityLevel_PUBLIC,
        }

        var ctx = context.Background()
        chaincodeInvocationSpec := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}

        tid := chaincodeInvocationSpec.ChaincodeSpec.ChaincodeID.Name

        // Now create the Transactions message and send to Peer.
        transaction, err := administrator.NewChaincodeQuery(chaincodeInvocationSpec, tid)
        if err != nil {
                return nil, fmt.Errorf("Error deploying chaincode: %s ", err)
        }

        ledger, err := ledger.GetLedger()
        ledger.BeginTxBatch("1")
        result, _, err := chaincode.Execute(ctx, chaincode.GetChain(chaincode.DefaultChain), transaction)
        if err != nil {
                return nil, fmt.Errorf("Error deploying chaincode: %s", err)
        }
        ledger.CommitTxBatch("1", []*pb.Transaction{transaction}, nil, nil)

        return result, err
}



//getED
func getCargoLocation(ContractNo string) ([]byte, error) {
        
        chaincodeInput := &pb.ChaincodeInput{Args: [][]byte{[]byte("getCargoLocation"), []byte(ContractNo)}}

        // Prepare spec and submit
        spec := &pb.ChaincodeSpec{
                Type:                 1,
                ChaincodeID:          &pb.ChaincodeID{Name: "mycc"},
                CtorMsg:              chaincodeInput,
                ConfidentialityLevel: pb.ConfidentialityLevel_PUBLIC,
        }

        var ctx = context.Background()
        chaincodeInvocationSpec := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}

        tid := chaincodeInvocationSpec.ChaincodeSpec.ChaincodeID.Name

        // Now create the Transactions message and send to Peer.
        transaction, err := administrator.NewChaincodeQuery(chaincodeInvocationSpec, tid)
        if err != nil {
                return nil, fmt.Errorf("Error deploying chaincode: %s ", err)
        }

        ledger, err := ledger.GetLedger()
        ledger.BeginTxBatch("1")
        result, _, err := chaincode.Execute(ctx, chaincode.GetChain(chaincode.DefaultChain), transaction)
        if err != nil {
                return nil, fmt.Errorf("Error deploying chaincode: %s", err)
        }
        ledger.CommitTxBatch("1", []*pb.Transaction{transaction}, nil, nil)

        return result, err
}

func listOfContracts() ([]byte, error) {
        
        chaincodeInput := &pb.ChaincodeInput{Args: [][]byte{[]byte("listOfContracts")}}

        // Prepare spec and submit
        spec := &pb.ChaincodeSpec{
                Type:                 1,
                ChaincodeID:          &pb.ChaincodeID{Name: "mycc"},
                CtorMsg:              chaincodeInput,
                ConfidentialityLevel: pb.ConfidentialityLevel_PUBLIC,
        }

        var ctx = context.Background()
        chaincodeInvocationSpec := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}

        tid := chaincodeInvocationSpec.ChaincodeSpec.ChaincodeID.Name

        // Now create the Transactions message and send to Peer.
        transaction, err := administrator.NewChaincodeQuery(chaincodeInvocationSpec, tid)
        if err != nil {
                return nil, fmt.Errorf("Error deploying chaincode: %s ", err)
        }

        ledger, err := ledger.GetLedger()
        ledger.BeginTxBatch("1")
        result, _, err := chaincode.Execute(ctx, chaincode.GetChain(chaincode.DefaultChain), transaction)
        if err != nil {
                return nil, fmt.Errorf("Error deploying chaincode: %s", err)
        }
        ledger.CommitTxBatch("1", []*pb.Transaction{transaction}, nil, nil)

        return result, err
}


func detailContractByContractID(ContractNo string) ([]byte, error) {
        
        chaincodeInput := &pb.ChaincodeInput{Args: [][]byte{[]byte("detailContractByContractID"), []byte(ContractNo)}}

        // Prepare spec and submit
        spec := &pb.ChaincodeSpec{
                Type:                 1,
                ChaincodeID:          &pb.ChaincodeID{Name: "mycc"},
                CtorMsg:              chaincodeInput,
                ConfidentialityLevel: pb.ConfidentialityLevel_PUBLIC,
        }

        var ctx = context.Background()
        chaincodeInvocationSpec := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}

        tid := chaincodeInvocationSpec.ChaincodeSpec.ChaincodeID.Name

        // Now create the Transactions message and send to Peer.
        transaction, err := administrator.NewChaincodeQuery(chaincodeInvocationSpec, tid)
        if err != nil {
                return nil, fmt.Errorf("Error deploying chaincode: %s ", err)
        }

        ledger, err := ledger.GetLedger()
        ledger.BeginTxBatch("1")
        result, _, err := chaincode.Execute(ctx, chaincode.GetChain(chaincode.DefaultChain), transaction)
        if err != nil {
                return nil, fmt.Errorf("Error deploying chaincode: %s", err)
        }
        ledger.CommitTxBatch("1", []*pb.Transaction{transaction}, nil, nil)

        return result, err
}

func submitSN(admCert crypto.CertificateHandler, ContractNo string, 
		UpdateTime string, SNSubmittedTime string) error {
        // Get a transaction handler to be used to submit the execute transaction
        // and bind the chaincode access control logic using the binding

        submittingCertHandler, err := administrator.GetTCertificateHandlerNext()
        if err != nil {
                return err
        }
        txHandler, err := submittingCertHandler.GetTransactionHandler()
        if err != nil {
                return err
        }
        binding, err := txHandler.GetBinding()
        if err != nil {
                return err
        }

       
        chaincodeInput := &pb.ChaincodeInput{Args: [][]byte{[]byte("submitSN"), []byte(ContractNo),[]byte(UpdateTime),
        []byte(SNSubmittedTime)}}

        chaincodeInputRaw, err := proto.Marshal(chaincodeInput)
        if err != nil {
                return err
        }

        // Access control. Administrator signs chaincodeInputRaw || binding to confirm his identity
        sigma, err := admCert.Sign(append(chaincodeInputRaw, binding...))
        if err != nil {
                return err
        }

        // Prepare spec and submit
        spec := &pb.ChaincodeSpec{
                Type:                 1,
                ChaincodeID:          &pb.ChaincodeID{Name: "mycc"},
                CtorMsg:              chaincodeInput,
                Metadata:             sigma, // Proof of identity
                ConfidentialityLevel: pb.ConfidentialityLevel_PUBLIC,
        }

        var ctx = context.Background()
        chaincodeInvocationSpec := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}

        tid := chaincodeInvocationSpec.ChaincodeSpec.ChaincodeID.Name

        // Now create the Transactions message and send to Peer.
        transaction, err := txHandler.NewChaincodeExecute(chaincodeInvocationSpec, tid)
        if err != nil {
                return fmt.Errorf("Error deploying chaincode: %s ", err)
        }

        ledger, err := ledger.GetLedger()
        ledger.BeginTxBatch("1")
        _, _, err = chaincode.Execute(ctx, chaincode.GetChain(chaincode.DefaultChain), transaction)
        if err != nil {
                return fmt.Errorf("Error deploying chaincode: %s", err)
        }
        ledger.CommitTxBatch("1", []*pb.Transaction{transaction}, nil, nil)

        return err
}


func rejectSN(admCert crypto.CertificateHandler, ContractNo string, rejectionComment string , UpdateTime string) error {
        // Get a transaction handler to be used to submit the execute transaction
        // and bind the chaincode access control logic using the binding

        submittingCertHandler, err := administrator.GetTCertificateHandlerNext()
        if err != nil {
                return err
        }
        txHandler, err := submittingCertHandler.GetTransactionHandler()
        if err != nil {
                return err
        }
        binding, err := txHandler.GetBinding()
        if err != nil {
                return err
        }

      
        chaincodeInput := &pb.ChaincodeInput{Args: [][]byte{[]byte("rejectSN"), []byte(ContractNo), []byte(rejectionComment), []byte(UpdateTime)}}
        chaincodeInputRaw, err := proto.Marshal(chaincodeInput)
        if err != nil {
                return err
        }

        // Access control. Administrator signs chaincodeInputRaw || binding to confirm his identity
        sigma, err := admCert.Sign(append(chaincodeInputRaw, binding...))
        if err != nil {
                return err
        }

        // Prepare spec and submit
        spec := &pb.ChaincodeSpec{
                Type:                 1,
                ChaincodeID:          &pb.ChaincodeID{Name: "mycc"},
                CtorMsg:              chaincodeInput,
                Metadata:             sigma, // Proof of identity
                ConfidentialityLevel: pb.ConfidentialityLevel_PUBLIC,
        }

        var ctx = context.Background()
        chaincodeInvocationSpec := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}

        tid := chaincodeInvocationSpec.ChaincodeSpec.ChaincodeID.Name

        // Now create the Transactions message and send to Peer.
        transaction, err := txHandler.NewChaincodeExecute(chaincodeInvocationSpec, tid)
        if err != nil {
                return fmt.Errorf("Error deploying chaincode: %s ", err)
        }

        ledger, err := ledger.GetLedger()
        ledger.BeginTxBatch("1")
        _, _, err = chaincode.Execute(ctx, chaincode.GetChain(chaincode.DefaultChain), transaction)
        if err != nil {
                return fmt.Errorf("Error deploying chaincode: %s", err)
        }
        ledger.CommitTxBatch("1", []*pb.Transaction{transaction}, nil, nil)

        return err
}


func getSN(ContractNo string) ([]byte, error) {
        
        chaincodeInput := &pb.ChaincodeInput{Args: [][]byte{[]byte("getSN"), []byte(ContractNo)}}

        // Prepare spec and submit
        spec := &pb.ChaincodeSpec{
                Type:                 1,
                ChaincodeID:          &pb.ChaincodeID{Name: "mycc"},
                CtorMsg:              chaincodeInput,
                ConfidentialityLevel: pb.ConfidentialityLevel_PUBLIC,
        }

        var ctx = context.Background()
        chaincodeInvocationSpec := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}

        tid := chaincodeInvocationSpec.ChaincodeSpec.ChaincodeID.Name

        // Now create the Transactions message and send to Peer.
        transaction, err := administrator.NewChaincodeQuery(chaincodeInvocationSpec, tid)
        if err != nil {
                return nil, fmt.Errorf("Error deploying chaincode: %s ", err)
        }

        ledger, err := ledger.GetLedger()
        ledger.BeginTxBatch("1")
        result, _, err := chaincode.Execute(ctx, chaincode.GetChain(chaincode.DefaultChain), transaction)
        if err != nil {
                return nil, fmt.Errorf("Error deploying chaincode: %s", err)
        }
        ledger.CommitTxBatch("1", []*pb.Transaction{transaction}, nil, nil)

        return result, err
}


func reSubmitPO(admCert crypto.CertificateHandler, ContractNo string, 
		RefNo string, ExporterName string, ImporterName string, Commodity string, UnitPrice string,
		 Amount string, Currency string, Quantity string, Weight string,TermsOfTrade string,
		 TermsOfInsurance string, TermsOfPayment string, PackingMethod string, WayOfTransportation string, TimeOfShipment string, 
		 PortOfShipment string, PortOfDischarge string, UpdateTime string, POCreatedTime string, POSubmittedTime string, 
		 CompanyIdOfExporter string, CompanyIdOfImporter string, PaymentDate string ) error {
        // Get a transaction handler to be used to submit the execute transaction
        // and bind the chaincode access control logic using the binding

        submittingCertHandler, err := administrator.GetTCertificateHandlerNext()
        if err != nil {
                return err
        }
        txHandler, err := submittingCertHandler.GetTransactionHandler()
        if err != nil {
                return err
        }
        binding, err := txHandler.GetBinding()
        if err != nil {
                return err
        }

       
        chaincodeInput := &pb.ChaincodeInput{Args: [][]byte{[]byte("reSubmitPO"), []byte(ContractNo),[]byte(RefNo),
        []byte(ExporterName),[]byte(ImporterName),[]byte(Commodity),[]byte(UnitPrice),[]byte(Amount),
        []byte(Currency),[]byte(Quantity),[]byte(Weight),[]byte(TermsOfTrade),[]byte(TermsOfInsurance),
        []byte(TermsOfPayment),[]byte(PackingMethod),[]byte(WayOfTransportation),[]byte(TimeOfShipment),[]byte(PortOfShipment),
        []byte(PortOfDischarge),[]byte(UpdateTime),[]byte(POCreatedTime),
        []byte(POSubmittedTime),[]byte(CompanyIdOfExporter),[]byte(CompanyIdOfImporter),[]byte(PaymentDate)}}

        chaincodeInputRaw, err := proto.Marshal(chaincodeInput)
        if err != nil {
                return err
        }

        // Access control. Administrator signs chaincodeInputRaw || binding to confirm his identity
        sigma, err := admCert.Sign(append(chaincodeInputRaw, binding...))
        if err != nil {
                return err
        }

        // Prepare spec and submit
        spec := &pb.ChaincodeSpec{
                Type:                 1,
                ChaincodeID:          &pb.ChaincodeID{Name: "mycc"},
                CtorMsg:              chaincodeInput,
                Metadata:             sigma, // Proof of identity
                ConfidentialityLevel: pb.ConfidentialityLevel_PUBLIC,
        }

        var ctx = context.Background()
        chaincodeInvocationSpec := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}

        tid := chaincodeInvocationSpec.ChaincodeSpec.ChaincodeID.Name

        // Now create the Transactions message and send to Peer.
        transaction, err := txHandler.NewChaincodeExecute(chaincodeInvocationSpec, tid)
        if err != nil {
                return fmt.Errorf("Error deploying chaincode: %s ", err)
        }

        ledger, err := ledger.GetLedger()
        ledger.BeginTxBatch("1")
        _, _, err = chaincode.Execute(ctx, chaincode.GetChain(chaincode.DefaultChain), transaction)
        if err != nil {
                return fmt.Errorf("Error deploying chaincode: %s", err)
        }
        ledger.CommitTxBatch("1", []*pb.Transaction{transaction}, nil, nil)

        return err
}


func getEPO(ContractNo string) ([]byte, error) {
        
        chaincodeInput := &pb.ChaincodeInput{Args: [][]byte{[]byte("getEPO"), []byte(ContractNo)}}

        // Prepare spec and submit
        spec := &pb.ChaincodeSpec{
                Type:                 1,
                ChaincodeID:          &pb.ChaincodeID{Name: "mycc"},
                CtorMsg:              chaincodeInput,
                ConfidentialityLevel: pb.ConfidentialityLevel_PUBLIC,
        }

        var ctx = context.Background()
        chaincodeInvocationSpec := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}

        tid := chaincodeInvocationSpec.ChaincodeSpec.ChaincodeID.Name

        // Now create the Transactions message and send to Peer.
        transaction, err := administrator.NewChaincodeQuery(chaincodeInvocationSpec, tid)
        if err != nil {
                return nil, fmt.Errorf("Error deploying chaincode: %s ", err)
        }

        ledger, err := ledger.GetLedger()
        ledger.BeginTxBatch("1")
        result, _, err := chaincode.Execute(ctx, chaincode.GetChain(chaincode.DefaultChain), transaction)
        if err != nil {
                return nil, fmt.Errorf("Error deploying chaincode: %s", err)
        }
        ledger.CommitTxBatch("1", []*pb.Transaction{transaction}, nil, nil)

        return result, err
}


func submitEarlyPaymentOffer(admCert crypto.CertificateHandler, ContractNo string , date1 string, date2 string,offerper string,disamount string, amount string, date3 string  ) error {
        // Get a transaction handler to be used to submit the execute transaction
        // and bind the chaincode access control logic using the binding

        submittingCertHandler, err := administrator.GetTCertificateHandlerNext()


        if err != nil {
                return err
        }
        txHandler, err := submittingCertHandler.GetTransactionHandler()
        if err != nil {
                return err
        }
        binding, err := txHandler.GetBinding()
        if err != nil {
                return err
        }

        
        chaincodeInput := &pb.ChaincodeInput{Args: [][]byte{[]byte("submitEarlyPaymentOffer"), []byte(ContractNo),[]byte(date1),[]byte(date2),[]byte(offerper),[]byte(disamount),[]byte(amount), []byte(date3)}}
        chaincodeInputRaw, err := proto.Marshal(chaincodeInput)
        if err != nil {
                return err
        }

        // Access control. Administrator signs chaincodeInputRaw || binding to confirm his identity
        sigma, err := admCert.Sign(append(chaincodeInputRaw, binding...))
        if err != nil {
                return err
        }

        // Prepare spec and submit
        spec := &pb.ChaincodeSpec{
                Type:                 1,
                ChaincodeID:          &pb.ChaincodeID{Name: "mycc"},
                CtorMsg:              chaincodeInput,
                Metadata:             sigma, // Proof of identity
                ConfidentialityLevel: pb.ConfidentialityLevel_PUBLIC,
        }

        var ctx = context.Background()
        chaincodeInvocationSpec := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}

        tid := chaincodeInvocationSpec.ChaincodeSpec.ChaincodeID.Name

        // Now create the Transactions message and send to Peer.
        transaction, err := txHandler.NewChaincodeExecute(chaincodeInvocationSpec, tid)
        if err != nil {
                return fmt.Errorf("Error deploying chaincode: %s ", err)
        }

        ledger, err := ledger.GetLedger()
        ledger.BeginTxBatch("1")
        _, _, err = chaincode.Execute(ctx, chaincode.GetChain(chaincode.DefaultChain), transaction)
        if err != nil {
                return fmt.Errorf("Error deploying chaincode: %s", err)
        }
        ledger.CommitTxBatch("1", []*pb.Transaction{transaction}, nil, nil)

        return err
}  




func AcceptEarlyPaymentOffer(admCert crypto.CertificateHandler, ContractNo string, date1 string) error {
        // Get a transaction handler to be used to submit the execute transaction
        // and bind the chaincode access control logic using the binding

        submittingCertHandler, err := administrator.GetTCertificateHandlerNext()
        if err != nil {
                return err
        }
        txHandler, err := submittingCertHandler.GetTransactionHandler()
        if err != nil {
                return err
        }
        binding, err := txHandler.GetBinding()
        if err != nil {
                return err
        }

      
        chaincodeInput := &pb.ChaincodeInput{Args: [][]byte{[]byte("acceptEarlyPaymentOffer"), []byte(ContractNo),[]byte(date1)}}
        chaincodeInputRaw, err := proto.Marshal(chaincodeInput)
        if err != nil {
                return err
        }

        // Access control. Administrator signs chaincodeInputRaw || binding to confirm his identity
        sigma, err := admCert.Sign(append(chaincodeInputRaw, binding...))
        if err != nil {
                return err
        }

        // Prepare spec and submit
        spec := &pb.ChaincodeSpec{
                Type:                 1,
                ChaincodeID:          &pb.ChaincodeID{Name: "mycc"},
                CtorMsg:              chaincodeInput,
                Metadata:             sigma, // Proof of identity
                ConfidentialityLevel: pb.ConfidentialityLevel_PUBLIC,
        }

        var ctx = context.Background()
        chaincodeInvocationSpec := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}

        tid := chaincodeInvocationSpec.ChaincodeSpec.ChaincodeID.Name

        // Now create the Transactions message and send to Peer.
        transaction, err := txHandler.NewChaincodeExecute(chaincodeInvocationSpec, tid)
        if err != nil {
                return fmt.Errorf("Error deploying chaincode: %s ", err)
        }

        ledger, err := ledger.GetLedger()
        ledger.BeginTxBatch("1")
        _, _, err = chaincode.Execute(ctx, chaincode.GetChain(chaincode.DefaultChain), transaction)
        if err != nil {
                return fmt.Errorf("Error deploying chaincode: %s", err)
        }
        ledger.CommitTxBatch("1", []*pb.Transaction{transaction}, nil, nil)

        return err
}
