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
	"context"
)

// RunOptions holds a set of configurations for Server.Run().
type RunOptions struct {
	gracefullCtx context.Context
}

// RunOption is an option of Server.Run().
type RunOption func(*RunOptions)

// WithGracefullContext accepts a context to shutdown a Server
// with care for existing client connections.
func WithGracefullContext(ctx context.Context) RunOption {
	return func(options *RunOptions) {
		options.gracefullCtx = ctx
	}
}
