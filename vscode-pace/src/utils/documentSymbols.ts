export interface TextPosition {
    line: number;
    character: number;
}

export interface TextRange {
    start: TextPosition;
    end: TextPosition;
}

interface TextLine {
    text: string;
}

interface TextDocumentLike {
    lineCount: number;
    lineAt(line: number): TextLine;
}

export interface PaceDocumentSymbols {
    tasks: string[];
    taskAliases: string[];
    hooks: string[];
}

export interface PaceDefinitionSymbol {
    type: 'task' | 'hook';
    name: string;
    alias?: string;
    range: TextRange;
    selectionRange: TextRange;
}

type Token =
    | { type: 'identifier'; value: string; range: TextRange }
    | { type: '{' | '}' | '[' | ']'; range: TextRange };

function position(line: number, character: number): TextPosition {
    return { line, character };
}

function makeRange(line: number, start: number, end: number): TextRange {
    return { start: position(line, start), end: position(line, end) };
}

function makeMultiRange(start: TextPosition, end: TextPosition): TextRange {
    return { start, end };
}

function tokenize(document: TextDocumentLike): Token[] {
    const tokens: Token[] = [];
    let inString = false;
    let inMultilineString = false;
    let escapeNext = false;

    const flushToken = (value: string, line: number, start: number, end: number) => {
        if (value) {
            tokens.push({ type: 'identifier', value, range: makeRange(line, start, end) });
        }
    };

    for (let lineNumber = 0; lineNumber < document.lineCount; lineNumber++) {
        const line = document.lineAt(lineNumber).text;
        let current = '';
        let currentStart = 0;

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
                flushToken(current, lineNumber, currentStart, i);
                current = '';
                break;
            }

            if (nextThree === '"""') {
                flushToken(current, lineNumber, currentStart, i);
                current = '';
                inMultilineString = true;
                i += 2;
                continue;
            }

            if (char === '"') {
                flushToken(current, lineNumber, currentStart, i);
                current = '';
                inString = true;
                continue;
            }

            if (/[A-Za-z0-9_-]/.test(char)) {
                if (!current) {
                    currentStart = i;
                }
                current += char;
                continue;
            }

            flushToken(current, lineNumber, currentStart, i);
            current = '';

            if (char === '{' || char === '}' || char === '[' || char === ']') {
                tokens.push({ type: char, range: makeRange(lineNumber, i, i + 1) });
            }
        }

        flushToken(current, lineNumber, currentStart, line.length);
    }

    return tokens;
}

function identifierAt(tokens: Token[], index: number): Extract<Token, { type: 'identifier' }> | undefined {
    const token = tokens[index];
    return token?.type === 'identifier' ? token : undefined;
}

function findMatchingBrace(tokens: Token[], openBraceIndex: number): Token | undefined {
    let depth = 0;
    for (let i = openBraceIndex; i < tokens.length; i++) {
        const token = tokens[i];
        if (token.type === '{') {
            depth++;
        } else if (token.type === '}') {
            depth--;
            if (depth === 0) {
                return token;
            }
        }
    }
    return undefined;
}

function lineRange(document: TextDocumentLike, line: number): TextRange {
    const text = document.lineAt(line).text;
    return makeRange(line, 0, text.length);
}

export class PaceDocumentSymbolExtractor {
    static extract(document: TextDocumentLike): PaceDocumentSymbols {
        const definitions = this.extractDefinitions(document);
        const taskAliases = new Set<string>();
        const tokens = tokenize(document);
        let depth = 0;

        for (let i = 0; i < tokens.length; i++) {
            const token = tokens[i];
            if (token.type === '{') {
                depth++;
                continue;
            }
            if (token.type === '}') {
                depth = Math.max(0, depth - 1);
                continue;
            }
            if (depth === 0 && token.type === 'identifier' && token.value === 'alias') {
                const alias = identifierAt(tokens, i + 1);
                const target = identifierAt(tokens, i + 2);
                if (alias && target) {
                    taskAliases.add(alias.value);
                }
            }
        }

        for (const definition of definitions) {
            if (definition.alias) {
                taskAliases.add(definition.alias);
            }
        }

        return {
            tasks: definitions.filter(symbol => symbol.type === 'task').map(symbol => symbol.name).sort(),
            taskAliases: [...taskAliases].sort(),
            hooks: definitions.filter(symbol => symbol.type === 'hook').map(symbol => symbol.name).sort()
        };
    }

    static extractDefinitions(document: TextDocumentLike): PaceDefinitionSymbol[] {
        const definitions: PaceDefinitionSymbol[] = [];
        const tokens = tokenize(document);
        let depth = 0;

        for (let i = 0; i < tokens.length; i++) {
            const token = tokens[i];

            if (token.type === '{') {
                depth++;
                continue;
            }

            if (token.type === '}') {
                depth = Math.max(0, depth - 1);
                continue;
            }

            if (depth !== 0 || token.type !== 'identifier') {
                continue;
            }

            if (token.value === 'task') {
                const taskName = identifierAt(tokens, i + 1);
                if (!taskName) {
                    continue;
                }

                let nextIndex = i + 2;
                let alias: string | undefined;
                if (tokens[nextIndex]?.type === '[') {
                    const aliasToken = identifierAt(tokens, nextIndex + 1);
                    if (aliasToken && tokens[nextIndex + 2]?.type === ']') {
                        alias = aliasToken.value;
                        nextIndex += 3;
                    }
                }

                if (tokens[nextIndex]?.type === '{') {
                    const closeBrace = findMatchingBrace(tokens, nextIndex);
                    const range = closeBrace
                        ? makeMultiRange(token.range.start, closeBrace.range.end)
                        : lineRange(document, token.range.start.line);
                    definitions.push({
                        type: 'task',
                        name: taskName.value,
                        alias,
                        range,
                        selectionRange: taskName.range
                    });
                }
                continue;
            }

            if (token.value === 'hook') {
                const hookName = identifierAt(tokens, i + 1);
                if (hookName && tokens[i + 2]?.type === '{') {
                    const closeBrace = findMatchingBrace(tokens, i + 2);
                    const range = closeBrace
                        ? makeMultiRange(token.range.start, closeBrace.range.end)
                        : lineRange(document, token.range.start.line);
                    definitions.push({
                        type: 'hook',
                        name: hookName.value,
                        range,
                        selectionRange: hookName.range
                    });
                }
            }
        }

        return definitions;
    }
}
