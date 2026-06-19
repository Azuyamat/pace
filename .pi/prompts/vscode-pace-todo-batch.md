Work through `@vscode-pace/TODO.md`.

Batch: {{batch}}

## Goal
Implement the requested TODO batch carefully and completely while keeping changes focused and reviewable.

## Rules
- Read `vscode-pace/TODO.md` first.
- Implement only the requested batch unless a dependency requires a small related change.
- Prefer source-of-truth correctness over cosmetic changes.
- Keep changes minimal, focused, and easy to review.
- Do not skip acceptance criteria silently.
- If an acceptance criterion cannot be completed, explain why and leave a clear note in the summary.
- After implementing, run the relevant build/test commands.
- At minimum, for extension changes run:
  - `cd vscode-pace && npm run compile`
- If Go code changes are made, also run relevant Go checks, usually:
  - `go test ./...`
- Update `vscode-pace/TODO.md` to mark completed tasks or add progress notes.
- Do not mark a task complete unless its acceptance criteria are satisfied.

## Suggested batch order
Use this mapping unless I specify a different batch:

- Batch 1: tasks 1-10 — grammar/snippet/README correctness fixes.
- Batch 2: tasks 11-15 — context detection and reference completions.
- Batch 3: tasks 16-19 — document symbols, CodeLens, hover docs.
- Batch 4: tasks 33-40 — tests, CI, packaging, README cleanup.
- Batch 5: tasks 21-23 and 50-52 — diagnostics, settings, validation workflow.
- Batch 6: tasks 24-26 and 44-49 — VS Code commands and project-aware workflows.
- Batch 7: tasks 27-32, 57, and 61-63 — Go-owned schema and generated VS Code artifacts.
- Batch 8: remaining polish tasks.

## Required final response
At the end, summarize:

1. Completed TODO items
2. Partially completed TODO items, if any
3. Files changed
4. Commands run and their results
5. Remaining recommended next batch
