import * as vscode from 'vscode';
import { snippetsConfig } from './snippets.config';
import { ContextDetector } from './utils/contextDetector';
import { CompletionFactory } from './utils/completionFactory';

export class PaceCompletionProvider implements vscode.CompletionItemProvider {
    provideCompletionItems(
        document: vscode.TextDocument,
        position: vscode.Position,
        token: vscode.CancellationToken,
        context: vscode.CompletionContext
    ): vscode.CompletionItem[] {
        if (!ContextDetector.hasTypedContent(document, position)) {
            return [];
        }

        const completions: vscode.CompletionItem[] = [];
        const docContext = ContextDetector.detectContext(document, position);

        if (docContext.isTopLevel) {
            completions.push(...CompletionFactory.createCompletionItems(snippetsConfig.topLevel));
        }

        if (docContext.inTaskBlock) {
            completions.push(...CompletionFactory.createCompletionItems(snippetsConfig.taskProperties));
        }

        if (docContext.inHookBlock) {
            completions.push(...CompletionFactory.createCompletionItems(snippetsConfig.hookProperties));
        }

        if (docContext.inArgsBlock) {
            completions.push(...CompletionFactory.createCompletionItems(snippetsConfig.argsProperties));
        }

        return completions;
    }
}
