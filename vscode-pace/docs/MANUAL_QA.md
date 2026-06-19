# Pace VS Code Extension Manual QA Checklist

Use the examples in `../examples/` and at least one real Pace workspace before publishing.

## Syntax highlighting

- [ ] `.pace` files and `config.pace` open with the Pace language mode.
- [ ] Top-level statements, task/hook properties, strings, comments, arrays, booleans, numbers, variables, and multiline strings highlight correctly.
- [ ] Braces and keywords inside comments/strings do not confuse highlighting or context-aware features.

## Completions and snippets

- [ ] Top-level completions appear on blank and partially typed lines.
- [ ] Task, hook, `args`, and `env` block contexts show only appropriate suggestions.
- [ ] Task and hook reference completions use names from the current document.
- [ ] Completion details, examples, deprecated markers, and replacement ranges behave correctly.

## Commands and CodeLens

- [ ] `Pace: Generate config.pace` works and asks before overwrite.
- [ ] `Pace: Run Current Task`, `Pace: Run Task`, and `Pace: Run Default Task` use the correct workspace folder.
- [ ] `Pace: Insert Generated Go Task` uses the Pace CLI and inserts valid syntax in Go workspaces.
- [ ] `Run task`/`Watch task` CodeLens appears only in Pace files and respects `pace.codeLens.enabled`.

## Diagnostics

- [ ] Valid examples have no diagnostics.
- [ ] `invalid-diagnostics.pace` reports parser/validation errors.
- [ ] Diagnostics update according to `pace.diagnostics.run` and stop when disabled.
- [ ] Missing Pace CLI failures are concise and logged to the Pace output channel.

## Schema compatibility

- [ ] `Pace: Show Schema Version` reports the bundled schema version and CLI schema version when available.
- [ ] `pace.schema.versionWarnings.enabled` controls mismatch warnings.

## Packaging and local install

- [ ] `npm ci`, `npm run compile`, `npm test`, and `npm run package` pass from a clean checkout.
- [ ] The generated VSIX installs and activates in a clean VS Code profile.
- [ ] The VSIX does not include source/test/dev-only files excluded by `.vscodeignore`.

## Platform notes

- [ ] Windows: commands work with workspace paths containing spaces; PowerShell users can install VSIX through the Command Palette if glob expansion fails.
- [ ] macOS/Linux: commands and terminals use the selected workspace folder.
- [ ] Multi-root: commands and diagnostics choose the active Pace document's workspace or prompt when needed.

## Telemetry

- [ ] No telemetry prompts or network telemetry are present; README states that no telemetry is collected.
