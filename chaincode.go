/**
@author: Arshad Sarfarz/sasi
@version: 1.0.0
@date: 10/04/2017
@Description: MedLab-Pharma chaincode v1
**/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strings"
)

const STATUS_SHIPPED = "shipped by manufacturer"
const STATUS_ACCEPTED_BY_DISTRIBUTOR = "accepted by distributor"
const STATUS_SHIPPED_BY_DISTRIBUTOR = "dispatched by distributor"
const STATUS_ACCEPTED_BY_LOGISTICS= "accepted by logistics"
const STATUS_SOLD_BY_RETAILER= "sold out"
const STATUS_PARTIALLY_SOLD_BY_RETAILER= "partially sold"
const STATUS_DISPATCH_IN_PROGRESS= "dispatch in progress by distributor"
const STATUS_ACCEPTED_BY_RETAILER = "accepted by retailer"
const STATUS_REJECTED_BY_RETAILER="rejected by retailer"
const STATUS_REJECTED_BY_CONSUMER="counterfeit  highlighted"
const STATUS_PARTIALLY_ACCEPTED_BY_DISTRIBUTOR = "partially accepted by distributor"
const STATUS_REJECTED_BY_LOGISTICS = "rejected by logistics"
const STATUS_REJECTED_BY_DISTRIBUTOR  = "rejected by distributor"
const STATUS_DISPATCHED_BY_LOGISTICS = "dispatched by logistics"
const UNIQUE_ID_COUNTER string = "UniqueIDCounter"
const CONTAINER_OWNER = "ContainerOwner"
//const RFC1123 = "Mon, 02 Jan 2006 15:04:05 MST"

type MedLabPharmaChaincode struct {
}

type UniqueIDCounter struct {
	ContainerMaxID int `json:"ContainerMaxID"`
	PalletMaxID    int `json:"PalletMaxID"`
}
type Shipment struct{
	ContainerList []Container `json:"container_list"`

}

type Container struct {
	ContainerId       string              `json:"container_id"`
	ParentContainerId string              `json:"parent_container_id"`
	ChildContainerId  []string            `json:"child_container_id"`
	Recipient         string              `json:"recipient_id"`
	Elements          ContainerElements   `json:"elements"`
	Provenance        ContainerProvenance `json:"provenance"`
	Repackagingstatus  []string            `json:"repackaged_pallets"`
	CertifiedBy       string              `json:"certified_by"`   ///New fields
	Address           string              `json:"address"`        ///New fields
	USN               string              `json:"usn"`            ///New fields
	ShipmentDate      string              `json:"shipment_date"`  ///New fields
	InvoiceNumber     string              `json:"invoice_number"` ///New fields
	Remarks           string              `json:"remarks"`        ///New fields
	ReceivedDate      string              `json:"recieved_date"` 
    SenderAddress     string              `json:"sender_address"`
}
type ContainerDrugNameCombination struct{
	ContainerId       string              `json:"container_id"`
	DrugName          string               `json:"drug_name"`
	GenericName	      string                `json:"generic_name"`
}
type UnitIDListJson struct {
	UnitID  []string            `json:"units_sold"`
	}
type Response struct {   
    Message string    `json:"message"`
}
type ContainerElements struct {
	Pallets []Pallet `json:"pallets"`
	Health string    `json:"container_health"`
	Remarks  string     `json:"container_remarks"`
}

type Pallet struct {
	PalletId string `json:"pallet_id"`
	Cases    []Case `json:"cases"`
	Health string    `json:"pallet_health"`
	Remarks  string     `json:"pallet_remarks"`
	}

type Case struct {
	CaseId string `json:"case_id"`
	Units  []Unit `json:"units"`
	Health string    `json:"case_health"`
	Remarks  string     `json:"case_remarks"`
	}

type Unit struct {
	DrugId       string `json:"drug_id"`
	DrugName     string `json:"drug_name"` ////New Fields
	UnitId       string `json:"unit_id"`
	ExpiryDate   string `json:"expiry_date"`
	HealthStatus string `json:"health_status"`
	BatchNumber  string `json:"batch_number"`
	LotNumber    string `json:"lot_number"`
	SaleStatus   string `json:"sale_status"`
	ConsumerName string `json:"consumer_name"`
	Health string    `json:"unit_health"`
	Remarks  string     `json:"unit_remarks"`
	GenericName string  `json:"Generic_Name"`
}

type ContainerProvenance struct {
	TransitStatus string          `json:transit_status`
	Sender        string          `json:sender`
	Receiver      string          `json:receiver`
	Supplychain   []ChainActivity `json:supplychain`
	ShipmentDate   string `json:"date"`
}

type ChainActivity struct {
	Sender   string `json:sender`
	Receiver string `json:receiver`
	Status   string `json:transit_status`
	ShipmentDate  string `json:"date"`
	Remarks   string  `json:remarks`
	}

type ContainerOwners struct {
	Owners []Owner `json:owners`
}

type Owner struct {
	OwnerId       string   `json:owner_id`
	ContainerList []string `json:container_id`
}

func main() {
	fmt.Println("Inside MedLabPharmaChaincode main function")
	err := shim.Start(new(MedLabPharmaChaincode))
	if err != nil {
		fmt.Printf("Error starting MedLabPharma chaincode: %s", err)
	}
}

// Init resets all the things
func (t *MedLabPharmaChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// Handle different functions
	if function == "init" {
		return t.init(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Invoke isur entry point to invoke a chaincode function
func (t *MedLabPharmaChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	user_byte,_ := t.GetUserAttribute(stub,"user_type")
		user_type := string(user_byte)
		if function == "ShipContainerUsingLogistics" {
		   if (user_type =="manufacturer"){
		     return t.ShipContainerUsingLogistics(stub, args[0], args[1], args[2], args[3], args[4],args[5])
		   }else{
                return nil, errors.New("User type: " + user_type+ "does not have privilege to execute chain code 'ShipContainerUsingLogistics'" )
		   }
	} else if function == "AcceptContainerbyLogistics"{
		  if (user_type =="logistics"){
			  return t.AcceptContainerbyLogistics(stub, args[0], args[1],args[2], args[3],args[4])
		  }else{
               return nil, errors.New("User type: " + user_type+ "does not have privilege to execute chain code 'AcceptContainerbyLogistics'" )
		  }	  
	}else if function == "DispatchContainer"{
		  if (user_type =="logistics"){
               return t.DispatchContainer(stub, args[0], args[1],args[2],args[3])	
		  }else{
                return nil, errors.New("User type: " + user_type+ "does not have privilege to execute chain code 'DispatchContainer'" )
		  }	  		
	}else if function == "UpdateContainerbyDistributor"{
		if (user_type =="distributor"){
		         return t.UpdateContainerbyDistributor(stub, args[0], args[1],args[2],args[3],args[4])		
		}else{
             return nil, errors.New("User type: " + user_type+ "does not have privilege to execute chain code 'UpdateContainerbyDistributor'" )
		}		   
	}else if function == "RejectContainerbyLogistics"{
		  if (user_type =="logistics"){
           	return t.RejectContainerbyLogistics(stub, args[0], args[1],args[2],args[3],args[4]) 
		}else{
              return nil, errors.New("User type: " + user_type+ "does not have privilege to execute chain code 'RejectContainerbyLogistics'" )
		}
	}else if function == "repackagingContainerbyDistributor"{
		if (user_type =="distributor"){
			return t.repackagingContainerbyDistributor(stub, args[0],args[1], args[2],args[3],args[4],args[5])		
		}else{
			return nil, errors.New("User type: " + user_type+ "does not have privilege to execute chain code 'repackagingContainerbyDistributor'" )
		}		   
	}else if function == "AcceptContainerbyRetailer"{
		if (user_type =="retailer"){
		         return t.AcceptContainerbyRetailer(stub, args[0],args[1], args[2],args[3])		
		}else{
            return nil, errors.New("User type: " + user_type+ "does not have privilege to execute chain code 'AcceptContainerbyRetailer'" )
		}		   
	}else if function == "SellingbyRetailer"{
		if (user_type =="retailer"){
		        return t.SellingbyRetailer(stub, args[0],args[1], args[2],args[3])		
		}else{
             return nil, errors.New("User type: " + user_type+ "does not have privilege to execute chain code 'SellingbyRetailer'" )
		}		   
	}else if function == "RejectContainerbyRetailer"{
		if (user_type =="retailer"){
		        return t.RejectContainerbyRetailer(stub, args[0],args[1], args[2],args[3])		
		}else{
             return nil, errors.New("User type: " + user_type+ "does not have privilege to execute chain code 'RejectContainerbyRetailer'" )
		}		   
	}else if function == "RejectingbyConsumer"{
		if (user_type =="consumer"){
		        return t.RejectingbyConsumer(stub, args[0],args[1], args[2],args[3])		
		}else{
             return nil, errors.New("User type: " + user_type+ "does not have privilege to execute chain code 'RejectingbyConsumer'" )
		}		   
	}				 
	fmt.Println("invoke did not find func: " + function)
	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *MedLabPharmaChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "GetContainerDetails" { //read a variable  
		return t.GetContainerDetails(stub, args[0])
	} else if function == "GetMaxIDValue" {
		return t.GetMaxIDValue(stub)
	} else if function == "GetEmptyContainer" {
		return t.GetEmptyContainer(stub)
	}  else if function == "GetContainerDetailsForOwner" {
		return t.GetContainerDetailsForOwner(stub, args[0])
	}else if function == "GetOwner" {
		return t.GetOwner(stub)
	}else if function == "GetUserAttribute" {
		return t.GetUserAttribute(stub, args[0])
	}else if function == "getProvenanceForContainer" {
		return t.getProvenanceForContainer(stub, args[0])
	}else if function == "SearchById" {
		return t.SearchById(stub, args[0])
	}else if function == "SearchByName" {
		return t.SearchByName(stub, args[0],args[1])
	}	
	fmt.Println("query did not find func: " + function)
	return nil, errors.New("Received unknown function query: " + function)
}

func (t *MedLabPharmaChaincode) init(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	maxIDCounter := UniqueIDCounter{
		ContainerMaxID: 0,
		PalletMaxID:    0}
	jsonVal, _ := json.Marshal(maxIDCounter)
	err := stub.PutState(UNIQUE_ID_COUNTER, []byte(jsonVal))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// write - invoke function to write key/value pair
func (t *MedLabPharmaChaincode) ShipContainerUsingLogistics(stub shim.ChaincodeStubInterface,
	senderID string, logisticsID string, receiverID string, remarks string, elementsJSON string,shipmentDate string) ([]byte, error) {
	var err error
	var containerId string
	shipment := Container{}
	json.Unmarshal([]byte(elementsJSON), &shipment)
	containerId=shipment.ContainerId
	valueAsbytes, err := stub.GetState(containerId)
	fmt.Println(string(valueAsbytes))	
	if(len(valueAsbytes)==0){
		fmt.Println("Validating duplicate containerID" + containerId)
		containerID, jsonValue := ShipContainerUsingLogistics_Internal(senderID, logisticsID, receiverID, remarks, elementsJSON,shipmentDate)
	    fmt.Println("running ShipContainerUsingLogistics.key:" + containerID)
	    fmt.Println(string(jsonValue))
	    valAsbytes, err := stub.GetState(containerID)
	    fmt.Println(string(valAsbytes))	
	    err = stub.PutState(containerID, jsonValue) //write the variable into the chaincode state
	    incrementCounter(stub) //increment the unique ids for container and Pallet
	    setCurrentOwner(stub, senderID, containerID)
	    setCurrentOwner(stub, logisticsID, containerID)
	   if err != nil {
		     return nil, err
	         }
	 }else{
		 fmt.Println("Container is already shipped cannot ship the same container ")
		 jsonResp := "{\"Error\":\"Container is already shipped cannot ship the same container \"}"
	     return nil, errors.New(jsonResp)
	}
	
	return nil, err

}
func (t *MedLabPharmaChaincode)DispatchContainer(stub shim.ChaincodeStubInterface,containerID string, receiverID string, remarks string,shipmentDate string) ([]byte, error) {
	var err error
		fmt.Println("running DispatchContainer:" + containerID)
     valAsbytes, err := stub.GetState(containerID)
	 if len(valAsbytes) == 0 {
		 	jsonResp := "{\"Error\":\"Failed to get state for Container id since there is no such container \"}"
		return nil, errors.New(jsonResp)
	 }
	 fmt.Println("json value from the container****************")
	 fmt.Println(valAsbytes)
	 if err != nil{
		jsonResp := "{\"Error\":\"Failed to get state for Container id \"}"
		return nil, errors.New(jsonResp)
	}
	 shipment := Container{}	  
	json.Unmarshal([]byte(valAsbytes), &shipment)
	shipment.Recipient = receiverID
	shipment.Remarks = ""
	conprov := shipment.Provenance  
    supplychain := conprov.Supplychain     
	chainActivity := ChainActivity{
		Sender:   shipment.Provenance.Receiver,//
		Receiver: receiverID,
		ShipmentDate :shipmentDate,
		Remarks: remarks,
		Status:   STATUS_DISPATCHED_BY_LOGISTICS ,		 
		}  
	supplychain = append(supplychain, chainActivity) 
	conprov.Supplychain = supplychain
   conprov.TransitStatus = STATUS_DISPATCHED_BY_LOGISTICS 
   conprov.Sender = shipment.Provenance.Receiver
   conprov.Receiver = receiverID
   shipment.Provenance = conprov
    shipment.ShipmentDate=shipmentDate
    jsonVal, _ := json.Marshal(shipment)
   	err = stub.PutState(containerID, jsonVal)//write the variable into the chaincode state
    if err != nil{
		jsonResp := "{\"Error\":\"Failed to put state for Container id \"}"
		return nil, errors.New(jsonResp)
	}
	fmt.Println("********DISPATCHED JSON***********")	
	fmt.Println("SHIPMENTDATE",shipment.Provenance.ShipmentDate)	
	fmt.Println(string(jsonVal))	
	setCurrentOwner(stub, receiverID, containerID)

	if err != nil {
		return nil, err
	}
	return nil, nil

}

// read - query function to read key/value pair
func (t *MedLabPharmaChaincode) GetContainerDetails(stub shim.ChaincodeStubInterface, container_id string) ([]byte, error) {
	fmt.Println("runnin GetContainerDetails ")
	var key, jsonResp string
	var err error

	if container_id == "" {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	fmt.Println("key:" + container_id)
	valAsbytes, err := stub.GetState(container_id)
	fmt.Println(valAsbytes)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

//Returns the maximum number used for ContainerID and PalletID in the format "ContainerMaxNumber, PalletMaxNumber"
func (t *MedLabPharmaChaincode) GetMaxIDValue(stub shim.ChaincodeStubInterface) ([]byte, error) {
	var jsonResp string
	var err error
	ConMaxAsbytes, err := stub.GetState(UNIQUE_ID_COUNTER)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for ContainerMaxNumber \"}"
		return nil, errors.New(jsonResp)
	}
	return ConMaxAsbytes, nil
}

func (t *MedLabPharmaChaincode) GetEmptyContainer(stub shim.ChaincodeStubInterface) ([]byte, error) {
	ConMaxAsbytes, err := stub.GetState(UNIQUE_ID_COUNTER)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for ContainerMaxNumber \"}"
		return nil, errors.New(jsonResp)
	}

	counter := UniqueIDCounter{}
	json.Unmarshal([]byte(ConMaxAsbytes), &counter)
	containerID := "CON" + strconv.Itoa(counter.ContainerMaxID+1)
	pallets := createPallet(containerID, counter.PalletMaxID+1)
	conelement := ContainerElements{Pallets: pallets}
	container := Container{
		ContainerId: containerID,
		Elements:    conelement}
	jsonVal, _ := json.Marshal(container)
	return jsonVal, nil
}

func ShipContainerUsingLogistics_Internal(senderID string,
	logisticsID string, receiverID string, remarks string, elementsJSON string,shipmentDate string) (string, []byte) {
		shipment := Container{}
	json.Unmarshal([]byte(elementsJSON), &shipment)
		chainActivity := ChainActivity{
		Sender:   senderID,
		Receiver: logisticsID,
		Status:   STATUS_SHIPPED,
		Remarks: remarks,
		ShipmentDate :shipmentDate,
		}
		var supplyChain []ChainActivity
	supplyChain = append(supplyChain, chainActivity)
	conprov := ContainerProvenance{
		TransitStatus: STATUS_SHIPPED,
		Sender:        senderID,
		Receiver:      logisticsID,
		Supplychain:   supplyChain}
	shipment.Recipient = receiverID
	shipment.ShipmentDate = shipmentDate
	shipment.Provenance = conprov
	//shipment.Remarks = remarks
	jsonVal, _ := json.Marshal(shipment)
	return shipment.ContainerId, jsonVal
}

func createUnit(caseID string) []Unit {
	units := make([]Unit, 3)

	for index := 0; index < 3; index++ {
		strIndex := strconv.Itoa(index + 1)
		unitid := caseID + "-UNIT" + strIndex
		units[index].UnitId = unitid
	}
	return units
}

func createCase(palletID string) []Case {
	cases := make([]Case, 3)

	for index := 0; index < 3; index++ {
		strIndex := strconv.Itoa(index + 1)
		caseid := palletID + "-CASE" + strIndex
		cases[index].CaseId = caseid
		cases[index].Units = createUnit(caseid)
	}
	return cases
}

func createPallet(containerID string, palletMaxID int) []Pallet {
	pallets := make([]Pallet, 3)
	for index := 0; index < 3; index++ {
		strMaxID := strconv.Itoa(palletMaxID)
		palletid := containerID + "-PAL" + strMaxID
		pallets[index].PalletId = palletid
		pallets[index].Cases = createCase(palletid)
		palletMaxID++
	}
	return pallets
}
// func RemoveIndex(s [3]string, index int)  []string  {
//            return append(s[:index], s[index+1:]...)
//             }

func All(vs [3]string, f func(string) bool) bool {
    for _, v := range vs {
        if !f(v) {
            return false
        }
    }
    return true
}

func validatePallet(shippedpallets []Pallet,dispatchedpallets []Pallet)([]Pallet, error) {
	var i, j int
	var t,u int
	fmt.Println("Am in validate pallet method")
	fmt.Println("mismatched list")
	var s int=0
    mismatchedlist:=[...]string{dispatchedpallets[s].PalletId,dispatchedpallets[s+1].PalletId,dispatchedpallets[s+2].PalletId}      		
    fmt.Println(mismatchedlist)		 
	for u=0; u < len(shippedpallets); u++ {	  	 
	  	for t=0; t < len(mismatchedlist); t++ {
		     if (shippedpallets[u].PalletId==mismatchedlist[t]){
 			 	mismatchedlist[t]=""
			  }
		  }	
	}  
	fmt.Println("Am printing using any in pallet")
	fmt.Println(mismatchedlist)
	flag:=All(mismatchedlist, func( v string) bool {        
	     return v == ""
    })
    fmt.Println(flag)
	if(flag==false){
			   fmt.Println(" {\"Error\":\"In valid Pallet IDs \"} ")
			   jsonResp := "{\"Error\":\"In valid Pallet IDs  \"}"
		       return nil, errors.New(jsonResp) 
		       } 
     if (len(shippedpallets)==len(dispatchedpallets)){
		  for i=0; i < len(shippedpallets); i++ {
                 for j=0; j < len(dispatchedpallets); j++ {
                        if (shippedpallets[i].PalletId==dispatchedpallets[j].PalletId){
						   	 fmt.Println(shippedpallets[i].PalletId)
						   	  fmt.Println(dispatchedpallets[j].PalletId)							  
							  flag1,_,count,counter:= validateCases(shippedpallets[i].Cases,dispatchedpallets[j].Cases)
							  if(counter>0){
								  fmt.Println("Invalid Container because of invalid Caseid")
								  jsonResp := "{\"Error\":\"Invalid Container because of  Invalid Case IDs \"}"
		                            return nil, errors.New(jsonResp)
							  }else if(count>0){
								  fmt.Println("Invalid Container because of invalid Unitid")
								  jsonResp := "{\"Error\":\"Invalid Container because of Invalid  Unit IDs \"}"
		                            return nil, errors.New(jsonResp)
							  }
							  fmt.Println(flag1)
							  fmt.Println("Test for cases")
							  fmt.Println(counter)
							  fmt.Println("Test for units")
							  fmt.Println(count)
						      if (flag1==false){
			                        fmt.Println(" {\"Error\":\"Invalid Container because of Palletid\"} ")
			                        jsonResp := "{\"Error\":\"Invalid Container because of Palletid \"}"
		                            return nil, errors.New(jsonResp) 
		                            }else if (dispatchedpallets[j].Health=="Healthy"){
                                      shippedpallets[i].Health="Healthy"
									  fmt.Println("pallet health is updated as Healthy")									  
						            }else if (dispatchedpallets[j].Health=="Partially Healthy"){
									  shippedpallets[i].Health="Partially Healthy"	
								      fmt.Println("pallet health is updated as Partially Healthy")
							         }else if (dispatchedpallets[j].Health=="UnHealthy"){
									  shippedpallets[i].Health="UnHealthy"	 
								      fmt.Println("pallet health is updated as un Healthy")
							         }		   
                       break
					  }						 
				}
		   }
		    
		   return shippedpallets,nil		  
	  }else{
		      fmt.Println("pallet lengths  are  not equal")
			  jsonResp := "{\"Error\":\"pallet lengths  are  not equal \"}"
		      return nil, errors.New(jsonResp)
	      }  		
}
func repackagedPallets(parentContainerId string,childContainerID string,acceptedpallets []Pallet,dispatchedpallets []Pallet)([]Pallet, error,bool) {
	var u,k int
	var find bool
	var find1 bool=true
	fmt.Println("Am in repackagedPallets")
	fmt.Println("dispatchedpallets")
	fmt.Println(dispatchedpallets)
	fmt.Println("acceptedpallets")
	fmt.Println(acceptedpallets)
    for	k=0;k<len(acceptedpallets);k++{
	       fmt.Println("checking if All the pallets in the parent conatiner are not repackaged in this child container")
           for u=0; u < len(dispatchedpallets); u++ {	
                  if(dispatchedpallets[u].Health=="Healthy"){	  	
	                     if(acceptedpallets[k].PalletId==dispatchedpallets[u].PalletId){
		                       fmt.Println(acceptedpallets[k].PalletId)
		                       fmt.Println(dispatchedpallets[u].PalletId)	   
		                       find = strings.Contains(dispatchedpallets[u].PalletId,parentContainerId)
		                       fmt.Println(dispatchedpallets[u].PalletId)
							   fmt.Println(parentContainerId)
		                       fmt.Println("Am printing the value of finds in repackagedpallets")
		                       fmt.Println(find)
		                       if(find){
                                   dispatchedpallets[u].PalletId=strings.Replace(dispatchedpallets[u].PalletId, parentContainerId+"-", childContainerID+"-", -1)
					               fmt.Println("after replacing it with the child container")
					               fmt.Println(dispatchedpallets[u].PalletId)
					               fmt.Println(acceptedpallets[k].PalletId)
		   	                       repackagedCases,_:=repackagedCases(parentContainerId,childContainerID,dispatchedpallets[u].Cases)
				                   fmt.Println("Cases after repackaging")
                                   fmt.Println(repackagedCases)
				                   } else{
					                      fmt.Println("match not found for container id")
				                          }
	                     }else{
							     fmt.Println("PalletIds doesnot match")
				                 
	                            }
	               }else{
                             find1=false
		                     fmt.Println("Unhealthy pallets cannot be repackaged")
							 return nil,nil,find1
		                  }
	          } 
     }	
	fmt.Println(parentContainerId)
	fmt.Println(childContainerID)
	fmt.Println(dispatchedpallets)
	return dispatchedpallets,nil,find1
}
func repackagedCases(parentContainerId string,childContainerID string,dispatchedCases []Case)([]Case, error) {
	var v int
	var find bool
	fmt.Println("Am in repackagedCases")
	for v=0; v < len(dispatchedCases); v++ {	
		       find = strings.Contains(dispatchedCases[v].CaseId,parentContainerId)
		       fmt.Println(dispatchedCases[v].CaseId)
		       fmt.Println(parentContainerId)
		       fmt.Println("Am printing the value of finds in repackagedCases")
		       fmt.Println(find)
		       if(find){
                  dispatchedCases[v].CaseId=strings.Replace(dispatchedCases[v].CaseId, parentContainerId+"-", childContainerID+"-", -1)
		   	      repackagedUnits,_:=repackagedUnits(parentContainerId,childContainerID,dispatchedCases[v].Units)
				  fmt.Println("Units after Repackaging ")
				  fmt.Println(repackagedUnits)	 
				} else{
					 fmt.Println("match not found for parent containerid")
				     }
	} 	
	fmt.Println(parentContainerId)
	fmt.Println(childContainerID)
	fmt.Println(dispatchedCases)
	return dispatchedCases,nil
}
func repackagedUnits(parentContainerId string,childContainerID string,dispatchedUnits []Unit)([]Unit, error) {
	var w int
	var find bool
	fmt.Println("Am in repackagedUnits")
	for w=0; w < len(dispatchedUnits); w++ {	
	     
		       find = strings.Contains(dispatchedUnits[w].UnitId,parentContainerId)
		       fmt.Println(dispatchedUnits[w].UnitId)
		       fmt.Println(parentContainerId)
		       fmt.Println("Am printing the value of finds in repackagedUnits")
		       fmt.Println(find)
		       if(find){
                 dispatchedUnits[w].UnitId=strings.Replace(dispatchedUnits[w].UnitId, parentContainerId+"-", childContainerID+"-", -1)
			   } else{
					 fmt.Println("match not found for parent containerid")
				 }
	} 	
	fmt.Println(parentContainerId)
	fmt.Println(childContainerID)
	fmt.Println(dispatchedUnits)
	return dispatchedUnits,nil
}
 func validateCases(shippedcases []Case,dispatchedcases []Case)(bool, error,int,int) {
    var k,l int
	var f,g int
	fmt.Println("Am in validate cases method")
	fmt.Println("mismatched list in cases")
	var v int=0
	var counter int=0
	var count int=0
    mismatchedlist:=[...]string{dispatchedcases[v].CaseId,dispatchedcases[v+1].CaseId,dispatchedcases[v+2].CaseId}      		
    fmt.Println(mismatchedlist)		 
	for f=0; f < len(shippedcases); f++ {	  	 
	  	for g=0; g < len(mismatchedlist); g++ {	 
		     if (shippedcases[f].CaseId==mismatchedlist[g]){
 			 	mismatchedlist[f]=""
			  }
		  }	
	}  
	fmt.Println("Am printing using any in cases")
    fmt.Println(mismatchedlist)
		     flag1:=All(mismatchedlist, func( v string) bool {        
		return v == ""
    })
	fmt.Println(flag1)
	if(flag1==false){
			   fmt.Println(" {\"Error\":\"In valid Case IDs \"} ")
			   counter++			   
		       return flag1, nil,count,counter
		       } 
	if (len(shippedcases)==len(dispatchedcases)){
		for k=0; k < len(shippedcases); k++ {
			for l=0; l < len(dispatchedcases); l++ {
               if (shippedcases[k].CaseId==dispatchedcases[l].CaseId){
				   fmt.Println(shippedcases[k].CaseId)
				   fmt.Println(dispatchedcases[l].CaseId)
				   flag2,_,count:= validateUnits(shippedcases[k].Units,dispatchedcases[l].Units)
				   fmt.Println("Testing units in Validate cases")
				   fmt.Println(flag2)
				   if (flag2==false){
					   flag1=flag2
					    return flag1,nil,count,counter
				   }
					if (dispatchedcases[l].Health=="Healthy"){
                         shippedcases[k].Health="Healthy"
						 fmt.Println("case health is updated as healthy")
					}else if(dispatchedcases[l].Health=="Partially Healthy"){
						shippedcases[k].Health="Partially Healthy"
						fmt.Println("case health is  updated as partially healthy")	   
					}else if(dispatchedcases[l].Health=="UnHealthy"){
						shippedcases[k].Health="UnHealthy"
						fmt.Println("case health is  updated as Unhealthy")	   
					}
				    break
     			}else{
						 fmt.Println("Case ids are not equal")
					}
			}
		}		
        return flag1,nil,count,counter
   }else{
		   fmt.Println("case lengths are not  equal")
		   jsonResp := "{\"Error\":\"case lengths are not  equal \"}"
		    return flag1, errors.New(jsonResp),count,counter
		  //return flag1,nil
	   }
    
}
func validateUnits(shippedunits []Unit,dispatchedunits []Unit)(bool, error,int) {
	var m,n int
	var y,z int
	fmt.Println("mismatched list unit list in Validate Units")
	var x int=0
	var count int=0
	fmt.Println("Am in validate units method")
    mismatchedlist:=[...]string{dispatchedunits[x].UnitId,dispatchedunits[x+1].UnitId,dispatchedunits[x+2].UnitId}      		
         fmt.Println(mismatchedlist)		 
	for y=0; y < len(shippedunits); y++ {	  	 
	  	for z=0; z < len(mismatchedlist); z++ {	 
		     if (shippedunits[y].UnitId==mismatchedlist[z]){
 			 	mismatchedlist[z]=""
			  }
		  }	
	}  
	fmt.Println("Am printing using any in Units")
    fmt.Println(mismatchedlist)
		     flag2:=All(mismatchedlist, func( v string) bool {        
		return v == ""
    })
	if(flag2==false){
			   fmt.Println(" {\"Error\":\"In valid Unit IDs \"} ")
			   count++
			    return flag2, nil,count 
		       } 
	if (len(shippedunits)==len(dispatchedunits)){	
		for m=0; m < len(dispatchedunits); m++ {	
			for n=0; n < len(dispatchedunits); n++ {
               if (shippedunits[m].UnitId==dispatchedunits[n].UnitId){
				   fmt.Println(shippedunits[m].UnitId)
				   fmt.Println(dispatchedunits[n].UnitId)
				     if (dispatchedunits[n].Health=="Healthy"){
                            shippedunits[m].Health="Healthy"
							fmt.Println("Unit health is updated as Healthy")
					   }else if (dispatchedunits[n].Health=="Pratially Healthy"){
						      shippedunits[m].Health="Pratially Healthy"
							fmt.Println("Unit health is updated as Partially Healthy")
					   }else if (dispatchedunits[n].Health=="UnHealthy"){
						      shippedunits[m].Health="UnHealthy"
							fmt.Println("Unit health is updated as UnHealthy")
					         }
					break
     			}else{
						   fmt.Println("Unit ids are not equal")
										 
					    }
		   }
	   }			  	
         return	flag2,nil,count
 }else{
		   fmt.Println("unit lengths are not  equal")
   		    return flag2, nil,count
	}  
	
}

func incrementCounter(stub shim.ChaincodeStubInterface) error {
	ConMaxAsbytes, err := stub.GetState(UNIQUE_ID_COUNTER)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for ContainerMaxNumber \"}"
		return errors.New(jsonResp)
	}
	counter := UniqueIDCounter{}
	json.Unmarshal([]byte(ConMaxAsbytes), &counter)
	counter.ContainerMaxID = counter.ContainerMaxID + 1
	counter.PalletMaxID = counter.PalletMaxID + 3
	jsonVal, _ := json.Marshal(counter)
	err = stub.PutState(UNIQUE_ID_COUNTER, []byte(jsonVal))
	if err != nil {
		return err
	}
	return nil
}

func (t *MedLabPharmaChaincode) SetCurrentOwnerTest(stub shim.ChaincodeStubInterface, ownerID string, containerID string) ([]byte, error) {
	err := setCurrentOwner(stub, ownerID, containerID)
	return []byte("success"), err
}

func (t *MedLabPharmaChaincode) GetContainerDetailsForOwner(stub shim.ChaincodeStubInterface, ownerID string) ([]byte, error) {

	fmt.Println("Fetching container details for Owner:" + ownerID)

	ConMaxAsbytes, err := stub.GetState(CONTAINER_OWNER)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for Container Owners \"}"
		return nil, errors.New(jsonResp)
	}
	ConOwners := ContainerOwners{}
	json.Unmarshal([]byte(ConMaxAsbytes), &ConOwners)

	var containerList []string
	var matchFound bool

	for index := range ConOwners.Owners {
		if ConOwners.Owners[index].OwnerId == ownerID {
			containerList = ConOwners.Owners[index].ContainerList
			matchFound = true
			break
		}
	}
	if matchFound {
		fmt.Println("MatchFound for Owner:" + ownerID)
		shipment := Shipment{}
	
		for _, containerID := range containerList {
			byteVal, _ := t.GetContainerDetails(stub, containerID)
			container := Container{}

			json.Unmarshal([]byte(byteVal), &container)
			shipment.ContainerList = append(shipment.ContainerList, container)
		}
		jsonVal, _ := json.Marshal(shipment)
		return jsonVal, nil
	} else {
		fmt.Println("Container details not found for Owner:" + ownerID)
		return nil, errors.New("Unable to get container details for Owner:" + ownerID)
	}
}
func (t *MedLabPharmaChaincode) GetOwner(stub shim.ChaincodeStubInterface) ([]byte, error) {

	ConMaxAsbytes, err := stub.GetState(CONTAINER_OWNER)
	fmt.Println("************Am in GET OWNER Method**********")
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for Container Owners \"}"
		return nil, errors.New(jsonResp)
	}
	return ConMaxAsbytes, nil
}
func (t *MedLabPharmaChaincode) AcceptContainerbyLogistics(stub shim.ChaincodeStubInterface,containerID string, logisticsID string, receiverID string, remarks string,date string) ([]byte, error) {

	fmt.Println("Accepting the  container by Logistics:" + logisticsID)
	fmt.Println("Accepting the  container by Logistics:" + containerID)
     valAsbytes, err := stub.GetState(containerID)
	 if len(valAsbytes) == 0 {
		 	jsonResp := "{\"Error\":\"Failed to get state for Container id since there is no such container \"}"
		return nil, errors.New(jsonResp)
	 }
	 fmt.Println("json value from the container****************")
	 fmt.Println(valAsbytes)
	 if err != nil{
		jsonResp := "{\"Error\":\"Failed to get state for Container id \"}"
		return nil, errors.New(jsonResp)
	}
	//timeLayOut := timePresent.Format(RFC1123)
	  shipment := Container{}	  
	json.Unmarshal([]byte(valAsbytes), &shipment)
	shipment.Recipient = receiverID
	shipment.Remarks=remarks
	shipment.ReceivedDate=date
	conprov := shipment.Provenance  
    supplychain := conprov.Supplychain     
	chainActivity := ChainActivity{
		Sender:   shipment.Provenance.Sender,
		Receiver: logisticsID,
		Remarks: remarks,
		ShipmentDate :date,
		Status:   STATUS_ACCEPTED_BY_LOGISTICS,		 
		}  
	supplychain = append(supplychain, chainActivity) 
	conprov.Supplychain = supplychain
   conprov.TransitStatus = STATUS_ACCEPTED_BY_LOGISTICS
   conprov.Sender = shipment.Provenance.Sender
   conprov.Receiver = logisticsID
   shipment.Provenance = conprov
   jsonVal, _ := json.Marshal(shipment)
   	err = stub.PutState(containerID, jsonVal)
    if err != nil{
		jsonResp := "{\"Error\":\"Failed to put state for Container id \"}"
		return nil, errors.New(jsonResp)
	}	
	fmt.Println(string(jsonVal))
	fmt.Println(string(shipment.Provenance.Sender))
	setCurrentOwner(stub, logisticsID, containerID)
	return nil, nil		
}
func (t *MedLabPharmaChaincode) RejectContainerbyLogistics(stub shim.ChaincodeStubInterface,containerID string, logisticsID string, receiverID string, remarks string,date string) ([]byte, error) {

	fmt.Println("Rejecting the  container by Logistics:" + logisticsID + containerID)
     valAsbytes, err := stub.GetState(containerID)
	 if len(valAsbytes) == 0 {
		 	jsonResp := "{\"Error\":\"Failed to get state for Container id since there is no such container \"}"
		return nil, errors.New(jsonResp)
	 }
	 fmt.Println("json value from the container****************")
	 fmt.Println(valAsbytes)
	 if err != nil{
		jsonResp := "{\"Error\":\"Failed to get state for Container id \"}"
		return nil, errors.New(jsonResp)
	}
	fmt.Println(remarks)
	if len(remarks) == 0 {
		 	jsonResp := "{\"Error\":\"Failed to have the remarks  for Container id since there is no input remarks \"}"
		return nil, errors.New(jsonResp)
	 }
	//timeLayOut := timePresent.Format(RFC1123)
	  shipment := Container{}	  
	json.Unmarshal([]byte(valAsbytes), &shipment)
	shipment.Recipient = receiverID
	shipment.Remarks = remarks
	shipment.ReceivedDate = date
	conprov := shipment.Provenance  
    supplychain := conprov.Supplychain     
	chainActivity := ChainActivity{
		Sender:   shipment.Provenance.Sender,
		Receiver: logisticsID,
		ShipmentDate :date,
		Remarks: remarks,
		Status:   STATUS_REJECTED_BY_LOGISTICS,		 
		}  
	supplychain = append(supplychain, chainActivity) 
	conprov.Supplychain = supplychain
   conprov.TransitStatus = STATUS_REJECTED_BY_LOGISTICS
   conprov.Sender = shipment.Provenance.Sender
   conprov.Receiver = logisticsID
   shipment.Provenance = conprov
   jsonVal, _ := json.Marshal(shipment)
   	err = stub.PutState(containerID, jsonVal)
    if err != nil{
		jsonResp := "{\"Error\":\"Failed to put state for Container id \"}"
		return nil, errors.New(jsonResp)
	}	
	fmt.Println(string(jsonVal))
	fmt.Println("SENDER",shipment.Provenance.Sender)
		setCurrentOwner(stub, logisticsID, containerID)
	return nil, nil		
}

func (t *MedLabPharmaChaincode) UpdateContainerbyDistributor(stub shim.ChaincodeStubInterface,containerID string, receiverID string, remarks string,elementsJSON string,date string) ([]byte, error) {
    var m int
	var count int=0
    fmt.Println("Running UpdateContainerbyDistributor ")
	fmt.Println("UpdateContainerbyDistributor:" + containerID)
     valAsbytes, err := stub.GetState(containerID)
	 if len(valAsbytes) == 0 {
		 	jsonResp := "{\"Error\":\"Failed to get state for Container id since there is no such container \"}"
		return nil, errors.New(jsonResp)
	 }
	 fmt.Println("json value from the container****************")
	 fmt.Println(valAsbytes)
	 if err != nil{
		jsonResp := "{\"Error\":\"Failed to get state for Container id \"}"
		return nil, errors.New(jsonResp)
	}		
	dispatchedJSON :=Container{}
	json.Unmarshal([]byte(elementsJSON), &dispatchedJSON)
	 dispatchedpallets :=dispatchedJSON.Elements.Pallets
	 shipment := Container{}	  
	json.Unmarshal([]byte(valAsbytes), &shipment)
	shippedPallets :=shipment.Elements.Pallets
	 updatedJSON :=Container{}
	 json.Unmarshal([]byte(elementsJSON), &updatedJSON)
	 updatedpallets,err :=validatePallet(shippedPallets,dispatchedpallets)
	 fmt.Println(" updatedpallets")
	 fmt.Println( updatedpallets)
     fmt.Println("begining")
	 fmt.Println( shippedPallets)
	 fmt.Println("dispatched pallets")
	 fmt.Println( dispatchedpallets)
	 fmt.Println("ending")
	 for m=0; m < len(updatedpallets); m++ {
		 fmt.Println("Am in update container by distributor and updating the container health")
		 if(updatedpallets[m].Health=="UnHealthy"){
		     count++
		}
	 }	
		fmt.Println(count)
	 if(count==0){
			shipment.Elements.Health="Healthy"
			fmt.Println("Am in update container by distributor and updated as healthy")
			shipment.Elements.Pallets=updatedpallets
	    //shipment.Recipient = receiverID
		shipment.Remarks=remarks
		shipment.ReceivedDate=date
	    conprov := shipment.Provenance  
        supplychain := conprov.Supplychain     
	    chainActivity := ChainActivity{
		Sender:   shipment.Provenance.Sender,
		Receiver: receiverID,
		ShipmentDate :date,
		Remarks: remarks,
		Status:   STATUS_ACCEPTED_BY_DISTRIBUTOR,		 
		}  
	supplychain = append(supplychain, chainActivity) 
	conprov.Supplychain = supplychain
    conprov.TransitStatus = STATUS_ACCEPTED_BY_DISTRIBUTOR
    conprov.Sender = shipment.Provenance.Sender//taking sender from the container to avoid inconsistency of sender from UI
    conprov.Receiver = receiverID  
    shipment.Provenance = conprov
	
		   }else if (count>=1)&&(count<3){            
			shipment.Elements.Health="PartialHealthy"
			fmt.Println("Am in update container by distributor and updated as PartialHealthy")
			shipment.Elements.Pallets=updatedpallets
	    //shipment.Recipient = receiverID
		shipment.Remarks=remarks
		shipment.ReceivedDate=date
	    conprov := shipment.Provenance  
        supplychain := conprov.Supplychain     
	    chainActivity := ChainActivity{
		Sender:   shipment.Provenance.Sender,
		Receiver: receiverID,
		ShipmentDate :date,
		Remarks: remarks,
		Status:   STATUS_PARTIALLY_ACCEPTED_BY_DISTRIBUTOR,		 
		}  
	supplychain = append(supplychain, chainActivity) 
	conprov.Supplychain = supplychain
    conprov.TransitStatus = STATUS_PARTIALLY_ACCEPTED_BY_DISTRIBUTOR
    conprov.Sender = shipment.Provenance.Sender//taking sender from the container to avoid inconsistency of sender from UI
    conprov.Receiver = receiverID  
    shipment.Provenance = conprov
	fmt.Println(shipment.Provenance)			
		}else if (count==3){         
			shipment.Elements.Health="UnHealthy"
			fmt.Println("Am in update container by distributor and updated as UnHealthy") 
			shipment.Elements.Pallets=updatedpallets
	    //shipment.Recipient = receiverID
		shipment.Remarks=remarks
		shipment.ReceivedDate=date
	    conprov := shipment.Provenance  
        supplychain := conprov.Supplychain     
	    chainActivity := ChainActivity{
		Sender:   shipment.Provenance.Sender,
		Receiver: receiverID,
		ShipmentDate :date,
		Remarks: remarks,
		Status:   STATUS_REJECTED_BY_DISTRIBUTOR,		 
		}  
	supplychain = append(supplychain, chainActivity) 
	conprov.Supplychain = supplychain
    conprov.TransitStatus = STATUS_REJECTED_BY_DISTRIBUTOR
    conprov.Sender = shipment.Provenance.Sender//taking sender from the container to avoid inconsistency of sender from UI
    conprov.Receiver = receiverID  
    shipment.Provenance = conprov
		   }	
   jsonVal, _ := json.Marshal(shipment)
   	err = stub.PutState(containerID, jsonVal)
    if err != nil{
		jsonResp := "{\"Error\":\"Failed to put state for Container id \"}"
		return nil, errors.New(jsonResp)
	}
	fmt.Println("JSON ACCEPTED BY Reciever")	
	fmt.Println(string(jsonVal))
	fmt.Println(receiverID)
	setCurrentOwner(stub, receiverID, containerID)	
	return jsonVal, nil		
}
func (t *MedLabPharmaChaincode) getProvenanceForContainer(stub shim.ChaincodeStubInterface, ContainerID string) ([]byte,error) {
	var m,s int
	var y int	
	var count int=0
	var count1 int=0
	//var count2 int=0
	fmt.Println("*****getProvenanceForContainer****** " + ContainerID)
	valAsbytes, err := stub.GetState(ContainerID)
	if len(valAsbytes) == 0 {
		 	jsonResp := "{\"Error\":\"Failed to get state for Container id since there is no such container \"}"
		return nil, errors.New(jsonResp)
	 }
	 fmt.Println("json value from the container****************")
	 fmt.Println(valAsbytes)
	 if err != nil{
		jsonResp := "{\"Error\":\"Failed to get state for Container id \"}"
		return nil, errors.New(jsonResp)
	}		
	 shipment := Container{}	  
	 json.Unmarshal([]byte(valAsbytes), &shipment)
	 if(len(shipment.ParentContainerId)!=0){
		         fmt.Println("It has parent provenance to be attached")
		         fmt.Println(shipment.ParentContainerId)
	 	         valueAsbytes, err := stub.GetState(shipment.ParentContainerId)
	             if len(valueAsbytes) == 0 {
		 	                jsonResp := "{\"Error\":\"Failed to get state for Container id since there is no such container \"}"
		                    return nil, errors.New(jsonResp)
	                        }
	            fmt.Println("json value from the container****************")
	            fmt.Println(valueAsbytes)
	            if err != nil{
		                   jsonResp := "{\"Error\":\"Failed to get state for Container id \"}"
		                   return nil, errors.New(jsonResp)
	                       }	
	            parentshipment := Container{}
	            json.Unmarshal([]byte(valueAsbytes), &parentshipment)
	            fmt.Println("Am Printing Parent shipment provenance")
	            fmt.Println(parentshipment.Provenance)
	            parentConProv := parentshipment.Provenance 
                parentSupplyChain := parentConProv.Supplychain 
	            childprov:=shipment.Provenance
	            for m=0; m < len(childprov.Supplychain); m++ {
		 	             parentSupplyChain=append(parentSupplyChain, childprov.Supplychain[m])
	                     }
	            childsupplychain:=parentSupplyChain
	            childprov.Supplychain=childsupplychain
	            fmt.Println("new conprov")
                fmt.Println(childprov)
	            fmt.Println("ending conprov")
	            fmt.Println(parentSupplyChain) 
                shipment.Provenance=childprov
 	            fmt.Println("Am Printing child shipment provenance")
	            fmt.Println(shipment.Provenance)	
	            jsonVal, _ := json.Marshal(shipment)
   	            fmt.Println(string(jsonVal)) 
	            return jsonVal, nil	 	    
	 }else if(len(shipment.ChildContainerId)!=0){
                 valuesAsbytes, err := stub.GetState(ContainerID)
	              if len(valuesAsbytes) == 0 {
		 	                   jsonResp := "{\"Error\":\"Failed to get state for child  Container id since there is no such container \"}"
		                        return nil, errors.New(jsonResp)
	                           }
	               fmt.Println("json value from the child container****************")
	               if err != nil{
		                        jsonResp := "{\"Error\":\"Failed to get state for child Container id \"}"
		                         return nil, errors.New(jsonResp)
	                           }
                   mainshipment := Container{}	  
	               json.Unmarshal([]byte(valuesAsbytes), &mainshipment)
                   mainConProv := mainshipment.Provenance                  
	     		   fmt.Println(len(shipment.ChildContainerId))
		           for y=0; y < len(shipment.ChildContainerId); y++ {		 
		                         newchild:=shipment.ChildContainerId[y]				   
		                         fmt.Println("new child")
		                         fmt.Println(newchild)				  
		                          valsAsbytes, err := stub.GetState(newchild)
	                              if len(valsAsbytes) == 0 {
		 	                            jsonResp := "{\"Error\":\"Failed to get state for child  Container id since there is no such container \"}"
		                                return nil, errors.New(jsonResp)
	                                    }
	         	                  if err != nil{
		                                 jsonResp := "{\"Error\":\"Failed to get state for child Container id \"}"
		                                  return nil, errors.New(jsonResp)
	                                     }
                                  childshipment := Container{}	  
	                              json.Unmarshal([]byte(valsAsbytes), &childshipment)
                                   parentSupplyChain1 := mainConProv.Supplychain 
	                               newConprov:=childshipment.Provenance
	                               newSupplyChain:=newConprov.Supplychain 								  
							       fmt.Println("Parent container has the following  children")		  	          
					                for s=0; s < len(newSupplyChain); s++ { 		
							                   newSupplyChain[s].Remarks="ChildContainerId: "+newchild+ " - " +newSupplyChain[s].Remarks	 
							                   parentSupplyChain1=append(parentSupplyChain1, newSupplyChain[s])
			                                   fmt.Println(newSupplyChain[s])									  
									           fmt.Println(newSupplyChain[s].Remarks)
                                      }
									  fmt.Println("printing count values")
									  fmt.Println(count)
									  fmt.Println(count1)
						  fmt.Println(parentSupplyChain1)						 
			                          mainConProv.Supplychain=parentSupplyChain1				
	                                  mainshipment.Provenance=mainConProv           
			            }				
	                fmt.Println("new parent mainConProv")
                    fmt.Println(mainConProv)
	                fmt.Println("ending child mainconprov")
	                jsonVal, _ := json.Marshal(mainshipment)
   	                fmt.Println(string(jsonVal))
			        return 	jsonVal,nil	
			 		
	           }else {
				   fmt.Println(" Provenance For the individual container without parent and child")
				           fmt.Println("*****getProvenanceForContainer****** " + ContainerID)
	                       valAsbytes, err := stub.GetState(ContainerID)
	                        if len(valAsbytes) == 0 {
		 	                          jsonResp := "{\"Error\":\"Failed to get state for Container id since there is no such container \"}"
		                              return nil, errors.New(jsonResp)
	                         }
	                     fmt.Println("json value from the container****************")
	                     fmt.Println(valAsbytes)
	                     if err != nil{
		                          jsonResp := "{\"Error\":\"Failed to get state for Container id \"}"
		                          return nil, errors.New(jsonResp)
	                              }		
	                    shipment := Container{}	  
	                    json.Unmarshal([]byte(valAsbytes), &shipment)
						jsonVal, _ := json.Marshal(shipment)
   	                    fmt.Println(string(jsonVal))
						return jsonVal,nil	
			   }
			  return nil,nil
}
func removeDuplicates(elements []string) []string{
    // Use map to record duplicates as we find them.
	fmt.Println("Removing Duplicates in Repackaging status")
    encountered := map[string]bool{}
    result := []string{}

    for v := range elements {
        if encountered[elements[v]] == true {
            // Do not add duplicate.
        } else {
            // Record this element as an encountered element.
            encountered[elements[v]] = true
            // Append to result slice.
            result = append(result, elements[v])
        }
    }
    // Return the new slice.
    return result
}

func (t *MedLabPharmaChaincode) repackagingContainerbyDistributor(stub shim.ChaincodeStubInterface,childContainerID string,containerID string, receiverID string, remarks string,elementsJSON string,shipmentDate string) ([]byte, error) {
	 var m,n int
	 var count int=0
	 var repackagingstatu1 []string
	 fmt.Println("Repackaging Container by Distributor:" + childContainerID)
	 valuAsbytes, err := stub.GetState(containerID)
	 shipment := Container{}	  
	 json.Unmarshal([]byte(valuAsbytes), &shipment)
	 shipment.ChildContainerId = append(shipment.ChildContainerId,childContainerID)
	  acceptedPallets :=shipment.Elements.Pallets
	  fmt.Println(acceptedPallets)
	  dispatchedshipment := Container{}	 
      json.Unmarshal([]byte(elementsJSON), &dispatchedshipment)
	  dispatchedshipment.ParentContainerId=containerID
	  dispatchedshipment.ContainerId=childContainerID
	 // dispatchedshipment.Recipient = receiverID
	  dispatchedPallets :=dispatchedshipment.Elements.Pallets
	  if(len(shipment.ChildContainerId)==1){
		     shipment.Recipient = receiverID
		    fmt.Println("This is the first child getting repackaged")
		    fmt.Println( shipment.ChildContainerId)
		    conprov1 := shipment.Provenance  
            supplychain1 := conprov1.Supplychain     
	    chainActivity1 := ChainActivity{
		Sender:   shipment.Provenance.Receiver,
		Receiver: receiverID,
		ShipmentDate :shipmentDate,
		Remarks: remarks,
		Status:   STATUS_DISPATCH_IN_PROGRESS,		 
		}  
	supplychain1 = append(supplychain1, chainActivity1) 
	conprov1.Supplychain = supplychain1
    conprov1.TransitStatus = STATUS_DISPATCH_IN_PROGRESS
    conprov1.Sender = shipment.Provenance.Receiver//taking sender from the container to avoid inconsistency of sender from UI
    conprov1.Receiver = receiverID  
    shipment.Provenance = conprov1
	shipment.ShipmentDate = shipmentDate
	shipment.Recipient = receiverID		
	  }	
	  for n=0; n < len(dispatchedPallets); n++ {		  
               repackagingstatu1=append(shipment.Repackagingstatus,dispatchedPallets[n].PalletId )
			   fmt.Println("before removing duplicates")
			   fmt.Println(repackagingstatu1)
			   shipment.Repackagingstatus= removeDuplicates(repackagingstatu1)
			   fmt.Println("after removing duplicates")
			   fmt.Println(shipment.Repackagingstatus)
	  }
	  fmt.Println("Printing Repackaged pallets in parent container")
	  fmt.Println(shipment.Repackagingstatus)
	 if ((len(valuAsbytes) == 0) || (err != nil)) {
		     fmt.Println("Failed to get state for  containerID ")
		 	 jsonResp := "{\"Error\":\"Failed to get state for Container id \"}"
	      	 return nil, errors.New(jsonResp)		    
	 }else{
		   fmt.Println("Parent container")
		   jsonVal, _ := json.Marshal(shipment)
		   fmt.Println("JSON updated for  parent Container")	
	       fmt.Println(string(jsonVal))
   	       err = stub.PutState(containerID, jsonVal)
           if err != nil{
		              jsonResp := "{\"Error\":\"Failed to put state for Container id \"}"
		              return nil, errors.New(jsonResp)
	                 }
	     }
		  fmt.Println(dispatchedPallets)
		  repackagedpallets,_,find2:=repackagedPallets(containerID,childContainerID,acceptedPallets,dispatchedPallets)
          fmt.Println("Repackaged Pallets")
		  fmt.Println(repackagedpallets)
          fmt.Println(find2)
		  if(!find2){
			  fmt.Println("Unhealthy Pallets cannot be repackaged")
			  jsonResp := "{\"Error\":\"Unhealthy Pallets cannot be repackaged \"}"
	 	      return nil, errors.New(jsonResp)
		  }else{
			 dispatchedshipment.Elements.Pallets=repackagedpallets
                   for m=0; m < len(dispatchedPallets); m++ {
		                   fmt.Println("Checking the pallet health in repackaging by distributor")
		                   if(repackagedpallets[m].Health=="Healthy"){
		                        count++
			                    
		                     }
			        }				 
		  }
		  if(len(repackagedpallets)==count){
			  dispatchedshipment.Elements.Health="Healthy"
               fmt.Println("updating pallet health in repackaging by distributor")
			   fmt.Println(dispatchedshipment.Elements.Health)
		  }
		 dispatchedshipment.Recipient = "" 
		 fmt.Println("Printing dispatchshipment provenance")
		 chainActivity := ChainActivity{
		            Sender:  shipment.Provenance.Sender,
		            Receiver: receiverID,
					ShipmentDate :shipmentDate,
					Remarks: remarks,
		            Status:   STATUS_SHIPPED_BY_DISTRIBUTOR,
		     }
		        var supplyChain []ChainActivity
	            supplyChain = append(supplyChain, chainActivity)
	            conprov := ContainerProvenance{
		        TransitStatus: STATUS_SHIPPED_BY_DISTRIBUTOR,
		        Sender:        shipment.Provenance.Sender,
		        Receiver:      receiverID,
		        Supplychain:   supplyChain}
		        conprov.Receiver = receiverID  
                dispatchedshipment.Provenance = conprov
				dispatchedshipment.ShipmentDate = shipmentDate
				dispatchedshipment.Recipient = receiverID
	            //dispatchedshipment.Remarks=remarks
		   jsonVall, _ := json.Marshal(dispatchedshipment)
   	       err = stub.PutState(childContainerID, jsonVall)
           if err != nil{
		        jsonResp := "{\"Error\":\"Failed to put state for child Container id \"}"
		         return nil, errors.New(jsonResp)
	       } 
		   fmt.Println("Final child container obtained")
		   fmt.Println(string(jsonVall))
		   fmt.Println(shipment.Provenance.Receiver)
		   fmt.Println(receiverID)
		  // setCurrentOwner(stub, shipment.Provenance.Receiver, childContainerID)
	       setCurrentOwner(stub, receiverID, childContainerID)
	      return jsonVall, nil		
	
}
func (t *MedLabPharmaChaincode) AcceptContainerbyRetailer(stub shim.ChaincodeStubInterface,containerID string, receiverID string, remarks string,date string) ([]byte, error) {

	fmt.Println("Accepting the  container by Retailer:" + containerID)
     valAsbytes, err := stub.GetState(containerID)
	 if len(valAsbytes) == 0 {
		 	jsonResp := "{\"Error\":\"Failed to get state for Container id since there is no such container \"}"
		return nil, errors.New(jsonResp)
	 }
	 fmt.Println("json value from the container****************")
	 fmt.Println(valAsbytes)
	 if err != nil{
		jsonResp := "{\"Error\":\"Failed to get state for Container id \"}"
		return nil, errors.New(jsonResp)
	}
	shipment := Container{}	  
	json.Unmarshal([]byte(valAsbytes), &shipment)
	shipment.Recipient = receiverID
	shipment.Remarks=remarks
	shipment.ReceivedDate=date
	conprov := shipment.Provenance  
    supplychain := conprov.Supplychain     
	chainActivity := ChainActivity{
		Sender:   shipment.Provenance.Sender,
		Receiver: receiverID,
		ShipmentDate :date,
		Remarks: remarks,
		Status:   STATUS_ACCEPTED_BY_RETAILER,		 
		}  
	supplychain = append(supplychain, chainActivity) 
	conprov.Supplychain = supplychain
   conprov.TransitStatus = STATUS_ACCEPTED_BY_RETAILER
   conprov.Sender = shipment.Provenance.Sender
   conprov.Receiver = receiverID
   shipment.Provenance = conprov
   jsonVal, _ := json.Marshal(shipment)
   	err = stub.PutState(containerID, jsonVal)
    if err != nil{
		jsonResp := "{\"Error\":\"Failed to put state for Container id \"}"
		return nil, errors.New(jsonResp)
	}	
	fmt.Println(string(jsonVal))
	fmt.Println(string(shipment.Provenance.Sender))
	fmt.Println(receiverID)
	setCurrentOwner(stub, receiverID, containerID)
	return jsonVal, nil		
}
func (t *MedLabPharmaChaincode)RejectContainerbyRetailer(stub shim.ChaincodeStubInterface,containerID string, receiverID string, remarks string,date string) ([]byte, error) {

	fmt.Println("Rejecting the  container by Retailer:" + containerID)
     valAsbytes, err := stub.GetState(containerID)
	 if len(valAsbytes) == 0 {
		 	jsonResp := "{\"Error\":\"Failed to get state for Container id since there is no such container \"}"
		return nil, errors.New(jsonResp)
	 }
	 fmt.Println("json value from the container****************")
	 fmt.Println(valAsbytes)
	 if err != nil{
		jsonResp := "{\"Error\":\"Failed to get state for Container id \"}"
		return nil, errors.New(jsonResp)
	}
	shipment := Container{}	  
	json.Unmarshal([]byte(valAsbytes), &shipment)
	shipment.Recipient = receiverID
	shipment.Remarks=remarks
	shipment.ReceivedDate=date
	conprov := shipment.Provenance  
    supplychain := conprov.Supplychain     
	chainActivity := ChainActivity{
		Sender:   shipment.Provenance.Sender,
		Receiver: receiverID,
		ShipmentDate :date,
		Remarks: remarks,
		Status:   STATUS_REJECTED_BY_RETAILER,		 
		}  
	supplychain = append(supplychain, chainActivity) 
	conprov.Supplychain = supplychain
   conprov.TransitStatus = STATUS_REJECTED_BY_RETAILER
   conprov.Sender = shipment.Provenance.Sender
   conprov.Receiver = receiverID
   shipment.Provenance = conprov
   jsonVal, _ := json.Marshal(shipment)
   	err = stub.PutState(containerID, jsonVal)
    if err != nil{
		jsonResp := "{\"Error\":\"Failed to put state for Container id \"}"
		return nil, errors.New(jsonResp)
	}	
	fmt.Println(string(jsonVal))
	fmt.Println(string(shipment.Provenance.Sender))
	fmt.Println(receiverID)
	setCurrentOwner(stub, receiverID, containerID)
	return jsonVal, nil		
}

func (t *MedLabPharmaChaincode)RejectingbyConsumer(stub shim.ChaincodeStubInterface,containerID string, receiverID string, remarks string,date string) ([]byte, error) {

	fmt.Println("Rejection  by Consumer:" + containerID)
     valAsbytes, err := stub.GetState(containerID)
	 if len(valAsbytes) == 0 {
		 	jsonResp := "{\"Error\":\"Failed to get state for Container id since there is no such container \"}"
		return nil, errors.New(jsonResp)
	 }
	 fmt.Println("json value from the container****************")
	 fmt.Println(valAsbytes)
	 if err != nil{
		jsonResp := "{\"Error\":\"Failed to get state for Container id \"}"
		return nil, errors.New(jsonResp)
	}
	shipment := Container{}	  
	json.Unmarshal([]byte(valAsbytes), &shipment)
	shipment.Recipient = receiverID
	shipment.Remarks=remarks
	shipment.ReceivedDate=date
	conprov := shipment.Provenance  
    supplychain := conprov.Supplychain     
	chainActivity := ChainActivity{
		Sender:   shipment.Provenance.Sender,
		Receiver: receiverID,
		ShipmentDate :date,
		Remarks: remarks,
		Status:   STATUS_REJECTED_BY_CONSUMER,		 
		}  
	supplychain = append(supplychain, chainActivity) 
	conprov.Supplychain = supplychain
   conprov.TransitStatus = STATUS_REJECTED_BY_CONSUMER
   conprov.Sender = shipment.Provenance.Sender
   conprov.Receiver = receiverID
   shipment.Provenance = conprov
   jsonVal, _ := json.Marshal(shipment)
   	err = stub.PutState(containerID, jsonVal)
    if err != nil{
		jsonResp := "{\"Error\":\"Failed to put state for Container id \"}"
		return nil, errors.New(jsonResp)
	}	
	fmt.Println(string(jsonVal))
	fmt.Println(string(shipment.Provenance.Sender))
	fmt.Println(receiverID)
	setCurrentOwner(stub, receiverID, containerID)
	return jsonVal, nil		
}
func searchedPalletsById(acceptedpallets []Pallet,ID string)([]Pallet, error,int) {
	 fmt.Println("In searchedPallets by Id")
	 fmt.Println(acceptedpallets)
	 fmt.Println(ID)
	 var a,count2 int
	 for a=0; a < len(acceptedpallets); a++ {
        checkCases,_,count:=searchedCasesById(acceptedpallets[a].Cases,ID)
		fmt.Println("after first pallet")
		fmt.Println(checkCases)
		fmt.Println(count)
		if(count==1){
			count2=1
			return acceptedpallets,nil,count
		}else{
			count2=2
			fmt.Println(acceptedpallets[a].PalletId)
		    }
	 }
	 fmt.Println("count in searched Pallets")
	 fmt.Println(count2)
	 	return acceptedpallets,nil,count2
}
func searchedCasesById(acceptedCases []Case,ID string)([]Case, error,int) {
	 fmt.Println("In searchedCases by Id")
	 fmt.Println(acceptedCases)
	 var b int
	 var count1 int
	 for b=0; b < len(acceptedCases); b++ {
        checkUnits,_,count:=searchedUnitsById(acceptedCases[b].Units,ID)
		fmt.Println(checkUnits)
		fmt.Println(count)
		if(count==1){
			count1=1
			return acceptedCases,nil,count1
		}else{
			count1=2	
			fmt.Println(acceptedCases[b].CaseId)	
		}
	 }
	 fmt.Println("count in searched cases")
	 fmt.Println(count1)
	return acceptedCases,nil,count1
}
func searchedUnitsById(acceptedUnits []Unit,ID string)([]Unit, error,int) {
	 fmt.Println("In searchedUnits by ID")
	 fmt.Println(acceptedUnits)
	 fmt.Println(ID)
	 var c,count int
	 for c=0; c < len(acceptedUnits); c++ {		 
              if(acceptedUnits[c].UnitId==ID){
                 fmt.Println("match occurred")
				 fmt.Println(acceptedUnits[c].SaleStatus)
				 fmt.Println(ID)
				 count =1
				 fmt.Println(count)
				 return acceptedUnits,nil,count
			     }else{
			            count=2
						fmt.Println(acceptedUnits[c].UnitId)			
		              }         
		      
	            }
				fmt.Println("count in searched units")
				fmt.Println(count)
		return acceptedUnits,nil,count
}

func (t *MedLabPharmaChaincode) SearchById(stub shim.ChaincodeStubInterface,ID string) ([]byte, error) {
    fmt.Println("This Method searches by ID" + ID)
    var string2 []string
    var containerID string
	var flag,flag1 bool
	var m,y,s,count,count1,countOfhyphen int
	fmt.Println("running by SearchById:" + ID)
	countOfhyphen=strings.Count("subCON5-PAL15-CASE2-UNIT3", "-")
	flag=strings.Contains(ID, "-")   
	if((flag==true)&&(countOfhyphen==3)){
          fmt.Println(strings.Contains(ID, "-"))
          string2= strings.Split(ID, "-")
          fmt.Println("My string is",string2[0])    
	      containerID=string2[0]
	      flag1=strings.Contains(containerID, "CON")
	      if(flag1){
                    valAsbytes, err := stub.GetState(containerID)
	                 if len(valAsbytes) == 0 {
		 	                    jsonResp := "{\"Error\":\"Failed to get state for Container id since there is no such container \"}"
		                        return nil, errors.New(jsonResp)
	                           }
	                 fmt.Println("json value from the container****************")
	                 fmt.Println(valAsbytes)
	                 if err != nil{
	            	 jsonResp := "{\"Error\":\"Failed to get state for Container id \"}"
		             return nil, errors.New(jsonResp)
	                 }
	  shipment := Container{}	  
	  json.Unmarshal([]byte(valAsbytes), &shipment)
	  acceptedpallets:=shipment.Elements.Pallets
	  tempConProv := shipment.Provenance 
      fmt.Println(tempConProv)
	  fmt.Println(tempConProv.TransitStatus)
	  //consumerName=tempConProv.Receiver
	  if((tempConProv.TransitStatus==STATUS_PARTIALLY_SOLD_BY_RETAILER)||(tempConProv.TransitStatus==STATUS_SOLD_BY_RETAILER)||(tempConProv.TransitStatus==STATUS_ACCEPTED_BY_RETAILER)){
	        searchedPallets,_,resultcount:=searchedPalletsById(acceptedpallets,ID)	
	        fmt.Println(searchedPallets)
	        if(len(shipment.ParentContainerId)!=0){
		         fmt.Println("It has parent provenance to be attached")
		         fmt.Println(shipment.ParentContainerId)
				 if(resultcount==1) { 
	                      fmt.Println("Checking if there is the unit in the repackaged container")
                          fmt.Println(resultcount)
	 	                  valueAsbytes, err := stub.GetState(shipment.ParentContainerId)
	                      if len(valueAsbytes) == 0 {
		 	                     jsonResp := "{\"Error\":\"Failed to get state for Container id since there is no such container \"}"
		                         return nil, errors.New(jsonResp)
	                             }
	     	              fmt.Println(valueAsbytes)
	                      if err != nil{
		                          jsonResp := "{\"Error\":\"Failed to get state for Container id \"}"
		                          return nil, errors.New(jsonResp)
	                             }	
	                      parentshipment := Container{}
	                      json.Unmarshal([]byte(valueAsbytes), &parentshipment)
	                      fmt.Println("Am Printing Parent shipment provenance")
	                      fmt.Println(parentshipment.Provenance)
	                      parentConProv := parentshipment.Provenance 
                          parentSupplyChain := parentConProv.Supplychain 
	                      childprov:=shipment.Provenance
	                      for m=0; m < len(childprov.Supplychain); m++ {
		 	                        parentSupplyChain=append(parentSupplyChain, childprov.Supplychain[m])
	                               }
						  childsupplychain:=parentSupplyChain
	                      childprov.Supplychain=childsupplychain
	                      fmt.Println("new conprov")
                          fmt.Println(childprov)
	                      fmt.Println("ending conprov")
	                      fmt.Println(parentSupplyChain) 
                          shipment.Provenance=childprov
 	                      fmt.Println("Am Printing child shipment provenance")
	                      fmt.Println(shipment.Provenance)	
	                      jsonVal, _ := json.Marshal(shipment)
   	                      fmt.Println(string(jsonVal)) 
	                      return jsonVal, nil	 	    
	          }else{
			             if(count!=1){
							  fmt.Println("Unit Id may  not be in the repackaged container or it might not be sold")
						      m := Response{"Unit Id may  not be in the repackaged container or it might not be sold"}
                              jsonval, _ := json.Marshal(m)
                              fmt.Println(jsonval)
                              fmt.Println(string(jsonval))
 		                      return jsonval,nil
		   	             }	                                    
		             }
	 }else if(len(shipment.ChildContainerId)!=0){
                 valuesAsbytes, err := stub.GetState(containerID)
	              if len(valuesAsbytes) == 0 {
		 	                   jsonResp := "{\"Error\":\"Failed to get state for child  Container id since there is no such container \"}"
		                        return nil, errors.New(jsonResp)
	                           }
	               fmt.Println("json value from the child container****************")
	               if err != nil{
		                        jsonResp := "{\"Error\":\"Failed to get state for child Container id \"}"
		                         return nil, errors.New(jsonResp)
	                           }
                   mainshipment := Container{}	  
	               json.Unmarshal([]byte(valuesAsbytes), &mainshipment)
                   mainConProv := mainshipment.Provenance                  
	     		   fmt.Println(len(shipment.ChildContainerId))
		           for y=0; y < len(shipment.ChildContainerId); y++ {		 
		                         newchild:=shipment.ChildContainerId[y]				   
		                         fmt.Println("new child")
		                         fmt.Println(newchild)				  
		                          valsAsbytes, err := stub.GetState(newchild)
	                              if len(valsAsbytes) == 0 {
		 	                            jsonResp := "{\"Error\":\"Failed to get state for child  Container id since there is no such container \"}"
		                                return nil, errors.New(jsonResp)
	                                    }
	         	                  if err != nil{
		                                 jsonResp := "{\"Error\":\"Failed to get state for child Container id \"}"
		                                  return nil, errors.New(jsonResp)
	                                     }
                                  childshipment := Container{}	  
	                              json.Unmarshal([]byte(valsAsbytes), &childshipment)
                                   parentSupplyChain1 := mainConProv.Supplychain 
	                               newConprov:=childshipment.Provenance
	                               newSupplyChain:=newConprov.Supplychain 								  
							       fmt.Println("Parent container has the following  children")		  	          
					                for s=0; s < len(newSupplyChain); s++ { 		
							                   newSupplyChain[s].Remarks="ChildContainerId: "+newchild+ " - " +newSupplyChain[s].Remarks	 
							                   parentSupplyChain1=append(parentSupplyChain1, newSupplyChain[s])
			                                   fmt.Println(newSupplyChain[s])									  
									           fmt.Println(newSupplyChain[s].Remarks)
                                      }
									  fmt.Println("printing count values")
									  fmt.Println(count)
									  fmt.Println(count1)
						              fmt.Println(parentSupplyChain1)						 
			                          mainConProv.Supplychain=parentSupplyChain1				
	                                  mainshipment.Provenance=mainConProv           
			            }				
	                fmt.Println("new parent mainConProv")
                    fmt.Println(mainConProv)
	                fmt.Println("ending child mainconprov")
	                jsonVal, _ := json.Marshal(mainshipment)
   	                fmt.Println(string(jsonVal))
			        return 	jsonVal,nil	
			 		
	           }else {
				   fmt.Println(" Provenance For the individual container without parent and child")
				           fmt.Println("*****getProvenanceForContainer****** " + containerID)
	                       valAsbytes, err := stub.GetState(containerID)
	                        if len(valAsbytes) == 0 {
		 	                          jsonResp := "{\"Error\":\"Failed to get state for Container id since there is no such container \"}"
		                              return nil, errors.New(jsonResp)
	                         }
	                     fmt.Println("json value from the container****************")
	                     fmt.Println(valAsbytes)
	                     if err != nil{
		                          jsonResp := "{\"Error\":\"Failed to get state for Container id \"}"
		                          return nil, errors.New(jsonResp)
	                              }		
	                    shipment := Container{}	  
	                    json.Unmarshal([]byte(valAsbytes), &shipment)
						jsonVal, _ := json.Marshal(shipment)
   	                    fmt.Println(string(jsonVal))
						return jsonVal,nil	
			   }
			  
			}else{
                  fmt.Println("Units which are not accepted by retailer cannot be seen")
				  m := Response{"Units which are not accepted by retailer cannot be seen"}
                  jsonval, _ := json.Marshal(m)
                  fmt.Println(jsonval)
                  fmt.Println(string(jsonval))
 		          return jsonval,nil
		      }	
	  }else{
		         fmt.Println("Though Id is seperated by '-' it doesnot contain Valid UnitID Format")
				 m := Response{"Though Id is seperated by '-' it doesnot contain Valid UnitID Format"}
                 jsonval, _ := json.Marshal(m)
                 fmt.Println(jsonval)
                 fmt.Println(string(jsonval))
 		         return jsonval,nil
	          }
	    
	}else{
		  fmt.Println("Entered ID is not in valid format")
		  m := Response{"Entered ID is not in valid format"}
          jsonval, _ := json.Marshal(m)
          fmt.Println(jsonval)
          fmt.Println(string(jsonval))
 		  return jsonval,nil
	    }
	return nil,nil
}
func searchByNameInPallets(pallets []Pallet, drugname  string,gename  string)([]Pallet, error,bool) {
	 fmt.Println("am in check pallets")
	 fmt.Println(pallets)
	 fmt.Println(drugname)
	  fmt.Println(gename)
	 var a int
	 //var result bool
	 var finalresult bool
	 for a=0; a < len(pallets); a++ {
        searchedpallets,_,result:=searchByNameInCases(pallets[a].Cases,drugname,gename)
		fmt.Println(searchedpallets)
		fmt.Println(result)
		if(result==false){
			finalresult=false
		}else{
		    finalresult=true
		}
	 }	
	 fmt.Println(finalresult)	
	return pallets,nil,finalresult
}
func searchByNameInCases(Cases []Case, drugname  string,gename  string)([]Case, error,bool) {
	 fmt.Println("am in check Cases")
	 fmt.Println(Cases)
	 var b int
	 var finalresult bool
	 for b=0; b < len(Cases); b++ {
        checkedUnits,_,result:=searchByNameInUnits(Cases[b].Units,drugname,gename)
		fmt.Println("searchByNameInCases first time")
		fmt.Println(checkedUnits)
		fmt.Println(result)
		if(result==false){
			finalresult=false
		}else{
			finalresult=true
		}
	 }
	 fmt.Println("searchByNameInCases")
	 fmt.Println(finalresult)
	return Cases,nil,finalresult
}
func searchByNameInUnits(Units []Unit,drugname  string,gename  string)([]Unit, error,bool) {
	 fmt.Println("am in check Units")
	 fmt.Println(drugname)
	 fmt.Println(gename)
	 var c int
	 var count int=0
	 var flag bool=true
	 for c=0; c < len(Units); c++ {		
              if((Units[c].DrugName==drugname)&&Units[c].GenericName==gename){
                    fmt.Println("match occurred in searchByNameInUnits")				 
				    fmt.Println(Units[c].UnitId)
				}else{
					 count++
				     fmt.Println("match doesnot occurred in searchByNameInUnits")		
			        }
	 }
	 
	 if(count==27){
          flag=false
	 }
	 fmt.Println(count)
	 fmt.Println(flag)
	return Units,nil,flag
}


func (t *MedLabPharmaChaincode) SearchByName(stub shim.ChaincodeStubInterface, drugname  string,gename  string) ([]byte, error) {
    var a,b,m,y,s,count,count1 int
	fmt.Println("running SearchByName ")
	fmt.Println(CONTAINER_OWNER)
	ConMaxAsbytes, err := stub.GetState(CONTAINER_OWNER)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for Container Owners \"}"
		return nil, errors.New(jsonResp)
	}
	ConOwners := ContainerOwners{}
	json.Unmarshal([]byte(ConMaxAsbytes), &ConOwners)
	fmt.Println(ConOwners)
	var containerlist []string
	var finalContainer []string
	for a=0; a < len(ConOwners.Owners); a++ {
			containerlist = ConOwners.Owners[a].ContainerList
			fmt.Println("ContainerList starts:")
			fmt.Println(a)
			fmt.Println(ConOwners.Owners[a].ContainerList)
			fmt.Println(ConOwners.Owners[a])
			for b=0; b < len(containerlist); b++ {
			     fmt.Println(containerlist[b])
			     finalContainer=append(finalContainer, containerlist[b])	
			}
			fmt.Println("finalContainer starts:")
			fmt.Println(finalContainer)
		             
		}	
		finConwitoutDuplicates:= removeDuplicates(finalContainer)
		fmt.Println(finConwitoutDuplicates)
		fmt.Println(finalContainer)
		containers := Shipment{}
		for _, containerID := range finConwitoutDuplicates {
			byteVal, _ := t.GetContainerDetails(stub, containerID)
			shipment := Container{}			
			json.Unmarshal([]byte(byteVal), &shipment)
			fmt.Println(shipment.ContainerId)
            if(len(shipment.ParentContainerId)!=0){
		         fmt.Println("It has parent provenance to be attached")
		         fmt.Println(shipment.ParentContainerId)
	 	         valueAsbytes, err := stub.GetState(shipment.ParentContainerId)
	             if len(valueAsbytes) == 0 {
		 	                jsonResp := "{\"Error\":\"Failed to get state for Container id since there is no such container \"}"
		                    return nil, errors.New(jsonResp)
	                        }
	            fmt.Println("json value from the container****************")
	            fmt.Println(valueAsbytes)
	            if err != nil{
		                   jsonResp := "{\"Error\":\"Failed to get state for Container id \"}"
		                   return nil, errors.New(jsonResp)
	                       }	
	            parentshipment := Container{}
	            json.Unmarshal([]byte(valueAsbytes), &parentshipment)
	            fmt.Println("Am Printing Parent shipment provenance")
	            fmt.Println(parentshipment.Provenance)
	            parentConProv := parentshipment.Provenance 
                parentSupplyChain := parentConProv.Supplychain 
	            childprov:=shipment.Provenance
	            for m=0; m < len(childprov.Supplychain); m++ {
		 	             parentSupplyChain=append(parentSupplyChain, childprov.Supplychain[m])
	                     }
	            childsupplychain:=parentSupplyChain
	            childprov.Supplychain=childsupplychain
	            fmt.Println("new conprov")
                fmt.Println(childprov)
	            fmt.Println("ending conprov")
	            fmt.Println(parentSupplyChain) 
                shipment.Provenance=childprov
 	            fmt.Println("Am Printing child shipment provenance")
	            fmt.Println(shipment.Provenance)	
	         }else if(len(shipment.ChildContainerId)!=0){
                 valuesAsbytes, err := stub.GetState(containerID)
	              if len(valuesAsbytes) == 0 {
		 	                   jsonResp := "{\"Error\":\"Failed to get state for child  Container id since there is no such container \"}"
		                        return nil, errors.New(jsonResp)
	                           }
	               fmt.Println("json value from the child container****************")
	               if err != nil{
		                        jsonResp := "{\"Error\":\"Failed to get state for child Container id \"}"
		                         return nil, errors.New(jsonResp)
	                           }
                   mainshipment := Container{}	  
	               json.Unmarshal([]byte(valuesAsbytes), &mainshipment)
                   mainConProv := mainshipment.Provenance                  
	     		   fmt.Println(len(shipment.ChildContainerId))
		           for y=0; y < len(shipment.ChildContainerId); y++ {		 
		                         newchild:=shipment.ChildContainerId[y]				   
		                         fmt.Println("new child")
		                         fmt.Println(newchild)				  
		                          valsAsbytes, err := stub.GetState(newchild)
	                              if len(valsAsbytes) == 0 {
		 	                            jsonResp := "{\"Error\":\"Failed to get state for child  Container id since there is no such container \"}"
		                                return nil, errors.New(jsonResp)
	                                    }
	         	                  if err != nil{
		                                 jsonResp := "{\"Error\":\"Failed to get state for child Container id \"}"
		                                  return nil, errors.New(jsonResp)
	                                     }
                                  childshipment := Container{}	  
	                              json.Unmarshal([]byte(valsAsbytes), &childshipment)
                                   parentSupplyChain1 := mainConProv.Supplychain 
	                               newConprov:=childshipment.Provenance
	                               newSupplyChain:=newConprov.Supplychain 								  
							       fmt.Println("Parent container has the following  children")		  	          
					                for s=0; s < len(newSupplyChain); s++ { 		
							                   newSupplyChain[s].Remarks="ChildContainerId: "+newchild+ " - " +newSupplyChain[s].Remarks	 
							                   parentSupplyChain1=append(parentSupplyChain1, newSupplyChain[s])
			                                   fmt.Println(newSupplyChain[s])									  
									           fmt.Println(newSupplyChain[s].Remarks)
                                      }
									  fmt.Println("printing count values")
									  fmt.Println(count)
									  fmt.Println(count1)
						              fmt.Println(parentSupplyChain1)						 
			                          mainConProv.Supplychain=parentSupplyChain1				
	                                  mainshipment.Provenance=mainConProv           
			            }				
	                fmt.Println("new parent mainConProv")
                    fmt.Println(mainConProv)
	                fmt.Println("ending child mainconprov")
					shipment=mainshipment
		           }
			pallets:=shipment.Elements.Pallets
            searchedpallets,_,result:= searchByNameInPallets(pallets,drugname,gename)
			fmt.Println(searchedpallets)
			fmt.Println(result)
			fmt.Println("container operation")
			fmt.Println(shipment)	
			if(result==true){	
				fmt.Println("appending containers")	
			    containers.ContainerList = append(containers.ContainerList, shipment)
		     }else{
				fmt.Println("not appending containers")	 
			 }
			fmt.Println(containers)
			 		
		}
		jsonVal, _ := json.Marshal(containers)
		fmt.Println("final containers obtained")
		fmt.Println(string(jsonVal))
		return jsonVal, nil		
}
func checkPallets(acceptedpallets []Pallet,soldunits []string,customerID string)([]Pallet, error,bool) {
	 fmt.Println("am in check pallets")
	 fmt.Println(acceptedpallets)
	 fmt.Println(soldunits)
	 var s bool
	 var a int
	 for a=0; a < len(acceptedpallets); a++ {
        checkCases,_,find2:=checkCases(acceptedpallets[a].Cases,soldunits,customerID)
		fmt.Println(checkCases)
		fmt.Println(find2)
	 }
	//return dispatchedpallets,nil,find1
	return acceptedpallets,nil,s
}
func checkCases(acceptedCases []Case,soldunits []string,customerID string)([]Case, error,bool) {
	 fmt.Println("am in check Cases")
	 fmt.Println(acceptedCases)
	 var s1 bool=false
	 var b int
	 for b=0; b < len(acceptedCases); b++ {
        checkUnits,_,find2:=checkUnits(acceptedCases[b].Units,soldunits,customerID)
		fmt.Println(checkUnits)
		fmt.Println(find2)
	 }
	//return dispatchedpallets,nil,find1
	return acceptedCases,nil,s1
}
func checkUnits(acceptedUnits []Unit,soldunits []string,customerID string)([]Unit, error,bool) {
	 fmt.Println("am in check Units")
	 fmt.Println(acceptedUnits)
	 fmt.Println(soldunits)
	 var s2 bool=false
	 var c,d int
	 for c=0; c < len(acceptedUnits); c++ {
		 for  d=0; d < len(soldunits); d++ {
              if(acceptedUnits[c].UnitId==soldunits[d]){
                 fmt.Println("match occurred")
				 acceptedUnits[c].SaleStatus=STATUS_SOLD_BY_RETAILER
				 acceptedUnits[c].ConsumerName=customerID				 
				 fmt.Println(acceptedUnits[c].SaleStatus)
				 fmt.Println(acceptedUnits[c].UnitId)
				 fmt.Println(soldunits[d])
			     }
		            
		  }
	 }
	//return dispatchedpallets,nil,find1
	return acceptedUnits,nil,s2
}

func (t *MedLabPharmaChaincode) SellingbyRetailer(stub shim.ChaincodeStubInterface,containerID string, customerID string,UnitJson string, remarks string) ([]byte, error) {
     var  s,UnitIdLength int
	 var flag,firstIndex,secondIndex,thirdIndex,fourthIndex bool
	 var string2 []string
	 fmt.Println("Invoking Sellingby Retailer")
	 fmt.Println(containerID)
	 fmt.Println(customerID)
	 fmt.Println(UnitJson)
	unitshipment := UnitIDListJson{}	  
	json.Unmarshal([]byte(UnitJson), &unitshipment) 
	fmt.Println(unitshipment)
    sellunits:=unitshipment.UnitID
	fmt.Println("sellunits")
	fmt.Println(sellunits)
    for s=0; s < len(sellunits); s++ {
	     flag=strings.Contains(sellunits[s],"-")
	      if(flag){
		       string2= strings.Split(sellunits[s], "-")    
		       UnitIdLength=len(string2)
	           firstIndex=strings.Contains(string2[0], "CON")
               secondIndex=strings.Contains(string2[1], "PAL")
		       thirdIndex=strings.Contains(string2[2], "CASE")
		       fourthIndex=strings.Contains(string2[3], "UNIT")
		       fmt.Println("My string is",firstIndex,secondIndex,thirdIndex,fourthIndex) 
	           if((firstIndex&&secondIndex&&thirdIndex&&fourthIndex)&&(UnitIdLength==4)){
				   fmt.Println("unit id is  in the valid format")
	                valAsbytes,err := stub.GetState(containerID)
					 if err != nil{
		                   jsonResp := "{\"Error\":\"Failed to get state for Container id \"}"
		                   return nil, errors.New(jsonResp)
	                       }
	                 if len(valAsbytes) == 0 {
		 	                  jsonResp := "{\"Error\":\"Failed to get state for Container id since there is no such container \"}"
		                       return nil, errors.New(jsonResp)
	                          }
					
					          shipment := Container{}	  
	                          json.Unmarshal([]byte(valAsbytes), &shipment) 
					          acceptedPallets :=shipment.Elements.Pallets
					          checkedpallets,_,find2:=checkPallets(acceptedPallets,sellunits,customerID)							  
							  fmt.Println(find2)
							  soldPallets,_,find3,count:=sellPallets(checkedpallets)
				              fmt.Println(count)
							  fmt.Println("finding if few pallets were not sold")
							  fmt.Println(find3)
				              if(count==len(acceptedPallets)&&(find3==true)){
                              fmt.Println("updating salestatus in validatesalestatus ")				 				 
				                       shipment.Elements.Health=STATUS_SOLD_BY_RETAILER 
				                       conprov:=shipment.Provenance
									   fmt.Println(conprov)
                                       supplychain := conprov.Supplychain     
	                                   chainActivity := ChainActivity{
		                                    Sender:   shipment.Provenance.Receiver,
		                                    Receiver: "",
		                                   // ShipmentDate :date,
		                                    Remarks: remarks,
		                                    Status:   STATUS_SOLD_BY_RETAILER,		 
		                                  }  
	                                  supplychain = append(supplychain, chainActivity) 
	                                  conprov.Supplychain = supplychain
                                      conprov.TransitStatus =STATUS_SOLD_BY_RETAILER 		
                                      conprov.Sender = shipment.Provenance.Receiver
                                      conprov.Receiver = customerID                                                                           		  
				                      shipment.Provenance = conprov
							  }else if(find3==false){
								     fmt.Println("only few units are sold out")
								     conprov := shipment.Provenance  
                                     supplychain := conprov.Supplychain     
	                                 chainActivity := ChainActivity{
		                             Sender:   shipment.Provenance.Receiver,
		                             Receiver: customerID,
		                           //  ShipmentDate :date,
		                             Remarks: remarks,
		                             Status:   STATUS_PARTIALLY_SOLD_BY_RETAILER,		 
		                             }  
	                                  supplychain = append(supplychain, chainActivity) 
	                                  conprov.Supplychain = supplychain
									  shipment.Provenance = conprov
									 }
							   shipment.Elements.Pallets=soldPallets	 
							  jsonVal, _ := json.Marshal(shipment)
  	                          fmt.Println(string(jsonVal))
							  err=stub.PutState(containerID, jsonVal)
                              if err != nil{
		                         jsonResp := "{\"Error\":\"Failed to put state for Container id \"}"
		                         return nil, errors.New(jsonResp)
	                             }		
	                         setCurrentOwner(stub, customerID, containerID)
	                         return jsonVal, nil
		               }else{
			                  fmt.Println("Unit id is not in the valid format")
		                      jsonResp := "{\"Error\":\"Unit id is not in the valid format \"}"
		                      return nil, errors.New(jsonResp)
		                     }	      
             }else{
			      fmt.Println("Unit id is not in the valid format")
		          jsonResp := "{\"Error\":\"Unit id is not in the valid format \"}"
		          return nil, errors.New(jsonResp)
		        }	

       }		

     return nil,nil
}

func sellPallets(acceptedpallets []Pallet)([]Pallet, error,bool,int) {
	 fmt.Println("am in sell Pallets")
	 fmt.Println(acceptedpallets)	 
	 var s bool=true
	 var a int
	 var count1 int=0
    for a=0; a < len(acceptedpallets); a++ {
        checkCases,_,find2,count:=sellCases(acceptedpallets[a].Cases)
		fmt.Println(checkCases)
		s=find2
		if(count==3){
		   acceptedpallets[a].Health=STATUS_SOLD_BY_RETAILER     
		   fmt.Println("updating salestatus in sellPallets")
		   fmt.Println(acceptedpallets[a].Health)		  
	    }
		if(acceptedpallets[a].Health==STATUS_SOLD_BY_RETAILER){
                 fmt.Println("updating salestatus in sellpallets as sold out")				 				 
				 fmt.Println(acceptedpallets[a].Health)
				 fmt.Println(acceptedpallets[a].PalletId)
				 count1++
			  }
		 fmt.Println("printing the values of the count from sell cases")
	     fmt.Println(count)
	     fmt.Println(count1)	     
	 }	
	return acceptedpallets,nil,s,count1
}

func sellCases(acceptedCases []Case)([]Case, error,bool,int) {
	var count1 int=0
	 fmt.Println("am in sell Cases")
	 fmt.Println(acceptedCases)
	 var s1 bool=true
	 var b int
	// var count int=0
	 for b=0; b < len(acceptedCases); b++ {
        checkUnits,_,count:=sellUnits(acceptedCases[b].Units)
		fmt.Println(checkUnits)
	    fmt.Println(count)
	   if(count==3){
		   acceptedCases[b].Health=STATUS_SOLD_BY_RETAILER     
		   fmt.Println("updating salestatus in sellCases")
		   fmt.Println(acceptedCases[b].Health)	    
		 }else{
			 fmt.Println("few units were not sold")
             s1=false
		 }
		 if(acceptedCases[b].Health==STATUS_SOLD_BY_RETAILER){
                 fmt.Println("updating salestatus in sellCases as sold out")				 				 
				 fmt.Println(acceptedCases[b].Health)
				 fmt.Println(acceptedCases[b].CaseId)
				 count1++
			  }
	 }
	  fmt.Println("sending count to sellpallets")	 
	  fmt.Println(count1)	 
	  return acceptedCases,nil,s1,count1
}

func sellUnits(acceptedUnits []Unit)([]Unit, error,int) {
	 fmt.Println("am in sell Units")
	 fmt.Println(acceptedUnits)
	 var c int
	 var count int=0
	 for c=0; c < len(acceptedUnits); c++ {		 
              if(acceptedUnits[c].SaleStatus==STATUS_SOLD_BY_RETAILER){
                 fmt.Println("match occurred")				 				 
				 fmt.Println(acceptedUnits[c].SaleStatus)
				 fmt.Println(acceptedUnits[c].UnitId)
				 count++
			  }
	 }
	return acceptedUnits,nil,count
}
// func (t *MedLabPharmaChaincode) validatesalestatus(stub shim.ChaincodeStubInterface,containerID string,) ([]byte, error) {
//     fmt.Println("running validatesalestatus")
// 	 valAsbytes, _ := stub.GetState(containerID)
// 	 if len(valAsbytes) == 0 {
// 		 	     jsonResp := "{\"Error\":\"Failed to get state for Container id since there is no such container \"}"
// 		         return nil, errors.New(jsonResp)
// 	            }
// 				shipment := Container{}	  
// 	            json.Unmarshal([]byte(valAsbytes), &shipment)
// 				jsonVal, _ := json.Marshal(shipment)
// 				fmt.Println("Before Validating sale status")
//   	            fmt.Println(string(jsonVal)) 
// 				acceptedPallets :=shipment.Elements.Pallets	
// 	            fmt.Println(acceptedPallets)
// 				soldPallets,_,find2,count:=sellPallets(acceptedPallets)
// 				fmt.Println(count)
// 				if(count==len(acceptedPallets)){
//                  fmt.Println("updating salestatus in validatesalestatust")				 				 
// 				  shipment.Elements.Health=STATUS_SOLD_BY_RETAILER 
// 				  conprov:=shipment.Provenance
//                   conprov.TransitStatus =STATUS_SOLD_BY_RETAILER
				  
// 				  shipment.Provenance = conprov
//  			  }
// 				fmt.Println(soldPallets)
// 				fmt.Println(find2)
// 				shipment.Elements.Pallets=soldPallets				
// 				fmt.Println("After Validating sale status")
// 				jsonVall, _ := json.Marshal(shipment)
//   	            fmt.Println(string(jsonVall))
// 				stub.PutState(containerID, jsonVall)  

//      return nil,nil
// }
// func (t *MedLabPharmaChaincode) SellingbyRetailer(stub shim.ChaincodeStubInterface,containerID string, customerID string,UnitID string, remarks string) ([]byte, error) {
//     var m,n int
// 	var o,l int
// 	var string2 []string
// 	var containerid string
// 	var flag,flag1 bool
// 	var flag2,flag3,flag4 bool
// 	flag=strings.Contains(UnitID,"-")
// 	if(flag){
// 		  string2= strings.Split(UnitID, "-")
// 		  l=len(string2)
//           fmt.Println("My string is",l)    
// 	      containerid=string2[0]
// 	      flag1=strings.Contains(string2[0], "CON")
//           flag2=strings.Contains(string2[1], "PAL")
// 		  flag3=strings.Contains(string2[2], "CASE")
// 		  flag4=strings.Contains(string2[3], "UNIT")
// 		  fmt.Println("My string is",flag1,flag2,flag3,flag4) 
// 	      if(flag1&&flag2&&flag3&&flag4){
// 			  	 fmt.Println("Selling the unit by Retailer:" + containerid)
// 	             valAsbytes, err := stub.GetState(containerid)
// 	             if len(valAsbytes) == 0 {
// 		 	            jsonResp := "{\"Error\":\"Failed to get state for Container id since there is no such container \"}"
// 		                return nil, errors.New(jsonResp)
// 	                   }
// 	            fmt.Println("json value from the container****************")
// 	            fmt.Println(valAsbytes)
// 	            if err != nil{
// 		             jsonResp := "{\"Error\":\"Failed to get state for Container id \"}"
// 		             return nil, errors.New(jsonResp)
// 	                 }
//           shipment := Container{}	  
// 	      json.Unmarshal([]byte(valAsbytes), &shipment)
// 	      for m=0; m < len(shipment.Elements.Pallets); m++ {
//               for n=0; n < len(shipment.Elements.Pallets[m].Cases); n++ {
//                    for o=0; o < len(shipment.Elements.Pallets[m].Cases[n].Units); o++ {
//                       if(shipment.Elements.Pallets[m].Cases[n].Units[o].UnitId==UnitID){
//                             shipment.Elements.Pallets[m].Cases[n].Units[o].SaleStatus=STATUS_SOLD_BY_RETAILER
// 					        fmt.Println(shipment.Elements.Pallets[m].Cases[n].Units[o].SaleStatus)
// 					        fmt.Println(shipment.Elements.Pallets[m].Cases[n].Units[o].UnitId)
// 					        fmt.Println(UnitID)
// 				          }else{
// 					                fmt.Println("Am not updating sale status")
// 				                 }
// 	                 }
// 	             }
// 	       }
// 	      jsonVal, _ := json.Marshal(shipment)
//           err = stub.PutState(containerID, jsonVal)
//           if err != nil{
// 		          jsonResp := "{\"Error\":\"Failed to put state for Container id \"}"
// 		          return nil, errors.New(jsonResp)
// 	              }	
// 	       fmt.Println(string(jsonVal))
// 	       setCurrentOwner(stub, customerID, containerID)
// 	       return jsonVal, nil	
// 	    }else{
//                   fmt.Println("Unit id is not in the valid format")
// 		          jsonResp := "{\"Error\":\"Unit id is not in the valid format \"}"
// 		          return nil, errors.New(jsonResp)
// 		     }
// 	}else{
// 		          fmt.Println("Unit id is not in the valid format")
// 		          jsonResp := "{\"Error\":\"Unit id is not in the valid format \"}"
// 		          return nil, errors.New(jsonResp)
// 	      }
	
// }
func (t *MedLabPharmaChaincode) GetUserAttribute(stub shim.ChaincodeStubInterface, attributeName string) ([]byte,error) {
	fmt.Println("***** Inside GetUserAttribute() func for attribute:" + attributeName)
	attributeValue, err := stub.ReadCertAttribute(attributeName)
	fmt.Println("attributeValue=" + string(attributeValue))
	
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get GetUserAttribute\"}"
		return nil, errors.New(jsonResp)
	}
	return attributeValue, nil
}

func setCurrentOwner(stub shim.ChaincodeStubInterface, ownerID string, containerID string) error {
	ConMaxAsbytes, err := stub.GetState(CONTAINER_OWNER)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for ContainerMaxNumber \"}"
		return errors.New(jsonResp)
	}
	ConOwners := ContainerOwners{}
	json.Unmarshal([]byte(ConMaxAsbytes), &ConOwners)

	var containerList []string
	var ownerIndex int
	var matchFound bool
	for index := range ConOwners.Owners {
		if ConOwners.Owners[index].OwnerId == ownerID {
			ownerIndex = index
			containerList = ConOwners.Owners[index].ContainerList
			matchFound = true
			break
		}
	}
	containerFound := false
	if matchFound {
		for index := range containerList {
			if containerList[index] == containerID {
				containerFound = true
				break
			}
		}
		if !containerFound {
			containerList = append(containerList, containerID)
			ConOwners.Owners[ownerIndex].ContainerList = containerList
		}
	} else {
		containerList := make([]string, 1)
		containerList[0] = containerID
		owner := Owner{OwnerId: ownerID, ContainerList: containerList}
		ConOwners.Owners = append(ConOwners.Owners, owner)
	}

	jsonVal, _ := json.Marshal(ConOwners)
	err = stub.PutState(CONTAINER_OWNER, []byte(jsonVal))
	if err != nil {
		return err
	}

	return nil
}
