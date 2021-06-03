package store

import (
	"context"
	"github.com/hamidds/restapi/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type WalletStore struct {
	db *mongo.Collection
}

func NewWalletStore(db *mongo.Collection) *WalletStore {
	return &WalletStore{db: db}
}

func (ws *WalletStore) Create(w *model.Wallet) error {
	_, err := ws.db.InsertOne(context.TODO(), w)
	return err
}

func (ws *WalletStore) Remove(field, value string) error {
	_, err := ws.db.DeleteOne(context.TODO(), bson.M{field: value})
	return err
}

func (ws *WalletStore) GetByName(name string) (*model.Wallet, error) {
	var w model.Wallet
	err := ws.db.FindOne(context.TODO(), bson.M{"name": name}).Decode(&w)
	return &w, err
}

func (ws *WalletStore) GetWalletsList() (*[]model.Wallet, error) {
	var wallets []model.Wallet
	query := bson.M{}
	res, err := ws.db.Find(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	if err = res.All(context.TODO(), &wallets); err != nil {
		return nil, err
	}
	return &wallets, err
}

func (ws *WalletStore) UpdateName(w *model.Wallet, name string) error {
	old, err := ws.GetByName(name)
	if err != nil {
		return err
	}
	_, err = ws.db.UpdateOne(context.TODO(),
		bson.M{"name": name},
		bson.M{"$set": bson.M{
			"name":         w.Name,
			"balance":      old.Balance,
			"coins":        old.Coins,
			"last_updated": time.Now().Format("2006-01-02 15:04:05"),
		},
		})
	return err
}


