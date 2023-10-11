package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/digihub/subscheck/simswapv2", handleSimSwapRequest)
	fmt.Println("Server listening on :8081")
	http.ListenAndServe(":8081", nil)
}

type SimSwapRequest struct {
	TransactionID string `json:"transaction_id"`
	ConsentID     int    `json:"consent_id"`
	ConsentName   string `json:"consent_name"`
	MSISDN        string `json:"msisdn"`
	Parameter     struct {
		PartnerName string `json:"partner_name"`
		PartnerID   string `json:"partner_id"`
	} `json:"parameter"`
}

type SimSwapResponse struct {
	TransactionID string `json:"transaction_id"`
	StatusCode    string `json:"status_code"`
	StatusDesc    string `json:"status_desc"`
	Score         string `json:"score,omitempty"`
}

func printRequestJSON(r map[string]interface{}) string {
	var jsonStr string
	for key, value := range r {
		if jsonStr != "" {
			jsonStr += ", "
		}
		switch v := value.(type) {
		case map[string]interface{}:
			jsonStr += "\"" + key + "\": {" + printRequestJSON(v) + "}"
		default:
			jsonStr += "\"" + key + "\": \"" + fmt.Sprintf("%v", v) + "\""
		}
	}
	return jsonStr
}

func handleSimSwapRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request:")

	// Print request method, URL, and headers
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)
	fmt.Println("Headers:")
	for key, values := range r.Header {
		fmt.Printf("%s: %s\n", key, values)
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var request map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Print request body in JSON format
	requestJSON := "{" + printRequestJSON(request) + "}"
	fmt.Println("Request Body (JSON):")
	fmt.Println(requestJSON)

	// Generate a random score between 1 and 4
	rand.Seed(time.Now().UnixNano())

	//response := SimSwapResponse{
	//	TransactionID: request["transaction_id"].(string), // Assuming "TransactionID" is a string key
	//	StatusCode:    "00000",
	//	StatusDesc:    "Success",
	//	Score:         fmt.Sprintf("%d", rand.Intn(4)+1),
	//}
	//
	//w.WriteHeader(http.StatusOK)

	response := SimSwapResponse{
		//TransactionID: request["transaction_id"].(string), // Assuming "TransactionID" is a string key
		StatusCode: "20005",
		StatusDesc: "Inactive MSISDN / MSISDN Not Found",
		//Score:         fmt.Sprintf("%d", randomScore),
	}

	w.WriteHeader(http.StatusBadRequest)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
