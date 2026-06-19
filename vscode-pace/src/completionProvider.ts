import * as vscode from 'vscode';
import { snippetsConfig } from './snippets.config';
import { ContextDetector } from './utils/contextDetector';
import { CompletionFactory } from './utils/completionFactory';
import { PaceDocumentSymbolExtractor } from './utils/documentSymbols';

const taskReferenceProperties = ['depends-on', 'dependencies'];
const hookReferenceProperties = ['requires', 'before', 'triggers', 'after', 'on_success', 'on_failure'];

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
        const referenceKind = this.getReferenceCompletionKind(document, position);

        if (referenceKind) {
            return CompletionFactory.applyReplacementRange(
                this.createReferenceCompletions(document, referenceKind, docContext.currentTaskName),
                document,
                position
            );
        }

        if (docContext.isTopLevel) {
            completions.push(...CompletionFactory.createCompletionItems(snippetsConfig.topLevel));
        }

        if (docContext.inEnvBlock) {
            return CompletionFactory.applyReplacementRange(completions, document, position);
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

        return CompletionFactory.applyReplacementRange(completions, document, position);
    }

    private getReferenceCompletionKind(document: vscode.TextDocument, position: vscode.Position): 'task' | 'hook' | undefined {
        const linePrefix = document.lineAt(position.line).text.substring(0, position.character);
        const commentIndex = linePrefix.indexOf('#');
        if (commentIndex !== -1) {
            return undefined;
        }

        if (/^\s*default\s+[A-Za-z0-9_-]*$/.test(linePrefix)) {
            return 'task';
        }

        if (/^\s*alias\s+[A-Za-z_][A-Za-z0-9_-]*\s+[A-Za-z0-9_-]*$/.test(linePrefix)) {
            return 'task';
        }

        const documentPrefix = document.getText(new vscode.Range(new vscode.Position(0, 0), position));
        if (this.isInReferenceArray(documentPrefix, taskReferenceProperties)) {
            return 'task';
        }

        if (this.isInReferenceArray(documentPrefix, hookReferenceProperties)) {
            return 'hook';
        }

        return undefined;
    }

    private isInReferenceArray(documentPrefix: string, properties: string[]): boolean {
        const tail = documentPrefix.slice(-1000);
        return properties.some(property => {
            const escapedProperty = property.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
            return new RegExp(`(?:^|\\s)${escapedProperty}\\s*\\[[^\\]]*$`, 's').test(tail);
        });
    }

    private createReferenceCompletions(document: vscode.TextDocument, kind: 'task' | 'hook', currentTaskName?: string): vscode.CompletionItem[] {
        const symbols = PaceDocumentSymbolExtractor.extract(document);
        const labels = kind === 'task'
            ? [...symbols.tasks, ...symbols.taskAliases]
            : symbols.hooks;

        return [...new Set(labels)].map(label => {
            const item = new vscode.CompletionItem(label, vscode.CompletionItemKind.Reference);
            item.detail = kind === 'task' ? 'Pace task reference' : 'Pace hook reference';
            item.insertText = label;
            if (kind === 'task' && label === currentTaskName) {
                item.sortText = `z-${label}`;
                item.detail = 'Pace task reference (current task)';
            }
            return item;
        });
    }
}
