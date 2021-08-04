package helm

import (
	"github.com/pkg/errors"
	"github.com/4paradigm/openaios-platform/src/internal/billingclient"
	"github.com/4paradigm/openaios-platform/src/pineapple/conf"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils"
)

type IPineappleInfo interface {
	SetValues(values map[string]interface{})
	GetName() string
	GetUserId() string
	GetPrefix() string
	GetValues() map[string]interface{}
	CreateChartValues() (map[string]interface{}, error)
}

type PineappleInfo struct {
	Name   string
	UserId string
	Prefix string
	Values map[string]interface{}
}

func NewPineappleInfo(name string, userId string, prefix string) (*PineappleInfo, error) {
	billingClient, err := billingclient.GetBillingClient(conf.GetBillingServerURL())
	if err != nil {
		return nil, errors.Wrap(err, utils.GetRuntimeLocation())
	}
	userBalance, err := billingclient.GetUserBalance(billingClient, userId)
	if err != nil {
		return nil, errors.WithMessage(err, "get GetUserBalance error: ")
	}
	if *userBalance <= 0.0 {
		return nil, errors.New("User does not have enough balance." + utils.GetRuntimeLocation())
	}
	pineappleInfo := PineappleInfo{
		Name:   name,
		UserId: userId,
		Prefix: prefix,
		Values: nil,
	}
	return &pineappleInfo, nil
}

func (p *PineappleInfo) GetName() string {
	return p.Name
}

func (p *PineappleInfo) GetUserId() string {
	return p.UserId
}

func (p *PineappleInfo) GetPrefix() string {
	return p.Prefix
}

func (p *PineappleInfo) GetValues() map[string]interface{} {
	return p.Values
}

func (p *PineappleInfo) SetValues(values map[string]interface{}) {
	p.Values = values
}

func (p *PineappleInfo) CreateChartValues() (map[string]interface{}, error) {
	var chartValues map[string]interface{}
	chartValues = p.Values

	// Create appConf Values
	appConf, err := conf.GetAppConf()
	if err != nil {
		errors.WithMessage(err, "GetAppConf error: ")
	}
	chartValues["appConf"] = appConf

	return chartValues, nil
}
