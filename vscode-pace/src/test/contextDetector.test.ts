import * as assert from 'assert';
import { ContextDetector } from '../utils/contextDetector';

type MockPosition = { line: number; character: number };

class MockDocument {
    private readonly lines: string[];

    constructor(text: string) {
        this.lines = text.split(/\r?\n/);
    }

    lineAt(line: number): { text: string } {
        return { text: this.lines[line] ?? '' };
    }
}

function positionOf(text: string, marker: string): MockPosition {
    const index = text.indexOf(marker);
    assert.notStrictEqual(index, -1, `marker not found: ${marker}`);
    const before = text.slice(0, index);
    const lines = before.split(/\r?\n/);
    return { line: lines.length - 1, character: lines[lines.length - 1].length };
}

function detect(text: string, marker = '<cursor>') {
    const position = positionOf(text, marker);
    const document = new MockDocument(text.replace(marker, '')) as any;
    return ContextDetector.detectContext(document, position as any);
}

{
    const context = detect(`task build {\n\tenv {\n\t\t<cursor>\n\t}\n}`);
    assert.strictEqual(context.inEnvBlock, true);
    assert.strictEqual(context.inTaskBlock, false);
    assert.strictEqual(context.currentTaskName, 'build');
    assert.deepStrictEqual(context.blockStack.map(block => block.type), ['task', 'env']);
}

{
    const context = detect(`task build {\n\tcommand "echo { not a block }"\n\t# } comment brace\n\t<cursor>\n}`);
    assert.strictEqual(context.inTaskBlock, true);
    assert.deepStrictEqual(context.blockStack.map(block => block.type), ['task']);
}

{
    const context = detect(`task build {\n\tcommand """\n\t{ not a block } # not a comment\n\t"""\n\targs {\n\t\t<cursor>\n\t}\n}`);
    assert.strictEqual(context.inArgsBlock, true);
    assert.deepStrictEqual(context.blockStack.map(block => block.type), ['task', 'args']);
}

console.log('contextDetector tests passed');
