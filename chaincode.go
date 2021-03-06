/**
@author: Arshad Sarfarz
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
)

const STATUS_SHIPPED = "shipped by manufacturer"
const STATUS_ACCEPTED_BY_DISTRIBUTOR = "accepted by distributor"
const STATUS_ACCEPTED_BY_LOGISTICS= "accepted by logistics"
const STATUS_PARTIALLY_ACCEPTED_BY_DISTRIBUTOR = "partiallly accepted by distributor"
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
	CertifiedBy       string              `json:"certified_by"`   ///New fields
	Address           string              `json:"address"`        ///New fields
	USN               string              `json:"usn"`            ///New fields
	ShipmentDate      string              `json:"shipment_date"`  ///New fields
	InvoiceNumber     string              `json:"invoice_number"` ///New fields
	Remarks           string              `json:"remarks"`        ///New fields
	ReceivedDate      string              `json:"recieved_date"` 
    SenderAddress     string              `json:"sender_address"`
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
		     return t.ShipContainerUsingLogistics(stub, args[0], args[1], args[2], args[3], args[4])
		   }
	} else if function == "AcceptContainerbyLogistics"{
		  if (user_type =="logistics"){
			  return t.AcceptContainerbyLogistics(stub, args[0], args[1],args[2], args[3])
		  }	  
	}else if function == "DispatchContainer"{
		  if (user_type =="logistics"){
               return t.DispatchContainer(stub, args[0], args[1],args[2],args[3])	
		  } 	  		
	}else if function == "UpdateContainerbyDistributor"{
		if (user_type =="distributor"){
		         return t.UpdateContainerbyDistributor(stub, args[0], args[1],args[2],args[3])		
		}		   
	}else if function == "RejectContainerbyLogistics"{
		  if (user_type =="logistics"){
           	return t.RejectContainerbyLogistics(stub, args[0], args[1],args[2],args[3]) 
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
	senderID string, logisticsID string, receiverID string, remarks string, elementsJSON string) ([]byte, error) {
	var err error
	var containerId string
	shipment := Container{}
	json.Unmarshal([]byte(elementsJSON), &shipment)
	containerId=shipment.ContainerId
	valueAsbytes, err := stub.GetState(containerId)
	fmt.Println(string(valueAsbytes))	
	if(len(valueAsbytes)==0){
		fmt.Println("Validating duplicate containerID" + containerId)
		containerID, jsonValue := ShipContainerUsingLogistics_Internal(senderID, logisticsID, receiverID, remarks, elementsJSON)
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
	conprov := shipment.Provenance  
    supplychain := conprov.Supplychain     
	chainActivity := ChainActivity{
		Sender:   shipment.Provenance.Receiver,//
		Receiver: receiverID,
		ShipmentDate :shipmentDate,
		Status:   STATUS_DISPATCHED_BY_LOGISTICS ,		 
		}  
	supplychain = append(supplychain, chainActivity) 
	conprov.Supplychain = supplychain
   conprov.TransitStatus = STATUS_DISPATCHED_BY_LOGISTICS 
   conprov.Sender = shipment.Provenance.Receiver
   conprov.Receiver = receiverID
   shipment.Provenance = conprov
    shipment.Provenance.ShipmentDate=shipmentDate
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
	logisticsID string, receiverID string, remarks string, elementsJSON string) (string, []byte) {
		chainActivity := ChainActivity{
		Sender:   senderID,
		Receiver: logisticsID,
		Status:   STATUS_SHIPPED,
		}
		var supplyChain []ChainActivity
	supplyChain = append(supplyChain, chainActivity)
	conprov := ContainerProvenance{
		TransitStatus: STATUS_SHIPPED,
		Sender:        senderID,
		Receiver:      logisticsID,
		Supplychain:   supplyChain}
	shipment := Container{}
	json.Unmarshal([]byte(elementsJSON), &shipment)
	shipment.Recipient = receiverID
	shipment.Provenance = conprov
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
func (t *MedLabPharmaChaincode) AcceptContainerbyLogistics(stub shim.ChaincodeStubInterface,containerID string, logisticsID string, receiverID string, remarks string) ([]byte, error) {

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
	conprov := shipment.Provenance  
    supplychain := conprov.Supplychain     
	chainActivity := ChainActivity{
		Sender:   shipment.Provenance.Sender,
		Receiver: logisticsID,
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
func (t *MedLabPharmaChaincode) RejectContainerbyLogistics(stub shim.ChaincodeStubInterface,containerID string, logisticsID string, receiverID string, remarks string) ([]byte, error) {

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
	
	conprov := shipment.Provenance  
    supplychain := conprov.Supplychain     
	chainActivity := ChainActivity{
		Sender:   shipment.Provenance.Sender,
		Receiver: logisticsID,
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

func (t *MedLabPharmaChaincode) UpdateContainerbyDistributor(stub shim.ChaincodeStubInterface,containerID string, receiverID string, remarks string,elementsJSON string) ([]byte, error) {
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
		 fmt.Println("Am in update container by distributor!!!!!!!!!!!!!!!!")
		 if(updatedpallets[m].Health=="Unhealthy"){
		 count++
			fmt.Println("Am in update container by distributor")
		}
		fmt.Println(count)
		    if(count==0){
			shipment.Elements.Health="Healthy"
			fmt.Println("Am in update container by distributor and updated as healthy")
		   }else if (count==1){            
			shipment.Elements.Health="PartialHealthy"
			fmt.Println("Am in update container by distributor and updated as PartialHealthy")
		}else if (count==2){         
			shipment.Elements.Health="UnHealthy"
			fmt.Println("Am in update container by distributor and updated as UnHealthy") 
		   }
	}
	if (shipment.Elements.Health=="Healthy"){
		shipment.Elements.Pallets=updatedpallets
	    shipment.Recipient = receiverID
	    conprov := shipment.Provenance  
        supplychain := conprov.Supplychain     
	    chainActivity := ChainActivity{
		Sender:   shipment.Provenance.Sender,
		Receiver: receiverID,
		Status:   STATUS_ACCEPTED_BY_DISTRIBUTOR,		 
		}  
	supplychain = append(supplychain, chainActivity) 
	conprov.Supplychain = supplychain
    conprov.TransitStatus = STATUS_ACCEPTED_BY_DISTRIBUTOR
    conprov.Sender = shipment.Provenance.Sender//taking sender from the container to avoid inconsistency of sender from UI
    conprov.Receiver = receiverID  
    shipment.Provenance = conprov
	}else if (shipment.Elements.Health=="Partially Healthy"){
		shipment.Elements.Pallets=updatedpallets
	    shipment.Recipient = receiverID
	    conprov := shipment.Provenance  
        supplychain := conprov.Supplychain     
	    chainActivity := ChainActivity{
		Sender:   shipment.Provenance.Sender,
		Receiver: receiverID,
		Status:   STATUS_PARTIALLY_ACCEPTED_BY_DISTRIBUTOR,		 
		}  
	supplychain = append(supplychain, chainActivity) 
	conprov.Supplychain = supplychain
    conprov.TransitStatus = STATUS_PARTIALLY_ACCEPTED_BY_DISTRIBUTOR
    conprov.Sender = shipment.Provenance.Sender//taking sender from the container to avoid inconsistency of sender from UI
    conprov.Receiver = receiverID  
    shipment.Provenance = conprov
	}else if (shipment.Elements.Health=="UnHealthy"){
		shipment.Elements.Pallets=updatedpallets
	    shipment.Recipient = receiverID
	    conprov := shipment.Provenance  
        supplychain := conprov.Supplychain     
	    chainActivity := ChainActivity{
		Sender:   shipment.Provenance.Sender,
		Receiver: receiverID,
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
	setCurrentOwner(stub, receiverID, containerID)
	return nil, nil		
}

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
