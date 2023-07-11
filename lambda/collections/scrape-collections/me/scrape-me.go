package me

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var (
	ctx = context.Background()

	client *http.Client
)

func getAndProcessListings(totalListing int, symbol string) {
	// https://api-mainnet.magiceden.dev/v2/collections/rox/listings?offset=0&limit=20
	var (
		limit  = 20
		offset = 0
	)

	for offset < totalListing {
		resp, err := client.Get(
			fmt.Sprintf(
				"https://api-mainnet.magiceden.dev/v2/collections/%s/listings?offset=%d&limit=20",
				symbol,
				offset,
			),
		)

		offset += limit
		if err != nil {
			log.Printf("Error sending request, err: %s", err)
			continue
		}

		// to avoid rate limiter
		time.Sleep(400 * time.Millisecond)

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading listing for collection %s, err: %s", symbol, err)
			continue
		}

		var listings []Listing
		err = json.Unmarshal(body, &listings)
		if err != nil {
			log.Printf("Error unmarshalling into json, err: %s", err)
			continue
		}

		for _, listing := range listings {
			var offer = shared.Offer{
				Key:         datastore.NameKey(shared.DATASTORE_OFFER_KIND, listing.PdaAddress, nil),
				Owner:       listing.Seller,
				Mint:        listing.TokenMint,
				Price:       int64(listing.Price * 1.e9),
				Marketplace: "Magic Eden",
				Contract:    contracts.Contracts["Magic Eden"][0].PublicKey,
				PubKey:      listing.PdaAddress,
			}

			// publish to pubsub
			if data, err := offer.ToJSON(); err == nil {
				log.Printf("addOffer: %s", data)
				results = append(results, addOfferTopic.Publish(ctx, &pubsub.Message{Data: data}))
			} else {
				log.Printf("Error marshalling, err: %s", err)
			}
		}
	}

	if len(results) > 0 {
		log.Printf("Published %d offer to %s", len(results), "addOffer")
		for _, r := range results {
			id, err := r.Get(ctx)
			if err != nil {
				log.Fatalf("Could not publish message with id %s: %v", id, err)
			}
		}
	}
}

func ScrapeCollections() {

	// Get API key
	err := godotenv.Load("apikey.env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	API_KEY := os.Getenv("ME_API_KEY")

	// 1) Get collections Stats

	var collectionStats StatsOfCollectionData

	url := "https://api-mainnet.magiceden.dev/v2/collections/vanguards/stats"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", API_KEY)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	err = json.Unmarshal(body, &collectionStats)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 2) collection holder stats

	var collectionHolderStats HolderStatsOfCollectionData

	url = "https://api-mainnet.magiceden.dev/v2/collections/vanguards/holder_stats"

	req, _ = http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", API_KEY)

	res, _ = http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ = io.ReadAll(res.Body)

	err = json.Unmarshal(body, &collectionHolderStats)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 3)

}

type CollectionData struct {
	Symbol      string   `json:"symbol"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Image       string   `json:"image"`
	Twitter     string   `json:"twitter"`
	Discord     string   `json:"discord"`
	Website     string   `json:"website"`
	Categories  []string `json:"categories"`
	IsBadged    bool     `json:"isBadged"`
}

type StatsOfCollectionData struct {
	FloorPrice   int     `json:"floorPrice"`
	ListedCount  int     `json:"listedCount"`
	AvgPrice24Hr float64 `json:"avgPrice24hr"`
	VolumeAll    float64 `json:"volumeAll"`
}

type ListingOfCollectionData struct {
	PdaAddress     string  `json:"pdaAddress"`
	AuctionHouse   string  `json:"auctionHouse"`
	TokenAddress   string  `json:"tokenAddress"`
	TokenMint      string  `json:"tokenMint"`
	Seller         string  `json:"seller"`
	SellerReferral string  `json:"sellerReferral"`
	TokenSize      int     `json:"tokenSize"`
	Price          float64 `json:"price"`
	Rarity         struct {
		Moonrank struct {
			Rank           int `json:"rank"`
			AbsoluteRarity int `json:"absolute_rarity"`
			Crawl          struct {
			} `json:"crawl"`
		} `json:"moonrank"`
	} `json:"rarity"`
	Extra struct {
		Img string `json:"img"`
	} `json:"extra"`
	Expiry int `json:"expiry"`
}
