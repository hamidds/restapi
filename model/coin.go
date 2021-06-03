package model

type Coin struct {
	Name   string  `json:"name"   bson:"name"   validate:"required" `
	Symbol string  `json:"symbol" bson:"symbol" validate:"required" `
	Amount float64 `json:"amount" bson:"amount" validate:"required" `
	Rate   float64 `json:"rate"   bson:"rate"   validate:"required" `
}

func NewCoin(name string, symbol string, amount float64, rate float64) *Coin {
	return &Coin{Name: name, Symbol: symbol, Amount: amount, Rate: rate}
}
