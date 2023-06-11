package chaincode

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
	constant "tutdeputraw.com/common/constants"
	"tutdeputraw.com/common/models"
)

func (s *RealEstateChaincode) RealEstateHistory_Create(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	var realEstateHistory = models.RealEstateHistoryModel{
		// Id:           args[0],
		OwnerID:      args[1],
		RealEstateId: args[2],
		DateTime:     args[3],
	}

	realEstateHistoryAsBytes, _ := json.Marshal(realEstateHistory)
	realEstateHistoryKey := constant.State_RealEstateHistory + realEstateHistory.RealEstateId + realEstateHistory.OwnerID
	// realEstateHistoryKey := constant.RealEstateHistoryKey + realEstateHistory.RealEstateId
	APIstub.PutState(
		realEstateHistoryKey,
		realEstateHistoryAsBytes,
	)

	return shim.Success(realEstateHistoryAsBytes)
}

func (s *RealEstateChaincode) RealEstateHistory_GetByRealEstateId(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments")
	}

	realEstateAsBytes, err := APIstub.GetState(constant.State_RealEstate + args[0])
	realEstate := models.RealEstateModel{}
	json.Unmarshal(realEstateAsBytes, &realEstate)
	realEstateId := constant.State_RealEstate + realEstate.RealEstateId

	ownerAndIdResultIterator, err := APIstub.GetStateByPartialCompositeKey(
		constant.Composite_GetOwnersByRealEstateKey,
		[]string{realEstateId},
	)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer ownerAndIdResultIterator.Close()

	var i int
	var id string

	var users []byte
	bArrayMemberAlreadyWritten := false

	users = append([]byte("["))

	for i = 0; ownerAndIdResultIterator.HasNext(); i++ {
		responseRange, err := ownerAndIdResultIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		// objectType, compositeKeyParts, err := APIstub.SplitCompositeKey(responseRange.Key)
		_, compositeKeyParts, err := APIstub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(err.Error())
		}

		id = compositeKeyParts[2]
		assetAsBytes, err := APIstub.GetState(id)

		fmt.Println("//==========//")
		fmt.Println("INFOOO")
		fmt.Println("compositeKeyParts\t: ", compositeKeyParts)
		fmt.Println("id\t\t\t: ", id)
		fmt.Println("assetAsBytes\t\t: ", string(assetAsBytes))

		if bArrayMemberAlreadyWritten == true {
			newBytes := append([]byte(","), assetAsBytes...)
			users = append(users, newBytes...)

		} else {
			// newBytes := append([]byte(","), carsAsBytes...)
			users = append(users, assetAsBytes...)
		}

		// fmt.Println("Found a asset\t: ", objectType)
		// fmt.Println("for index\t: ", compositeKeyParts[0])
		// fmt.Println("asset id\t: ", compositeKeyParts[1])

		bArrayMemberAlreadyWritten = true
	}

	users = append(users, []byte("]")...)

	fmt.Println("//==========//")
	fmt.Println("RESULT\t\t: ", string(users))

	return shim.Success(users)
}

func (s *RealEstateChaincode) GetRealEstatesHistoriesDatasetByFilePath(filePath string) []models.RealEstateHistoryModel {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading JSON file: %s", err)
	}

	var realEstates []models.RealEstateHistoryModel

	err = json.Unmarshal(fileContent, &realEstates)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %s", err)
	}

	// for _, realEstate := range realEstates {
	// 	fmt.Println("AAIU RealEstateId: %s, OwnerId: %s", realEstate.RealEstateId, realEstate.OwnerId)
	// }

	return realEstates
}
