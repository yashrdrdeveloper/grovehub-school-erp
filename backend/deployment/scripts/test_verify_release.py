import hashlib
import importlib.util
import io
import json
from pathlib import Path
import tarfile
import tempfile
import unittest

spec = importlib.util.spec_from_file_location("verify_release", Path(__file__).with_name("verify-release.py"))
verify = importlib.util.module_from_spec(spec)
spec.loader.exec_module(verify)
COMMIT = "a" * 40


class VerifyReleaseTest(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.addCleanup(self.tmp.cleanup)
        self.root = Path(self.tmp.name)
        self.archive = self.root / "release.tar.gz"
        self.dest = self.root / "release"

    def make_archive(self, extra=None, bad_checksum=False, project="school-backend"):
        files = {
            "manifest.json": json.dumps({"project": project, "commit": COMMIT, "platform": "linux-amd64"}).encode(),
            f"payload/{project}": b"test executable",
        }
        checks = "".join(f"{hashlib.sha256(data).hexdigest()}  {name}\n" for name, data in files.items())
        if bad_checksum:
            checks = checks.replace(hashlib.sha256(b"test executable").hexdigest(), "0" * 64)
        files["SHA256SUMS"] = checks.encode()
        with tarfile.open(self.archive, "w:gz") as out:
            for name, data in files.items():
                item = tarfile.TarInfo(name)
                item.size = len(data)
                out.addfile(item, io.BytesIO(data))
            if extra:
                item, data = extra
                out.addfile(item, io.BytesIO(data) if data is not None else None)

    def unpack(self, checksum=None):
        verify.unpack(self.archive, self.dest, "school-backend", COMMIT, checksum or verify.digest(self.archive))

    def test_valid_archive(self):
        self.make_archive()
        self.unpack()
        self.assertEqual((self.dest / "payload/school-backend").read_bytes(), b"test executable")
        self.assertEqual((self.dest / "payload/school-backend").stat().st_mode & 0o777, 0o755)

    def test_outer_checksum(self):
        self.make_archive()
        with self.assertRaises(ValueError):
            self.unpack("0" * 64)
        self.assertFalse(self.dest.exists())

    def test_inner_checksum(self):
        self.make_archive(bad_checksum=True)
        with self.assertRaises(ValueError):
            self.unpack()
        self.assertFalse(self.dest.exists())

    def test_wrong_project(self):
        self.make_archive(project="platform")
        with self.assertRaises(ValueError):
            self.unpack()

    def test_traversal(self):
        item = tarfile.TarInfo("payload/../../escape")
        item.size = 1
        self.make_archive((item, b"x"))
        with self.assertRaises(ValueError):
            self.unpack()
        self.assertFalse((self.root / "escape").exists())

    def test_symlink(self):
        item = tarfile.TarInfo("payload/link")
        item.type = tarfile.SYMTYPE
        item.linkname = "/etc/passwd"
        self.make_archive((item, None))
        with self.assertRaises(ValueError):
            self.unpack()

    def test_unlisted_file(self):
        item = tarfile.TarInfo("payload/unlisted")
        item.size = 1
        self.make_archive((item, b"x"))
        with self.assertRaises(ValueError):
            self.unpack()

    def test_duplicate_entry(self):
        item = tarfile.TarInfo("payload/school-backend")
        item.size = 1
        self.make_archive((item, b"x"))
        with self.assertRaises(ValueError):
            self.unpack()


if __name__ == "__main__":
    unittest.main()
