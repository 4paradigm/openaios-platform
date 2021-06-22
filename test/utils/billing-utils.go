package utils

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/4paradigm/openaios-platform/test/openapi/billing/apigen/billingclient"
)

func getBillingClient() (*billingclient.Client, error) {
	billingServerURL := GetConfig("../test-cicd.toml").Env.BillingServerURL
	return billingclient.NewClient(billingServerURL)
}

func GetUserBalance(userID string) (*float64, error) {
	client, err := getBillingClient()
	if err != nil {
		return nil, errors.New("cannot connect to billing server.")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.GetAccountUseridBalance(ctx, userID)
	if err != nil {
		return nil, errors.New("cannot get user account balance")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("cannot get user account balance")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("cannot get user account balance")
	}
	var balance float64
	err = json.Unmarshal(body, &balance)
	if err != nil {
		return nil, errors.New("cannot get user account balance")
	}
	return &balance, nil
}

func UpdateUserBalance(userID string, balance float64) error {
	client, err := getBillingClient()
	if err != nil {
		return errors.New("cannot connect to billing server.")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var reqBody billingclient.PutAccountUseridJSONRequestBody
	reqBody.Balance = &balance
	resp, err := client.PutAccountUserid(ctx, userID, reqBody)
	if err != nil {
		return errors.New("cannot put user account balance")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("cannot put user account balance")
	}
	return nil
}

func RechargeUserBalance(userID string, balance float64) error {
	client, err := getBillingClient()
	if err != nil {
		return errors.New("cannot connect to billing server.")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var params billingclient.PostAccountUseridBalanceParams
	params.BuyBalance = balance
	resp, err := client.PostAccountUseridBalance(ctx, userID, &params)
	if err != nil {
		return errors.New("cannot put user account balance")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("cannot put user account balance")
	}
	return nil
}
