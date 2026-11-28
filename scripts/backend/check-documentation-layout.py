#!/usr/bin/env python3
"""Check tracked backend documentation and Platform service entrypoints."""

from pathlib import PurePosixPath
import subprocess
import sys


def main():
    tracked = subprocess.check_output(["git", "ls-files", "-z", "backend"]).decode().split("\0")
    invalid = []
    for name in filter(None, tracked):
        path = PurePosixPath(name)
        parts = path.parts
        if path.suffix.lower() in {".md", ".mdx", ".rst", ".adoc"} or "docs" in parts:
            if not any(parts[i:i + 3] == ("internal", "resources", "docs") for i in range(len(parts) - 2)):
                invalid.append(name)
        if len(parts) > 5 and parts[:2] == ("backend", "projects") and parts[3] == "services":
            if parts[5:6] == ("tests",) or parts[5:7] == ("internal", "tests"):
                invalid.append(f"{name}: service tests must live under internal/resources/tests")
    for owner in ("audit", "foundation", "identity", "infrastructure", "integration"):
        prefix = PurePosixPath("backend/projects/platform/services") / owner
        root_go = sorted(
            PurePosixPath(name).name for name in filter(None, tracked)
            if PurePosixPath(name).parent == prefix and name.endswith(".go")
        )
        if root_go != ["main.go"]:
            invalid.append(f"{prefix}: expected main.go only, found {root_go}")
    if invalid:
        print("Backend architecture violations:", *invalid, sep="\n", file=sys.stderr)
        return 1
    print("Backend documentation, service tests and service entrypoints follow the required layout")
    return 0


if __name__ == "__main__":
    sys.exit(main())
