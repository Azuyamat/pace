import * as cp from 'child_process';
import * as vscode from 'vscode';
import { paceSchema } from './generated/paceSchema';
import { log, debugLog } from './output';

interface RuntimeSchemaInfo {
    version?: number;
}

function configuredPacePath(): string {
    return vscode.workspace.getConfiguration('pace').get<string>('path', 'pace') || 'pace';
}

function versionWarningsEnabled(): boolean {
    return vscode.workspace.getConfiguration('pace').get<boolean>('schema.versionWarnings.enabled', true);
}

function readRuntimeSchemaVersion(command: string): Promise<number | undefined> {
    return new Promise(resolve => {
        cp.execFile(command, ['schema', '--json'], { windowsHide: true, timeout: 5000, maxBuffer: 1024 * 1024 }, (error, stdout, stderr) => {
            if (error) {
                debugLog(`Unable to read Pace CLI schema version: ${stderr || error.message}`);
                resolve(undefined);
                return;
            }
            try {
                const parsed = JSON.parse(stdout) as RuntimeSchemaInfo;
                resolve(typeof parsed.version === 'number' ? parsed.version : undefined);
            } catch (parseError) {
                debugLog(`Unable to parse Pace CLI schema JSON: ${String(parseError)}`);
                resolve(undefined);
            }
        });
    });
}

function major(version: number): number {
    // Schema versions are currently integers. Reserve the major calculation for
    // future semver-like values while keeping v1 compatible with existing CLIs.
    return Math.trunc(version);
}

export async function checkSchemaVersionCompatibility(): Promise<void> {
    const bundledVersion = paceSchema.version;
    log(`Bundled Pace language schema version: ${bundledVersion}`);

    if (!versionWarningsEnabled()) {
        return;
    }

    const runtimeVersion = await readRuntimeSchemaVersion(configuredPacePath());
    if (runtimeVersion === undefined) {
        return;
    }

    log(`Pace CLI language schema version: ${runtimeVersion}`);
    if (major(runtimeVersion) !== major(bundledVersion)) {
        vscode.window.showWarningMessage(
            `Pace CLI schema version ${runtimeVersion} differs from bundled extension schema version ${bundledVersion}. Some language features may not match.`,
            'Show Pace Logs'
        ).then(selection => {
            if (selection === 'Show Pace Logs') {
                vscode.commands.executeCommand('workbench.action.output.toggleOutput');
            }
        });
    }
}

export async function showSchemaVersionCommand(): Promise<void> {
    const bundledVersion = paceSchema.version;
    const runtimeVersion = await readRuntimeSchemaVersion(configuredPacePath());
    const runtimeText = runtimeVersion === undefined ? 'unavailable (Pace CLI not found or too old)' : String(runtimeVersion);
    vscode.window.showInformationMessage(`Pace schema versions — bundled: ${bundledVersion}, CLI: ${runtimeText}`);
}
