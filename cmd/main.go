package main

import (
	"context"
	"fmt"
	"github.com/Tabernol/inforce-go-task/pkg/rarible"
	"log"
)

func main() {
	cfg := rarible.NewConfigFromEnv()
	client := rarible.NewClient(cfg)

	// ---- Test 1: Ownership ----
	ownershipID := "ETHEREUM%3A0xb66a603f4cfe17e3d27b87a8bfcad319856518b8%3A32292934596187112148346015918544186536963932779440027682601542850818403729410%3A0x4765273c477c2dc484da4f1984639e943adccfeb"

	resp, err := client.GetOwnershipByID(context.Background(), ownershipID)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("Ownership result:\n%+v\n", resp)

	// ---- Test 2: Trait rarity ----
	req := rarible.TraitRarityRequest{
		CollectionId: "ETHEREUM:0x60e4d786628fea6478f785a6d7e704777c86a7c6",
		Properties: []rarible.TraitProperty{
			{Key: "Hat", Value: "Halo"},
		},
		Limit: 12,
	}

	rarityResp, err := client.QueryTraitsWithRarity(context.Background(), req)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("Trait rarity result:\n%+v\n", rarityResp)
}
