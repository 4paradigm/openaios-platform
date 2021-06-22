package integration

import (
	// "time"

	"fmt"
	"time"

	"github.com/4paradigm/openaios-platform/test/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var testEnvName, bearerToken string

var _ = Describe("收费功能测试", func() {
	BeforeSuite(func() {
		bearerToken = utils.GetToken()
		userID, _ := utils.GetUserID(bearerToken)
		utils.UpdateUserBalance(userID, 50)
		envName, err := utils.CreateBasicEnvironment(bearerToken, "single-core", "env/pytorch", "20.12-py3", "public")
		Expect(err).To(BeNil())
		testEnvName = envName
		fmt.Println(fmt.Sprintf("test env %s is created!", testEnvName))
	})

	AfterSuite(func() {
		err := utils.DeleteEnvironment(bearerToken, testEnvName)
		if err != nil {
			fmt.Println(err)
			Expect(err).To(BeNil())
		}
		fmt.Println(fmt.Sprintf("test env %s is deleted!", testEnvName))
	})

	Describe("更新用户余额为50,并取得用户余额", func() {
		Context("等待一分钟后", func() {
			It("用户余额已经正确减去每分钟花费", func() {
				userID, _ := utils.GetUserID(bearerToken)
				utils.UpdateUserBalance(userID, 50.0)
				getBalance, _ := utils.GetUserBalance(userID)
				Expect(*getBalance).To(Equal(50.0))
				time.Sleep(60 * time.Second)
				cost, err := utils.GetUserCostPerMinute(bearerToken)
				Expect(err).To(BeNil())
				newGetBalance, _ := utils.GetUserBalance(userID)
				Expect(*newGetBalance).To(Equal(*getBalance - cost))
			})
		})
	})

	Describe("更新用户余额为-0.1", func() {
		Context("等待90秒后", func() {
			It("因为余额不足，用户之前创建的环境会进入Unknown状态，用户每分钟花费应变为0", func() {
				userID, _ := utils.GetUserID(bearerToken)
				utils.UpdateUserBalance(userID, -0.1)
				getBalance, _ := utils.GetUserBalance(userID)
				Expect(*getBalance).To(Equal(-0.1))
				time.Sleep(90 * time.Second)
				cost, err := utils.GetUserCostPerMinute(bearerToken)
				Expect(err).To(BeNil())
				Expect(cost).To(Equal(0.0))
			})
		})
	})
})
