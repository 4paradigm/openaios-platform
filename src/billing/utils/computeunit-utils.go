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

package utils

import (
	"context"
	"github.com/4paradigm/openaios-platform/src/internal/mongodb"
	"github.com/4paradigm/openaios-platform/src/internal/response"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

type ComputeunitInfo struct {
	ID    string   `bson:"id,omitempty"`
	Price *float64 `bson:"price,omitempty"`
}

type ComputeunitGroup struct {
	GroupName       string   `bson:"groupName,omitempty"`
	ComputeunitList []string `bson:"computeunitList,omitempty"`
	Avl             bool     `bson:"avl,omitempty"`
}

func ModifyUserComputeunitGroup(client *mongo.Client, userID string, groups []string) error {
	uniqueKey := AccountInfo{UserID: userID}
	modifyCount, err := mongodb.UpdateOneDocument(client, database, userColl, uniqueKey,
		mongodb.MongodbOperation{Operator: "$set", Document: AccountInfo{ComputeunitGroup: groups}})
	if err != nil {
		return errors.Wrap(err, response.GetRuntimeLocation())
	} else if modifyCount == 0 {
		log.Warnf("cannot find such user %s.", userID)
	}
	return nil
}

func GetUserComputeunitGroup(client *mongo.Client, userID string) ([]string, error) {
	uniqueKey := AccountInfo{UserID: userID}
	document := mongodb.FindOneDocument(client, database, userColl, uniqueKey)
	if document == nil {
		return nil, errors.New("Cannot find user. " + response.GetRuntimeLocation())
	} else {
		var accountInfo AccountInfo
		err := document.Decode(&accountInfo)
		if err != nil {
			return nil, errors.Wrap(err, response.GetRuntimeLocation())
		}
		return accountInfo.ComputeunitGroup, nil
	}
}

func GetComputeunitGroupList(client *mongo.Client) ([]string, error) {
	cursor, err := mongodb.FindDocuments(client, database, computeunitGroupColl, "")
	if err != nil {
		return nil, errors.Wrap(err, "cannot find documents "+response.GetRuntimeLocation())
	}
	var computeunitGroup ComputeunitGroup
	var groupList = []string{}
	for cursor.Next(context.Background()) {
		if err = cursor.Decode(&computeunitGroup); err != nil {
			log.Warn(errors.Wrap(err, response.GetRuntimeLocation()))
			continue
		}
		groupList = append(groupList, computeunitGroup.GroupName)
	}
	return groupList, nil
}

func GetComputeunitInfoByGroup(client *mongo.Client, groupName string) ([]map[string]interface{}, bool, error) {
	uniqueKey := ComputeunitGroup{GroupName: groupName}
	group := mongodb.FindOneDocument(client, database, computeunitGroupColl, uniqueKey)
	if group == nil {
		return nil, false, errors.New("cannot find such group " + response.GetRuntimeLocation())
	}
	var computeunitGroup ComputeunitGroup
	err := group.Decode(&computeunitGroup)
	if err != nil {
		return nil, false, errors.Wrap(err, response.GetRuntimeLocation())
	}
	var computeunitList []map[string]interface{}
	var uniqueKeyList []ComputeunitInfo
	if len(computeunitGroup.ComputeunitList) == 0 {
		return computeunitList, computeunitGroup.Avl, nil
	}
	for _, id := range computeunitGroup.ComputeunitList {
		uniqueKeyList = append(uniqueKeyList, ComputeunitInfo{ID: id})
	}
	cursor, err := mongodb.FindDocumentsByMultiKey(client, database, computeunitColl,
		"$or", uniqueKeyList)
	if err != nil {
		return nil, false, errors.Wrap(err, response.GetRuntimeLocation())
	}
	for cursor.Next(context.Background()) {
		var currentMap map[string]interface{}
		if err = bson.Unmarshal(cursor.Current, &currentMap); err != nil {
			log.Error(errors.Wrap(err, response.GetRuntimeLocation()))
			continue
		}
		//id := currentMap["id"].(string)
		delete(currentMap, "_id")
		computeunitList = append(computeunitList, currentMap)
	}
	return computeunitList, computeunitGroup.Avl, nil
}

func GetComputeunitInGroupByID(client *mongo.Client, groupName string, computeunitID string) (map[string]interface{}, bool, error) {
	uniqueKey := ComputeunitGroup{GroupName: groupName}
	group := mongodb.FindOneDocument(client, database, computeunitGroupColl, uniqueKey)
	if group == nil {
		return nil, false, errors.New("cannot find such group " + response.GetRuntimeLocation())
	}
	var computeunitGroup ComputeunitGroup
	err := group.Decode(&computeunitGroup)
	if err != nil {
		return nil, false, errors.Wrap(err, response.GetRuntimeLocation())
	}
	for _, id := range computeunitGroup.ComputeunitList {
		if id == computeunitID {
			oneUniqueKey := ComputeunitInfo{ID: computeunitID}
			document := mongodb.FindOneDocument(client, database, computeunitColl, oneUniqueKey)
			if document == nil {
				return nil, false, errors.New("cannot find such computeunit. " + response.GetRuntimeLocation())
			}
			var result map[string]interface{}
			if err = document.Decode(&result); err != nil {
				return nil, false, errors.Wrap(err, response.GetRuntimeLocation())
			}
			delete(result, "_id")
			return result, computeunitGroup.Avl, nil
		}
	}
	return nil, false, errors.New("cannot find such computeunit. " + response.GetRuntimeLocation())
}

func GetPriceMap(client *mongo.Client) (map[string]float64, error) {
	cursor, err := mongodb.FindDocuments(client, database, computeunitColl, "")
	if err != nil {
		return nil, errors.Wrap(err, "Cannot get price list "+response.GetRuntimeLocation())
	}
	var computeunitInfo ComputeunitInfo
	var priceMap = map[string]float64{}
	for cursor.Next(context.Background()) {
		if err = cursor.Decode(&computeunitInfo); err != nil {
			log.Warn(errors.Wrap(err, response.GetRuntimeLocation()))
			continue
		}
		if computeunitInfo.Price != nil {
			priceMap[computeunitInfo.ID] = *computeunitInfo.Price
		}
	}
	return priceMap, nil
}

func In(haystack interface{}, needle interface{}) (bool, error) {
	sVal := reflect.ValueOf(haystack)
	kind := sVal.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		for i := 0; i < sVal.Len(); i++ {
			if sVal.Index(i).Interface() == needle {
				return true, nil
			}
		}
		return false, nil
	}
	return false, errors.New("ErrUnSupportHaystack " + response.GetRuntimeLocation())
}
