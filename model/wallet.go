package model

import (
	"time"
)

type Wallet struct {
	Name        string    `json:"name"         bson:"name"         validate:"required"`
	Balance     float64   `json:"balance"      bson:"balance"      validate:"isdefault"`
	Coins       []Coin    `json:"coins"        bson:"coins"        validate:"isdefault"`
	LastUpdated time.Time `json:"last_updated" bson:"last_updated" validate:"isdefault"`
	//Others      interface{} `json:"-" validate:"isdefault"`
}

func NewWallet(name string, balance float64, coins []Coin, lastUpdated time.Time) *Wallet {
	return &Wallet{Name: name, Balance: balance, Coins: coins, LastUpdated: lastUpdated}
}

func (cw *Wallet) CalculateBalance() {
	var balance float64
	for _, coin := range cw.Coins {
		balance += coin.Amount * coin.Rate
	}
	cw.Balance = balance
}
