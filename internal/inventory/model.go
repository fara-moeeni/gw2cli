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
