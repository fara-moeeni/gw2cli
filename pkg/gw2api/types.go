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

type MaterialStorageEntry struct {
	ID       int    `json:"id"`
	Category int    `json:"category"`
	Binding  string `json:"binding,omitempty"`
	Count    int    `json:"count"`
}

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

type CommerceDelivery struct {
	Coins int `json:"coins"`
	Items []struct {
		ID    int `json:"id"`
		Count int `json:"count"`
	} `json:"items"`
}

type CommercePrice struct {
	ID    int `json:"id"`
	Whitelisted bool `json:"whitelisted"`
	Buys  struct {
		Quantity  int `json:"quantity"`
		UnitPrice int `json:"unit_price"`
	} `json:"buys"`
	Sells struct {
		Quantity  int `json:"quantity"`
		UnitPrice int `json:"unit_price"`
	} `json:"sells"`
}

type CommerceTransaction struct {
	ID        int       `json:"id"`
	ItemID    int       `json:"item_id"`
	Price     int       `json:"price"`
	Quantity  int       `json:"quantity"`
	Created   string    `json:"created"`
	Purchased string    `json:"purchased,omitempty"`
}

type CommerceExchange struct {
	CoinsPerGem int `json:"coins_per_gem"`
	Quantity    int `json:"quantity"`
}

type LegendaryArmoryItem struct {
	ID    int `json:"id"`
	Count int `json:"count"`
}
