package chaincode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
	constant "tutdeputraw.com/common/constants"
	mock "tutdeputraw.com/common/mocks"
	"tutdeputraw.com/common/models"
)

func (s *RealEstateChaincode) GetUsersDataset() []models.UserModel {
	fileContent, err := ioutil.ReadFile("../common/dataset/lite/users.json")
	if err != nil {
		log.Fatalf("Error reading JSON file: %s", err)
	}

	var users []models.UserModel

	err = json.Unmarshal(fileContent, &users)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %s", err)
	}

	// for _, user := range users {
	// 	fmt.Println("AAIU Name: %s, Id: %s, Email: %s", user.Name, user.Id, user.Email)
	// }

	return users
}

func (s *RealEstateChaincode) User_InitTest(APIstub shim.ChaincodeStubInterface) sc.Response {
	users := mock.Mock_Users

	i := 0
	usersLen := len(users)
	for i < usersLen {
		realEstateAsBytes, _ := json.Marshal(users[i])
		APIstub.PutState(constant.State_User+users[i].Id, realEstateAsBytes)
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *RealEstateChaincode) User_Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	users := s.GetUsersDataset()

	i := 0
	usersLen := len(users)
	for i < usersLen {
		realEstateAsBytes, _ := json.Marshal(users[i])
		APIstub.PutState(constant.State_User+users[i].Id, realEstateAsBytes)
		i = i + 1
		fmt.Println("AAI1:", string(realEstateAsBytes))
	}

	return shim.Success(nil)
}

func (s *RealEstateChaincode) User_CheckIfUserExist(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	fmt.Println("IDID: ", constant.State_User+args[0])

	bytes, _ := APIstub.GetState(constant.State_User + args[0])
	if bytes != nil {
		return shim.Success([]byte(strconv.FormatBool(true)))
	}

	return shim.Success([]byte(strconv.FormatBool(false)))
}

func (s *RealEstateChaincode) User_Create(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	// check if user has not registered yet
	responseIsUserRegistered := s.User_CheckIfUserExist(APIstub, []string{args[0]})
	userIsRegistered, _ := strconv.ParseBool(string(responseIsUserRegistered.Payload))
	if userIsRegistered {
		return shim.Error("user has already registered")
	}

	bytes, _ := APIstub.GetState(args[0])
	if bytes != nil {
		return shim.Error("user is already registered")
	}

	var user = models.UserModel{
		Id:          args[0],
		Name:        args[1],
		NPWPNumber:  args[2],
		PhoneNumber: args[3],
		Email:       args[4],
	}

	userAsBytes, _ := json.Marshal(user)
	APIstub.PutState(constant.State_User+args[0], userAsBytes)

	return shim.Success(userAsBytes)
}

func (s *RealEstateChaincode) User_GetById(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	realEstateAsBytes, _ := APIstub.GetState(constant.State_User + args[0])
	// fmt.Printf("- queryRealEstateById:\n%s\n", string(realEstateAsBytes))
	return shim.Success(realEstateAsBytes)
}

func (s *RealEstateChaincode) User_GetAllWithPagination(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	startKey := constant.State_User + "0"
	endKey := constant.State_User + "999"
	pageSize, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		return shim.Error(err.Error())
	}
	bookmark := args[1]

	resultsIterator, responseMetadata, err := APIstub.GetStateByRangeWithPagination(startKey, endKey, int32(pageSize), bookmark)
	if err != nil {
		return shim.Error("failed to get state by range with pagination")
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIteratorQueryWithPagination(resultsIterator)
	if err != nil {
		return shim.Error("failed to construct query response from iterator")
	}

	bufferWithPaginationInfo := addPaginationMetadataToQueryResults(buffer, responseMetadata)
	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", bufferWithPaginationInfo.String())

	return shim.Success(buffer.Bytes())
}

func (s *RealEstateChaincode) User_GetAll(APIstub shim.ChaincodeStubInterface) sc.Response {
	startKey := constant.State_User + "0"
	endKey := constant.State_User + "999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	// fmt.Printf("- queryAllUsers:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
