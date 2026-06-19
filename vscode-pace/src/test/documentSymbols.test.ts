import * as assert from 'assert';
import { PaceDocumentSymbolExtractor } from '../utils/documentSymbols';

class MockDocument {
    private readonly lines: string[];
    readonly lineCount: number;

    constructor(text: string) {
        this.lines = text.split(/\r?\n/);
        this.lineCount = this.lines.length;
    }

    lineAt(line: number): { text: string } {
        return { text: this.lines[line] ?? '' };
    }
}

const document = new MockDocument(`# task ignored {\n# hook ignored {\nvar msg = "task ignored {"\ntask build [b] {\n\tcommand "go build"\n}\n\ntask test\n{\n\tcommand "go test"\n\tenv {\n\t\tNOT_A_TASK = "task nested {"\n\t}\n}\n\nalias t test\n\nhook setup {\n\tcommand "echo setup"\n}\n\ncommand """\nhook ignored_multiline {\n"""\n`) as any;

const symbols = PaceDocumentSymbolExtractor.extract(document);
assert.deepStrictEqual(symbols.tasks, ['build', 'test']);
assert.deepStrictEqual(symbols.taskAliases, ['b', 't']);
assert.deepStrictEqual(symbols.hooks, ['setup']);

console.log('documentSymbols tests passed');
