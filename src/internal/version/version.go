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

// Package version provides version functions
package version

import (
    "flag"
    "fmt"
    "os"
)

var (
    version string
    printVersion = flag.Bool("version", false, "print version")
)

func GetVersion() string {
    return version
}

// CheckVersionFlag if version flag == true, then print version info and exit
func CheckVersionFlag() {
    if *printVersion {
        fmt.Printf("version: %v\n", GetVersion())
        os.Exit(0)
    }
}
