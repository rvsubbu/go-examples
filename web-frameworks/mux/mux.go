package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

// This is the response struct that will be serialized and sent back
type StatusResponse struct {
	Status     string `json:"status"`
	ActionItem string `json:"action_item"`
}

func ActionItemGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	w.Header().Set("Content-Type", "application/json")

	body := StatusResponse{
		Status:     "Hello world from mux!",
		ActionItem: vars["ai"],
	}

	serializedBody, _ := json.Marshal(body)
	_, _ = w.Write(serializedBody)
}

type RequestBody struct {
	ActionItem string `json:"ai"`
}

func ActionItemPostHandler(w http.ResponseWriter, r *http.Request) {
	// Read complete request body
	rawRequestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Transform into RequestBody struct
	requestBody := &RequestBody{}
	err = json.Unmarshal(rawRequestBody, requestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	body := StatusResponse{
		Status:     "Hello world from mux!",
		ActionItem: requestBody.ActionItem,
	}

	serializedBody, _ := json.Marshal(body)
	_, _ = w.Write(serializedBody)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/ais/{ai}", ActionItemGetHandler).Methods("GET")
	r.HandleFunc("/ais", ActionItemPostHandler).Methods("POST")

	log.Println("Listening on :8004")
	log.Fatal(http.ListenAndServe(":8004", r))
}
