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
