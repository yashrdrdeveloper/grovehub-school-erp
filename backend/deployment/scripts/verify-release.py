#!/usr/bin/env python3
"""Verify and extract a CI artifact into a new, private staging directory."""
import argparse
import hashlib
import json
import os
from pathlib import Path, PurePosixPath
import re
import shutil
import tarfile

PROJECTS = {"platform", "school-backend", "internal-tool-backend",
            "school-frontend", "internal-tool-frontend"}
MAX_BYTES = 1024 * 1024 * 1024
MAX_FILES = 20000


def digest(path):
    with path.open("rb") as stream:
        return hashlib.file_digest(stream, "sha256").hexdigest()


def unpack(archive, destination, project, commit, expected_sha256):
    if project not in PROJECTS or not re.fullmatch(r"[0-9a-f]{40}", commit):
        raise ValueError("invalid release identity")
    if not re.fullmatch(r"[0-9a-f]{64}", expected_sha256):
        raise ValueError("invalid trusted archive checksum")
    if digest(archive) != expected_sha256:
        raise ValueError("archive checksum mismatch")
    destination.mkdir(mode=0o700, parents=False, exist_ok=False)
    try:
        seen = set()
        size = 0
        with tarfile.open(archive, "r:gz") as source:
            for item in source:
                name = item.name.rstrip("/")
                parts = PurePosixPath(name).parts
                if (not parts or name.startswith("/") or ".." in parts or
                        name != "/".join(parts) or name in seen or
                        len(seen) >= MAX_FILES or
                        any(not re.fullmatch(r"[A-Za-z0-9_.@+ -]+", p) for p in parts)):
                    raise ValueError("unsafe archive path")
                if not (name in {"manifest.json", "SHA256SUMS", "payload"} or
                        name.startswith("payload/")):
                    raise ValueError("unexpected archive entry")
                if any(p.startswith(".env") or p in {".git", ".ssh"} or
                       p.endswith((".pem", ".key")) for p in parts):
                    raise ValueError("secret-like archive entry")
                seen.add(name)
                target = destination.joinpath(*parts)
                if item.isdir():
                    target.mkdir(mode=0o755, parents=True, exist_ok=True)
                elif item.isfile():
                    size += item.size
                    if size > MAX_BYTES:
                        raise ValueError("archive size limit exceeded")
                    target.parent.mkdir(mode=0o755, parents=True, exist_ok=True)
                    with source.extractfile(item) as reader, target.open("xb") as writer:
                        shutil.copyfileobj(reader, writer)
                    target.chmod(0o644)
                else:
                    raise ValueError("links and special files are forbidden")
        manifest = json.loads((destination / "manifest.json").read_text())
        if manifest != {"project": project, "commit": commit, "platform": "linux-amd64"}:
            raise ValueError("manifest mismatch")
        checks = {}
        for line in (destination / "SHA256SUMS").read_text().splitlines():
            match = re.fullmatch(r"([0-9a-f]{64})  (.+)", line)
            if not match:
                raise ValueError("invalid checksum manifest")
            checksum, name = match.groups()
            if name in checks:
                raise ValueError("duplicate checksum")
            checks[name] = checksum
        files = {p.relative_to(destination).as_posix() for p in destination.rglob("*")
                 if p.is_file() and p.relative_to(destination).as_posix() != "SHA256SUMS"}
        if set(checks) != files:
            raise ValueError("checksum coverage mismatch")
        for name, checksum in checks.items():
            if digest(destination / name) != checksum:
                raise ValueError("payload checksum mismatch")
        if project.endswith("-frontend"):
            if "payload/index.html" not in files:
                raise ValueError("missing frontend index")
        else:
            if files != {"manifest.json", f"payload/{project}"}:
                raise ValueError("unexpected backend payload")
            (destination / "payload" / project).chmod(0o755)
        destination.chmod(0o755)
    except BaseException:
        shutil.rmtree(destination)
        raise


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("archive", type=Path)
    parser.add_argument("destination", type=Path)
    parser.add_argument("project", choices=sorted(PROJECTS))
    parser.add_argument("commit")
    parser.add_argument("sha256")
    args = parser.parse_args()
    unpack(args.archive, args.destination, args.project, args.commit, args.sha256)
