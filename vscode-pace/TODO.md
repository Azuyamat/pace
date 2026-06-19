# vscode-pace TODO

This TODO tracks improvements for the Pace VS Code extension. Each task includes a title, description, and acceptance criteria.

---

## 1. ✅ Align TextMate grammar keywords with the Go parser

### Description
Update `syntaxes/pace.tmLanguage.json` so top-level keywords match the real Pace parser definitions in `internal/config/parsing/statement_parser.go`.

Current grammar includes unsupported or questionable keywords such as `set` and `globals`, and misses supported statements such as `var` and `import`.

### Acceptance criteria
- `var` is highlighted as a top-level keyword.
- `import` is highlighted as a top-level keyword.
- Unsupported keywords are removed unless they are intentionally supported by the parser.
- Grammar keyword list matches the parser-supported top-level statements: `task`, `hook`, `var`, `default`, `alias`, `import`.
- Existing syntax highlighting still works for tasks, hooks, strings, comments, variables, booleans, numbers, arrays, and braces.

---

## 2. ✅ Align task property highlighting with the Go parser

### Description
Update task property highlighting to include every supported task property from `internal/config/parsing/property_parser.go`.

Supported task properties include:

- `command`
- `inputs`
- `outputs`
- `depends-on`
- `dependencies`
- `env`
- `cache`
- `working_dir`
- `requires`
- `before`
- `triggers`
- `after`
- `description`
- `watch`
- `parallel`
- `silent`
- `continue_on_error`
- `timeout`
- `retry`
- `retry_delay`
- `on_success`
- `on_failure`
- `when`
- `args`

### Acceptance criteria
- Every parser-supported task property is highlighted.
- Deprecated/alias properties are still highlighted if supported by the parser.
- No unsupported task properties are documented or highlighted unless intentionally supported.
- Grammar and README property lists agree with the Go parser.

---

## 3. ✅ Align hook property highlighting with the Go parser

### Description
Update hook property highlighting to include every supported hook property from `internal/config/parsing/property_parser.go`.

Supported hook properties:

- `command`
- `env`
- `working_dir`
- `description`

### Acceptance criteria
- Every parser-supported hook property is highlighted.
- Unsupported hook properties are removed from docs/snippets unless intentionally supported.
- Hook property list in grammar, README, snippets, and Go parser are consistent.

---

## 4. ✅ Fix `env` snippet syntax

### Description
The current `env` snippet uses invalid syntax:

```pace
env {
    "KEY" "value"
}
```

The parser expects assignment syntax:

```pace
env {
    KEY = "value"
}
```

### Acceptance criteria
- Task `env` snippet inserts valid Pace syntax.
- Hook `env` snippet inserts valid Pace syntax.
- Snippet allows identifier keys and string values.
- README examples use valid `env` syntax.
- Manual test confirms generated snippet parses successfully with Pace.

---

## 5. ✅ Improve boolean property snippets

### Description
Boolean property completions currently insert only the property name, such as:

```pace
cache
```

The parser expects a boolean value:

```pace
cache true
```

Update boolean snippets to include a `true`/`false` choice.

### Acceptance criteria
- `cache` completion inserts `cache ${1|true,false|}` or equivalent.
- `parallel` completion inserts a boolean value.
- `silent` completion inserts a boolean value.
- `watch` completion inserts a boolean value.
- `continue_on_error` completion inserts a boolean value.
- Generated snippets are valid Pace syntax.

---

## 6. ✅ Add missing top-level completions

### Description
Add completions for all parser-supported top-level statements.

Currently completions are missing at least `import`.

### Acceptance criteria
- Completion list includes `task`.
- Completion list includes `hook`.
- Completion list includes `var`.
- Completion list includes `default`.
- Completion list includes `alias`.
- Completion list includes `import`.
- `alias` documentation indicates the inline alias syntax is preferred if standalone aliases remain deprecated.
- Top-level snippets produce valid Pace syntax.

---

## 7. ✅ Add missing task property completions

### Description
Add completions for all parser-supported task properties that are currently missing from `src/snippets.config.ts`.

Missing or incomplete candidates include:

- `depends-on`
- `requires`
- `triggers`
- `working_dir`
- `when`

### Acceptance criteria
- Completion list inside task blocks includes every supported task property.
- Alias properties such as `before`/`requires`, `after`/`triggers`, and `depends-on`/`dependencies` are documented clearly.
- Snippets produce valid Pace syntax.
- README examples match the snippets.

---

## 8. ✅ Add missing hook property completions

### Description
Add completions for all parser-supported hook properties.

Currently `working_dir` is missing from hook snippets.

### Acceptance criteria
- Completion list inside hook blocks includes `description`, `command`, `env`, and `working_dir`.
- Snippets produce valid Pace syntax.
- Hook README section matches actual completions.

---

## 9. ✅ Add args block completions

### Description
Improve completions inside `args` blocks.

Supported args properties:

- `required`
- `optional`

### Acceptance criteria
- Inside `args {}`, completions include `required` and `optional`.
- Snippets produce valid array syntax.
- Completions do not incorrectly include task-level properties inside `args` blocks.

---

## 10. ✅ Relax or remove `hasTypedContent`

### Description
`ContextDetector.hasTypedContent` currently prevents completions on empty or indented lines. This hurts editing inside empty task/hook blocks because users do not get suggestions until they type something.

### Acceptance criteria
- Users can trigger completions on a blank line inside a task block.
- Users can trigger completions on a blank line inside a hook block.
- Users can trigger completions on a blank top-level line.
- Completions are still suppressed in comments and strings if that behavior is implemented.


### Batch 1 progress
Completed in this batch: grammar keywords/properties aligned with the Go parser; task/hook/env/args snippets corrected and completed; boolean snippets now include values; blank-line completions are enabled; args/env contexts no longer receive task/hook property completions; README examples and supported syntax were updated for parser-valid syntax.

---

## 11. ✅ Replace brittle context detection with a block stack scanner

### Description
`ContextDetector.isInBlock` currently scans lines using simple `includes` calls and brace counts. This can break when braces appear in strings/comments or when blocks are nested.

Replace it with a lightweight scanner that:

- ignores comments
- ignores strings and multiline strings
- tracks `{` and `}`
- tracks named blocks such as `task`, `hook`, `args`, and `env`
- returns the current block stack at the cursor position

### Acceptance criteria
- Context detection works with nested blocks, e.g. `task { env { ... } }`.
- Braces inside strings do not affect context detection.
- Braces inside comments do not affect context detection.
- `args` block is correctly detected inside tasks.
- `env` block is correctly detected inside tasks and hooks.
- Unit tests cover common and edge cases.

---

## 12. ✅ Use env context or remove unused `inEnvBlock`

### Description
`DocumentContext` includes `inEnvBlock`, but the completion provider does not currently use it.

Decide whether to support env-specific completions or remove the unused field.

### Acceptance criteria
- If env completions are supported, `inEnvBlock` is used by the completion provider.
- If env completions are not supported, `inEnvBlock` is removed to reduce dead code.
- TypeScript compile passes.

---

## 13. ✅ Add task and hook name extraction from the current document

### Description
Implement a parser/scanner in the extension that extracts task and hook names from the open `.pace` document.

This will power smarter completions for references.

### Acceptance criteria
- Extension can list all task names in the current document.
- Extension can list all task aliases in the current document.
- Extension can list all hook names in the current document.
- Extraction ignores commented-out tasks/hooks.
- Extraction handles simple whitespace variations.
- Unit tests cover extraction behavior.

---

## 14. ✅ Add task reference completions

### Description
Use extracted task names to provide completions in task-reference contexts.

Relevant properties:

- `depends-on`
- `dependencies`
- `default`
- `alias` target argument

### Acceptance criteria
- Inside `depends-on [...]`, task names are suggested.
- Inside `dependencies [...]`, task names are suggested.
- After `default `, task names are suggested.
- In standalone `alias short target`, task names are suggested for the target.
- The current task is not suggested as its own dependency, or is deprioritized.

---

## 15. ✅ Add hook reference completions

### Description
Use extracted hook names to provide completions in hook-reference contexts.

Relevant properties:

- `requires`
- `before`
- `triggers`
- `after`
- `on_success`
- `on_failure`

### Acceptance criteria
- Inside `requires [...]`, hook names are suggested.
- Inside `before [...]`, hook names are suggested.
- Inside `triggers [...]`, hook names are suggested.
- Inside `after [...]`, hook names are suggested.
- Inside `on_success [...]`, hook names are suggested.
- Inside `on_failure [...]`, hook names are suggested.

---

### Batch 2 progress
Completed in this batch: context detection now uses a string/comment-aware block stack scanner; `env` context is explicitly handled by completions; task names, task aliases, and hook names are extracted from the current document; task and hook reference completions are provided in supported reference contexts; unit tests cover scanner and extraction behavior.

---

## 16. ✅ Add document symbols

### Description
Implement a `DocumentSymbolProvider` so VS Code's Outline panel shows Pace tasks and hooks.

### Acceptance criteria
- Each `task` appears as a symbol.
- Each `hook` appears as a symbol.
- Symbols include accurate ranges.
- Symbols are nested or categorized clearly where appropriate.
- Outline works for valid `.pace` files.
- Malformed files do not crash the extension.

---

## 17. ✅ Add CodeLens actions for tasks

### Description
Add CodeLens actions above task definitions for common operations.

Potential actions:

- Run task
- Watch task
- Run without cache

### Acceptance criteria
- Each task has a `Run task` CodeLens.
- `Run task` executes `pace run <task>` in a VS Code terminal.
- Tasks with `watch true` can expose a `Watch task` CodeLens, or all tasks can expose it if Pace supports watch execution.
- CodeLens can be disabled by configuration.
- CodeLens does not appear in non-Pace files.

---

## 18. ✅ Add CodeLens actions for hooks if useful

### Description
Evaluate whether hooks should expose CodeLens actions. Hooks may not be directly runnable depending on Pace semantics.

### Acceptance criteria
- Decision is documented.
- If hooks are runnable, add appropriate CodeLens actions.
- If hooks are not runnable, no misleading hook CodeLens is shown.

---

## 19. ✅ Add hover documentation for statements and properties

### Description
Implement hover docs for Pace keywords and properties using a shared schema.

Examples:

- Hovering `cache` explains smart caching and input requirements.
- Hovering `depends-on` explains task dependencies.
- Hovering `requires` explains before-task hooks.

### Acceptance criteria
- Hover docs exist for every top-level statement.
- Hover docs exist for every task property.
- Hover docs exist for every hook property.
- Hover docs exist for args properties.
- Docs match parser behavior.
- Hover provider does not show irrelevant docs inside strings/comments.

---

### Batch 3 progress
Completed in this batch: document symbols now list task and hook definitions with ranges; task CodeLens actions run `pace run <task>` and watch-enabled tasks expose `Watch task`; CodeLens can be disabled with `pace.codeLens.enabled`; hook CodeLens was intentionally not added because hooks are task-referenced and not directly runnable; hover docs now cover top-level statements, task properties, hook properties, and args properties while suppressing hovers in strings/comments.

---

## 20. Add basic diagnostics from extension-side parsing

### Description
Before integrating the real Pace parser, add lightweight diagnostics for obvious document issues.

Potential diagnostics:

- unmatched braces
- unknown top-level statement
- unknown task property
- unknown hook property
- duplicate task names
- duplicate hook names

### Acceptance criteria
- Diagnostics appear in VS Code Problems panel.
- Diagnostics update when the document changes.
- Diagnostics are scoped to `.pace` files.
- Diagnostics do not produce false positives for valid examples from README.
- Extension does not crash on malformed files.

---

## 21. ✅ Add `pace validate --json` to the Go CLI

### Description
Add a machine-readable validation command to the Pace CLI so editors can use the real parser and validator.

Proposed command:

```bash
pace validate --json config.pace
```

The JSON output should include file, range, severity, message, and optional hint.

### Acceptance criteria
- CLI validates a Pace config without running tasks.
- CLI can output JSON diagnostics.
- JSON includes enough location information for VS Code diagnostics.
- Validation catches parser errors.
- Validation catches semantic errors from the existing validator.
- Non-JSON output remains human-friendly.

---

## 22. ✅ Integrate VS Code diagnostics with `pace validate --json`

### Description
Use the real Pace CLI validation command from the extension to provide accurate diagnostics.

### Acceptance criteria
- Extension detects the workspace Pace binary or uses configured path.
- Extension runs validation on save or debounce after edits.
- Parser/validation errors appear in the Problems panel.
- Diagnostics include line/column ranges when available.
- Users can disable CLI diagnostics in settings.
- Missing Pace binary produces a non-intrusive warning, not repeated errors.

---

## 23. ✅ Add extension settings

### Description
Add VS Code configuration settings for extension behavior.

Suggested settings:

- `pace.path`
- `pace.diagnostics.enabled`
- `pace.diagnostics.run`
- `pace.codeLens.enabled`
- `pace.completions.projectAware`

### Acceptance criteria
- Settings are contributed in `package.json`.
- Settings have descriptions and defaults.
- Extension reads settings dynamically.
- Settings are documented in `README.md`.

---

## 24. ✅ Add `Pace: Generate config.pace` command

### Description
Add a VS Code command that initializes a Pace config for the current workspace.

Preferred implementation: call the Pace CLI's existing `pace init` functionality rather than duplicating project detection in TypeScript.

### Acceptance criteria
- Command appears in Command Palette.
- If `config.pace` does not exist, command can generate one.
- If `config.pace` exists, user is asked before overwriting.
- Command uses workspace root as the working directory.
- Command surfaces success/failure messages in VS Code.
- Command works for Go projects using existing Pace project detection.

---

## 25. ✅ Add project-aware Go task snippets

### Description
Provide commands or snippets that infer useful Go task definitions from the current Go project.

Potential inferred values:

- module name from `go.mod`
- main package from `main.go` or `cmd/*/main.go`
- output binary path
- presence of `.golangci` config
- presence of `//go:generate`

### Acceptance criteria
- Extension can suggest a Go build task with correct main package.
- Extension can suggest a Go run task with correct main package.
- Extension can suggest Go test, vet, tidy tasks.
- Extension suggests lint only when a known lint config exists.
- Logic is sourced from Pace CLI if possible, or clearly kept in sync with Go implementation.

---

## 26. ✅ Avoid duplicating Go project detection logic in TypeScript

### Description
The Go project generator already exists in `internal/template/generator/go.go`. Avoid reimplementing and drifting from it in the extension.

Instead, expose a CLI command that emits generated config or suggestions as JSON.

Example:

```bash
pace init --type go --print-json
```

or:

```bash
pace schema project-template --type go --json
```

### Acceptance criteria
- VS Code extension can retrieve generated project/task suggestions from the Pace CLI.
- Go detection logic remains primarily in Go.
- TypeScript code does not duplicate full project detection behavior.
- Generated results are stable and testable.

---

## 27. ✅ Create a Go-owned Pace language schema

### Description
Introduce a single source of truth for Pace language statements and properties in Go.

This schema should describe:

- top-level statements
- task properties
- hook properties
- args properties
- property types
- aliases/deprecated forms
- snippets
- documentation

### Acceptance criteria
- Schema exists in a Go package, for example `internal/config/schema`.
- Parser registries can be generated from or directly use the schema.
- VS Code snippets can be generated from the schema.
- Documentation can be generated from the schema.
- Schema covers every currently supported parser feature.

---

## 28. ✅ Generate VS Code snippets from the Go schema

### Description
Replace manually maintained `src/snippets.config.ts` with generated data from the Go-owned schema.

### Acceptance criteria
- A generation command creates TypeScript snippet configuration.
- Generated snippets include all supported statements/properties.
- Generated snippets produce valid Pace syntax.
- Manual edits to generated files are discouraged with a header comment.
- `npm run compile` uses generated snippets.

---

## 29. ✅ Generate grammar keyword lists from the Go schema

### Description
Use the Go-owned schema to generate keyword/property lists for TextMate grammar.

### Acceptance criteria
- Generated grammar keyword lists include all supported statements and properties.
- Generated output is consumed by or merged into `pace.tmLanguage.json`.
- Keyword drift between parser and grammar is eliminated.
- Generation is documented.

---

## 30. ✅ Generate documentation from the Go schema

### Description
Use the Go-owned schema to generate language reference documentation for README and/or docs site.

### Acceptance criteria
- Generated docs include every top-level statement.
- Generated docs include every task property.
- Generated docs include every hook property.
- Generated docs include examples for property types.
- Docs indicate aliases and deprecated syntax.
- Generated docs match parser behavior.

---

## 31. ✅ Add a generator command for VS Code artifacts

### Description
Add a script or Go command to generate extension artifacts.

Potential command:

```bash
go run ./cmd/pace dev generate-vscode
```

or:

```bash
go generate ./...
```

### Acceptance criteria
- One command regenerates snippets and grammar keyword data.
- Command is documented in `vscode-pace/README.md` or `CONTRIBUTING.md`.
- Generated files are deterministic.
- CI can verify generated files are up to date.

---

## 32. ✅ Add generated-file freshness check to CI

### Description
Ensure generated extension artifacts are committed and up to date, or generated during build.

### Acceptance criteria
- CI runs the generator.
- CI fails if generated files differ from committed files, if committed generated files are required.
- Alternatively, CI builds from generated files without requiring generated artifacts to be committed.
- Developer workflow is documented.

---

## 33. ✅ Add extension unit tests

### Description
Add automated tests for extension utility logic and completion behavior.

Test targets:

- context detection
- snippet generation/factory
- task/hook extraction
- completion provider behavior
- diagnostics parser if implemented

### Acceptance criteria
- Test framework is configured.
- Tests run with `npm test`.
- Completion tests cover top-level, task, hook, args, and env contexts.
- Context tests cover nested blocks, comments, and strings.
- CI runs the tests.

---

## 34. ✅ Add compile/package CI for `vscode-pace`

### Description
Add a GitHub Actions job for the VS Code extension.

### Acceptance criteria
- CI runs `npm ci` inside `vscode-pace`.
- CI runs `npm run compile`.
- CI runs tests when available.
- CI optionally runs `npm run package`.
- CI fails on TypeScript errors.

---

## 35. ✅ Add `@vscode/vsce` as a dev dependency

### Description
The README currently instructs installing `vsce` globally. Add it as a dev dependency for reproducible packaging.

### Acceptance criteria
- `@vscode/vsce` is listed in `devDependencies`.
- `npm run package` packages the extension without global tools.
- README no longer requires global `vsce` installation unless optional.

---

## 36. ✅ Add package/install/dev scripts to `package.json`

### Description
Improve extension development workflow with scripts.

Suggested scripts:

```json
"package": "vsce package",
"install:local": "code --install-extension pace-language-*.vsix --force",
"dev": "npm run compile && code --extensionDevelopmentPath=."
```

### Acceptance criteria
- `npm run package` creates a `.vsix`.
- `npm run install:local` installs the generated extension locally.
- `npm run dev` launches an extension development host or documents how to do so.
- Scripts work on supported development platforms, or platform caveats are documented.

---

## 37. ✅ Add root VS Code tasks for extension development

### Description
Add workspace tasks for common extension actions.

Potential tasks:

- compile extension
- watch extension
- package extension
- install extension locally

### Acceptance criteria
- Tasks are available from VS Code's task runner.
- Tasks run from the repository root.
- Tasks use `vscode-pace` as the working directory where appropriate.
- Tasks are documented or clearly named.

---

## 38. ✅ Improve `.vscodeignore`

### Description
Reduce packaged extension size and avoid shipping unnecessary files.

Candidates to exclude:

- `src/**`
- `tsconfig.json`
- `*.map`
- `.vscode/**`
- test files
- local config files

Only exclude source maps if debugging packaged extensions is not needed.

### Acceptance criteria
- VSIX contains only runtime-required files.
- Extension still works after packaging.
- Packaged size is reduced or justified.
- `.vscodeignore` is documented enough to avoid accidental breakage.

---

## 39. ✅ Decide whether to commit `build/`

### Description
`build/` exists but is ignored by `.gitignore`. Decide whether generated JS should be committed.

Options:

1. Do not commit `build/`; compile during packaging.
2. Commit `build/`; ensure it is always up to date.

### Acceptance criteria
- Decision is documented.
- Repository state matches the decision.
- CI enforces the decision.
- Packaging process works from a clean clone.

---

## 40. ✅ Improve README accuracy

### Description
Update `vscode-pace/README.md` so it accurately reflects extension features and the real Pace language.

Current README contains some drift, including properties not fully supported by snippets/grammar and examples that may not match parser behavior.

### Acceptance criteria
- README lists only supported syntax.
- README examples parse successfully with Pace.
- README describes autocomplete, not just syntax highlighting.
- Installation instructions use local `npm run package` after `@vscode/vsce` is added.
- README documents extension settings and commands when implemented.


### Batch 4 progress
Completed in this batch: extension unit tests now include completion behavior across top-level, task, hook, args, env, and reference contexts; CI installs, compiles, tests, and packages the VS Code extension; `@vscode/vsce` is a dev dependency; package/install/dev scripts and root VS Code tasks were added; `.vscodeignore` now excludes source/test/dev-only files from VSIX packaging; `build/` remains uncommitted by design and is compiled during CI/packaging; README installation/development instructions now use local npm scripts and accurately describe current features.

---

## 41. ✅ Add a language reference section or link

### Description
Add a concise language reference to the extension README, or link to generated Pace docs.

### Acceptance criteria
- Users can find supported statements.
- Users can find supported task properties.
- Users can find supported hook properties.
- Users can find examples for arrays, maps, booleans, strings, multiline strings, and variables.
- Reference is generated or clearly tied to the Go schema to prevent drift.

---

## 42. Add snippets for common full task templates

### Description
Add higher-level snippets for common complete task definitions.

Potential snippets:

- basic task
- cached task
- Go build task
- Go test task
- task with args
- task with hooks

### Acceptance criteria
- Snippets produce valid Pace syntax.
- Snippets are not noisy in completion results.
- Snippets include helpful placeholders.
- Snippets are documented.

---

## 43. Add formatter support

### Description
Implement a basic document formatter for `.pace` files.

Formatting behavior may include:

- normalize indentation
- align block indentation
- keep braces consistent
- preserve comments
- normalize simple arrays if safe

### Acceptance criteria
- `Format Document` works for `.pace` files.
- Formatter preserves comments.
- Formatter preserves strings and multiline strings.
- Formatter does not change semantic content.
- Malformed files are handled gracefully.
- Tests cover representative formatting cases.

---

## 44. ✅ Add command to run selected/current task

### Description
Add a command that runs the task under the cursor.

### Acceptance criteria
- Command appears in Command Palette.
- If cursor is inside a task block, it runs that task.
- If cursor is not inside a task, user is prompted to select a task from the document.
- Task is executed in a VS Code terminal.
- Terminal uses the workspace folder as current directory.

---

## 45. ✅ Add command to list and run tasks

### Description
Add a command that lists all tasks in the current Pace config and lets the user run one.

### Acceptance criteria
- Command appears as `Pace: Run Task`.
- Quick Pick lists task names.
- Quick Pick includes descriptions when available.
- Selecting a task runs `pace run <task>` in a terminal.
- Aliases are displayed when available.

---

## 46. ✅ Add support for default task execution

### Description
Add a command to run the default Pace task.

### Acceptance criteria
- Command appears as `Pace: Run Default Task`.
- Command executes `pace run` in a terminal.
- If no workspace folder exists, user receives a useful error.
- Terminal is reused or named clearly.

---

## 47. ✅ Add problem matcher for Pace tasks

### Description
If Pace emits structured or recognizable errors, add a VS Code problem matcher for task execution.

### Acceptance criteria
- Problem matcher captures Pace config parse errors when possible.
- Problem matcher captures validation errors when possible.
- Problem matcher is documented.
- It does not produce misleading diagnostics for normal command output.

---

## 48. ✅ Add multi-root workspace handling

### Description
Ensure commands and diagnostics work correctly in VS Code multi-root workspaces.

### Acceptance criteria
- Commands choose the correct workspace folder for the active `.pace` file.
- Diagnostics are scoped per workspace folder.
- `Pace: Generate config.pace` prompts for folder if needed.
- Terminals run in the correct folder.

---

## 49. ✅ Add Windows path and shell considerations

### Description
Pace is cross-platform, and the extension should avoid Unix-only assumptions.

### Acceptance criteria
- Commands work on Windows, macOS, and Linux.
- Terminal command construction handles paths with spaces.
- Local install/package scripts document platform limitations or provide cross-platform alternatives.
- Generated examples avoid hardcoded OS-specific shell commands unless intentional.

---

### Batch 6 progress
Completed in this batch: added VS Code commands to generate `config.pace`, run the current/selected/default task, and insert Go task templates; task execution now validates task names and uses safer platform-aware shell quoting; command workspace selection handles multi-root workspaces; a `$pace` problem matcher is contributed; the Pace CLI initializer now supports non-interactive `--yes`, `--force`, `--stdout`, and `--json` modes so editor integrations can source Go project detection from Go instead of duplicating it in TypeScript; Go project scanning now skips large generated/vendor directories.

---

## 50. ✅ Add cancellation and debounce for expensive operations

### Description
Diagnostics and CLI calls should not run too often or continue after newer edits.

### Acceptance criteria
- Document validation is debounced.
- In-flight validation can be ignored or cancelled when a newer edit happens.
- Extension remains responsive on large files.
- CLI failures do not spam users.

---

## 51. ✅ Add logging/output channel

### Description
Add a dedicated VS Code output channel for Pace extension logs.

### Acceptance criteria
- Extension creates a `Pace` output channel.
- CLI command failures are logged there.
- Debug information can be enabled via setting.
- User-facing notifications remain concise.

---

## 52. ✅ Add activation events explicitly if needed

### Description
Review extension activation behavior. Modern VS Code may infer activation from language contributions, but explicit activation events can improve clarity and compatibility.

### Acceptance criteria
- `package.json` activation behavior is reviewed.
- If needed, `activationEvents` includes `onLanguage:pace` and Pace commands.
- Extension activates only when useful.
- Commands work after activation.


### Batch 5 progress
Completed in this batch: added `pace validate --json` with parser/validator diagnostics; VS Code diagnostics now call the CLI with debounce and stale-run cancellation; extension settings now cover CLI path, diagnostics, CodeLens, project-aware completions, and debug logging; a `Pace` output channel logs CLI and debug details; activation events explicitly cover Pace files and commands.

---

## 53. ✅ Improve completion item metadata

### Description
Make completion results more useful with details, docs, sorting, and replacement ranges.

### Acceptance criteria
- Completion items include `detail` such as `Pace task property`.
- Completion docs include examples.
- Deprecated items are marked deprecated where applicable.
- Completion ordering prioritizes common properties.
- Replacement range avoids duplicating partially typed text.

---

## 54. Add completion support for variables

### Description
Provide completions for variables defined with `var` and environment variables where appropriate.

### Acceptance criteria
- Variables defined in the current document are extracted.
- Inside strings, `${var}` completions are offered.
- `$var` completions are offered where appropriate.
- Built-in/context variables are included if Pace has any.
- Completions do not trigger in irrelevant contexts.

---

## 55. Add completion support for file globs

### Description
Provide helpful snippets or completions for common input/output glob patterns.

Examples:

- `**/*.go`
- `go.mod`
- `go.sum`
- `cmd/**/*.go`
- `internal/**/*.go`

### Acceptance criteria
- Glob snippets appear in `inputs` and `outputs` array contexts.
- Go-specific globs are only suggested when a Go project is detected or project-aware completions are enabled.
- Generic globs are available for all projects.
- Suggestions are not noisy in non-array contexts.

---

## 56. ✅ Add semantic tokens if TextMate becomes limiting

### Description
Evaluate whether semantic tokenization would improve highlighting beyond TextMate grammar.

Potential improvements:

- distinguish task names from hook names
- distinguish property names from arbitrary identifiers
- distinguish references from declarations

### Acceptance criteria
- Decision is documented.
- If implemented, semantic tokens enhance existing grammar without breaking themes.
- Tokenization handles malformed files gracefully.
- Tests or manual test cases cover representative syntax.

---

## 57. ✅ Add snippets/config generation tests on the Go side

### Description
If a Go schema/generator is introduced, test that generated extension artifacts are complete and stable.

### Acceptance criteria
- Go tests verify all parser-supported properties exist in the schema.
- Go tests verify generated snippets contain valid syntax patterns.
- Go tests verify generated docs include all schema entries.
- Tests fail when parser registries and schema drift.

---

## 58. Refactor parser registries to avoid reflection where possible

### Description
The Go parser uses reflection in `PropertyParser.setFieldValue`. This is flexible but weakens compile-time safety.

Consider replacing reflection-based assignment with typed setter functions in the schema/registry.

### Acceptance criteria
- Property parsing remains behavior-compatible.
- Invalid field names are impossible or caught at compile time.
- Tests cover all task/hook properties.
- Parser remains easy to extend.

---

## 59. ✅ Fix or review `pace init` unknown project prompt logic

### Description
In `internal/command/init.go`, when project type is unknown and the user is prompted to specify a type, the code appears to check `answer != "y"` after parsing the project type. This may cancel unexpectedly unless the user typed `y`.

This is not directly part of the VS Code extension, but it affects `Pace: Generate config.pace` if the extension calls `pace init`.

### Acceptance criteria
- Unknown project prompt flow is reviewed.
- User can type a supported project type and continue.
- User can intentionally cancel.
- Tests or manual verification cover unknown project initialization.

---

## 60. ✅ Add non-interactive Pace init mode for editor integrations

### Description
Editor integrations should not rely on interactive CLI prompts.

Add flags such as:

```bash
pace init --type go --yes
pace init --type go --force
pace init --type go --stdout
```

### Acceptance criteria
- Extension can generate config non-interactively.
- Existing interactive behavior remains available for terminal users.
- CLI refuses to overwrite existing config unless `--force` is used.
- `--stdout` or equivalent can print generated config without writing.

---

## 61. ✅ Add JSON schema export command

### Description
Expose the Pace language schema as JSON for editor integrations and docs tooling.

Potential command:

```bash
pace schema --json
```

### Acceptance criteria
- Command outputs top-level statements.
- Command outputs task properties.
- Command outputs hook properties.
- Command outputs args properties.
- Output includes property types, aliases, docs, and snippets.
- Output is versioned for compatibility.

---

## 62. ✅ Consume JSON schema in the VS Code extension

### Description
Instead of hardcoding all language metadata in TypeScript, consume generated or bundled JSON schema from the Go CLI/schema generator.

### Acceptance criteria
- Extension completion data is derived from schema JSON.
- Hover docs are derived from schema JSON.
- Diagnostics can use schema JSON for known properties.
- TypeScript code has minimal duplicated language metadata.

---

## 63. ✅ Add bundled fallback schema

### Description
If runtime CLI schema loading is unavailable, the extension should use a bundled schema generated at package time.

### Acceptance criteria
- Extension works without Pace CLI installed.
- Bundled schema is generated from Go source before packaging.
- Users with newer Pace CLI can optionally use runtime schema if enabled.
- Version mismatch behavior is documented.

---

### Batch 7 progress
Completed in this batch: introduced a Go-owned Pace language schema; parser registries now derive statement/property metadata from that schema; added `pace schema --json` and `pace schema --generate-vscode`; generated bundled schema JSON, TypeScript schema/snippet data, grammar keyword data, and a generated language reference; VS Code completions and hover docs now consume generated schema-backed data; package generation refreshes bundled schema artifacts; CI verifies generated artifacts are fresh; Go tests cover parser/schema drift and generated docs/keyword completeness.

## 64. ✅ Add version compatibility handling

### Description
The installed extension and installed Pace CLI may support different language versions.

### Acceptance criteria
- Extension can show bundled schema version.
- Extension can detect Pace CLI version when available.
- Major mismatches produce a helpful warning.
- Users can disable version warnings.

---

## 65. ✅ Add marketplace readiness metadata

### Description
Review `package.json` for marketplace quality.

Potential additions:

- keywords
- bugs URL
- homepage
- gallery banner
- better category list

### Acceptance criteria
- `package.json` includes useful keywords.
- Repository metadata is correct.
- Icon path is valid.
- License metadata is correct.
- Extension can be packaged with no marketplace warnings, or warnings are documented.

---

## 66. ✅ Add changelog

### Description
Add a `CHANGELOG.md` for extension releases.

### Acceptance criteria
- Changelog exists in `vscode-pace/CHANGELOG.md`.
- Current version has an entry.
- Release process updates changelog.
- Marketplace package includes changelog if desired.

---

## 67. ✅ Add contributing/development docs

### Description
Add a short development guide for working on the extension.

### Acceptance criteria
- Docs explain install dependencies.
- Docs explain compile/watch.
- Docs explain launching extension development host.
- Docs explain packaging and local install.
- Docs explain generated files, if introduced.

---

## 68. ✅ Add examples directory for `.pace` files

### Description
Add example Pace configs that can be used for manual extension testing.

Suggested examples:

- minimal config
- Go project config
- config with hooks
- config with args
- config with env
- invalid config for diagnostics testing

### Acceptance criteria
- Examples exist under `vscode-pace/examples` or shared repo docs.
- Examples parse successfully except intentionally invalid ones.
- README references examples.
- Manual testing checklist uses the examples.

---

## 69. ✅ Add manual QA checklist

### Description
Create a checklist for validating extension behavior before release.

### Acceptance criteria
- Checklist covers syntax highlighting.
- Checklist covers completions.
- Checklist covers snippets.
- Checklist covers commands.
- Checklist covers diagnostics.
- Checklist covers packaging and local install.
- Checklist includes Windows/macOS/Linux notes.

---

## 70. ✅ Revisit language file extension and filename associations

### Description
Currently the extension contributes `.pace`. The main repo uses `config.pace`. Consider whether to add filename-specific associations or icons.

### Acceptance criteria
- `.pace` files are recognized.
- `config.pace` is recognized.
- Any additional associations are intentional.
- No unrelated files are claimed by the extension.

---

## 71. ✅ Add file icon/theme contribution if desired

### Description
Consider adding a Pace file icon for `.pace` files using the existing Pace logo assets.

### Acceptance criteria
- Decision is documented.
- If implemented, file icon displays for `.pace` files.
- Icon works in light and dark themes.
- Packaging includes required assets.

---

## 72. Improve multiline string support in grammar/context detection

### Description
The grammar supports triple-quoted strings, but context detection and future formatter/diagnostics must also handle them correctly.

### Acceptance criteria
- Braces inside multiline strings do not affect context detection.
- Comments inside multiline strings are not treated as comments.
- Formatter preserves multiline strings.
- Tests cover multiline strings.

---

## 73. Support comments in scanners and future parsers

### Description
Any TypeScript scanner used by the extension should understand Pace comments.

### Acceptance criteria
- `#` starts a comment outside strings.
- Comments are ignored by context detection.
- Comments are ignored by symbol extraction.
- Comments are preserved by formatter.
- Tests cover comments with braces and keywords.

---

## 74. ✅ Add safe terminal command construction

### Description
Commands like `pace run <task>` should avoid shell injection or quoting bugs.

### Acceptance criteria
- Task names are validated before command construction.
- Terminal commands quote paths safely.
- Workspace paths with spaces work.
- Invalid task names from malformed files are not executed without confirmation.

---

## 75. ✅ Add extension telemetry decision

### Description
Decide whether the extension will collect telemetry. Prefer no telemetry unless explicitly needed.

### Acceptance criteria
- Decision is documented.
- If no telemetry, README states that no telemetry is collected.
- If telemetry is added, it follows VS Code marketplace policies and has settings to disable it.

---


### Batch 8 progress
Completed in this batch: added README language-reference/examples coverage and telemetry/schema compatibility notes; improved completion metadata with details, snippet examples, deprecated tagging, stable ordering, and replacement ranges; documented the semantic-token and file-icon decisions; added schema version compatibility checks and the `Pace: Show Schema Version` command; improved marketplace metadata and filename association; added changelog, contributing guide, examples, and a manual QA checklist; confirmed existing non-interactive init and safe terminal command work are complete.

## Suggested implementation order

1. Fix correctness drift in grammar/snippets/README.
2. Improve completion context detection.
3. Add task/hook extraction, reference completions, and document symbols.
4. Add dev scripts, tests, CI, and packaging improvements.
5. Add CLI/editor integration commands such as `Pace: Generate config.pace` and `Pace: Run Task`.
6. Introduce Go-owned schema and generated VS Code artifacts.
7. Add real diagnostics via `pace validate --json`.
8. Add advanced editor features: hover docs, CodeLens, formatter, semantic tokens.
