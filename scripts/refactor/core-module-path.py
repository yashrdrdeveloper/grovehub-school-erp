#!/usr/bin/env python3
"""Migrate tracked Core module references; no network or Git writes."""

import argparse
import pathlib
import subprocess

OLD = "github.com/yashrdrdeveloper/grovehub-school-erp/" + "libs/" + "backend/core"
NEW = "grovehub/libs/core"


def tracked_paths():
    raw = subprocess.check_output(["git", "ls-files", "-z"])
    return [pathlib.Path(p.decode()) for p in raw.split(b"\x00") if p]


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--apply", action="store_true")
    args = parser.parse_args()
    root = subprocess.check_output(["git", "rev-parse", "--show-toplevel"], text=True).strip()
    import os
    os.chdir(root)
    changed = []
    for path in tracked_paths():
        if path.is_symlink() or not path.is_file() or path.as_posix() == "scripts/refactor/core-module-path.py":
            continue
        data = path.read_bytes()
        if b"\x00" in data:
            continue
        try:
            source = data.decode("utf-8")
        except UnicodeDecodeError:
            continue
        # Also handle escaped regexp literals and legacy package fragments in guards.
        updated = source.replace(OLD, NEW)
        updated = updated.replace(OLD.replace(".", "\\."), NEW)
        updated = updated.replace(OLD.replace("/", "\\/").replace(".", "\\."), NEW.replace("/", "\\/"))
        updated = updated.replace("libs/" + "backend/core", "libs/core")
        if source == updated:
            continue
        changed.append(path.as_posix())
        if args.apply:
            path.write_bytes(updated.encode("utf-8"))
    print(f"{len(changed)} tracked text files {'updated' if args.apply else 'require migration'}")
    for path in changed:
        print(path)


if __name__ == "__main__":
    main()
