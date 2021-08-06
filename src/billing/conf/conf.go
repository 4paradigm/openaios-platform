/*
 * Copyright © 2021 peizhaoyou <peizhaoyou@4paradigm.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package conf implements methods to get mongodb configs from flags or environment variables.
package conf

import (
	"flag"
	"os"
)

var (
	mongodbURL = flag.String("mongodb-url", os.Getenv("PINEAPPLE_MONGODB_URL"),
		"mongodb url")
	mongodbDatabase = flag.String("mongodb-database", os.Getenv("PINEAPPLE_MONGODB_DATABASE"),
		"mongodb database")
)

// GetmongodbURL returns mongodbURL string
func GetmongodbURL() string {
	return *mongodbURL
}

// GetMongodbDatabase returns mongodbDatabase string
func GetMongodbDatabase() string {
	return *mongodbDatabase
}
