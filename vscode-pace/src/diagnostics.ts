import * as vscode from 'vscode';
import { execFile, ChildProcess } from 'child_process';
import { log, debugLog } from './output';

interface PaceCliDiagnostic {
    file?: string;
    range?: {
        start?: { line?: number; column?: number };
        end?: { line?: number; column?: number };
    };
    severity?: string;
    message?: string;
    hint?: string;
}

export class PaceDiagnosticsManager implements vscode.Disposable {
    private readonly collection = vscode.languages.createDiagnosticCollection('pace');
    private readonly disposables: vscode.Disposable[] = [];
    private readonly timers = new Map<string, NodeJS.Timeout>();
    private readonly versions = new Map<string, number>();
    private readonly running = new Map<string, ChildProcess>();
    private warnedMissingBinary = false;

    constructor() {
        this.disposables.push(this.collection);
        this.disposables.push(vscode.workspace.onDidOpenTextDocument(document => this.schedule(document)));
        this.disposables.push(vscode.workspace.onDidChangeTextDocument(event => {
            if (this.runMode(event.document) === 'onType') {
                this.schedule(event.document);
            }
        }));
        this.disposables.push(vscode.workspace.onDidSaveTextDocument(document => {
            if (this.runMode(document) === 'onSave') {
                this.schedule(document, 0);
            }
        }));
        this.disposables.push(vscode.workspace.onDidCloseTextDocument(document => {
            this.clear(document.uri);
        }));
        this.disposables.push(vscode.workspace.onDidChangeConfiguration(event => {
            if (event.affectsConfiguration('pace.diagnostics') || event.affectsConfiguration('pace.path')) {
                this.validateOpenDocuments();
            }
        }));
        this.validateOpenDocuments();
    }

    dispose(): void {
        for (const timer of this.timers.values()) {
            clearTimeout(timer);
        }
        for (const process of this.running.values()) {
            process.kill();
        }
        this.disposables.forEach(disposable => disposable.dispose());
    }

    private validateOpenDocuments(): void {
        vscode.workspace.textDocuments.forEach(document => this.schedule(document));
    }

    private schedule(document: vscode.TextDocument, delay = 400): void {
        if (!this.shouldValidate(document)) {
            this.clear(document.uri);
            return;
        }

        const key = document.uri.toString();
        const existing = this.timers.get(key);
        if (existing) {
            clearTimeout(existing);
        }

        const timer = setTimeout(() => {
            this.timers.delete(key);
            this.validate(document);
        }, delay);
        this.timers.set(key, timer);
    }

    private shouldValidate(document: vscode.TextDocument): boolean {
        return document.languageId === 'pace'
            && document.uri.scheme === 'file'
            && vscode.workspace.getConfiguration('pace', document.uri).get<boolean>('diagnostics.enabled', true);
    }

    private runMode(document: vscode.TextDocument): 'onType' | 'onSave' {
        return vscode.workspace.getConfiguration('pace', document.uri).get<'onType' | 'onSave'>('diagnostics.run', 'onType');
    }

    private validate(document: vscode.TextDocument): void {
        const key = document.uri.toString();
        const version = (this.versions.get(key) ?? 0) + 1;
        this.versions.set(key, version);

        const existing = this.running.get(key);
        if (existing) {
            existing.kill();
        }

        const pacePath = vscode.workspace.getConfiguration('pace', document.uri).get<string>('path', 'pace') || 'pace';
        const cwd = vscode.workspace.getWorkspaceFolder(document.uri)?.uri.fsPath;
        debugLog(`Running diagnostics: ${pacePath} validate --json ${document.uri.fsPath}`, document.uri);

        const child = execFile(pacePath, ['validate', '--json', document.uri.fsPath], { cwd }, (error, stdout, stderr) => {
            this.running.delete(key);
            if (this.versions.get(key) !== version) {
                return;
            }

            if (error && (error as NodeJS.ErrnoException).code === 'ENOENT') {
                this.collection.delete(document.uri);
                if (!this.warnedMissingBinary) {
                    this.warnedMissingBinary = true;
                    vscode.window.showWarningMessage('Pace diagnostics are disabled because the Pace CLI was not found. Set pace.path to the CLI location.');
                }
                log(`Pace CLI not found: ${pacePath}`, document.uri);
                return;
            }

            if (stderr.trim()) {
                log(`pace validate stderr: ${stderr.trim()}`, document.uri);
            }

            try {
                const parsed = JSON.parse(stdout || '[]') as PaceCliDiagnostic[];
                this.collection.set(document.uri, parsed.map(item => this.toDiagnostic(item)));
            } catch (parseError) {
                this.collection.delete(document.uri);
                log(`Failed to parse pace validate JSON: ${parseError instanceof Error ? parseError.message : String(parseError)}`, document.uri);
            }
        });

        this.running.set(key, child);
    }

    private toDiagnostic(item: PaceCliDiagnostic): vscode.Diagnostic {
        const startLine = Math.max((item.range?.start?.line ?? 1) - 1, 0);
        const startColumn = Math.max((item.range?.start?.column ?? 1) - 1, 0);
        const endLine = Math.max((item.range?.end?.line ?? item.range?.start?.line ?? 1) - 1, startLine);
        const endColumn = Math.max((item.range?.end?.column ?? item.range?.start?.column ?? 2) - 1, startColumn + 1);
        const diagnostic = new vscode.Diagnostic(
            new vscode.Range(startLine, startColumn, endLine, endColumn),
            item.hint ? `${item.message ?? 'Pace validation error'}\n${item.hint}` : item.message ?? 'Pace validation error',
            this.toSeverity(item.severity)
        );
        diagnostic.source = 'pace';
        return diagnostic;
    }

    private toSeverity(severity?: string): vscode.DiagnosticSeverity {
        switch (severity) {
            case 'warning':
                return vscode.DiagnosticSeverity.Warning;
            case 'info':
                return vscode.DiagnosticSeverity.Information;
            case 'hint':
                return vscode.DiagnosticSeverity.Hint;
            default:
                return vscode.DiagnosticSeverity.Error;
        }
    }

    private clear(uri: vscode.Uri): void {
        const key = uri.toString();
        const timer = this.timers.get(key);
        if (timer) {
            clearTimeout(timer);
            this.timers.delete(key);
        }
        const process = this.running.get(key);
        if (process) {
            process.kill();
            this.running.delete(key);
        }
        this.collection.delete(uri);
    }
}
