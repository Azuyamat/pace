import * as vscode from 'vscode';
import { docsByContext, LanguageDoc } from './languageMetadata';
import { ContextDetector } from './utils/contextDetector';

const wordPattern = /[A-Za-z_][A-Za-z0-9_-]*/g;

function getWordAtPosition(document: vscode.TextDocument, position: vscode.Position): { word: string; range: vscode.Range } | undefined {
    const line = document.lineAt(position.line).text;
    for (const match of line.matchAll(wordPattern)) {
        const start = match.index ?? 0;
        const end = start + match[0].length;
        if (position.character >= start && position.character <= end) {
            return {
                word: match[0],
                range: new vscode.Range(new vscode.Position(position.line, start), new vscode.Position(position.line, end))
            };
        }
    }
    return undefined;
}

function markdownFor(doc: LanguageDoc): vscode.MarkdownString {
    const markdown = new vscode.MarkdownString();
    markdown.appendMarkdown(`**${doc.label}**\n\n${doc.documentation}`);
    return markdown;
}

export class PaceHoverProvider implements vscode.HoverProvider {
    provideHover(document: vscode.TextDocument, position: vscode.Position): vscode.Hover | undefined {
        if (!ContextDetector.isPositionInCode(document, position)) {
            return undefined;
        }

        const word = getWordAtPosition(document, position);
        if (!word) {
            return undefined;
        }

        const context = ContextDetector.detectContext(document, word.range.start);
        let doc: LanguageDoc | undefined;

        if (context.isTopLevel) {
            doc = docsByContext.topLevel.get(word.word);
        }

        if (!doc && context.inArgsBlock) {
            doc = docsByContext.args.get(word.word);
        }

        if (!doc && context.inTaskBlock) {
            doc = docsByContext.task.get(word.word);
        }

        if (!doc && context.inHookBlock) {
            doc = docsByContext.hook.get(word.word);
        }

        return doc ? new vscode.Hover(markdownFor(doc), word.range) : undefined;
    }
}
