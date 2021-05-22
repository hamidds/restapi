package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/hamidds/restapi/model"
	"log"
	"net/http"
)

func main() {
	// Init Router
	r := mux.NewRouter()

	model.Wallets = append(model.Wallets, model.Wallet{Name: "a", Coins: nil, Balance: 16.6})

	// Route Handlers
	r.HandleFunc("/wallets", createWallet).Methods("POST")
	r.HandleFunc("/wallets", getWallets).Methods("GET")
	r.HandleFunc("/wallets/{wname}", updateWallet).Methods("PUT")
	r.HandleFunc("/wallets/{wname}", deleteWallet).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}

func deleteWallet(writer http.ResponseWriter, request *http.Request) {

}

func updateWallet(writer http.ResponseWriter, request *http.Request) {

}

func getWallets(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(model.Wallets)
}

func createWallet(writer http.ResponseWriter, request *http.Request) {

}
