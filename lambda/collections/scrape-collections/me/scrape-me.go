package me

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func ScrapeCollections() {
	err := godotenv.Load("apikey.env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	API_KEY := os.Getenv("ME_API_KEY")

}
