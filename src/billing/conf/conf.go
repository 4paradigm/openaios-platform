package conf

import (
	"flag"
	"os"
)

var (
	mongodbUrl = flag.String("mongodb-url", os.Getenv("PINEAPPLE_MONGODB_URL"),
		"mongodb url")
	mongodbDatabase = flag.String("mongodb-database", os.Getenv("PINEAPPLE_MONGODB_DATABASE"),
		"mongodb database")
)

func GetMongodbUrl() string {
	return *mongodbUrl
}

func GetMongodbDatabase() string {
	return *mongodbDatabase
}
