import * as vscode from 'vscode';

export interface SnippetDefinition {
    label: string;
    snippet: string;
    documentation: string;
    kind?: vscode.CompletionItemKind;
}

export interface SnippetsConfig {
    topLevel: SnippetDefinition[];
    taskProperties: SnippetDefinition[];
    hookProperties: SnippetDefinition[];
    argsProperties: SnippetDefinition[];
}

export const snippetsConfig: SnippetsConfig = {
    topLevel: [
        {
            label: 'var',
            snippet: 'var ${1:VAR_NAME} = "${2:value}"',
            documentation: 'Define a variable'
        },
        {
            label: 'default',
            snippet: 'default ${1:task_name}',
            documentation: 'Set default task'
        },
        {
            label: 'alias',
            snippet: 'alias ${1:short} ${2:task_name}',
            documentation: 'Create task alias'
        },
        {
            label: 'globals',
            snippet: 'globals {\n\t"${1:KEY}" "${2:value}"\n}',
            documentation: 'Define global environment variables'
        },
        {
            label: 'hook',
            snippet: 'hook "${1:hook-name}" {\n\tdescription "${2:description}"\n\tcommand "${3:command}"\n}',
            documentation: 'Define a hook'
        },
        {
            label: 'task',
            snippet: 'task ${1:task_name} {\n\tdescription "${2:description}"\n\tcommand "${3:command}"\n}',
            documentation: 'Define a task'
        }
    ],
    taskProperties: [
        {
            label: 'description',
            snippet: 'description "${1:task description}"',
            documentation: 'Task description'
        },
        {
            label: 'command',
            snippet: 'command "${1:command to run}"',
            documentation: 'Command to execute'
        },
        {
            label: 'dependencies',
            snippet: 'dependencies [${1:"dep1", "dep2"}]',
            documentation: 'Task dependencies'
        },
        {
            label: 'before',
            snippet: 'before [${1:"hook1"}]',
            documentation: 'Hooks to run before task'
        },
        {
            label: 'after',
            snippet: 'after [${1:"hook1"}]',
            documentation: 'Hooks to run after task'
        },
        {
            label: 'on_success',
            snippet: 'on_success [${1:"hook1"}]',
            documentation: 'Hooks to run on success'
        },
        {
            label: 'on_failure',
            snippet: 'on_failure [${1:"hook1"}]',
            documentation: 'Hooks to run on failure'
        },
        {
            label: 'inputs',
            snippet: 'inputs [${1:"src/**/*.go"}]',
            documentation: 'Input file patterns'
        },
        {
            label: 'outputs',
            snippet: 'outputs [${1:"build/output"}]',
            documentation: 'Output file patterns'
        },
        {
            label: 'env',
            snippet: 'env {\n\t"${1:KEY}" "${2:value}"\n}',
            documentation: 'Environment variables'
        },
        {
            label: 'args',
            snippet: 'args {\n\trequired [${1:"arg1"}]\n}',
            documentation: 'Command arguments'
        },
        {
            label: 'cache',
            snippet: 'cache',
            documentation: 'Enable caching (true/false)',
            kind: vscode.CompletionItemKind.Keyword
        },
        {
            label: 'parallel',
            snippet: 'parallel',
            documentation: 'Run dependencies in parallel (true/false)',
            kind: vscode.CompletionItemKind.Keyword
        },
        {
            label: 'silent',
            snippet: 'silent',
            documentation: 'Suppress output (true/false)',
            kind: vscode.CompletionItemKind.Keyword
        },
        {
            label: 'watch',
            snippet: 'watch',
            documentation: 'Enable watch mode (true/false)',
            kind: vscode.CompletionItemKind.Keyword
        },
        {
            label: 'continue_on_error',
            snippet: 'continue_on_error',
            documentation: 'Continue on error (true/false)',
            kind: vscode.CompletionItemKind.Keyword
        },
        {
            label: 'timeout',
            snippet: 'timeout "${1:5m}"',
            documentation: 'Execution timeout'
        },
        {
            label: 'retry',
            snippet: 'retry ${1:2}',
            documentation: 'Number of retries'
        },
        {
            label: 'retry_delay',
            snippet: 'retry_delay "${1:3s}"',
            documentation: 'Delay between retries'
        }
    ],
    hookProperties: [
        {
            label: 'description',
            snippet: 'description "${1:hook description}"',
            documentation: 'Hook description'
        },
        {
            label: 'command',
            snippet: 'command "${1:command to run}"',
            documentation: 'Command to execute'
        },
        {
            label: 'env',
            snippet: 'env {\n\t"${1:KEY}" "${2:value}"\n}',
            documentation: 'Environment variables'
        }
    ],
    argsProperties: [
        {
            label: 'required',
            snippet: 'required [${1:"arg1"}]',
            documentation: 'Required arguments'
        },
        {
            label: 'optional',
            snippet: 'optional [${1:"arg1"}]',
            documentation: 'Optional arguments'
        }
    ]
};
