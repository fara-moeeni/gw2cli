package gw2api

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type InventorySlot struct {
	ID      int    `json:"id"`
	Count   int    `json:"count"`
	Binding string `json:"binding,omitempty"`
}

type Bag struct {
	ID        int              `json:"id"`
	Size      int              `json:"size"`
	Inventory []*InventorySlot `json:"inventory"`
}

type EquipmentSlot struct {
	ID   int    `json:"id"`
	Slot string `json:"slot"`
}

type Character struct {
	Name       string          `json:"name"`
	Race       string          `json:"race"`
	Gender     string          `json:"gender"`
	Profession string          `json:"profession"`
	Level      int             `json:"level"`
	Age        int             `json:"age"`
	Created    string          `json:"created"`
	Bags       []*Bag          `json:"bags"`
	Equipment  []EquipmentSlot `json:"equipment"`
}

type AccountInventory []*InventorySlot

type WalletCurrency struct {
	ID    int `json:"id"`
	Value int `json:"value"`
}

type Currency struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Order       int    `json:"order"`
	Icon        string `json:"icon"`
}
