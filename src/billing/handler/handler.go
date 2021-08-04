package handler

import (
	"github.com/labstack/gommon/log"
	"github.com/4paradigm/openaios-platform/src/billing/conf"
	"github.com/4paradigm/openaios-platform/src/billing/utils"
	"github.com/4paradigm/openaios-platform/src/internal/mongodb"
)

type Handler struct{}

var mongodbUrl = conf.GetMongodbUrl()

func NewHandler() Handler {
	return Handler{}
}

func InitBillingServer() {
	// init mongodb collection
	client, err := mongodb.GetMongodbClient(mongodbUrl)
	defer mongodb.KillMongodbClient(client)
	if err != nil {
		log.Error(err.Error())
		return
	}
	if err = utils.InitColl(client); err != nil {
		log.Error(err.Error())
	}
}
