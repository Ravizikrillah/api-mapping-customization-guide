package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/example/target", handleExampleRequest)
	fmt.Println("Server listening on :8083")
	http.ListenAndServe(":8083", nil)
}

type ExampleRequest struct {
	PhoneNumber  string `json:"phoneNumber"`
	MaxAge       int    `json:"maxAge"`
	CountryCode  string `json:"countryCode"`
	CategoryCode string `json:"categoryCode"`
}

type ExampleResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    Data   `json:"data"`
}

type Data struct {
	NewPhoneNumber string `json:"newPhoneNumber"`
	Age            int    `json:"age"`
	Score          int    `json:"score"`
	Location       string `json:"location"`
	ProductCodes   string `json:"productCodes"`
	Description    string `json:"description"`
	Enabled        bool   `json:"enabled"`
	Items          []Item `json:"items"`
}

type Item struct {
	Location    string `json:"location"`
	Product     string `json:"product"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

func handleExampleRequest(w http.ResponseWriter, r *http.Request) {
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

	var request ExampleRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Print request body in JSON format
	requestJSON, _ := json.Marshal(request)
	fmt.Println("Request Body (JSON):")
	fmt.Println(string(requestJSON))

	// Dummy response
	response := ExampleResponse{
		Status:  "success",
		Message: "Request successful",
		Data: Data{
			NewPhoneNumber: request.PhoneNumber,
			Age:            request.MaxAge,
			Score:          2,
			Location:       request.CountryCode,
			ProductCodes:   request.CategoryCode,
			Description:    "Sample Description",
			Enabled:        true,
			Items: []Item{
				{
					Location:    "Item Location 1",
					Product:     "Item Product 1",
					Description: "Item Description 1",
					Enabled:     true,
				},
				{
					Location:    "Item Location 2",
					Product:     "Item Product 2",
					Description: "Item Description 2",
					Enabled:     false,
				},
			},
		},
	}

	fmt.Println(response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
