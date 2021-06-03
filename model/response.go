package model

type WalletsResponse struct {
	Size    int       `json:"size"`
	Wallets *[]Wallet `json:"wallets"`
	Code    int       `json:"code"`
	Message string    `json:"message"`
}

func NewWalletsResponse(size int, wallets *[]Wallet, code int, message string) WalletsResponse {
	return WalletsResponse{Size: size, Wallets: wallets, Code: code, Message: message}
}

type WalletResponse struct {
	Name        string    `json:"name"          validate:"required"`
	Balance     float64   `json:"balance"       validate:"isdefault"`
	Coins       []Coin    `json:"coins"         validate:"isdefault"`
	LastUpdated string `json:"last_updated"  validate:"isdefault"`
	Code        int       `json:"code"`
	Message     string    `json:"message"`
}

func NewWalletResponse(name string, balance float64, coins []Coin, lastUpdated string, code int, message string) WalletResponse {
	return WalletResponse{Name: name, Balance: balance, Coins: coins, LastUpdated: lastUpdated, Code: code, Message: message}
}

type CoinResponse struct {
	Name    string  `json:"name"    validate:"required"`
	Symbol  string  `json:"symbol"  validate:"required"`
	Amount  float64 `json:"amount"  validate:"required"`
	Rate    float64 `json:"rate"    validate:"required"`
	Code    int     `json:"code"`
	Message string  `json:"message"`
}

func NewCoinResponse(name string, symbol string, amount float64, rate float64, code int, message string) CoinResponse {
	return CoinResponse{Name: name, Symbol: symbol, Amount: amount, Rate: rate, Code: code, Message: message}
}
