package me

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	ctx = context.Background()

	client http.Client
)

// Retrieve Top collections
func getTopCollections() []Collection {

	var topCollections []Collection

	resp, err := client.Get("https://api-mainnet.magiceden.dev/v2/marketplace/popular_collections?timeRange=1d")
	if err != nil {
		log.Printf("Error sending request, err: %s", err)
		return topCollections
	}

	// to avoid rate limiter
	// time.Sleep(400 * time.Millisecond)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error getting collections, err: %s", err)
		return topCollections
	}

	err = json.Unmarshal(body, &topCollections)
	if err != nil {
		log.Printf("Error unmarshalling into json, err: %s", err)
		return topCollections
	}

	return topCollections

}

// Check if the collections exist in the backend
func getMissingCollection(collecions []Collection) []Collection {

	return collecions
}

// Scrape collecion stats
func getCollectionStats(collectionSymbol string) CollectionStats {

	var collectionStats CollectionStats

	resp, err := client.Get(fmt.Sprintf("https://api-mainnet.magiceden.dev/v2/collections/%s/stats", collectionSymbol))
	if err != nil {
		log.Printf("Error sending request, err: %s", err)
		return collectionStats
	}

	// to avoid rate limiter
	// time.Sleep(400 * time.Millisecond)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error getting stat for collection %s, err: %s", collectionSymbol, err)
		return collectionStats
	}

	err = json.Unmarshal(body, &collectionStats)
	if err != nil {
		log.Printf("Error unmarshalling into json, err: %s", err)
		return collectionStats
	}

	return collectionStats

}

func getCollectionHolderStats(collectionSymbol string) CollectionHolderStats {

	var collectionHolderStats CollectionHolderStats

	resp, err := client.Get(fmt.Sprintf("https://api-mainnet.magiceden.dev/v2/collections/%s/holder_stats", collectionSymbol))
	if err != nil {
		log.Printf("Error sending request, err: %s", err)
		return collectionHolderStats
	}

	// to avoid rate limiter
	// time.Sleep(400 * time.Millisecond)

	body, err := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(body))

	if err != nil {
		log.Printf("Error getting stat for collection %s, err: %s", collectionSymbol, err)
		return collectionHolderStats
	}

	err = json.Unmarshal(body, &collectionHolderStats)
	if err != nil {
		log.Printf("Get collcetionholderstats Error unmarshalling into json, err: %s", err)
		return collectionHolderStats
	}

	return collectionHolderStats
}

// Scrape collection Listing
func getCollectionListing(collectionSymbol string, listedCount int) []CollectionListing {
	var (
		limit             = 20
		offset            = 0
		collectionListing []CollectionListing
	)

	for offset < listedCount {
		resp, err := client.Get(
			fmt.Sprintf(
				"https://api-mainnet.magiceden.dev/v2/collections/%s/listings?offset=%d&limit=20",
				collectionSymbol,
				offset,
			),
		)

		offset += limit
		if err != nil {
			log.Printf("Error sending request, err: %s", err)
			continue
		}

		// to avoid rate limiter
		// time.Sleep(400 * time.Millisecond)

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading listing for collection %s, err: %s", collectionSymbol, err)
			continue
		}

		var listings []CollectionListing
		err = json.Unmarshal(body, &listings)
		if err != nil {
			log.Printf("Error unmarshalling into json, err: %s", err)
			continue
		}

		collectionListing = append(collectionListing, listings...)
		break

	}

	return collectionListing

}

func ScrapeMagicEden() ([]MagicEdenCollectionData, map[string]string) {

	var magicEdenCollectionData []MagicEdenCollectionData

	// Scrape top collection
	var collections []Collection = getTopCollections()

	println("collection")
	fmt.Println(collections[0])
	println("\n\n")

	// find collections thats not in kyzzen backend
	var missingCollections []Collection = getMissingCollection(collections)

	collectionSymbolToTokenMintMap := make(map[string]string)

	// Iterate though the missing collections
	for index, collection := range missingCollections {

		//Scrape collection stats
		var collectionStats CollectionStats = getCollectionStats(collection.Symbol)

		println("collection stats")
		fmt.Println(collectionStats)
		println("\n\n")

		//scrape collection holder stats
		// var collectionHolderStats CollectionHolderStats = getCollectionHolderStats(collection.Symbol)

		//scrape collection listing
		var collectionListing []CollectionListing = getCollectionListing(collection.Symbol, collectionStats.ListedCount)
		println("collection listing")
		fmt.Println(collectionListing[0])
		println("\n\n")

		meCollectionData := MagicEdenCollectionData{
			Collection:      collection,
			CollectionStats: collectionStats,
			// CollectionHolderStats: collectionHolderStats,
			CollectionListing: collectionListing,
		}

		magicEdenCollectionData = append(magicEdenCollectionData, meCollectionData)

		collectionSymbolToTokenMintMap[collection.Symbol] = collectionListing[0].TokenMint

		break

		if index == 4 {
			break
		}

	}

	return magicEdenCollectionData, collectionSymbolToTokenMintMap

}
