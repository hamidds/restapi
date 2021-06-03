package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/hamidds/restapi/model"
	"github.com/hamidds/restapi/store"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//var (
//	walletsDb   *mongo.Collection
//	coinsDb     *mongo.Collection
//	WalletStore *store.WalletStore
//)

var WalletStore *store.WalletStore

func DeleteWallet(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	var response model.WalletResponse

	currentWallet, err := WalletStore.GetByName(params["wname"])
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		response = model.NewWalletResponse("", 0, nil, getTime(time.Now()), 400, "Wallet doesn't exist!")
		json.NewEncoder(writer).Encode(response)
		return
	}

	err = WalletStore.Remove("name", params["wname"])
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		response = model.NewWalletResponse("", 0, nil, getTime(time.Now()), 400, err.Error())
		json.NewEncoder(writer).Encode(response)
		return
	}
	writer.WriteHeader(http.StatusOK)
	response = model.NewWalletResponse(currentWallet.Name, currentWallet.Balance, currentWallet.Coins, getTime(currentWallet.LastUpdated), 200, "Wallet deleted (logged out) successfully!")
	json.NewEncoder(writer).Encode(response)
}

func UpdateWallet(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)

	reqBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		response := model.NewWalletResponse("", 0, []model.Coin{}, getTime(time.Now()), 400, err.Error())
		json.NewEncoder(writer).Encode(response)
		return
	}

	var response model.WalletResponse
	var wallet model.Wallet
	err = json.Unmarshal(reqBody, &wallet)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		response = model.NewWalletResponse("", 0, nil, getTime(time.Now()), 400, "Wallet doesn't exist!")
		json.NewEncoder(writer).Encode(response)
		return
	}

	if _, err := WalletStore.GetByName(wallet.Name); err == nil {
		writer.WriteHeader(http.StatusBadRequest)
		response = model.NewWalletResponse(wallet.Name, wallet.Balance, wallet.Coins, getTime(wallet.LastUpdated), 400, "Username is already taken!")
		json.NewEncoder(writer).Encode(response)
		return
	}

	err = WalletStore.UpdateName(&wallet, params["wname"])
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		response = model.NewWalletResponse(wallet.Name, wallet.Balance, wallet.Coins, getTime(wallet.LastUpdated), 400, err.Error())
		json.NewEncoder(writer).Encode(response)
		return
	}

	writer.WriteHeader(http.StatusOK)
	updatedWallet, _ := WalletStore.GetByName(wallet.Name)
	response = model.NewWalletResponse(updatedWallet.Name, updatedWallet.Balance, updatedWallet.Coins, getTime(updatedWallet.LastUpdated), 200, "Wallet name changed successfully!")
	json.NewEncoder(writer).Encode(response)
}

func GetWallets(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var wallets *[]model.Wallet

	wallets, err := WalletStore.GetWalletsList()
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		response := model.NewWalletResponse("", 0, []model.Coin{}, getTime(time.Now()), 400, err.Error())
		json.NewEncoder(writer).Encode(response)
		return
	}

	writer.WriteHeader(http.StatusOK)
	response := model.NewWalletsResponse(len(*wallets), wallets, 200, "All wallets received successfully!")
	json.NewEncoder(writer).Encode(response)
}

func CreateWallet(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	reqBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		response := model.NewWalletResponse("", 0, []model.Coin{}, getTime(time.Now()), 400, err.Error())
		json.NewEncoder(writer).Encode(response)
		return
	}

	var response model.WalletResponse
	var wallet model.Wallet

	err = json.Unmarshal(reqBody, &wallet)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		response = model.NewWalletResponse(wallet.Name, wallet.Balance, wallet.Coins, getTime(wallet.LastUpdated), 400, err.Error())
		json.NewEncoder(writer).Encode(response)
		return
	}

	validate, trans := setUpValidator()
	filedErrors := validate.Struct(wallet)
	if filedErrors != nil {
		var message string
		for _, e := range filedErrors.(validator.ValidationErrors) {
			fmt.Println(e.Translate(trans))
			message = e.Translate(trans)
		}
		writer.WriteHeader(http.StatusBadRequest)
		response = model.NewWalletResponse(wallet.Name, wallet.Balance, wallet.Coins, getTime(wallet.LastUpdated), 400, message)
		json.NewEncoder(writer).Encode(response)
		return
	}

	if _, err := WalletStore.GetByName(wallet.Name); err == nil {
		writer.WriteHeader(http.StatusBadRequest)
		response = model.NewWalletResponse(wallet.Name, wallet.Balance, wallet.Coins, getTime(wallet.LastUpdated), 400, "Username is already taken!")
		json.NewEncoder(writer).Encode(response)
		return
	}

	newWallet := model.NewWallet(wallet.Name, 0.0, []model.Coin{}, time.Now())
	err = WalletStore.Create(newWallet)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		response = model.NewWalletResponse(wallet.Name, wallet.Balance, wallet.Coins, getTime(wallet.LastUpdated), 400, err.Error())
		json.NewEncoder(writer).Encode(response)
		return
	}

	writer.WriteHeader(http.StatusOK)
	response = model.NewWalletResponse(wallet.Name, 0.0, []model.Coin{}, getTime(time.Now()), 200, "wallet is created successfully!")
	json.NewEncoder(writer).Encode(response)
}

func setUpValidator() (*validator.Validate, ut.Translator) {
	validate := validator.New()
	translator := en.New()
	uni := ut.New(translator, translator)

	trans, found := uni.GetTranslator("en")
	if !found {
		log.Fatal("translator not found")
	}

	_ = validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is required!", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())

		return t
	})

	_ = validate.RegisterTranslation("isdefault", trans, func(ut ut.Translator) error {
		return ut.Add("isdefault", "{0} should not be sent from user!", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("isdefault", fe.Field())

		return t
	})

	return validate, trans
}
