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

package handler

import (
	"context"
	"encoding/json"
	"github.com/4paradigm/openaios-platform/src/internal/billingclient"
	"github.com/4paradigm/openaios-platform/src/internal/mongodb"
	"github.com/4paradigm/openaios-platform/src/internal/response"
	"github.com/4paradigm/openaios-platform/src/pineapple/apigen"
	"github.com/4paradigm/openaios-platform/src/pineapple/conf"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type UserProfile struct {
	UserID  string                 `bson:"userID,omitempty"`
	RegTime time.Time              `bson:"regTime,omitempty"`
	Profile map[string]interface{} `bson:"profile,omitempty"`
	Inviter string                 `bson:"inviter,omitempty"`
}

type CompetitionInfo struct {
	Name             string     `bson:"name,omitempty"`
	ID               string     `bson:"id,omitempty"`
	Beginning        *time.Time `bson:"beginning,omitempty"`
	Deadline         *time.Time `bson:"deadline,omitempty"`
	ComputeunitGroup string     `bson:"computeunitGroup,omitempty"`
	BaseParticipants int64      `bson:"baseParticipants,omitempty"`
}

var mongodbCompetitionsColl = "competitions"

func (handler *Handler) GetCompetitionCompetitionID(ctx echo.Context, competitionID string) error {
	userID := ctx.Get("userID").(string)
	mongodbURL := conf.GetMongodbURL()
	database := conf.GetMongodbDatabase()
	client, err := mongodb.GetMongodbClient(mongodbURL)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "check competition failed.").SetInternal(
			errors.Wrap(err, "cannot connect to mongodb. "+utils.GetRuntimeLocation()))
	}
	defer mongodb.KillMongodbClient(client)

	document := mongodb.FindOneDocument(client, database, competitionID, UserProfile{UserID: userID})
	if document == nil {
		return ctx.String(http.StatusOK, strconv.FormatBool(false))
	}
	var competition CompetitionInfo
	err = document.Decode(&competition)
	if err != nil {
		ctx.Logger().Warn(errors.Wrap(err, utils.GetRuntimeLocation()))
		return ctx.String(http.StatusOK, strconv.FormatBool(false))
	}
	return ctx.String(http.StatusOK, strconv.FormatBool(true))
}

func (handler *Handler) PostCompetitionCompetitionID(ctx echo.Context, competitionID string,
	params apigen.PostCompetitionCompetitionIDParams) error {
	userID := ctx.Get("userID").(string)
	mongodbURL := conf.GetMongodbURL()
	database := conf.GetMongodbDatabase()
	client, err := mongodb.GetMongodbClient(mongodbURL)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "join competition failed.").SetInternal(
			errors.Wrap(err, "cannot connect to mongodb. "+utils.GetRuntimeLocation()))
	}
	defer mongodb.KillMongodbClient(client)

	// decode body
	var body = map[string]interface{}{}
	err = json.NewDecoder(ctx.Request().Body).Decode(&body)
	if err != nil {
		ctx.Logger().Warn(errors.Wrap(err, "cannot get request body. "+utils.GetRuntimeLocation()))
	}

	document := mongodb.FindOneDocument(client, database, mongodbCompetitionsColl,
		CompetitionInfo{ID: competitionID})
	if document == nil {
		return response.BadRequestWithMessage(ctx, "competition not exists.")
	}
	var competition CompetitionInfo
	err = document.Decode(&competition)
	if err != nil {
		return response.BadRequestWithMessage(ctx, "competition not exists.")
	}
	if competition.Beginning != nil && competition.Beginning.After(time.Now()) {
		return response.BadRequestWithMessage(ctx, "competition not begin")
	}
	if competition.Deadline != nil && competition.Deadline.Before(time.Now()) {
		return response.BadRequestWithMessage(ctx, "competition is over")
	}

	// set index
	err = mongodb.CreateUniqueIndex(client, database, competitionID, "userID")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "cannot join competition").SetInternal(
			errors.Wrap(err, response.GetRuntimeLocation()))
	}
	err = mongodb.CreateIndex(client, database, competitionID, "inviter")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "cannot join competition").SetInternal(
			errors.Wrap(err, response.GetRuntimeLocation()))
	}

	// insert user information
	var inviter string
	if params.Inviter == nil || *params.Inviter == userID {
		inviter = ""
	} else {
		inviter = *params.Inviter
	}
	_, err = mongodb.InsertOneDocument(client, database, competitionID,
		UserProfile{UserID: userID, Profile: body, RegTime: time.Now(), Inviter: inviter})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key error") {
			return response.BadRequestWithMessage(ctx, "already joined this competition.")
		} else {
			return echo.NewHTTPError(
				http.StatusInternalServerError, "cannot join competition").SetInternal(
				errors.Wrap(err, "cannot insert document "+utils.GetRuntimeLocation()))
		}
	}

	// insert user to computeunit group
	billingClient, err := billingclient.GetBillingClient(conf.GetBillingServerURL())
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "cannot connect to billing server").SetInternal(
			errors.Wrap(err, utils.GetRuntimeLocation()))
	}
	err = billingclient.AddComputeunitGroupToUser(billingClient, userID, competition.ComputeunitGroup)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "cannot join competition").SetInternal(
			errors.Wrap(err, "cannot add computeunit group "+utils.GetRuntimeLocation()))
	}
	return response.StatusOKNoContent(ctx)
}

func (handler *Handler) GetCompetition(ctx echo.Context, params apigen.GetCompetitionParams) error {
	mongodbURL := conf.GetMongodbURL()
	database := conf.GetMongodbDatabase()
	client, err := mongodb.GetMongodbClient(mongodbURL)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "check competition failed.").SetInternal(
			errors.Wrap(err, "cannot connect to mongodb. "+utils.GetRuntimeLocation()))
	}
	defer mongodb.KillMongodbClient(client)

	cursor, err := mongodb.FindDocuments(client, database, mongodbCompetitionsColl, "")
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "cannot find competitions.").SetInternal(
			errors.Wrap(err, "cannot find competition "+utils.GetRuntimeLocation()))
	}

	var competitionInfo CompetitionInfo
	var competitionList []apigen.CompetitionInfo
	billingClient, err := billingclient.GetBillingClient(conf.GetBillingServerURL())
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "cannot connect to billing server").SetInternal(
			errors.Wrap(err, utils.GetRuntimeLocation()))
	}
	for cursor.Next(context.Background()) {
		if err = cursor.Decode(&competitionInfo); err != nil {
			ctx.Logger().Warn(errors.Wrap(err, "cannot decode document "+utils.GetRuntimeLocation()))
			continue
		}
		beginning := *competitionInfo.Beginning
		deadline := *competitionInfo.Deadline
		id := competitionInfo.ID
		name := competitionInfo.Name
		computeunitList := []apigen.ComputeUnitSpec{}
		avl := true
		computeunitGroup, err := billingclient.GetComputeUnitListByGroupName(billingClient, competitionInfo.ComputeunitGroup)
		if err != nil {
			ctx.Logger().Warn(errors.Wrap(err, "cannot get computeunit group "+utils.GetRuntimeLocation()))
		}
		if computeunitGroup == nil {
			avl = false
			if !params.Beginning.IsZero() && competitionInfo.Deadline != nil && params.Beginning.After(*competitionInfo.Deadline) {
				continue
			}
			if params.Beginning.IsZero() && time.Now().After(*competitionInfo.Deadline) {
				continue
			}
			if !params.End.IsZero() && competitionInfo.Beginning != nil && params.End.Before(*competitionInfo.Beginning) {
				continue
			}
		}
		for _, unit := range computeunitGroup {
			ID := apigen.ComputeUnitId(*unit.Id)
			computeunitList = append(computeunitList, apigen.ComputeUnitSpec{
				Id: &ID, Price: unit.Price, Description: unit.Description})
		}
		participant, err := mongodb.CountDocuments(client, database, id, "")
		if err != nil {
			ctx.Logger().Warn(errors.Wrap(err, "cannot get competition participant "+utils.GetRuntimeLocation()))
		}
		participant += competitionInfo.BaseParticipants

		competitionList = append(competitionList,
			apigen.CompetitionInfo{Id: &id, Name: &name, Beginning: &beginning, Avl: &avl, Participant: &participant,
				Deadline: &deadline, ComputingResource: &computeunitList})
	}
	return ctx.JSON(http.StatusOK, competitionList)
}

func (handler *Handler) GetCompetitionCompetitionIDInvitation(ctx echo.Context, competitionID string) error {
	mongodbURL := conf.GetMongodbURL()
	database := conf.GetMongodbDatabase()
	client, err := mongodb.GetMongodbClient(mongodbURL)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "check competition failed.").SetInternal(
			errors.Wrap(err, "cannot connect to mongodb. "+utils.GetRuntimeLocation()))
	}
	defer mongodb.KillMongodbClient(client)

	userID := ctx.Get("userID").(string)
	operator := mongodb.ComparisonQueryOperator{Operation: "$eq", Value: userID}
	cursor, err := mongodb.FindDocuments(client, database, competitionID, "inviter", operator)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "check competition failed.").SetInternal(
			errors.Wrap(err, "find document failed "+utils.GetRuntimeLocation()))
	}
	var count int64 = 0
	var userProfile UserProfile
	for cursor.Next(context.Background()) {
		err = cursor.Decode(&userProfile)
		if err != nil {
			continue
		}
		count = count + 1
	}
	return ctx.String(http.StatusOK, strconv.FormatInt(count, 10))
}
