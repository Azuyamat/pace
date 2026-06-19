import * as vscode from 'vscode';
import { snippetsConfig } from './snippets.config';
import { paceSchema } from './generated/paceSchema';

export interface LanguageDoc {
    label: string;
    documentation: string;
    kind?: vscode.CompletionItemKind;
}

export const topLevelDocs: LanguageDoc[] = paceSchema.topLevel.map(statement => ({
    label: statement.name,
    documentation: statement.description
}));

export const taskPropertyDocs: LanguageDoc[] = snippetsConfig.taskProperties.map(property => ({
    label: property.label,
    documentation: property.documentation,
    kind: property.kind
}));

export const hookPropertyDocs: LanguageDoc[] = snippetsConfig.hookProperties.map(property => ({
    label: property.label,
    documentation: property.documentation,
    kind: property.kind
}));

export const argsPropertyDocs: LanguageDoc[] = snippetsConfig.argsProperties.map(property => ({
    label: property.label,
    documentation: property.documentation,
    kind: property.kind
}));

export const docsByContext = {
    topLevel: new Map(topLevelDocs.map(doc => [doc.label, doc])),
    task: new Map(taskPropertyDocs.map(doc => [doc.label, doc])),
    hook: new Map(hookPropertyDocs.map(doc => [doc.label, doc])),
    args: new Map(argsPropertyDocs.map(doc => [doc.label, doc]))
};
