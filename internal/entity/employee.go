package entity

type Employee struct {
	ID           int
	Balance      int
	Username     string
	PasswordHash string
}

type EmployeeInfo struct {
	Coins       int
	Inventory   []*InventoryItem
	CoinHistory *CoinHistory
}
