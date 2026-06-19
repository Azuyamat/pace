import * as cp from 'child_process';
import * as fs from 'fs';
import * as path from 'path';
import * as vscode from 'vscode';
import { runPaceTask, runPaceDefaultTask } from './codeLensProvider';
import { debugLog, log } from './output';
import { PaceDocumentSymbolExtractor, PaceDefinitionSymbol } from './utils/documentSymbols';

interface GeneratedConfig {
    Tasks?: Record<string, GeneratedTask>;
}

interface GeneratedTask {
    Name?: string;
    Alias?: string;
    Command?: string;
    Inputs?: string[];
    Outputs?: string[];
    DependsOn?: string[];
    Cache?: boolean;
    Watch?: boolean;
    Description?: string;
}

function getPacePath(uri?: vscode.Uri): string {
    return vscode.workspace.getConfiguration('pace', uri).get<string>('path', 'pace') || 'pace';
}

async function pickWorkspaceFolder(): Promise<vscode.WorkspaceFolder | undefined> {
    const activeDocument = vscode.window.activeTextEditor?.document;
    if (activeDocument) {
        const folder = vscode.workspace.getWorkspaceFolder(activeDocument.uri);
        if (folder) {
            return folder;
        }
    }

    const folders = vscode.workspace.workspaceFolders || [];
    if (folders.length === 0) {
        vscode.window.showErrorMessage('Open a workspace folder before running Pace commands.');
        return undefined;
    }
    if (folders.length === 1) {
        return folders[0];
    }

    return vscode.window.showWorkspaceFolderPick({ placeHolder: 'Select the workspace folder to use for Pace.' });
}

function execFile(command: string, args: string[], cwd: string, uri?: vscode.Uri): Promise<{ stdout: string; stderr: string }> {
    debugLog(`Running ${command} ${args.join(' ')} in ${cwd}`, uri);
    return new Promise((resolve, reject) => {
        cp.execFile(command, args, { cwd, windowsHide: true }, (error, stdout, stderr) => {
            if (error) {
                log(`Command failed: ${command} ${args.join(' ')}\n${stderr || error.message}`);
                reject(new Error(stderr || error.message));
                return;
            }
            resolve({ stdout, stderr });
        });
    });
}

export async function generateConfigCommand(): Promise<void> {
    const folder = await pickWorkspaceFolder();
    if (!folder) {
        return;
    }

    const configPath = path.join(folder.uri.fsPath, 'config.pace');
    const exists = fs.existsSync(configPath);
    if (exists) {
        const answer = await vscode.window.showWarningMessage(
            `${folder.name} already contains config.pace. Overwrite it?`,
            { modal: true },
            'Overwrite'
        );
        if (answer !== 'Overwrite') {
            return;
        }
    }

    try {
        const args = ['init', '--yes'];
        if (exists) {
            args.push('--force');
        }
        await execFile(getPacePath(folder.uri), args, folder.uri.fsPath, folder.uri);
        vscode.window.showInformationMessage(`Generated ${configPath}`);
    } catch (error) {
        vscode.window.showErrorMessage(`Failed to generate config.pace: ${String(error instanceof Error ? error.message : error)}`);
    }
}

function taskAtCursor(document: vscode.TextDocument, position: vscode.Position): PaceDefinitionSymbol | undefined {
    return PaceDocumentSymbolExtractor.extractDefinitions(document)
        .filter(definition => definition.type === 'task')
        .find(definition => {
            const range = new vscode.Range(
                new vscode.Position(definition.range.start.line, definition.range.start.character),
                new vscode.Position(definition.range.end.line, definition.range.end.character)
            );
            return range.contains(position);
        });
}

function descriptionForTask(document: vscode.TextDocument, definition: PaceDefinitionSymbol): string | undefined {
    const range = new vscode.Range(
        new vscode.Position(definition.range.start.line, definition.range.start.character),
        new vscode.Position(definition.range.end.line, definition.range.end.character)
    );
    const match = document.getText(range).match(/(?:^|\s)description\s+"([^"]*)"/);
    return match?.[1];
}

async function pickTask(document: vscode.TextDocument): Promise<string | undefined> {
    const definitions = PaceDocumentSymbolExtractor.extractDefinitions(document).filter(definition => definition.type === 'task');
    if (definitions.length === 0) {
        vscode.window.showErrorMessage('No Pace tasks found in the current document.');
        return undefined;
    }

    const selected = await vscode.window.showQuickPick(definitions.map(definition => ({
        label: definition.name,
        description: definition.alias ? `alias: ${definition.alias}` : undefined,
        detail: descriptionForTask(document, definition)
    })), { placeHolder: 'Select a Pace task to run' });

    return selected?.label;
}

export async function runSelectedOrCurrentTaskCommand(): Promise<void> {
    const editor = vscode.window.activeTextEditor;
    const document = editor?.document;
    if (!editor || !document || document.languageId !== 'pace') {
        vscode.window.showErrorMessage('Open a Pace file before running the current task.');
        return;
    }

    const currentTask = taskAtCursor(document, editor.selection.active);
    const taskName = currentTask?.name || await pickTask(document);
    if (taskName) {
        runPaceTask(taskName, document);
    }
}

export async function runTaskCommand(): Promise<void> {
    const document = vscode.window.activeTextEditor?.document;
    if (!document || document.languageId !== 'pace') {
        vscode.window.showErrorMessage('Open a Pace file before running a task.');
        return;
    }

    const taskName = await pickTask(document);
    if (taskName) {
        runPaceTask(taskName, document);
    }
}

export function runDefaultTaskCommand(): void {
    const document = vscode.window.activeTextEditor?.document;
    runPaceDefaultTask(document?.languageId === 'pace' ? document : undefined);
}

function quote(value: string): string {
    return `"${value.replace(/\\/g, '\\\\').replace(/"/g, '\\"')}"`;
}

function arrayProperty(name: string, values: string[] | undefined): string | undefined {
    return values && values.length > 0 ? `    ${name} [${values.map(quote).join(', ')}]` : undefined;
}

function generatedTaskSnippet(name: string, task: GeneratedTask): string {
    const lines = [`task ${task.Name || name}${task.Alias ? ` [${task.Alias}]` : ''} {`];
    if (task.Command) {
        lines.push(`    command ${quote(task.Command)}`);
    }
    if (task.Description) {
        lines.push(`    description ${quote(task.Description)}`);
    }
    for (const line of [arrayProperty('inputs', task.Inputs), arrayProperty('outputs', task.Outputs), arrayProperty('depends-on', task.DependsOn)]) {
        if (line) {
            lines.push(line);
        }
    }
    if (task.Cache) {
        lines.push('    cache true');
    }
    if (task.Watch) {
        lines.push('    watch true');
    }
    lines.push('}');
    return lines.join('\n');
}

export async function insertGoTaskTemplateCommand(): Promise<void> {
    const editor = vscode.window.activeTextEditor;
    const folder = await pickWorkspaceFolder();
    if (!editor || !folder) {
        return;
    }

    try {
        const result = await execFile(getPacePath(folder.uri), ['init', '--type', 'go', '--json', '--stdout', '--yes'], folder.uri.fsPath, folder.uri);
        const config = JSON.parse(result.stdout) as GeneratedConfig;
        const tasks = Object.entries(config.Tasks || {}).filter(([name]) => ['build', 'run', 'test', 'vet', 'tidy', 'lint'].includes(name));
        const selected = await vscode.window.showQuickPick(tasks.map(([name, task]) => ({
            label: name,
            detail: task.Command,
            description: task.Description,
            task
        })), { placeHolder: 'Select a generated Go task to insert' });
        if (!selected) {
            return;
        }
        await editor.insertSnippet(new vscode.SnippetString(generatedTaskSnippet(selected.label, selected.task)));
    } catch (error) {
        vscode.window.showErrorMessage(`Failed to load generated Go tasks from Pace CLI: ${String(error instanceof Error ? error.message : error)}`);
    }
}
