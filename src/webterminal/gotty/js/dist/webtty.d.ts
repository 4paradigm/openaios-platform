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

export declare const protocols: string[];
export declare const msgInputUnknown = "0";
export declare const msgInput = "1";
export declare const msgPing = "2";
export declare const msgResizeTerminal = "3";
export declare const msgUnknownOutput = "0";
export declare const msgOutput = "1";
export declare const msgPong = "2";
export declare const msgSetWindowTitle = "3";
export declare const msgSetPreferences = "4";
export declare const msgSetReconnect = "5";
export interface Terminal {
    info(): {
        columns: number;
        rows: number;
    };
    output(data: string): void;
    showMessage(message: string, timeout: number): void;
    removeMessage(): void;
    setWindowTitle(title: string): void;
    setPreferences(value: object): void;
    onInput(callback: (input: string) => void): void;
    onResize(callback: (colmuns: number, rows: number) => void): void;
    reset(): void;
    deactivate(): void;
    close(): void;
}
export interface Connection {
    open(): void;
    close(): void;
    send(data: string): void;
    isOpen(): boolean;
    onOpen(callback: () => void): void;
    onReceive(callback: (data: string) => void): void;
    onClose(callback: () => void): void;
}
export interface ConnectionFactory {
    create(): Connection;
}
export declare class WebTTY {
    term: Terminal;
    connectionFactory: ConnectionFactory;
    args: string;
    authToken: string;
    reconnect: number;
    constructor(term: Terminal, connectionFactory: ConnectionFactory, args: string, authToken: string);
    open(): () => void;
}
