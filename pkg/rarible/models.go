package rarible

// Ownership response model
type Ownership struct {
	Id       string `json:"id"`
	Owner    string `json:"owner"`
	Contract string `json:"contract"`
	TokenId  string `json:"tokenId"`
	Value    string `json:"value"`
}

// Trait rarity request
type TraitRarityRequest struct {
	CollectionId string          `json:"collectionId"`
	Properties   []TraitProperty `json:"properties"`
	Limit        int             `json:"limit,omitempty"`
}

// Trait rarity response
type TraitRarityResponse struct {
	Items      []TraitRarity `json:"items"`
	NextCursor string        `json:"nextCursor"`
}

// Individual trait rarity entry
type TraitRarity struct {
	TraitType string  `json:"traitType"`
	Value     string  `json:"value"`
	Rarity    float64 `json:"rarity"`
}

type TraitProperty struct {
	Name   string   `json:"name,omitempty"`
	Values []string `json:"values,omitempty"`
}
