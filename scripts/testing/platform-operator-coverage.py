#!/usr/bin/env python3
"""Verify the relocated Platform operator code, not the empty cmd package."""

import json
import os
from pathlib import Path
import shutil
import subprocess
import sys


def main():
    root = Path(__file__).resolve().parents[2]
    project = root / "backend/projects/platform"
    generated = project / ".generated/unit-coverage"
    generated.mkdir(parents=True, exist_ok=True)
    raw = generated / "operator-raw.out"
    profile = generated / "operator.out"
    overlay = json.loads((project / ".generated/overlay.json").read_text())
    owners = ("foundation", "infrastructure", "integration", "identity", "audit")
    production = [
        "./launcher/internal/src/operations/command/...",
        "./launcher/internal/src/strategies/migrate",
        "./launcher/internal/src/strategies/serve",
        "./services/audit/internal/app/strategies/command/...",
        "./services/identity/internal/app/strategies/command/...",
    ]
    providers = [
        path
        for owner in owners
        for path in (f"./services/{owner}", f"./services/{owner}/internal/app/bootstrap/impl")
    ]
    tests = [
        "./launcher/internal/resources/tests/unit-test",
        "./launcher/internal/resources/tests/unit-test/operations/command/...",
        "./launcher/internal/resources/tests/unit-test/strategies/migrate",
        "./launcher/internal/resources/tests/unit-test/strategies/serve",
        "./services/audit/internal/tests/unit-test/strategies/command/...",
        "./services/identity/internal/tests/unit-test/strategies/command/...",
    ]

    # Go coverage instrumentation must see the generated logical source paths.
    # Never overwrite maintained source; remove only the copies made here.
    created = []
    try:
        for logical, replacement in sorted(overlay["Replace"].items()):
            destination = Path(logical)
            destination.parent.mkdir(parents=True, exist_ok=True)
            with destination.open("xb") as output:
                created.append(destination)
                with Path(replacement).open("rb") as source:
                    shutil.copyfileobj(source, output)

        subprocess.run(
            [
                "go", "test", "-count=1", "-vet=off", "-covermode=atomic",
                "-coverpkg=" + ",".join(production + providers),
                "-coverprofile=" + str(raw),
                *production, *providers, *tests,
            ],
            cwd=project, check=True, env={**os.environ, "GOWORK": "off"},
        )

        facade_lines = {}
        for owner in owners:
            relative = f"services/{owner}/main.go"
            lines = (project / relative).read_text().splitlines()
            start = next(i for i, line in enumerate(lines, 1) if line.startswith("func Bootstrap("))
            end = next(i for i in range(start + 1, len(lines) + 1) if lines[i - 1] == "}")
            facade_lines["grovehub/" + relative] = (start, end)

        def in_scope(block):
            filename, location = block.rsplit(":", 1)
            if filename in facade_lines:
                start, end = facade_lines[filename]
                return start <= int(location.split(".", 1)[0]) <= end
            return (
                "/launcher/internal/src/operations/command/" in filename
                or "/launcher/internal/src/strategies/migrate/" in filename
                or "/launcher/internal/src/strategies/serve/" in filename
                or "/internal/app/strategies/command/" in filename
                or filename.endswith("/internal/app/bootstrap/impl/operator.go")
            )

        blocks = {}
        for line in raw.read_text().splitlines()[1:]:
            block, statements, hits = line.split()
            if in_scope(block):
                count, previous = blocks.get(block, (int(statements), 0))
                blocks[block] = (count, max(previous, int(hits)))

        with profile.open("w") as output:
            output.write("mode: atomic\n")
            for block, (statements, hits) in sorted(blocks.items()):
                output.write(f"{block} {statements} {hits}\n")

        total = sum(statements for statements, _ in blocks.values())
        covered = sum(statements for statements, hits in blocks.values() if hits)
        if not total:
            raise RuntimeError("no operator source statements instrumented")
        subprocess.run(["go", "tool", "cover", "-func=" + str(profile)], cwd=project, check=True)
        print(f"Operator coverage: {covered}/{total} statements ({100 * covered / total:.2f}%)")
        if covered != total:
            for block, (statements, hits) in sorted(blocks.items()):
                if not hits:
                    print(f"uncovered: {block} ({statements} statements)", file=sys.stderr)
            return 1
        return 0
    finally:
        for path in reversed(created):
            path.unlink(missing_ok=True)


if __name__ == "__main__":
    sys.exit(main())
