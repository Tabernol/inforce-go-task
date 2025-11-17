package rarible

import (
	"context"
)

type Client interface {
	// Get NFT ownership by ID
	GetOwnershipByID(ctx context.Context, ownershipID string) (*Ownership, error)

	// Query trait rarities
	QueryTraitsWithRarity(ctx context.Context, req TraitRarityRequest) (*TraitRarityResponse, error)
}
