package dto

type SendCoinRequest struct {
	ToUser string `json:"toUser" validate:"required"`
	Amount int    `json:"amount" validate:"required,min=1"`
}

type CoinHistory struct {
	Received []CoinTransaction `json:"received"`
	Sent     []CoinTransaction `json:"sent"`
}

type CoinTransaction struct {
	User   string `json:"user"`
	Amount int    `json:"amount"`
}
