package utils

import (
	"math/rand"
	"time"

	"github.com/4paradigm/openaios-platform/test/openapi/main/apigen/restclient"
)

func GenRandEnvName(n int, m int) restclient.EnvironmentName {
	rand.Seed(time.Now().Unix())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890")
	b := make([]rune, n+m+1)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	b[0] = rune('a')
	b[n] = rune('-')
	return restclient.EnvironmentName(b)
}
