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
	Traits []TraitRarity `json:"traits"`
}

// Individual trait rarity entry
type TraitRarity struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Rarity string `json:"rarity"`
}

type TraitProperty struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
