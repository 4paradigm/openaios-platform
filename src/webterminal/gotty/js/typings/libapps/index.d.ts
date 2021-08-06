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

export declare namespace hterm {
    export class Terminal {
        io: IO;
        onTerminalReady: () => void;

        constructor();
        getPrefs(): Prefs;
        decorate(HTMLElement);
        installKeyboard(): void;
        uninstallKeyboard(): void;
        setWindowTitle(title: string): void;
        reset(): void;
        softReset(): void;
    }

    export class IO {
        writeUTF8: ((data: string) => void);
        writeUTF16: ((data: string) => void);
        onVTKeystroke: ((data: string) => void) | null;
        sendString: ((data: string) => void) | null;
        onTerminalResize: ((columns: number, rows: number) => void) | null;

        push(): IO;
        writeUTF(data: string);
        showOverlay(message: string, timeout: number | null);
    }

    export class Prefs {
        set(key: string, value: string): void;
    }

    export var defaultStorage: lib.Storage;
}

export declare namespace lib {
    export interface Storage {
    }

    export interface Memory {
        new (): Storage;
        Memory(): Storage
    }

    export var Storage: {
        Memory: Memory
    }

    export class UTF8Decoder {
        decode(str: string)
    }
}
