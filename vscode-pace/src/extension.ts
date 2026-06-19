import * as vscode from 'vscode';
import { PaceCodeLensProvider, runPaceTask } from './codeLensProvider';
import { PaceCompletionProvider } from './completionProvider';
import { PaceHoverProvider } from './hoverProvider';
import { PaceDocumentSymbolProvider } from './documentSymbolProvider';
import { PaceDiagnosticsManager } from './diagnostics';
import { log } from './output';
import { checkSchemaVersionCompatibility, showSchemaVersionCommand } from './versionCompatibility';
import {
    generateConfigCommand,
    insertGoTaskTemplateCommand,
    runDefaultTaskCommand,
    runSelectedOrCurrentTaskCommand,
    runTaskCommand
} from './commands';

export function activate(context: vscode.ExtensionContext) {
    log('Pace language extension is now active');

    const completionProvider = vscode.languages.registerCompletionItemProvider(
        'pace',
        new PaceCompletionProvider(),
        ' ', '"'
    );

    const documentSymbolProvider = vscode.languages.registerDocumentSymbolProvider(
        'pace',
        new PaceDocumentSymbolProvider()
    );

    const codeLensProvider = vscode.languages.registerCodeLensProvider(
        'pace',
        new PaceCodeLensProvider()
    );

    const hoverProvider = vscode.languages.registerHoverProvider(
        'pace',
        new PaceHoverProvider()
    );

    const runTaskFromCodeLensCommand = vscode.commands.registerCommand(
        'pace.runTaskFromCodeLens',
        (taskName: string, document?: vscode.TextDocument) => runPaceTask(taskName, document)
    );

    const watchTaskCommand = vscode.commands.registerCommand(
        'pace.watchTaskFromCodeLens',
        (taskName: string, document?: vscode.TextDocument) => runPaceTask(taskName, document, 'watch')
    );

    const generateConfig = vscode.commands.registerCommand('pace.generateConfig', generateConfigCommand);
    const runCurrentTask = vscode.commands.registerCommand('pace.runCurrentTask', runSelectedOrCurrentTaskCommand);
    const runTask = vscode.commands.registerCommand('pace.runTask', runTaskCommand);
    const runDefaultTask = vscode.commands.registerCommand('pace.runDefaultTask', runDefaultTaskCommand);
    const insertGoTaskTemplate = vscode.commands.registerCommand('pace.insertGoTaskTemplate', insertGoTaskTemplateCommand);
    const showSchemaVersion = vscode.commands.registerCommand('pace.showSchemaVersion', showSchemaVersionCommand);

    void checkSchemaVersionCompatibility();

    context.subscriptions.push(
        completionProvider,
        documentSymbolProvider,
        codeLensProvider,
        hoverProvider,
        new PaceDiagnosticsManager(),
        runTaskFromCodeLensCommand,
        watchTaskCommand,
        generateConfig,
        runCurrentTask,
        runTask,
        runDefaultTask,
        insertGoTaskTemplate,
        showSchemaVersion
    );
}

export function deactivate() {}
