import * as vscode from 'vscode';
import { PaceDocumentSymbolExtractor, TextPosition, TextRange } from './utils/documentSymbols';

function toPosition(position: TextPosition): vscode.Position {
    return new vscode.Position(position.line, position.character);
}

export function toRange(range: TextRange): vscode.Range {
    return new vscode.Range(toPosition(range.start), toPosition(range.end));
}

export class PaceDocumentSymbolProvider implements vscode.DocumentSymbolProvider {
    provideDocumentSymbols(document: vscode.TextDocument): vscode.DocumentSymbol[] {
        try {
            return PaceDocumentSymbolExtractor.extractDefinitions(document).map(definition => {
                return new vscode.DocumentSymbol(
                    definition.name,
                    definition.type === 'task' && definition.alias ? `task alias: ${definition.alias}` : definition.type,
                    definition.type === 'task' ? vscode.SymbolKind.Function : vscode.SymbolKind.Event,
                    toRange(definition.range),
                    toRange(definition.selectionRange)
                );
            });
        } catch {
            return [];
        }
    }
}
