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

package webtty

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// Option is an option for WebTTY.
type Option func(*WebTTY) error

// WithPermitWrite sets a WebTTY to accept input from slaves.
func WithPermitWrite() Option {
	return func(wt *WebTTY) error {
		wt.permitWrite = true
		return nil
	}
}

// WithFixedColumns sets a fixed width to TTY master.
func WithFixedColumns(columns int) Option {
	return func(wt *WebTTY) error {
		wt.columns = columns
		return nil
	}
}

// WithFixedRows sets a fixed height to TTY master.
func WithFixedRows(rows int) Option {
	return func(wt *WebTTY) error {
		wt.rows = rows
		return nil
	}
}

// WithWindowTitle sets the default window title of the session
func WithWindowTitle(windowTitle []byte) Option {
	return func(wt *WebTTY) error {
		wt.windowTitle = windowTitle
		return nil
	}
}

// WithReconnect enables reconnection on the master side.
func WithReconnect(timeInSeconds int) Option {
	return func(wt *WebTTY) error {
		wt.reconnect = timeInSeconds
		return nil
	}
}

// WithMasterPreferences sets an optional configuration of master.
func WithMasterPreferences(preferences interface{}) Option {
	return func(wt *WebTTY) error {
		prefs, err := json.Marshal(preferences)
		if err != nil {
			return errors.Wrapf(err, "failed to marshal preferences as JSON")
		}
		wt.masterPrefs = prefs
		return nil
	}
}
