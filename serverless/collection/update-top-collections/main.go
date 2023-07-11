package main

import (
	"fmt"
	"serverless/collection/update-top-collections/me"
	"serverless/collection/update-top-collections/tensor"
)

func test() {
	meData, _ := me.ScrapeMagicEden()
	// var cl CollectionListing
	// err := json.Unmarshal(meData[0].CollectionListing, &cl)
	// if err != nil {
	// 	log.Printf("Error unmarshalling into json, err: %s", err)
	// 	// return topCollections
	// }
	fmt.Printf("PdaAddress: %s\n", meData[4].CollectionListing[0].TokenMint)
	// tokenmint := meData[4].CollectionListing[0].TokenMint
	// tensor.ScrapeTensor(tokenmint)
}

func main() {

	magicEdenData, collectionSymbolToTokenMintMap := me.ScrapeMagicEden()
	tensorData := tensor.ScrapeTensor(collectionSymbolToTokenMintMap)

	fmt.Println(magicEdenData)
	fmt.Println(tensorData)

}

// type CollectionData

// Issues
// 1) there is inconsistent return data from Magic Eden collecion stats
// 2) Tensor collection stats api keeps throwing errors
