# Contributing to the Pace VS Code Extension

## Setup

```bash
cd vscode-pace
npm ci
```

## Development loop

- `npm run compile` compiles TypeScript to `build/`.
- `npm run watch` recompiles on change.
- `npm test` runs extension utility and completion tests.
- `npm run dev` compiles and launches VS Code with this folder as the extension development path. You can also press `F5` from VS Code.

## Generated files

Language metadata is owned by Go in `internal/config/schema`. Regenerate extension artifacts after schema changes:

```bash
cd vscode-pace
npm run generate:schema
```

This updates:

- `schemas/pace.schema.json`
- `src/generated/paceSchema.ts`
- `src/generated/schemaSnippets.ts`
- `syntaxes/pace.generated-keywords.json`
- `syntaxes/pace.tmLanguage.json` keyword patterns
- `LANGUAGE_REFERENCE.generated.md`

Do not hand-edit generated files except when changing the generator itself.

## Packaging and local install

```bash
npm run package
npm run install:local
```

`npm run install:local` uses the VS Code `code` command and shell glob expansion. If your shell does not expand `*.vsix`, install the generated VSIX through the Command Palette.

## Examples and QA

Manual testing examples live in `examples/`. Use `docs/MANUAL_QA.md` before publishing a release.

## Semantic tokens decision

The extension currently relies on the generated TextMate grammar plus context-aware providers rather than a semantic token provider. TextMate is sufficient for the supported highlighting surface today, and avoiding semantic tokens keeps malformed-file handling simple. Revisit this if themes need declaration/reference-specific coloring that TextMate cannot express.

## File icon decision

The extension currently contributes the Pace logo as the marketplace/icon asset only. It does not contribute a file icon theme because VS Code file icons are normally supplied by user-selected icon themes; adding a custom icon theme would be a separate, intentional feature.
