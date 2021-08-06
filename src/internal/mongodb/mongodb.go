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

package mongodb

import (
	"context"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"github.com/4paradigm/openaios-platform/src/internal/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"strings"
	"time"
)

type MongodbOperation struct {
	Operator string
	Document interface{}
}

type ComparisonQueryOperator struct {
	Operation string
	Value     interface{}
}

func GetMongodbClient(mongodbUrl string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodbUrl))
	if err != nil {
		return nil, err
	}
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Error(err.Error())
	}
	log.Debug("Successfully connected and pinged.")
	return client, nil
}

func KillMongodbClient(client *mongo.Client) {
	if client == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Disconnect(ctx); err != nil {
		log.Error(err.Error())
	}
}

func CreateIndex(client *mongo.Client, database string, collection string, key string) error {
	if client == nil {
		return errors.New("mongodb client is nil " + response.GetRuntimeLocation())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coll := client.Database(database).Collection(collection)
	mod := mongo.IndexModel{Keys: bson.M{key: bsonx.Int32(1)}}
	_, err := coll.Indexes().CreateOne(ctx, mod)
	if err != nil {
		return errors.Wrap(err, response.GetRuntimeLocation())
	}
	return nil
}

func CreateUniqueIndex(client *mongo.Client, database string, collection string, keys ...string) error {
	if client == nil {
		return errors.New("mongodb client is nil.")
	}
	db := client.Database(database).Collection(collection)
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	indexView := db.Indexes()
	keysDoc := bsonx.Doc{}

	for _, key := range keys {
		if strings.HasPrefix(key, "-") {
			keysDoc = keysDoc.Append(strings.TrimLeft(key, "-"), bsonx.Int32(-1))
		} else {
			keysDoc = keysDoc.Append(key, bsonx.Int32(1))
		}
	}

	result, err := indexView.CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    keysDoc,
			Options: options.Index().SetUnique(true),
		},
		opts,
	)
	if result == "" || err != nil {
		log.Error("EnsureIndex error", err)
		return errors.New("EnsureIndex error")
	}
	return nil
}

func CountDocuments(client *mongo.Client, database string, collection string,
	key string, operators ...ComparisonQueryOperator) (int64, error) {
	if client == nil {
		return 0, errors.New("mongodb client is nil.")
	}
	db := client.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	operatorList := bson.M{}
	for _, operator := range operators {
		operatorList[operator.Operation] = operator.Value
	}
	if key != "" {
		return db.CountDocuments(ctx, bson.M{key: operatorList})
	} else {
		return db.CountDocuments(ctx, bson.D{})
	}
}

func InsertOneDocument(client *mongo.Client, database string, collection string,
	document interface{}) (*mongo.InsertOneResult, error) {
	if client == nil {
		return nil, errors.New("mongodb client is nil.")
	}
	db := client.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := db.InsertOne(ctx, document)
	return result, err
}

func DeleteOneDocument(client *mongo.Client, database string,
	collection string, uniqueKey interface{}) error {
	if client == nil {
		return errors.New("mongodb client is nil.")
	}
	db := client.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := db.DeleteOne(ctx, uniqueKey)
	return err
}

// warning: one operator only appear once
// return modify_count
func UpdateOneDocument(client *mongo.Client, database string, collection string,
	uniqueKey interface{}, operations ...MongodbOperation) (int64, error) {
	if client == nil {
		return 0, errors.New("mongodb client is nil.")
	}
	operationList := bson.M{}
	for _, operation := range operations {
		operationList[operation.Operator] = operation.Document
	}
	db := client.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	updateResult, err := db.UpdateOne(ctx, uniqueKey, operationList)
	if err != nil {
		return 0, err
	} else {
		return updateResult.ModifiedCount, nil
	}
}

func FindOneDocument(client *mongo.Client, database string,
	collection string, uniqueKey interface{}) *mongo.SingleResult {
	if client == nil {
		return nil
	}
	db := client.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return db.FindOne(ctx, uniqueKey)
}

// warning: only support single key
func FindDocuments(client *mongo.Client, database string, collection string,
	key string, operators ...ComparisonQueryOperator) (*mongo.Cursor, error) {
	if client == nil {
		return nil, errors.New("mongodb client is nil.")
	}
	operatorList := bson.M{}
	for _, operator := range operators {
		operatorList[operator.Operation] = operator.Value
	}
	db := client.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if key != "" {
		return db.Find(ctx, bson.M{key: operatorList})
	} else {
		return db.Find(ctx, bson.D{})
	}
}

func FindDocumentsByMultiKey(client *mongo.Client, database string, collection string,
	condition string, keys interface{}) (*mongo.Cursor, error) {
	if client == nil {
		return nil, errors.New("mongodb client is nil.")
	}
	db := client.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return db.Find(ctx, bson.M{condition: keys})
}

// warning: different operations cannot have same operator
func UpdateOrInsertOneDocument(client *mongo.Client, database string, collection string,
	uniqueKey interface{}, operations ...MongodbOperation) error {
	if client == nil {
		return errors.New("mongodb client is nil.")
	}
	operationList := bson.M{}
	for _, operation := range operations {
		operationList[operation.Operator] = operation.Document
	}
	db := client.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := db.UpdateOne(ctx, uniqueKey, operationList, options.Update().SetUpsert(true))
	return err
}

func InsertedIdToObjectId(insertedId interface{}) (string, error) {
	if oid, ok := insertedId.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	} else {
		return "", errors.New("Cannot convert inserted id to object id.")
	}
}

//func CheckCollectionExists(client *mongo.Client, database string, collection string) (bool, error) {
//	if client == nil {
//		return false, errors.New("mongodb client is nil.")
//	}
//	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
//	defer cancel()
//
//	names, err := client.Database(database).ListCollectionNames(ctx, bson.D{})
//	if err != nil {
//		return false, err
//	}
//	// Simply search in the names slice, e.g.
//	for _, name := range names {
//		if name == collection {
//			return true, nil
//		}
//	}
//	return false, nil
//}
