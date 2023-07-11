package me

type Collection struct {
	Symbol      string  `json:"symbol"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	FloorPrice  int64   `json:"floorPrice"`
	VolumeAll   float64 `json:"volumeAll"`
}

type CollectionStats struct {
	FloorPrice   int     `json:"floorPrice"`
	ListedCount  int     `json:"listedCount"`
	AvgPrice24Hr float64 `json:"avgPrice24hr"`
	VolumeAll    float64 `json:"volumeAll"`
}

type CollectionHolderStats struct {
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

type CollectionListing struct {
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

type MagicEdenCollectionData struct {
	Collection      Collection      `json:"collection"`
	CollectionStats CollectionStats `json:"collectionStats"`
	// CollectionHolderStats CollectionHolderStats `json:"collectionHolderStats"` // The data returned is different for some collections eg mad_lads
	CollectionListing []CollectionListing `json:"collectionListing"`
}
