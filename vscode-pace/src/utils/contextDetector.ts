import * as vscode from 'vscode';

export interface DocumentContext {
    inTaskBlock: boolean;
    inHookBlock: boolean;
    inEnvBlock: boolean;
    inArgsBlock: boolean;
    isTopLevel: boolean;
}

export class ContextDetector {
    static detectContext(document: vscode.TextDocument, position: vscode.Position): DocumentContext {
        const inTaskBlock = this.isInBlock(document, position, 'task');
        const inHookBlock = this.isInBlock(document, position, 'hook');
        const inEnvBlock = this.isInBlock(document, position, 'env');
        const inArgsBlock = this.isInBlock(document, position, 'args');

        return {
            inTaskBlock,
            inHookBlock,
            inEnvBlock,
            inArgsBlock,
            isTopLevel: !inTaskBlock && !inHookBlock && !inEnvBlock && !inArgsBlock
        };
    }

    private static isInBlock(
        document: vscode.TextDocument,
        position: vscode.Position,
        blockType: string
    ): boolean {
        let openBraces = 0;
        let inTargetBlock = false;

        for (let i = position.line; i >= 0; i--) {
            const line = document.lineAt(i).text;
            
            if (line.includes('}')) openBraces--;
            if (line.includes('{')) {
                openBraces++;
                if (openBraces > 0 && line.includes(blockType + ' ')) {
                    inTargetBlock = true;
                    break;
                }
            }
        }

        return inTargetBlock && openBraces > 0;
    }
    
    static hasTypedContent(document: vscode.TextDocument, position: vscode.Position): boolean {
        const lineText = document.lineAt(position.line).text;
        const textBeforeCursor = lineText.substring(0, position.character).trim();
        return textBeforeCursor.length > 0;
    }
}
