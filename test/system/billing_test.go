package system

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/4paradigm/openaios-platform/test/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("billing测试", func() {
	Describe("使用PUT接口更新用户余额为一个随机浮点数", func() {
		Context("当成功更新后", func() {
			It("GET接口取得用户的余额应为更新后的值", func() {
				rand.Seed(time.Now().UTC().UnixNano())
				randBalance := rand.Float64() * 10000
				userID, _ := utils.GetUserID(utils.GetToken())
				fmt.Println(userID)
				getBalance, _ := utils.GetUserBalance(userID)
				utils.UpdateUserBalance(userID, randBalance)
				getBalance, _ = utils.GetUserBalance(userID)
				Expect(*getBalance).To(Equal(randBalance))
			})
		})
	})
	Describe("使用POST接口新增用户余额为一个随机浮点数", func() {
		Context("当成功更新后", func() {
			It("GET接口取得用户的余额应为充值后的值", func() {
				rand.Seed(time.Now().UTC().UnixNano())
				randBalance := rand.Float64() * 10000
				userID, _ := utils.GetUserID(utils.GetToken())
				balance, _ := utils.GetUserBalance(userID)
				utils.RechargeUserBalance(userID, randBalance)
				getBalance, _ := utils.GetUserBalance(userID)
				newBalance := *balance + randBalance
				Expect(*getBalance).To(Equal(newBalance))
			})
		})
	})
})
