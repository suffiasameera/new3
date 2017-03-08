package main

import "fmt"
import "errors"
import "encoding/json"
import "github.com/hyperledger/fabric/core/chaincode/shim"

type SampleChaincode struct {
}

type Incident struct {
IncidentID string `json:"iid"`
IName string `json:"iname"`
Desc string `json:"desc"`
Orig string `json:"orig"`
Status string `json:"status"`
}


func Create(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    fmt.Println("Entering Create function")

    if len(args) != 5 {
        fmt.Println("Invalid number of args")
        return nil, errors.New("Expected at least two arguments for loan application creation")
    }

    var Id = args[0]
    var In = args[1]
    var De = args[2]
    var Or = args[3]
    var St = args[4]

    var incidentInfo Incident
    incidentInfo = Incident{Id, In, De, Or, St}
     piBytes, err := json.Marshal(&incidentInfo)
     if err != nil {
       fmt.Println("Error in marshaling ",err)
       return nil, err
     }


    err = stub.PutState(Id, piBytes)
    if err != nil {
        fmt.Println("Could not save changes", err)
        return nil, err
    }

    fmt.Println("Successfully saved changes")
    return nil, nil
}

func Get(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    fmt.Println("Entering Get function")

    if len(args) < 1 {
        fmt.Println("Invalid number of arguments")
        return nil, errors.New("Missing ID")
    }

    var Id = args[0]
    piBytes, err := stub.GetState(Id)
    if err != nil {
        fmt.Println("Could not fetch data with id "+Id+" from ledger", err)
        return nil, err
    }
    piBytes2 := piBytes
    var incidentInfo Incident
    err = json.Unmarshal(piBytes, &incidentInfo)
    if err != nil {
      fmt.Println("Error in unmarshaling",err)
      return nil, err
    }
    fmt.Println(incidentInfo.IName)
    fmt.Println(incidentInfo.Desc)
    fmt.Println(incidentInfo.Orig)
    fmt.Println(incidentInfo.Status)

    piBytes2, err = json.Marshal(&incidentInfo)
    if err != nil {
      fmt.Println("Error in unmarshaling",err)
      return nil, err
    }
/*
    Bytes, err := stub.GetState(Id)
    if err != nil {
        fmt.Println("Not found id "+Id+" from ledger", err)
        return nil, err
    }
*/
    return piBytes, nil
}

func Update(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering Update")

	if len(args) < 2 {
		fmt.Println("Invalid number of args")
		return nil, errors.New("Expected atleast two arguments for update")
	}

	var Id = args[0]
	var status = args[1]

	laBytes, err := stub.GetState(Id)
	if err != nil {
		fmt.Println("Could not fetch data from ledger", err)
		return nil, err
	}

	var incidentInfo Incident
	err = json.Unmarshal(laBytes, &incidentInfo)
	incidentInfo.Status = status

	laBytes, err = json.Marshal(&incidentInfo)
	if err != nil {
		fmt.Println("Could not marshal ", err)
		return nil, err
	}

	err = stub.PutState(Id, laBytes)
	if err != nil {
		fmt.Println("Could not save update", err)
		return nil, err
	}

	fmt.Println("Successfully updated changes")
	return nil, nil

}


func (t *SampleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

   /* if len(args) != 4 {
     return nil, errors.New("Incorrect number of arguments. Expecting 4")
   }
   */

     var incidentInfo Incident
      incidentInfo = Incident{"1", "One", "d1", "o1", "s1"}
      bytes, err := json.Marshal (&incidentInfo)
      if err != nil {
             fmt.Println("Could not marshal incident info object", err)
             return nil, err
      }

   /*bytes, err := json.Marshal (&args[0])
   if err != nil {
          fmt.Println("Could not marshal incident info object", err)
          return nil, err
   }
  */

   err = stub.PutState("1", bytes)
   if err != nil {
     fmt.Println("Could not save ", err)
     return nil, err
   }

    return nil, nil
}

func (t *SampleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
fmt.Println("Inside Query function")
  if function == "get" {
    return Get(stub, args)
  }
  fmt.Println("Query could not find: " + function)
    return nil, nil
}

func (t *SampleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
/*  if function == "init" {
		return Init(stub, "init", args)
  }
  */
  if function == "create" {
    return Create(stub, args)
  }
  if function=="update" {
    return Update(stub,args)
  }
  fmt.Println("Invoke did not find func: " + function)
    return nil, nil
}



func main() {
    err := shim.Start(new(SampleChaincode))
    if err != nil {
        fmt.Println("Could not start SampleChaincode")
    } else {
        fmt.Println("SampleChaincode successfully started")
    }

}
