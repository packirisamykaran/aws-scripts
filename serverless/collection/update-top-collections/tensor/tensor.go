package tensor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const API_KEY = "eff00cb4-025d-4bf5-84d5-eb170054af72"
const ENDPOINT = "https://api.tensor.so/graphql"

type TensorCollectionData struct {
	MintList []string `json:"mintList"`
	// CollectionStats CollectionStats `json:"collectionStats"`
}

type CollectionStats struct {
	Floor1h    float64 `json:"floor1h"`
	Floor24h   float64 `json:"floor24h"`
	Floor7d    float64 `json:"floor7d"`
	FloorPrice string  `json:"floorPrice"`
	NumListed  int     `json:"numListed"`
	NumMints   int     `json:"numMints"`
	PriceUnit  string  `json:"priceUnit"`
	Sales1h    int     `json:"sales1h"`
	Sales24h   int     `json:"sales24h"`
	Sales7d    int     `json:"sales7d"`
	Volume1h   string  `json:"volume1h"`
	Volume24h  string  `json:"volume24h"`
	Volume7d   string  `json:"volume7d"`
}

type Mints struct {
	Data struct {
		Mints []struct {
			Slug string `json:"slug"`
		} `json:"mints"`
	} `json:"data"`
}

type MintListResponse struct {
	Data struct {
		MintList []string `json:"mintList"`
	} `json:"data"`
}

func getCollectionSlugByMint(mint string) string {

	client := &http.Client{}

	payload := fmt.Sprintf(
		`{"query":"query Mints($tokenMints: [String!]!) {\n  mints(tokenMints: $tokenMints) {\n    slug\n  }\n}","variables":{"tokenMints":["%s"]},"operationName":"Mints"}`,
		mint,
	)

	req, err := http.NewRequest(
		"POST",
		ENDPOINT,
		strings.NewReader(payload),
	)

	if err != nil {

		return ""
	}

	req.Header.Set("content-type", "application/json")
	req.Header.Set("referer", "https://studio.apollographql.com/")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("X-Tensor-Api-Key", API_KEY)

	resp, err := client.Do(req)
	if err != nil {
		// log.Error("error request", zap.Error(err))
		return ""
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	// log.Info("resp", zap.ByteString("body", body))

	var mints Mints
	if err := json.Unmarshal(body, &mints); err != nil {
		// log.Error("unmarshal body", zap.Error(err))
		return ""
	}

	if len(mints.Data.Mints) > 0 {

		println("Slug")
		fmt.Println(mints.Data.Mints[0].Slug)
		println("\n\n")

		return mints.Data.Mints[0].Slug

	}

	return ""
}

func getMintList(slug string) []string {
	// GraphQL endpoint and payload

	payload := []byte(`{"query":"query MintList($slug: String!) { mintList(slug: $slug) }", "variables": {"slug": "` + slug + `"}}`)

	// Create a new HTTP client
	client := &http.Client{}

	var mintList []string

	// Create an HTTP request with the payload
	req, err := http.NewRequest("POST", ENDPOINT, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return mintList
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-TENSOR-API-KEY", API_KEY)

	// Send the request
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Request failed:", err)
		return mintList
	}

	body, _ := io.ReadAll(response.Body)
	// log.Info("resp", zap.ByteString(
	defer response.Body.Close()

	var mintListResponse MintListResponse

	// Read the response body
	if err := json.Unmarshal(body, &mintListResponse); err != nil {
		// log.Error("unmarshal body", zap.Error(err))
		return mintList
	}

	mintList = mintListResponse.Data.MintList

	// println("mint list")
	// fmt.Println(mintList)
	// println("\n\n")

	return mintList
}

func getCollectionStats(slug string) CollectionStats {
	// GraphQL endpoint and payload

	payload := []byte(`{
		"query": "query CollectionStats($slug: String!) {
			instrumentTV2(slug: $slug) {
				id
				slug
				slugMe
				slugDisplay
				statsOverall {
					floor1h
					floor24h
					floor7d
					floorPrice
					numListed
					numMints
					priceUnit
					sales1h
					sales24h
					sales7d
					volume1h
					volume24h
					volume7d
				}
				statsSwap {
					buyNowPrice
					sellNowPrice
				}
				statsTSwap {
					buyNowPrice
					nftsForSale
					numMints
					priceUnit
					sales7d
					sellNowPrice
					solDeposited
					volume7d
				}
				statsHSwap {
					buyNowPrice
					nftsForSale
					priceUnit
					sales7d
					sellNowPrice
					solDeposited
					volume7d
				}
				tswapTVL
				firstListDate
				name
			}
		}",
		"variables": {
			"slug": "2c6c34a6-a621-4c5c-a47c-af0aac421579"
		}
	}`)

	// Create a new HTTP client
	client := &http.Client{}

	var collectionStats CollectionStats

	// Create an HTTP request with the payload
	req, err := http.NewRequest("POST", ENDPOINT, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return collectionStats
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-TENSOR-API-KEY", API_KEY)

	// Send the request
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Request failed:", err)
		return collectionStats
	}

	body, _ := io.ReadAll(response.Body)
	// log.Info("resp", zap.ByteString(
	defer response.Body.Close()

	var collectionStatsResponse CollectionStatsResponse

	fmt.Println(string(body))

	// Read the response body
	if err := json.Unmarshal(body, &collectionStatsResponse); err != nil {
		// log.Error("unmarshal body", zap.Error(err))
		return collectionStats
	}

	collectionStats = collectionStatsResponse.Data.CollectionOverallStats.Stats

	println("Tensor Collection Stats")
	fmt.Println(body)
	println("\n\n")

	return collectionStats
}

func ScrapeTensor(collectionSymbolToTokenMintMap map[string]string) map[string]TensorCollectionData {

	tensorCollectionData := make(map[string]TensorCollectionData)

	for symbol, tokenMint := range collectionSymbolToTokenMintMap {

		slug := getCollectionSlugByMint(tokenMint)

		mintList := getMintList(slug)

		tensorCollectionData[symbol] = TensorCollectionData{
			MintList: mintList,
			// CollectionStats: collectionStats,
		}

	}

	return tensorCollectionData

}

type ColletionOverallStats struct {
	ID          string `json:"id"`
	Slug        string `json:"slug"`
	SlugMe      string `json:"slugMe"`
	SlugDisplay string `json:"slugDisplay"`
	Stats       struct {
		Floor1h    float64 `json:"floor1h"`
		Floor24h   float64 `json:"floor24h"`
		Floor7d    float64 `json:"floor7d"`
		FloorPrice string  `json:"floorPrice"`
		NumListed  int     `json:"numListed"`
		NumMints   int     `json:"numMints"`
		PriceUnit  string  `json:"priceUnit"`
		Sales1h    int     `json:"sales1h"`
		Sales24h   int     `json:"sales24h"`
		Sales7d    int     `json:"sales7d"`
		Volume1h   string  `json:"volume1h"`
		Volume24h  string  `json:"volume24h"`
		Volume7d   string  `json:"volume7d"`
	} `json:"stats"`
	StatsSwap struct {
		BuyNowPrice  string `json:"buyNowPrice"`
		SellNowPrice string `json:"sellNowPrice"`
	} `json:"statsSwap"`
	StatsTSwap struct {
		BuyNowPrice  string `json:"buyNowPrice"`
		NftsForSale  int    `json:"nftsForSale"`
		NumMints     int    `json:"numMints"`
		PriceUnit    string `json:"priceUnit"`
		Sales7d      int    `json:"sales7d"`
		SellNowPrice string `json:"sellNowPrice"`
		SolDeposited string `json:"solDeposited"`
		Volume7d     string `json:"volume7d"`
	} `json:"statsTSwap"`
	StatsHSwap struct {
		BuyNowPrice  *string `json:"buyNowPrice"`
		NftsForSale  int     `json:"nftsForSale"`
		PriceUnit    string  `json:"priceUnit"`
		Sales7d      int     `json:"sales7d"`
		SellNowPrice *string `json:"sellNowPrice"`
		SolDeposited string  `json:"solDeposited"`
		Volume7d     string  `json:"volume7d"`
	} `json:"statsHSwap"`
	TswapTVL      string `json:"tswapTVL"`
	FirstListDate int64  `json:"firstListDate"`
	Name          string `json:"name"`
}

type CollectionStatsResponse struct {
	Data struct {
		CollectionOverallStats ColletionOverallStats `json:"collectionOverallStats"`
	} `json:"data"`
}
