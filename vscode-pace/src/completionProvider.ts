import * as vscode from 'vscode';

export class PaceCompletionProvider implements vscode.CompletionItemProvider {
    
    provideCompletionItems(
        document: vscode.TextDocument,
        position: vscode.Position,
        token: vscode.CancellationToken,
        context: vscode.CompletionContext
    ): vscode.CompletionItem[] {
        
        const linePrefix = document.lineAt(position).text.substr(0, position.character);
        const completions: vscode.CompletionItem[] = [];

        // Check if we're inside a block
        const inTaskBlock = this.isInBlock(document, position, 'task');
        const inHookBlock = this.isInBlock(document, position, 'hook');
        const inEnvBlock = this.isInBlock(document, position, 'env');
        const inArgsBlock = this.isInBlock(document, position, 'args');

        // Top-level keywords
        if (!inTaskBlock && !inHookBlock && !inEnvBlock && !inArgsBlock) {
            completions.push(...this.getTopLevelCompletions());
        }

        // Task/Hook properties
        if (inTaskBlock) {
            completions.push(...this.getTaskPropertyCompletions());
        }

        if (inHookBlock) {
            completions.push(...this.getHookPropertyCompletions());
        }

        if (inArgsBlock) {
            completions.push(...this.getArgsPropertyCompletions());
        }

        return completions;
    }

    private isInBlock(document: vscode.TextDocument, position: vscode.Position, blockType: string): boolean {
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

    private getTopLevelCompletions(): vscode.CompletionItem[] {
        return [
            this.createSnippet('set', 'set "${1:VAR_NAME}" "${2:value}"', 'Define a variable'),
            this.createSnippet('default', 'default "${1:task-name}"', 'Set default task'),
            this.createSnippet('alias', 'alias "${1:short}" "${2:task-name}"', 'Create task alias'),
            this.createSnippet('globals', 'globals {\n\t"${1:KEY}" "${2:value}"\n}', 'Define global environment variables'),
            this.createSnippet('hook', 'hook "${1:hook-name}" {\n\tdescription "${2:description}"\n\tcommand "${3:command}"\n}', 'Define a hook'),
            this.createSnippet('task', 'task "${1:task-name}" {\n\tdescription "${2:description}"\n\tcommand "${3:command}"\n}', 'Define a task'),
        ];
    }

    private getTaskPropertyCompletions(): vscode.CompletionItem[] {
        return [
            this.createSnippet('description', 'description "${1:task description}"', 'Task description'),
            this.createSnippet('command', 'command "${1:command to run}"', 'Command to execute'),
            this.createSnippet('dependencies', 'dependencies [${1:"dep1", "dep2"}]', 'Task dependencies'),
            this.createSnippet('before', 'before [${1:"hook1"}]', 'Hooks to run before task'),
            this.createSnippet('after', 'after [${1:"hook1"}]', 'Hooks to run after task'),
            this.createSnippet('on_success', 'on_success [${1:"hook1"}]', 'Hooks to run on success'),
            this.createSnippet('on_failure', 'on_failure [${1:"hook1"}]', 'Hooks to run on failure'),
            this.createSnippet('inputs', 'inputs [${1:"src/**/*.go"}]', 'Input file patterns'),
            this.createSnippet('outputs', 'outputs [${1:"build/output"}]', 'Output file patterns'),
            this.createSnippet('env', 'env {\n\t"${1:KEY}" "${2:value}"\n}', 'Environment variables'),
            this.createSnippet('args', 'args {\n\trequired [${1:"arg1"}]\n}', 'Command arguments'),
            this.createKeyword('cache', 'Enable caching (true/false)'),
            this.createKeyword('parallel', 'Run dependencies in parallel (true/false)'),
            this.createKeyword('silent', 'Suppress output (true/false)'),
            this.createKeyword('watch', 'Enable watch mode (true/false)'),
            this.createKeyword('continue_on_error', 'Continue on error (true/false)'),
            this.createSnippet('timeout', 'timeout "${1:5m}"', 'Execution timeout'),
            this.createSnippet('retry', 'retry ${1:2}', 'Number of retries'),
            this.createSnippet('retry_delay', 'retry_delay "${1:3s}"', 'Delay between retries'),
        ];
    }

    private getHookPropertyCompletions(): vscode.CompletionItem[] {
        return [
            this.createSnippet('description', 'description "${1:hook description}"', 'Hook description'),
            this.createSnippet('command', 'command "${1:command to run}"', 'Command to execute'),
            this.createSnippet('env', 'env {\n\t"${1:KEY}" "${2:value}"\n}', 'Environment variables'),
        ];
    }

    private getArgsPropertyCompletions(): vscode.CompletionItem[] {
        return [
            this.createSnippet('required', 'required [${1:"arg1"}]', 'Required arguments'),
            this.createSnippet('optional', 'optional [${1:"arg1"}]', 'Optional arguments'),
        ];
    }

    private createSnippet(label: string, snippet: string, documentation: string): vscode.CompletionItem {
        const item = new vscode.CompletionItem(label, vscode.CompletionItemKind.Snippet);
        item.insertText = new vscode.SnippetString(snippet);
        item.documentation = new vscode.MarkdownString(documentation);
        return item;
    }

    private createKeyword(label: string, documentation: string): vscode.CompletionItem {
        const item = new vscode.CompletionItem(label, vscode.CompletionItemKind.Keyword);
        item.documentation = new vscode.MarkdownString(documentation);
        return item;
    }
}
