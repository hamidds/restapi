package store

import (
	"context"
	"errors"
	"github.com/hamidds/restapi/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type CoinStore struct {
	db *mongo.Collection
}

func NewCoinStore(db *mongo.Collection) *CoinStore {
	return &CoinStore{db: db}
}

func (ws *WalletStore) AddCoin(wallet *model.Wallet, coin *model.Coin) error {
	wallet.Coins = append(wallet.Coins, *coin)
	wallet.CalculateBalance()
	now, _ := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	_, err := ws.db.UpdateOne(context.TODO(), bson.M{"name": wallet.Name}, bson.M{"$set": bson.M{"coins": wallet.Coins, "balance": wallet.Balance, "last_updated": now}})
	if err != nil {
		return err
	}
	return nil
}

func (ws *WalletStore) GetCoinBySymbol(w *model.Wallet, symbol string) (model.Coin, error) {
	for _, o := range w.Coins {
		if o.Symbol == symbol {
			return o, nil
		}
	}
	return model.Coin{}, errors.New("coin doesn't exist")
}

func (ws *WalletStore) GetCoinByName(w *model.Wallet, name string) (model.Coin, error) {
	for _, o := range w.Coins {
		if o.Name == name {
			return o, nil
		}
	}
	return model.Coin{}, errors.New("coin doesn't exist")
}

func (ws *WalletStore) RemoveCoin(wallet *model.Wallet, symbol string) error {
	newCoins := &[]model.Coin{}
	for _, o := range wallet.Coins {
		if o.Symbol != symbol {
			*newCoins = append(*newCoins, o)
		}
	}
	wallet.Coins = *newCoins
	wallet.CalculateBalance()
	now, _ := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	_, err := ws.db.UpdateOne(context.TODO(), bson.M{"name": wallet.Name}, bson.M{"$set": bson.M{"coins": newCoins, "balance": wallet.Balance, "last_updated": now}})
	if err != nil {
		return err
	}
	wallet.Coins = *newCoins
	return nil
}

func (ws *WalletStore) UpdateCoin(w *model.Wallet, newCoin model.Coin, symbol string) error {
	newCoins := &[]model.Coin{}
	for _, o := range w.Coins {
		if o.Symbol != symbol {
			*newCoins = append(*newCoins, o)
		} else {
			*newCoins = append(*newCoins, newCoin)
		}
	}
	w.Coins = *newCoins
	w.CalculateBalance()
	now, _ := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	_, err := ws.db.UpdateOne(context.TODO(),
		bson.M{"name": w.Name},
		bson.M{"$set": bson.M{
			"name":         w.Name,
			"balance":      w.Balance,
			"coins":        newCoins,
			"last_updated": now,
		},
		})
	return err
}
