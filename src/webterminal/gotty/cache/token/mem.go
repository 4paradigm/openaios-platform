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

package token

import (
	"log"
	"time"

	"github.com/patrickmn/go-cache"
)

//MemCache use memory to store token and TtyParameter
type MemCache struct {
	cache *cache.Cache
}

//NewMemCache new MemCache
func NewMemCache() *MemCache {
	return &MemCache{
		cache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

//Get token param from memory
func (r *MemCache) Get(token string) *TtyParameter {
	obj, exit := r.cache.Get(token)
	if !exit {
		return nil
	}
	param, ok := obj.(TtyParameter)
	if ok {
		return &param
	}
	log.Printf("get token %s from mem obj is not tty param", token)
	return nil
}

//Delete token from memory
func (r *MemCache) Delete(token string) error {
	r.cache.Delete(token)
	return nil
}

//Add token to memory
func (r *MemCache) Add(token string, param *TtyParameter, d time.Duration) error {
	return r.cache.Add(token, *param, d)
}
