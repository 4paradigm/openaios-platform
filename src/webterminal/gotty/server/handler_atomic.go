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

package server

import (
	"sync"
	"time"
)

type counter struct {
	duration    time.Duration
	zeroTimer   *time.Timer
	wg          sync.WaitGroup
	connections int
	mutex       sync.Mutex
}

func newCounter(duration time.Duration) *counter {
	zeroTimer := time.NewTimer(duration)

	// when duration is 0, drain the expire event here
	// so that user will never get the event.
	if duration == 0 {
		<-zeroTimer.C
	}

	return &counter{
		duration:  duration,
		zeroTimer: zeroTimer,
	}
}

func (counter *counter) add(n int) int {
	counter.mutex.Lock()
	defer counter.mutex.Unlock()

	if counter.duration > 0 {
		counter.zeroTimer.Stop()
	}
	counter.wg.Add(n)
	counter.connections += n

	return counter.connections
}

func (counter *counter) done() int {
	counter.mutex.Lock()
	defer counter.mutex.Unlock()

	counter.connections--
	counter.wg.Done()
	if counter.connections == 0 && counter.duration > 0 {
		counter.zeroTimer.Reset(counter.duration)
	}

	return counter.connections
}

func (counter *counter) count() int {
	counter.mutex.Lock()
	defer counter.mutex.Unlock()

	return counter.connections
}

func (counter *counter) wait() {
	counter.wg.Wait()
}

func (counter *counter) timer() *time.Timer {
	return counter.zeroTimer
}
