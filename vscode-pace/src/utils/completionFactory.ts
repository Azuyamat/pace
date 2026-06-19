import * as vscode from 'vscode';
import { SnippetDefinition } from '../snippets.config';

export class CompletionFactory {
    static createCompletionItem(definition: SnippetDefinition): vscode.CompletionItem {
        const kind = definition.kind || vscode.CompletionItemKind.Snippet;
        const item = new vscode.CompletionItem(definition.label, kind);
        
        item.insertText = new vscode.SnippetString(definition.snippet);
        item.detail = definition.detail;
        item.documentation = new vscode.MarkdownString([
            definition.documentation,
            '',
            'Example:',
            '',
            '```pace',
            definition.snippet,
            '```'
        ].join('\n'));
        if (definition.deprecated && vscode.CompletionItemTag?.Deprecated !== undefined) {
            item.tags = [vscode.CompletionItemTag.Deprecated];
        }
        return item;
    }

    static createCompletionItems(definitions: SnippetDefinition[]): vscode.CompletionItem[] {
        return definitions.map((def, index) => {
            const item = this.createCompletionItem(def);
            item.sortText = `${String(index).padStart(3, '0')}-${def.label}`;
            return item;
        });
    }

    static applyReplacementRange(items: vscode.CompletionItem[], document: vscode.TextDocument, position: vscode.Position): vscode.CompletionItem[] {
        if (typeof document.getWordRangeAtPosition !== 'function') {
            return items;
        }
        const wordRange = document.getWordRangeAtPosition(position, /[A-Za-z0-9_-]+/);
        if (wordRange) {
            for (const item of items) {
                item.range = wordRange;
            }
        }
        return items;
    }
}
