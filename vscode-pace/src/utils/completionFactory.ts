import * as vscode from 'vscode';
import { SnippetDefinition } from '../snippets.config';

export class CompletionFactory {
    static createCompletionItem(definition: SnippetDefinition): vscode.CompletionItem {
        const kind = definition.kind || vscode.CompletionItemKind.Snippet;
        const item = new vscode.CompletionItem(definition.label, kind);
        
        if (kind === vscode.CompletionItemKind.Snippet) {
            item.insertText = new vscode.SnippetString(definition.snippet);
        } else {
            item.insertText = definition.snippet;
        }
        
        item.documentation = new vscode.MarkdownString(definition.documentation);
        return item;
    }

    static createCompletionItems(definitions: SnippetDefinition[]): vscode.CompletionItem[] {
        return definitions.map(def => this.createCompletionItem(def));
    }
}
