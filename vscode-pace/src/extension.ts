import * as vscode from 'vscode';
import { PaceCompletionProvider } from './completionProvider';

export function activate(context: vscode.ExtensionContext) {
    console.log('Pace language extension is now active');

    const completionProvider = vscode.languages.registerCompletionItemProvider(
        'pace',
        new PaceCompletionProvider(),
        ' ', '"'
    );

    context.subscriptions.push(completionProvider);
}

export function deactivate() {}
