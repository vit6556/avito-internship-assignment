package dto

type EmployeeInfoResponse struct {
	Coins       int              `json:"coins"`
	Inventory   []*InventoryItem `json:"inventory"`
	CoinHistory *CoinHistory     `json:"coinHistory"`
}

type InventoryItem struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}
