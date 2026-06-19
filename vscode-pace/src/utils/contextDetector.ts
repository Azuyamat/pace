import * as vscode from 'vscode';

export interface BlockFrame {
    type: string;
    name?: string;
}

export interface DocumentContext {
    inTaskBlock: boolean;
    inHookBlock: boolean;
    inEnvBlock: boolean;
    inArgsBlock: boolean;
    isTopLevel: boolean;
    blockStack: BlockFrame[];
    currentTaskName?: string;
}

const namedBlocks = new Set(['task', 'hook', 'env', 'args']);

export class ContextDetector {
    static detectContext(document: vscode.TextDocument, position: vscode.Position): DocumentContext {
        const blockStack = this.getBlockStack(document, position);
        const currentBlock = blockStack[blockStack.length - 1];
        const currentTask = [...blockStack].reverse().find(block => block.type === 'task');

        return {
            inTaskBlock: !!currentBlock && currentBlock.type === 'task',
            inHookBlock: !!currentBlock && currentBlock.type === 'hook',
            inEnvBlock: !!currentBlock && currentBlock.type === 'env',
            inArgsBlock: !!currentBlock && currentBlock.type === 'args',
            isTopLevel: blockStack.length === 0,
            blockStack,
            currentTaskName: currentTask?.name
        };
    }

    static getBlockStack(document: vscode.TextDocument, position: vscode.Position): BlockFrame[] {
        const stack: BlockFrame[] = [];
        let pendingBlock: BlockFrame | undefined;
        let expectBlockName = false;
        let inString = false;
        let inMultilineString = false;
        let escapeNext = false;

        const handleToken = (token: string) => {
            if (!token) {
                return;
            }

            if (expectBlockName && pendingBlock && (pendingBlock.type === 'task' || pendingBlock.type === 'hook')) {
                pendingBlock.name = token;
                expectBlockName = false;
                return;
            }

            if (namedBlocks.has(token)) {
                pendingBlock = { type: token };
                expectBlockName = token === 'task' || token === 'hook';
            }
        };

        for (let lineNumber = 0; lineNumber <= position.line; lineNumber++) {
            const fullLine = document.lineAt(lineNumber).text;
            const line = lineNumber === position.line ? fullLine.substring(0, position.character) : fullLine;
            let token = '';

            const flushToken = () => {
                handleToken(token);
                token = '';
            };

            for (let i = 0; i < line.length; i++) {
                const char = line[i];
                const nextThree = line.substring(i, i + 3);

                if (inMultilineString) {
                    if (nextThree === '"""') {
                        inMultilineString = false;
                        i += 2;
                    }
                    continue;
                }

                if (inString) {
                    if (escapeNext) {
                        escapeNext = false;
                    } else if (char === '\\') {
                        escapeNext = true;
                    } else if (char === '"') {
                        inString = false;
                    }
                    continue;
                }

                if (char === '#') {
                    flushToken();
                    break;
                }

                if (nextThree === '"""') {
                    flushToken();
                    inMultilineString = true;
                    i += 2;
                    continue;
                }

                if (char === '"') {
                    flushToken();
                    inString = true;
                    continue;
                }

                if (/[A-Za-z0-9_-]/.test(char)) {
                    token += char;
                    continue;
                }

                flushToken();

                if (char === '{') {
                    stack.push(pendingBlock ?? { type: 'anonymous' });
                    pendingBlock = undefined;
                    expectBlockName = false;
                } else if (char === '}') {
                    stack.pop();
                    pendingBlock = undefined;
                    expectBlockName = false;
                }
            }

            flushToken();
        }

        return stack;
    }
    
    static hasTypedContent(document: vscode.TextDocument, position: vscode.Position): boolean {
        return true;
    }

    static isPositionInCode(document: vscode.TextDocument, position: vscode.Position): boolean {
        let inString = false;
        let inMultilineString = false;
        let escapeNext = false;

        for (let lineNumber = 0; lineNumber <= position.line; lineNumber++) {
            const fullLine = document.lineAt(lineNumber).text;
            const line = lineNumber === position.line ? fullLine.substring(0, position.character) : fullLine;

            for (let i = 0; i < line.length; i++) {
                const char = line[i];
                const nextThree = line.substring(i, i + 3);

                if (inMultilineString) {
                    if (nextThree === '"""') {
                        inMultilineString = false;
                        i += 2;
                    }
                    continue;
                }

                if (inString) {
                    if (escapeNext) {
                        escapeNext = false;
                    } else if (char === '\\') {
                        escapeNext = true;
                    } else if (char === '"') {
                        inString = false;
                    }
                    continue;
                }

                if (char === '#') {
                    if (lineNumber === position.line) {
                        return false;
                    }
                    break;
                }

                if (nextThree === '"""') {
                    inMultilineString = true;
                    i += 2;
                    continue;
                }

                if (char === '"') {
                    inString = true;
                }
            }
        }

        return !inString && !inMultilineString;
    }
}
