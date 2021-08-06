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

export class ConnectionFactory {
    url: string;
    protocols: string[];

    constructor(url: string, protocols: string[]) {
        this.url = url;
        this.protocols = protocols;
    };

    create(): Connection {
        return new Connection(this.url, this.protocols);
    };
}

export class Connection {
    bare: WebSocket;


    constructor(url: string, protocols: string[]) {
        this.bare = new WebSocket(url, protocols);
    }

    open() {
        // nothing todo for websocket
    };

    close() {
        this.bare.close();
    };

    send(data: string) {
        this.bare.send(data);
    };

    isOpen(): boolean {
        if (this.bare.readyState == WebSocket.CONNECTING ||
            this.bare.readyState == WebSocket.OPEN) {
            return true
        }
        return false
    }

    onOpen(callback: () => void) {
        this.bare.onopen = (event) => {
            callback();
        }
    };

    onReceive(callback: (data: string) => void) {
        this.bare.onmessage = (event) => {
            callback(event.data);
        }
    };

    onClose(callback: () => void) {
        this.bare.onclose = (event) => {
            callback();
        };
    };
}
