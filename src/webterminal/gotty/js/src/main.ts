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

import { Hterm } from "./hterm";
import { Xterm } from "./xterm";
import { Terminal, WebTTY, protocols } from "./webtty";
import { ConnectionFactory } from "./websocket";

// @TODO remove these
declare var gotty_auth_token: string;
declare var gotty_term: string;

const elem = document.getElementById("terminal")

if (elem !== null) {
    var term: Terminal;
    if (gotty_term == "hterm") {
        term = new Hterm(elem);
    } else {
        term = new Xterm(elem);
    }
    const httpsEnabled = window.location.protocol == "https:";
    const url = (httpsEnabled ? 'wss://' : 'ws://') + window.location.host + window.location.pathname + 'ws';
    const args = window.location.search;
    const factory = new ConnectionFactory(url, protocols);
    const wt = new WebTTY(term, factory, args, gotty_auth_token);
    const closer = wt.open();

    window.addEventListener("unload", () => {
        closer();
        term.close();
    });
};
