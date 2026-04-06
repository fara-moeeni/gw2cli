package inventory

// ItemLocation represents where a specific instance of an item is found.
type ItemLocation struct {
	Source string // e.g., "Bank", "Shared", "CharacterName"
	Detail string // e.g., "Tab 1", "Bag 2", "Equipped: Helm"
	Count  int
}

// ItemDetail combines the API data with our location data.
type ItemDetail struct {
	ID        int
	Name      string
	Type      string
	Locations []ItemLocation
}

// TotalCount calculates the total number of this item across all locations.
func (i *ItemDetail) TotalCount() int {
	total := 0
	for _, loc := range i.Locations {
		total += loc.Count
	}
	return total
}

type RecipeIngredientDetail struct {
	ItemID int
	Name   string
	Count  int
}

type RecipeDetail struct {
	ID          int
	OutputName  string
	Discipline  string
	Rating      int
	Ingredients []RecipeIngredientDetail
}

type CollectionItem struct {
	Name string
	Type string
}

type CollectionSummary struct {
	Skins     int
	Dyes      int
	Minis     int
	Mounts    int
	Outfits   int
	Novelties int
	Finishers int
}

type DailyStatus struct {
	Name      string
	Completed bool
}

type FractalDaily struct {
	Name      string
	Tier      string
	Level     int
	Completed bool
}

type WizardsVaultStatus struct {
	Title        string
	ProgressCur  int
	ProgressGoal int
	Acclaim      int
	Completed    bool
}

type RaidWingStatus struct {
	Name   string
	Events []DailyStatus
}

type AchievementProgress struct {
	ID           int
	Name         string
	Description  string
	Requirement  string
	Current      int
	Max          int
	Points       int
	Done         bool
	TierStatus   string // e.g. "2/4 tiers"
	StatusSymbol string // [✓], [~], [ ]
}

type CategorySummary struct {
	Name      string
	Completed int
	Total     int
	AP        int
}

type MasteryRegion struct {
	Name   string
	Spent  int
	Earned int
}

type MasterySummary struct {
	Regions []MasteryRegion
	Luck    int
}
