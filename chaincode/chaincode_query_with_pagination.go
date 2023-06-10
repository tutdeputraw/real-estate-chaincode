package chaincode

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
)

func (s *RealEstateChaincode) QueryWithPagination(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	queryString := args[0]
	//return type of ParseInt is int64
	pageSize, err := strconv.ParseInt(args[1], 10, 32)
	if err != nil {
		return shim.Error(err.Error())
	}
	bookmark := args[2]

	queryResults, err := getQueryResultForQueryStringWithPagination(stub, queryString, int32(pageSize), bookmark)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func getQueryResultForQueryStringWithPagination(stub shim.ChaincodeStubInterface, queryString string, pageSize int32, bookmark string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, responseMetadata, err := stub.GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIteratorQueryWithPagination(resultsIterator)
	if err != nil {
		return nil, err
	}

	bufferWithPaginationInfo := addPaginationMetadataToQueryResults(buffer, responseMetadata)

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", bufferWithPaginationInfo.String())

	return buffer.Bytes(), nil
}

// func (s *RealEstateChaincode) PokokString() sc.Response {
// 	result := bytes.NewBufferString("")
// 	response := bytes.NewBufferString("[{}, {}, {}]")
// 	var responseMetaCountRecord int32 = 10
// 	responseMetaBookmark := "selanjute"

// 	result.WriteString(fmt.Sprintf(`{
// 		"response": %v,
// 		"metadata": {
// 			"RecordsCount": %v,
// 			"Bookmark": %v
// 		}
// 	}`, response, responseMetaCountRecord, responseMetaBookmark))

// 	// str := "Hello, World!"
// 	// fmt.Fprintf(result, "%s", str)
// 	fmt.Println("XXY:", result.String())
// 	return shim.Success(nil)
// }

func addPaginationMetadataToQueryResults(buffer *bytes.Buffer, responseMetadata *sc.QueryResponseMetadata) *bytes.Buffer {
	// buffer.WriteString("{\"ResponseMetadata\":{\"RecordsCount\":")
	// buffer.WriteString("\"")
	// buffer.WriteString(fmt.Sprintf("%v", responseMetadata.FetchedRecordsCount))
	// buffer.WriteString("\"")
	// buffer.WriteString(", \"Bookmark\":")
	// buffer.WriteString("\"")
	// buffer.WriteString(responseMetadata.Bookmark)
	// buffer.WriteString("\"}}")
	buffer.WriteString("\"ResponseMetadata\":{\"RecordsCount\":")
	buffer.WriteString("\"")
	buffer.WriteString(fmt.Sprintf("%v", responseMetadata.FetchedRecordsCount))
	buffer.WriteString("\"")
	buffer.WriteString(", \"Bookmark\":")
	buffer.WriteString("\"")
	buffer.WriteString(responseMetadata.Bookmark)
	buffer.WriteString("\"}}")

	// result := bytes.NewBufferString("")

	// result.WriteString(fmt.Sprintf(`{
	// 	"response": %v,
	// 	"metadata": {
	// 		"RecordsCount": %v,
	// 		"Bookmark": %v
	// 	}
	// }`, buffer, responseMetadata.FetchedRecordsCount, responseMetadata.Bookmark))

	return buffer
}

func constructQueryResponseFromIteratorQueryWithPagination(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("{")
	buffer.WriteString("\"data\":")
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
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
	buffer.WriteString("],")

	return &buffer, nil
}
