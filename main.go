package main

import (
	"github.com/gorilla/mux"
	"github.com/hamidds/restapi/db"
	"github.com/hamidds/restapi/handler"
	"github.com/hamidds/restapi/store"
	"log"
	"net/http"
)

func main() {
	// Init Router
	r := mux.NewRouter()

	mongoClient, err := db.GetMongoClient()
	if err != nil {
		log.Fatal(err)
	}
	walletsDb := db.SetupWalletsDb(mongoClient)
	handler.WalletStore = store.NewWalletStore(walletsDb)
	//coinsDb := db.SetupCoinsDb(mongoClient)
	//handler.CoinStore = store.NewCoinStore(coinsDb)

	// Wallet Handlers
	r.HandleFunc("/wallets", handler.CreateWallet).Methods("POST")
	r.HandleFunc("/wallets", handler.GetWallets).Methods("GET")
	r.HandleFunc("/wallets/{wname}", handler.UpdateWallet).Methods("PUT")
	r.HandleFunc("/wallets/{wname}", handler.DeleteWallet).Methods("DELETE")

	// Coin Handlers
	r.HandleFunc("/{wname}/coins", handler.CreateCoin).Methods("POST")
	r.HandleFunc("/{wname}", handler.GetCoins).Methods("GET")
	r.HandleFunc("/{wname}/{symbol}", handler.UpdateCoin).Methods("PUT")
	r.HandleFunc("/{wname}/{symbol}", handler.DeleteCoin).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}


