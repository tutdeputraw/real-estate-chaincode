package chaincode

import (
	"encoding/json"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
	"tutdeputraw.com/common/models"
)

func (s *RealEstateChaincode) SalesAndPurchaseAgreement_Create(APIstub shim.ChaincodeStubInterface, agreementJSON string) sc.Response {
	var agreement models.Agreement
	err := json.Unmarshal([]byte(agreementJSON), &agreement)
	if err != nil {
		return shim.Error("failed to unmarshal agreement JSON:" + err.Error())
	}

	existing, err := APIstub.GetState(agreement.ID)
	if err != nil {
		return shim.Error("failed to read from world state:" + err.Error())
	}
	if existing != nil {
		return shim.Error("agreement with ID " + agreement.ID + " already exists")
	}

	agreementBytes, err := json.Marshal(agreement)
	if err != nil {
		return shim.Error("failed to marshal agreement:" + err.Error())
	}

	err = APIstub.PutState(agreement.ID, agreementBytes)
	if err != nil {
		return shim.Error("failed to put agreement state:" + err.Error())
	}

	return shim.Success(nil)
}

func (s *RealEstateChaincode) SalesAndPurchaseAgreement_Update(APIstub shim.ChaincodeStubInterface, agreementJSON string) sc.Response {
	var agreement models.Agreement
	err := json.Unmarshal([]byte(agreementJSON), &agreement)
	if err != nil {
		return shim.Error("failed to unmarshal agreement JSON:" + err.Error())
	}

	existing, err := APIstub.GetState(agreement.ID)
	if err != nil {
		return shim.Error("failed to read from world state:" + err.Error())
	}
	if existing == nil {
		return shim.Error("agreement with ID " + agreement.ID + " does not exist")
	}

	agreementBytes, err := json.Marshal(agreement)
	if err != nil {
		return shim.Error("failed to marshal agreement:" + err.Error())
	}

	err = APIstub.PutState(agreement.ID, agreementBytes)
	if err != nil {
		return shim.Error("failed to put agreement state:" + err.Error())
	}

	return shim.Success(nil)
}

func (s *RealEstateChaincode) SalesAndPurchaseAgreement_Delete(APIstub shim.ChaincodeStubInterface, agreementID string) sc.Response {
	existing, err := APIstub.GetState(agreementID)
	if err != nil {
		return shim.Error("failed to read from world state:" + err.Error())
	}
	if existing == nil {
		return shim.Error("agreement with ID " + agreementID + " does not exist")
	}

	err = APIstub.DelState(agreementID)
	if err != nil {
		return shim.Error("failed to delete agreement state:" + err.Error())
	}

	return shim.Success(nil)
}
