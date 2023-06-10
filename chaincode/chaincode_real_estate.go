package chaincode

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"

	constant "tutdeputraw.com/common/constants"
	mock "tutdeputraw.com/common/mocks"
	"tutdeputraw.com/common/models"
)

//==========[CORE]==========//

func (s *RealEstateChaincode) RealEstate_Init_Test(APIstub shim.ChaincodeStubInterface) sc.Response {
	// realEstates := mock.Mock_RealEstates
	realEstates := mock.Mock_RealEstates_TransactionHistory

	i := 0
	realEstatesLen := len(realEstates)
	for i < realEstatesLen {
		realEstateAsBytes, _ := json.Marshal(realEstates[i])

		realEstate := models.RealEstateModel{}
		json.Unmarshal(realEstateAsBytes, &realEstate)

		APIstub.PutState(constant.State_RealEstate+realEstates[i].RealEstateId, realEstateAsBytes)

		// create real estate history
		s.RealEstateHistory_Create(APIstub, []string{
			realEstate.RealEstateId, // id
			realEstate.OwnerId,
			realEstate.RealEstateId,
			time.Now().Local().String(),
		})

		// write the user to the real estate ownership history
		createCompositeRealEstatesByOwner := s.RealEstate_CreateCompositeRealEstatesByOwner(APIstub, realEstate)
		createCompositeRealEstatesByOwnerIsSuccess, _ := strconv.ParseBool(string(createCompositeRealEstatesByOwner.Payload))
		if !createCompositeRealEstatesByOwnerIsSuccess {
			return shim.Error("failed to create composite real estates by owner")
		}

		createCompositeOwnersByRealEstate := s.RealEstate_CreateCompositeOwnersByRealEstate(APIstub, realEstate)
		createCompositeOwnersByRealEstateIsSuccess, _ := strconv.ParseBool(string(createCompositeOwnersByRealEstate.Payload))
		if !createCompositeOwnersByRealEstateIsSuccess {
			return shim.Error("failed to create composite owners by real estate")
		}

		i = i + 1
	}

	return shim.Success(nil)
}

func (s *RealEstateChaincode) RealEstate_Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	realEstates := s.GetRealEstatesUniqueDataset()

	i := 0
	realEstatesLen := len(realEstates)
	for i < realEstatesLen {
		realEstateAsBytes, _ := json.Marshal(realEstates[i])

		realEstate := models.RealEstateModel{}
		json.Unmarshal(realEstateAsBytes, &realEstate)

		APIstub.PutState(constant.State_RealEstate+realEstates[i].RealEstateId, realEstateAsBytes)

		// create real estate history
		s.RealEstateHistory_Create(APIstub, []string{
			realEstate.RealEstateId, // id
			realEstate.OwnerId,
			realEstate.RealEstateId,
			time.Now().Local().String(),
		})

		// write the user to the real estate ownership history
		createCompositeRealEstatesByOwner := s.RealEstate_CreateCompositeRealEstatesByOwner(APIstub, realEstate)
		createCompositeRealEstatesByOwnerIsSuccess, _ := strconv.ParseBool(string(createCompositeRealEstatesByOwner.Payload))
		if !createCompositeRealEstatesByOwnerIsSuccess {
			return shim.Error("failed to create composite real estates by owner")
		}

		createCompositeOwnersByRealEstate := s.RealEstate_CreateCompositeOwnersByRealEstate(APIstub, realEstate)
		createCompositeOwnersByRealEstateIsSuccess, _ := strconv.ParseBool(string(createCompositeOwnersByRealEstate.Payload))
		if !createCompositeOwnersByRealEstateIsSuccess {
			return shim.Error("failed to create composite owners by real estate")
		}

		// // write the user to the real estate ownership history
		// compositeKey := constant.Composite_GetRealEstatesByOwnerKey
		// ownerId := constant.State_User + realEstates[i].OwnerId
		// realEstateId := constant.State_RealEstate + realEstates[i].RealEstateId
		// ownerRealEstateIdIndexKey, err := APIstub.CreateCompositeKey(
		// 	compositeKey,
		// 	[]string{ownerId, realEstateId},
		// )
		// if err != nil {
		// 	return shim.Error(err.Error())
		// }
		// value := []byte{0x00}
		// APIstub.PutState(ownerRealEstateIdIndexKey, value)

		i = i + 1
	}

	return shim.Success(nil)
}

func (s *RealEstateChaincode) RealEstate_CheckIfRealEstateHasAlreadyRegistered(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	bytes, _ := APIstub.GetState(constant.State_RealEstate + args[0])
	if bytes != nil {
		return shim.Success([]byte(strconv.FormatBool(true)))
	}

	return shim.Success([]byte(strconv.FormatBool(false)))
}

func (s *RealEstateChaincode) RealEstate_Create(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 13 {
		return shim.Error("Incorrect number of arguments. Expecting 13")
	}

	var realEstate = models.RealEstateModel{
		RealEstateId: args[0],
		OwnerId:      args[1],
		Price:        args[2],
		Bed:          args[3],
		Bath:         args[4],
		AcreLot:      args[5],
		FullAddress:  args[6],
		Street:       args[7],
		City:         args[8],
		State:        args[9],
		ZipCode:      args[10],
		HouseSize:    args[11],
		IsOpenToSell: args[12],
	}

	realEstateAsBytes, _ := json.Marshal(realEstate)
	APIstub.PutState(constant.State_RealEstate+args[0], realEstateAsBytes)

	return shim.Success(realEstateAsBytes)
}

func (s *RealEstateChaincode) RealEstate_RegisterNewRealEstate(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 13 {
		return shim.Error("Incorrect number of arguments. Expecting 13")
	}

	var realEstate = models.RealEstateModel{
		RealEstateId: args[0],
		OwnerId:      args[1],
		Price:        args[2],
		Bed:          args[3],
		Bath:         args[4],
		AcreLot:      args[5],
		FullAddress:  args[6],
		Street:       args[7],
		City:         args[8],
		State:        args[9],
		ZipCode:      args[10],
		HouseSize:    args[11],
	}

	// return false if a user is not registered first
	responseIsUserRegistered := s.User_CheckIfUserExist(APIstub, []string{realEstate.OwnerId})
	userIsRegistered, _ := strconv.ParseBool(string(responseIsUserRegistered.Payload))
	userIsNotRegistered := !userIsRegistered
	if userIsNotRegistered {
		return shim.Error("user is not registered")
	}

	// return false if the real estate is already registered
	responseCheckIfRealEstateHasAlreadyRegistered := s.RealEstate_CheckIfRealEstateHasAlreadyRegistered(APIstub, args)
	realEstateIsRegistered, _ := strconv.ParseBool(string(responseCheckIfRealEstateHasAlreadyRegistered.Payload))
	if !realEstateIsRegistered {
		// create real estate
		s.RealEstate_Create(APIstub, args)
	}

	// create real estate history
	s.RealEstateHistory_Create(APIstub, []string{
		realEstate.RealEstateId, // id
		realEstate.OwnerId,
		realEstate.RealEstateId,
		time.Now().Local().String(),
	})

	// write the user to the real estate ownership history
	createCompositeRealEstatesByOwner := s.RealEstate_CreateCompositeRealEstatesByOwner(APIstub, realEstate)
	createCompositeRealEstatesByOwnerIsSuccess, _ := strconv.ParseBool(string(createCompositeRealEstatesByOwner.Payload))
	if !createCompositeRealEstatesByOwnerIsSuccess {
		return shim.Error("failed to create composite real estates by owner")
	}

	createCompositeOwnersByRealEstate := s.RealEstate_CreateCompositeOwnersByRealEstate(APIstub, realEstate)
	createCompositeOwnersByRealEstateIsSuccess, _ := strconv.ParseBool(string(createCompositeOwnersByRealEstate.Payload))
	if !createCompositeOwnersByRealEstateIsSuccess {
		return shim.Error("failed to create composite owners by real estate")
	}

	bytesRealEstate, _ := json.Marshal(realEstate)
	return shim.Success(bytesRealEstate)
}

func (s *RealEstateChaincode) RealEstate_GetById(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	realEstateKey := constant.State_RealEstate + args[0]

	realEstateAsBytes, _ := APIstub.GetState(realEstateKey)

	return shim.Success(realEstateAsBytes)
}

func (s *RealEstateChaincode) RealEstate_GetAllWithPagination(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	startKey := constant.State_RealEstate + "0"
	endKey := constant.State_RealEstate + "999"
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

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return shim.Error("failed to construct query response from iterator")
	}

	bufferWithPaginationInfo := addPaginationMetadataToQueryResults(buffer, responseMetadata)
	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", bufferWithPaginationInfo.String())

	return shim.Success(buffer.Bytes())
}

func (s *RealEstateChaincode) RealEstate_GetAll(APIstub shim.ChaincodeStubInterface) sc.Response {
	startKey := constant.State_RealEstate + "0"
	endKey := constant.State_RealEstate + "999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var realEstates []models.RealEstateModel
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		var realEstate models.RealEstateModel
		err = json.Unmarshal(queryResponse.Value, &realEstate)
		if err != nil {
			return shim.Error(fmt.Sprintf("Failed to unmarshal state: %s", err.Error()))
		}

		// fmt.Println("realEstate: ", realEstate)

		realEstates = append(realEstates, realEstate)
	}

	realEstatesBytes, err := json.Marshal(realEstates)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to marshal states: %s", err.Error()))
	}

	return shim.Success(realEstatesBytes)
}

func (s *RealEstateChaincode) NYOBAK(APIstub shim.ChaincodeStubInterface) sc.Response {
	// result, _ := APIstub.GetCreator()
	result, _ := APIstub.GetSignedProposal()
	println("KIKI: ", string(result.ProposalBytes))
	println("KIKI: ", string(result.Signature))

	APIstub.GetSignedProposal()
	return shim.Success(result.ProposalBytes)
}

func (s *RealEstateChaincode) RealEstate_GetByOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments")
	}
	owner := args[0]

	ownerAndIdResultIterator, err := APIstub.GetStateByPartialCompositeKey(
		constant.Composite_GetRealEstatesByOwnerKey,
		[]string{constant.State_User + owner},
	)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer ownerAndIdResultIterator.Close()

	var i int
	var id string

	var realestates []byte
	bArrayMemberAlreadyWritten := false

	realestates = append([]byte("["))

	for i = 0; ownerAndIdResultIterator.HasNext(); i++ {
		responseRange, err := ownerAndIdResultIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		objectType, compositeKeyParts, err := APIstub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(err.Error())
		}

		id = compositeKeyParts[1]
		assetAsBytes, err := APIstub.GetState(id)

		fmt.Println("INFOOOUSER")
		fmt.Println("compositeKeyParts\t: ", compositeKeyParts)
		fmt.Println("id\t\t\t: ", id)
		fmt.Println("assetAsBytes\t\t: \n", string(assetAsBytes))

		if bArrayMemberAlreadyWritten == true {
			newBytes := append([]byte(","), assetAsBytes...)
			realestates = append(realestates, newBytes...)

		} else {
			// newBytes := append([]byte(","), carsAsBytes...)
			realestates = append(realestates, assetAsBytes...)
		}

		fmt.Println("Found a asset\t: ", objectType)
		fmt.Println("for index\t: ", compositeKeyParts[0])
		fmt.Println("asset id\t: ", compositeKeyParts[1])
		bArrayMemberAlreadyWritten = true
	}

	realestates = append(realestates, []byte("]")...)

	return shim.Success(realestates)
}

func (s *RealEstateChaincode) RealEstate_CreateCompositeRealEstatesByOwner(APIstub shim.ChaincodeStubInterface, realEstate models.RealEstateModel) sc.Response {
	compositeKey := constant.Composite_GetRealEstatesByOwnerKey
	ownerId := constant.State_User + realEstate.OwnerId
	realEstateId := constant.State_RealEstate + realEstate.RealEstateId

	ownerRealEstateIdIndexKey, err := APIstub.CreateCompositeKey(
		compositeKey,
		[]string{ownerId, realEstateId},
	)
	if err != nil {
		return shim.Success([]byte(strconv.FormatBool(false)))
	}

	value := []byte{0x00}
	fmt.Println("JJJA1: ", ownerRealEstateIdIndexKey)
	APIstub.PutState(ownerRealEstateIdIndexKey, value)

	return shim.Success([]byte(strconv.FormatBool(true)))
}

func (s *RealEstateChaincode) RealEstate_CreateCompositeOwnersByRealEstate(APIstub shim.ChaincodeStubInterface, realEstate models.RealEstateModel) sc.Response {
	compositeKey := constant.Composite_GetOwnersByRealEstateKey
	ownerId := constant.State_User + realEstate.OwnerId
	realEstateId := constant.State_RealEstate + realEstate.RealEstateId
	realEstateHistoryKey := constant.State_RealEstateHistory + realEstate.RealEstateId + realEstate.OwnerId

	ownerRealEstateIdIndexKey, err := APIstub.CreateCompositeKey(
		compositeKey,
		[]string{realEstateId, ownerId, realEstateHistoryKey},
	)
	if err != nil {
		return shim.Success([]byte(strconv.FormatBool(false)))
	}

	value := []byte{0x00}
	fmt.Println("JJJA2: ", ownerRealEstateIdIndexKey)
	APIstub.PutState(ownerRealEstateIdIndexKey, value)

	return shim.Success([]byte(strconv.FormatBool(true)))
}

func (s *RealEstateChaincode) RealEstate_BuyRealEstateOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	// create spa
	// change real estate owner
	// s.SalesAndPurchaseAgreement_Create(APIstub,[]str)
	s.RealEstate_ChangeRealEstateOwner(APIstub, []string{})
	return shim.Success(nil)
}

func (s *RealEstateChaincode) RealEstate_ChangeRealEstateOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	realEstateKey := constant.State_RealEstate + args[0]
	newOwnerId := args[1]

	realEstateAsBytes, _ := APIstub.GetState(realEstateKey)
	realEstate := models.RealEstateModel{}
	json.Unmarshal(realEstateAsBytes, &realEstate)

	print("realEstate.IsOpenToSell: ", realEstate.IsOpenToSell, "\t", realEstate.RealEstateId)

	if realEstate.IsOpenToSell == "false" {
		return shim.Error("= Real estate is not open to sell")
	}

	//==========[delete old key]==========//
	resulta, _ := APIstub.GetState("realestate~owneriduser-2real-estate-3")
	fmt.Println("PIPIR: " + string(resulta))

	compositeKey := constant.Composite_GetRealEstatesByOwnerKey
	ownerId := constant.State_User + realEstate.OwnerId
	realEstateId := constant.State_RealEstate + realEstate.RealEstateId

	ownerCaridIndexKey, err := APIstub.CreateCompositeKey(
		compositeKey,
		[]string{ownerId, realEstateId},
	)
	fmt.Println("CAKCIKCUK: ", ownerCaridIndexKey)
	length := len(ownerCaridIndexKey)
	fmt.Println("joss " + strconv.Itoa(length))
	fmt.Println("JASJUS: " + strconv.FormatBool(strings.Contains(ownerCaridIndexKey, "")))
	if err != nil {
		return shim.Error(err.Error())
	}

	state, _ := APIstub.GetState(ownerCaridIndexKey)
	fmt.Println("SAKJANESS: " + string(state) + "\t" + strconv.FormatBool(strings.Contains(string(state), "")))
	if state == nil {
		fmt.Println("GAONOSS")
		return shim.Success([]byte(strconv.FormatBool(true)))
	}

	err = APIstub.DelState(ownerCaridIndexKey)
	if err != nil {
		return shim.Error("Failed to delete state:" + string(err.Error()))
	}
	//----------[delete old key]----------//

	//==========[update realestate state]==========//
	realEstate.OwnerId = newOwnerId
	realEstate.IsOpenToSell = "false"

	realEstateAsBytes, _ = json.Marshal(realEstate)

	APIstub.PutState(realEstateKey, realEstateAsBytes)
	//----------[update realestate state]----------//

	//==========[create real estate history]==========//
	s.RealEstateHistory_Create(APIstub, []string{
		realEstate.RealEstateId, // id
		realEstate.OwnerId,
		realEstate.RealEstateId,
		time.Now().Local().String(),
	})
	//----------[create real estate history]----------//

	//==========[add new user into the composite of real estate history]==========//
	createCompositeRealEstatesByOwner := s.RealEstate_CreateCompositeRealEstatesByOwner(APIstub, realEstate)
	createCompositeRealEstatesByOwnerIsSuccess, _ := strconv.ParseBool(string(createCompositeRealEstatesByOwner.Payload))
	if !createCompositeRealEstatesByOwnerIsSuccess {
		return shim.Error("failed to create composite real estates by owner")
	}

	createCompositeOwnersByRealEstate := s.RealEstate_CreateCompositeOwnersByRealEstate(APIstub, realEstate)
	createCompositeOwnersByRealEstateIsSuccess, _ := strconv.ParseBool(string(createCompositeOwnersByRealEstate.Payload))
	if !createCompositeOwnersByRealEstateIsSuccess {
		return shim.Error("failed to create composite owners by real estate")
	}
	//----------[add new user into the composite of real estate history]----------//

	return shim.Success(nil)
}

func (s *RealEstateChaincode) RealEstate_ChangeRealEstateSellStatus(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	realEstateId := args[0]
	status := args[1]

	if status != "true" && status != "false" {
		return shim.Error("args[1] value is not in range true or false, value type must be string")
	}

	realEstateAsBytes, err := APIstub.GetState(constant.State_RealEstate + realEstateId)
	if err != nil {
		return shim.Error(err.Error())
	}

	realEstate := models.RealEstateModel{}
	json.Unmarshal(realEstateAsBytes, &realEstate)

	realEstate.IsOpenToSell = status

	realEstateAsBytes, _ = json.Marshal(realEstate)

	APIstub.PutState(constant.State_RealEstate+realEstateId, realEstateAsBytes)

	//==========[create sale history]==========//
	if status == "true" {
		s.RealEstateSalesRecord_Create(APIstub, []string{
			realEstate.RealEstateId,
			realEstate.OwnerId,
		})
		createCompositeRealEstateRecordByRealEstate := s.RealEstateSalesRecord_CreateComposite(APIstub, realEstate)
		createCompositeRealEstateRecordByRealEstateIsSuccess, _ := strconv.ParseBool(string(createCompositeRealEstateRecordByRealEstate.Payload))
		if !createCompositeRealEstateRecordByRealEstateIsSuccess {
			return shim.Error("failed to create composite real estate sale record by real estate")
		}

		// s.SalesAndPurchaseAgreement_Create(APIstub,)
	} else {
		fmt.Println("MELBUO KENE")
		s.RealEstateSalesRecord_Delete(APIstub, []string{
			realEstate.RealEstateId + realEstate.OwnerId,
		})
		deleteCompositeKey := s.RealEstateSalesRecord_DeleteComposite(APIstub, realEstate)
		deleteCompositeKeyIsSuccess, _ := strconv.ParseBool(string(deleteCompositeKey.Payload))
		if !deleteCompositeKeyIsSuccess {
			return shim.Error("failed to delete the composite key real estate sales record by real estate")
		}
	}
	//----------[]----------//

	return shim.Success(realEstateAsBytes)
}

func (s *RealEstateChaincode) GetRealEstatesUniqueDataset() []models.RealEstateModel {
	fileContent, err := ioutil.ReadFile("../common/dataset/lite/real-estates-unique.json")
	if err != nil {
		log.Fatalf("Error reading JSON file: %s", err)
	}

	var realEstates []models.RealEstateModel

	err = json.Unmarshal(fileContent, &realEstates)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %s", err)
	}

	// for _, realEstate := range realEstates {
	// 	fmt.Println("AAIU RealEstateId: %s, OwnerId: %s", realEstate.RealEstateId, realEstate.OwnerId)
	// }

	return realEstates
}

func (s *RealEstateChaincode) GetRealEstatesDataset() []models.RealEstateModel {
	fileContent, err := ioutil.ReadFile("../common/dataset/lite/real-estates.json")
	if err != nil {
		log.Fatalf("Error reading JSON file: %s", err)
	}

	var realEstates []models.RealEstateModel

	err = json.Unmarshal(fileContent, &realEstates)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %s", err)
	}

	// for _, realEstate := range realEstates {
	// 	fmt.Println("AAIU RealEstateId: %s, OwnerId: %s", realEstate.RealEstateId, realEstate.OwnerId)
	// }

	return realEstates
}

func (s *RealEstateChaincode) GetRealEstatesDatasetByFilePath(filePath string) []models.RealEstateModel {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading JSON file: %s", err)
	}

	var realEstates []models.RealEstateModel

	err = json.Unmarshal(fileContent, &realEstates)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %s", err)
	}

	// for _, realEstate := range realEstates {
	// 	fmt.Println("AAIU RealEstateId: %s, OwnerId: %s", realEstate.RealEstateId, realEstate.OwnerId)
	// }

	return realEstates
}
