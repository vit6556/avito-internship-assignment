package entity

type MerchItem struct {
	ID    int
	Name  string
	Price int
}

type InventoryItem struct {
	Type     string
	Quantity int
}
