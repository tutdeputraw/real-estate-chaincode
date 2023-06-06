package models

import (
	"time"
)

type Party struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
}

type Property struct {
	Address     string    `json:"address"`
	Description string    `json:"description"`
	Price       string    `json:"price"`
	ClosingDate time.Time `json:"closing_date"`
}

type Terms struct {
	DepositAmount              string `json:"deposit_amount"`
	InspectionPeriod           string `json:"inspection_period"`
	FinancingRequired          string `json:"financing_required"`
	FinancingContingencyPeriod string `json:"financing_contingency_period"`
	TitleCompany               string `json:"title_company"`
	EscrowAgent                string `json:"escrow_agent"`
}

type Agreement struct {
	ID              string           `json:"id"`
	Parties         map[string]Party `json:"parties"`
	Property        Property         `json:"property"`
	Terms           Terms            `json:"terms"`
	AdditionalTerms []string         `json:"additional_terms"`
	Signatures      struct {
		SellerSignature string    `json:"seller_signature"`
		BuyerSignature  string    `json:"buyer_signature"`
		Date            time.Time `json:"date"`
	} `json:"signatures"`
}

// func main() {
// 	agreement := Agreement{
// 		Parties: map[string]Party{
// 			"seller": {
// 				Name:    "John Doe",
// 				Address: "123 Main Street, City, State",
// 				Email:   "john.doe@example.com",
// 				Phone:   "555-1234",
// 			},
// 			"buyer": {
// 				Name:    "Jane Smith",
// 				Address: "456 Elm Street, City, State",
// 				Email:   "jane.smith@example.com",
// 				Phone:   "555-5678",
// 			},
// 		},
// 		Property: Property{
// 			Address:     "789 Oak Avenue, City, State",
// 			Description: "A beautiful single-family home with three bedrooms and two bathrooms.",
// 			Price:       250000,
// 			ClosingDate: time.Date(2023, time.July, 1, 0, 0, 0, 0, time.UTC),
// 		},
// 		Terms: Terms{
// 			DepositAmount:              50000,
// 			InspectionPeriod:           14,
// 			FinancingRequired:          true,
// 			FinancingContingencyPeriod: 30,
// 			TitleCompany:               "ABC Title Company",
// 			EscrowAgent:                "XYZ Escrow Services",
// 		},
// 		AdditionalTerms: []string{
// 			"The seller agrees to provide a one-year home warranty to the buyer.",
// 			"The buyer acknowledges that the property is being sold 'as is' with no warranties.",
// 		},
// 		Signatures: struct {
// 			SellerSignature string    `json:"seller_signature"`
// 			BuyerSignature  string    `json:"buyer_signature"`
// 			Date            time.Time `json:"date"`
// 		}{
// 			SellerSignature: "John Doe",
// 			BuyerSignature:  "Jane Smith",
// 			Date:            time.Now(),
// 		},
// 	}

// 	jsonData, err := json.MarshalIndent(agreement, "", "  ")
// 	if err != nil {
// 		fmt.Println("Error marshaling JSON:", err)
// 		return
// 	}

// 	fmt.Println(string(jsonData))
// }
