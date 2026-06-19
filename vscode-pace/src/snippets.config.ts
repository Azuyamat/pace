import * as vscode from 'vscode';

export interface SnippetDefinition {
    label: string;
    snippet: string;
    documentation: string;
    kind?: vscode.CompletionItemKind;
    detail?: string;
    deprecated?: boolean;
}

export interface SnippetsConfig {
    topLevel: SnippetDefinition[];
    taskProperties: SnippetDefinition[];
    hookProperties: SnippetDefinition[];
    argsProperties: SnippetDefinition[];
}

// Snippet data is generated from the Go-owned schema. Regenerate with:
//   go run ./cmd/pace schema --generate-vscode
// Do not add language metadata here by hand.
export { generatedSnippetsConfig as snippetsConfig } from './generated/schemaSnippets';
