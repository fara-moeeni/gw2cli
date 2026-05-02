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

type Ingredient struct {
	ItemID int `json:"item_id"`
	Count  int `json:"count"`
}

type Recipe struct {
	ID           int          `json:"id"`
	OutputItemID int          `json:"output_item_id"`
	MinRating    int          `json:"min_rating"`
	Disciplines  []string     `json:"disciplines"`
	Ingredients  []Ingredient `json:"ingredients"`
}

type NamedEntity struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Skin struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type DailyAchievement struct {
	ID        int      `json:"id"`
	Level     struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"level"`
	RequiredAccess []string `json:"required_access"`
}

type DailyAchievements struct {
	PVE      []DailyAchievement `json:"pve"`
	PVP      []DailyAchievement `json:"pvp"`
	WVW      []DailyAchievement `json:"wvw"`
	Fractals []DailyAchievement `json:"fractals"`
	Special  []DailyAchievement `json:"special"`
}

type Achievement struct {
	ID          int               `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Requirement string            `json:"requirement"`
	LockedText  string            `json:"locked_text"`
	Type        string            `json:"type"`
	Flags       []string          `json:"flags"`
	Tiers       []AchievementTier `json:"tiers"`
}

type Dungeon struct {
	ID    string `json:"id"`
	Paths []struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"paths"`
}

type Raid struct {
	ID    string `json:"id"`
	Wings []struct {
		ID     string `json:"id"`
		Events []struct {
			ID   string `json:"id"`
			Type string `json:"type"`
		} `json:"events"`
	} `json:"wings"`
}

type WizardsVaultObjective struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Track         string `json:"track"`
	Acclaim       int    `json:"acclaim"`
	ProgressCur   int    `json:"progress_current"`
	ProgressGoal  int    `json:"progress_goal"`
	Claimed       bool   `json:"claimed"`
}

type WizardsVaultResponse struct {
	Objectives []WizardsVaultObjective `json:"objectives"`
}

type AchievementTier struct {
	Count  int `json:"count"`
	Points int `json:"points"`
}

type AchievementCategory struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Order        int    `json:"order"`
	Achievements []int  `json:"achievements"`
}

type AchievementGroup struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Order        int    `json:"order"`
	Categories   []int  `json:"categories"`
}

type AccountAchievement struct {
	ID        int  `json:"id"`
	Current   int  `json:"current"`
	Max       int  `json:"max"`
	Done      bool `json:"done"`
	Repeated  int  `json:"repeated,omitempty"`
	Bits      []int `json:"bits,omitempty"`
}

type Mastery struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Region string `json:"region"`
}

type MasteryPointSummary struct {
	Totals []struct {
		Region string `json:"region"`
		Spent  int    `json:"spent"`
		Earned int    `json:"earned"`
	} `json:"totals"`
}

type Luck struct {
	ID    string `json:"id"`
	Value int    `json:"value"`
}

type AccountMastery struct {
	ID    int `json:"id"`
	Level int `json:"level"`
}

type Account struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Age          int      `json:"age"`
	World        int      `json:"world"`
	Guilds       []string `json:"guilds"`
	Created      string   `json:"created"`
	Access       []string `json:"access"`
	Commander    bool     `json:"commander"`
	FractalLevel int      `json:"fractal_level"`
	DailyAP      int      `json:"daily_ap"`
	MonthlyAP    int      `json:"monthly_ap"`
	WvpRank      int      `json:"wvw_rank"`
}
