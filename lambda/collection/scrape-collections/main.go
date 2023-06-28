package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {

	url := "https://api-mainnet.magiceden.dev/v2/collections"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var objects []interface{}

	// Unmarshal the JSON byte array into the array of interfaces
	err := json.Unmarshal(body, &objects)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the decoded objects
	for index, obj := range objects {
		fmt.Println(obj)
		fmt.Println("\n\n")

		if index == 3 {
			break
		}

	}

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


type StatsOfCollectionData struct{
	FloorPrice   int     `json:"floorPrice"`
	ListedCount  int     `json:"listedCount"`
	AvgPrice24Hr float64 `json:"avgPrice24hr"`
	VolumeAll    float64 `json:"volumeAll"`

}


type ListingOfCollectionData struct{
		
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


type HolderStatsOfCollectionData struct {
	Symbol         string `json:"symbol"`
	TotalSupply    int    `json:"totalSupply"`
	UniqueHolders  int    `json:"uniqueHolders"`
	TokenHistogram struct {
		Bars []struct {
			LVal  int `json:"l_val"`
			Hight int `json:"hight"`
		} `json:"bars"`
	} `json:"tokenHistogram"`
	TopHolders []struct {
		Tokens int    `json:"tokens"`
		Owner  string `json:"owner"`
		Buy7D  struct {
			Count  int   `json:"count"`
			Volume int64 `json:"volume"`
		} `json:"buy7d"`
		Sell7D struct {
			Count  int   `json:"count"`
			Volume int64 `json:"volume"`
		} `json:"sell7d"`
	} `json:"topHolders"`
}



id

banner_url
tags
creator_id


disputed_message
documentation
endpoint
facebook

instagram
is_curated
✓is_derivative
is_nsfw
links
medium
mint_list_url
roadmap
solo_image_url
solo_username
solo_verified
state
team
telegram
thumbnail_url
tiktok

verifeyed
volume_past_24h
volume_past_7d
candy_machine_addresses
whitepaper
youtube
created_at
modified_at
volume_modified_at
items

