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
