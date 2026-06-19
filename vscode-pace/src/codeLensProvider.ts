import * as vscode from 'vscode';
import { toRange } from './documentSymbolProvider';
import { PaceDocumentSymbolExtractor, PaceDefinitionSymbol, TextPosition } from './utils/documentSymbols';

function quoteShellArg(value: string): string {
    if (process.platform === 'win32') {
        return `"${value.replace(/"/g, '\\"')}"`;
    }
    return `'${value.replace(/'/g, `'\\''`)}'`;
}

function toPosition(position: TextPosition): vscode.Position {
    return new vscode.Position(position.line, position.character);
}

function getWorkspaceFolder(document?: vscode.TextDocument): vscode.WorkspaceFolder | undefined {
    if (document) {
        const folder = vscode.workspace.getWorkspaceFolder(document.uri);
        if (folder) {
            return folder;
        }
    }
    return vscode.workspace.workspaceFolders?.[0];
}

export function runPaceTask(taskName: string, document?: vscode.TextDocument, command: 'run' | 'watch' = 'run'): void {
    if (!/^[A-Za-z_][A-Za-z0-9_-]*$/.test(taskName)) {
        vscode.window.showErrorMessage(`Invalid Pace task name: ${taskName}`);
        return;
    }

    const workspaceFolder = getWorkspaceFolder(document);
    if (!workspaceFolder) {
        vscode.window.showErrorMessage('Open a workspace folder before running Pace tasks.');
        return;
    }

    const pacePath = vscode.workspace.getConfiguration('pace', document?.uri).get<string>('path', 'pace') || 'pace';
    const terminal = vscode.window.createTerminal({
        name: `Pace: ${taskName}`,
        cwd: workspaceFolder.uri.fsPath
    });
    terminal.sendText(`${quoteShellArg(pacePath)} ${command} ${quoteShellArg(taskName)}`);
    terminal.show();
}

export function runPaceDefaultTask(document?: vscode.TextDocument): void {
    const workspaceFolder = getWorkspaceFolder(document);
    if (!workspaceFolder) {
        vscode.window.showErrorMessage('Open a workspace folder before running Pace tasks.');
        return;
    }

    const pacePath = vscode.workspace.getConfiguration('pace', document?.uri).get<string>('path', 'pace') || 'pace';
    const terminal = vscode.window.createTerminal({
        name: 'Pace: default task',
        cwd: workspaceFolder.uri.fsPath
    });
    terminal.sendText(`${quoteShellArg(pacePath)} run`);
    terminal.show();
}

export class PaceCodeLensProvider implements vscode.CodeLensProvider {
    provideCodeLenses(document: vscode.TextDocument): vscode.CodeLens[] {
        const enabled = vscode.workspace.getConfiguration('pace', document.uri).get<boolean>('codeLens.enabled', true);
        if (!enabled || document.languageId !== 'pace') {
            return [];
        }

        try {
            return PaceDocumentSymbolExtractor.extractDefinitions(document)
                .filter(definition => definition.type === 'task')
                .flatMap(definition => this.createTaskCodeLenses(document, definition));
        } catch {
            return [];
        }
    }

    private createTaskCodeLenses(document: vscode.TextDocument, definition: PaceDefinitionSymbol): vscode.CodeLens[] {
        const start = toPosition(definition.selectionRange.start);
        const range = new vscode.Range(start, start);
        const codeLenses = [
            new vscode.CodeLens(range, {
                title: 'Run task',
                command: 'pace.runTaskFromCodeLens',
                arguments: [definition.name, document]
            })
        ];

        if (this.hasWatchEnabled(document, definition)) {
            codeLenses.push(new vscode.CodeLens(range, {
                title: 'Watch task',
                command: 'pace.watchTaskFromCodeLens',
                arguments: [definition.name, document]
            }));
        }

        return codeLenses;
    }

    private hasWatchEnabled(document: vscode.TextDocument, definition: PaceDefinitionSymbol): boolean {
        const text = document.getText(toRange(definition.range));
        return /(?:^|\s)watch\s+true(?:\s|$)/.test(text);
    }
}
