import * as vscode from 'vscode';

let channel: vscode.OutputChannel | undefined;

export function getPaceOutputChannel(): vscode.OutputChannel {
    if (!channel) {
        channel = vscode.window.createOutputChannel('Pace');
    }
    return channel;
}

export function log(message: string, document?: vscode.Uri): void {
    const prefix = new Date().toISOString();
    const suffix = document ? ` (${document.fsPath})` : '';
    getPaceOutputChannel().appendLine(`[${prefix}] ${message}${suffix}`);
}

export function debugLog(message: string, document?: vscode.Uri): void {
    const enabled = vscode.workspace.getConfiguration('pace', document).get<boolean>('logging.debug', false);
    if (enabled) {
        log(message, document);
    }
}
