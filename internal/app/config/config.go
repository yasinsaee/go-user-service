package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

var MongoClient *mongo.Client
var MongoDB *mongo.Database

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ No .env file found. Using system env variables.")
	}
}

// مقدار env با fallback
func GetEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
