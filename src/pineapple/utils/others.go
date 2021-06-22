package utils

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
)

const codes = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const codesLength = len(codes)

func RandCode(n int) string {
	b := make([]byte, n)

	for i := range b {
		b[i] = codes[rand.Intn(codesLength)]
	}
	return string(b)
}

func GetEnvDefault(env, defaultValue string) string {
	v, exist := os.LookupEnv(env)
	if exist {
		return v
	} else {
		return defaultValue
	}
}

func GetRuntimeLocation() string {
	_, fn, line, _ := runtime.Caller(1)
	return fmt.Sprintf("%s:%d", fn, line)
}
