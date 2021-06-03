package handler

import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/hamidds/restapi/model"
	"github.com/hamidds/restapi/store"
	"io/ioutil"
	"net/http"
	"time"
)

var CoinStore *store.CoinStore

func CreateCoin(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)

	reqBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		response := model.NewCoinResponse("", "", 0, 0, 400, err.Error())
		json.NewEncoder(writer).Encode(response)
		return
	}

	var coin model.Coin
	err = json.Unmarshal(reqBody, &coin)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		response := model.NewCoinResponse(coin.Name, coin.Symbol, coin.Amount, coin.Rate, 400, err.Error())
		json.NewEncoder(writer).Encode(response)
		return
	}

	validate, trans := setUpValidator()
	filedErrors := validate.Struct(coin)
	if filedErrors != nil {
		var message string
		for _, e := range filedErrors.(validator.ValidationErrors) {
			fmt.Println(e.Translate(trans))
			message = e.Translate(trans)
		}
		writer.WriteHeader(http.StatusBadRequest)
		response := model.NewCoinResponse(coin.Name, coin.Symbol, coin.Amount, coin.Rate, 400, message)
		json.NewEncoder(writer).Encode(response)
		return
	}

	currentWallet, err := WalletStore.GetByName(params["wname"])
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		response := model.NewWalletResponse(params["wname"], 0, []model.Coin{}, getTime(time.Now()), 400, "Wallet doesn't exist!")
		json.NewEncoder(writer).Encode(response)
		return
	}

	if _, err := WalletStore.GetCoinByName(currentWallet, coin.Name); err == nil {
		writer.WriteHeader(http.StatusBadRequest)
		response := model.NewCoinResponse(coin.Name, coin.Symbol, coin.Amount, coin.Rate, 400, "This coin (name) is already added!")
		json.NewEncoder(writer).Encode(response)
		return
	}

	if _, err := WalletStore.GetCoinBySymbol(currentWallet, coin.Symbol); err == nil {
		writer.WriteHeader(http.StatusBadRequest)
		response := model.NewCoinResponse(coin.Name, coin.Symbol, coin.Amount, coin.Rate, 400, "This coin (symbol) is already added!")
		json.NewEncoder(writer).Encode(response)
		return
	}

	writer.WriteHeader(http.StatusOK)
	NewCoin := model.NewCoin(coin.Name, coin.Symbol, coin.Amount, coin.Rate)
	WalletStore.AddCoin(currentWallet, NewCoin)
	response := model.NewCoinResponse(coin.Name, coin.Symbol, coin.Amount, coin.Rate, 200, "Coin added successfully!")
	json.NewEncoder(writer).Encode(response)
}

func GetCoins(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)

	wallet, err := WalletStore.GetByName(params["wname"])
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		response := model.NewWalletResponse(params["wname"], 0, []model.Coin{}, getTime(time.Now()), 400, "Wallet doesn't exist!")
		json.NewEncoder(writer).Encode(response)
		return
	}

	writer.WriteHeader(http.StatusOK)
	response := model.NewWalletResponse(wallet.Name, wallet.Balance, wallet.Coins, getTime(wallet.LastUpdated), 200, "All coins received successfully!")
	json.NewEncoder(writer).Encode(response)
}

func UpdateCoin(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)

	reqBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		response := model.NewWalletResponse(params["wname"], 0, []model.Coin{}, getTime(time.Now()), 400, err.Error())
		json.NewEncoder(writer).Encode(response)
		return
	}

	var coin model.Coin
	err = json.Unmarshal(reqBody, &coin)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		response := model.NewCoinResponse(coin.Name, params["symbol"], coin.Amount, coin.Rate, 400, err.Error())
		json.NewEncoder(writer).Encode(response)
		return
	}

	validate, trans := setUpValidator()
	filedErrors := validate.Struct(coin)
	if filedErrors != nil {
		var message string
		for _, e := range filedErrors.(validator.ValidationErrors) {
			fmt.Println(e.Translate(trans))
			message = e.Translate(trans)
		}
		writer.WriteHeader(http.StatusBadRequest)
		response := model.NewCoinResponse(coin.Name, coin.Symbol, coin.Amount, coin.Rate, 400, message)
		json.NewEncoder(writer).Encode(response)
		return
	}

	wallet, err := WalletStore.GetByName(params["wname"])
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		response := model.NewWalletResponse(params["wname"], 0, []model.Coin{}, getTime(time.Now()), 400, "Wallet doesn't exist!")
		json.NewEncoder(writer).Encode(response)
		return
	}

	_, err = WalletStore.GetCoinBySymbol(wallet, params["symbol"])
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		response := model.NewCoinResponse("", params["symbol"], 0, 0, 400, "Coin doesn't exist!")
		json.NewEncoder(writer).Encode(response)
		return
	}

	if existed, err := WalletStore.GetCoinByName(wallet, coin.Name); err == nil && existed.Name != coin.Name{
		writer.WriteHeader(http.StatusBadRequest)
		response := model.NewCoinResponse(existed.Name, existed.Symbol, existed.Amount, existed.Rate, 400, "This coin (name) exists!")
		json.NewEncoder(writer).Encode(response)
		return
	}

	if existed, err := WalletStore.GetCoinBySymbol(wallet, coin.Symbol); err == nil && params["symbol"] != coin.Symbol {
		writer.WriteHeader(http.StatusBadRequest)
		response := model.NewCoinResponse(existed.Name, existed.Symbol, existed.Amount, existed.Rate, 400, "This coin (name) exists!")
		json.NewEncoder(writer).Encode(response)
		return
	}

	err = WalletStore.UpdateCoin(wallet, coin, params["symbol"])
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		response := model.NewCoinResponse(coin.Name, params["symbol"], coin.Amount, coin.Rate, 400, err.Error())
		json.NewEncoder(writer).Encode(response)
		return
	}

	writer.WriteHeader(http.StatusOK)
	response := model.NewCoinResponse(coin.Name, coin.Symbol, coin.Amount, coin.Rate, 200, "Coin updated successfully!")
	json.NewEncoder(writer).Encode(response)
}

func DeleteCoin(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)

	wallet, err := WalletStore.GetByName(params["wname"])
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		response := model.NewWalletResponse(params["wname"], 0, []model.Coin{}, getTime(time.Now()), 400, "Wallet doesn't exist!")
		json.NewEncoder(writer).Encode(response)
		return
	}

	coin, err := WalletStore.GetCoinBySymbol(wallet, params["symbol"])
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		response := model.NewCoinResponse("", params["symbol"], 0, 0, 400, "Coin doesn't exist!")
		json.NewEncoder(writer).Encode(response)
		return
	}

	err = WalletStore.RemoveCoin(wallet, params["symbol"])
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		response := model.NewCoinResponse(coin.Name, params["symbol"], coin.Amount, coin.Rate, 400, err.Error())
		json.NewEncoder(writer).Encode(response)
		return
	}

	writer.WriteHeader(http.StatusOK)
	response := model.NewCoinResponse(coin.Name, coin.Symbol, coin.Amount, coin.Rate, 200, "Coin deleted successfully!")
	json.NewEncoder(writer).Encode(response)
}

func getTime(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}
