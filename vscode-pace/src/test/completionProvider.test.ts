import * as assert from 'assert';
import Module = require('module');

class Position {
    constructor(public line: number, public character: number) {}
}

class Range {
    constructor(public start: Position, public end: Position) {}
}

class CompletionItem {
    insertText?: unknown;
    documentation?: unknown;
    detail?: string;
    sortText?: string;

    constructor(public label: string, public kind?: number) {}
}

class SnippetString {
    constructor(public value: string) {}
}

class MarkdownString {
    constructor(public value: string) {}
}

const vscodeMock = {
    Position,
    Range,
    CompletionItem,
    SnippetString,
    MarkdownString,
    CompletionItemKind: {
        Text: 0,
        Method: 1,
        Function: 2,
        Constructor: 3,
        Field: 4,
        Variable: 5,
        Class: 6,
        Interface: 7,
        Module: 8,
        Property: 9,
        Unit: 10,
        Value: 11,
        Enum: 12,
        Keyword: 13,
        Snippet: 14,
        Color: 15,
        File: 16,
        Reference: 17
    }
};

const originalLoad = (Module as any)._load;
(Module as any)._load = function(request: string, parent: unknown, isMain: boolean) {
    if (request === 'vscode') {
        return vscodeMock;
    }
    return originalLoad.apply(this, arguments as any);
};

const { PaceCompletionProvider } = require('../completionProvider');
const { CompletionFactory } = require('../utils/completionFactory');
const { snippetsConfig } = require('../snippets.config');

class MockDocument {
    private readonly lines: string[];
    readonly lineCount: number;

    constructor(private readonly text: string) {
        this.lines = text.split(/\r?\n/);
        this.lineCount = this.lines.length;
    }

    lineAt(line: number): { text: string } {
        return { text: this.lines[line] ?? '' };
    }

    getText(range?: Range): string {
        if (!range) {
            return this.text;
        }

        const endOffset = this.offsetAt(range.end);
        const startOffset = this.offsetAt(range.start);
        return this.text.slice(startOffset, endOffset);
    }

    private offsetAt(position: Position): number {
        let offset = 0;
        for (let i = 0; i < position.line; i++) {
            offset += (this.lines[i] ?? '').length + 1;
        }
        return offset + position.character;
    }
}

function positionOf(text: string, marker = '<cursor>'): Position {
    const index = text.indexOf(marker);
    assert.notStrictEqual(index, -1, `marker not found: ${marker}`);
    const before = text.slice(0, index);
    const lines = before.split(/\r?\n/);
    return new Position(lines.length - 1, lines[lines.length - 1].length);
}

function labelsFor(textWithMarker: string): string[] {
    const position = positionOf(textWithMarker);
    const document = new MockDocument(textWithMarker.replace('<cursor>', ''));
    const provider = new PaceCompletionProvider();
    const items = provider.provideCompletionItems(document as any, position as any, {} as any, {} as any) as CompletionItem[];
    return items.map(item => item.label).sort();
}

{
    const cacheSnippet = snippetsConfig.taskProperties.find((definition: { label: string }) => definition.label === 'cache');
    assert.ok(cacheSnippet, 'cache snippet exists');
    const item = CompletionFactory.createCompletionItem(cacheSnippet);
    assert.strictEqual(item.label, 'cache');
    assert.strictEqual((item.insertText as SnippetString).value, 'cache ${1|true,false|}');
}

{
    const labels = labelsFor('<cursor>');
    for (const expected of ['task', 'hook', 'var', 'default', 'alias', 'import']) {
        assert.ok(labels.includes(expected), `top-level completions include ${expected}`);
    }
}

{
    const labels = labelsFor('task build {\n    <cursor>\n}');
    for (const expected of ['command', 'depends-on', 'requires', 'args', 'when']) {
        assert.ok(labels.includes(expected), `task completions include ${expected}`);
    }
}

{
    const labels = labelsFor('hook setup {\n    <cursor>\n}');
    assert.deepStrictEqual(labels, ['command', 'description', 'env', 'working_dir'].sort());
}

{
    const labels = labelsFor('task build {\n    args {\n        <cursor>\n    }\n}');
    assert.deepStrictEqual(labels, ['optional', 'required']);
}

{
    const labels = labelsFor('task build {\n    env {\n        <cursor>\n    }\n}');
    assert.deepStrictEqual(labels, []);
}

{
    const labels = labelsFor('task build {\n    depends-on [<cursor>\n}\n\ntask test {\n    command "go test"\n}');
    assert.ok(labels.includes('test'), 'task reference completions include task names');
}

console.log('completionProvider tests passed');
