package entity

type CoinTransaction struct {
	User   string
	Amount int
}

type CoinHistory struct {
	Received []CoinTransaction
	Sent     []CoinTransaction
}
