"""Exercise the real deployment snippets through a disposable local Nginx."""
import json
import pathlib
import socket
import subprocess
import tempfile
import threading
import time
import urllib.error
import urllib.request
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer

ROOT = pathlib.Path(__file__).resolve().parents[1]

class Upstream(BaseHTTPRequestHandler):
    def do_GET(self):
        body = json.dumps({
            "authorization": self.headers.get("Authorization", ""),
            "proof": self.headers.get("X-GH-Identity-Proof", ""),
            "path": self.path,
        }).encode()
        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    do_POST = do_GET

    def log_message(self, *args):
        pass

def free_port():
    with socket.socket() as sock:
        sock.bind(("127.0.0.1", 0))
        return sock.getsockname()[1]

def main():
    servers = [ThreadingHTTPServer(("127.0.0.1", port), Upstream) for port in (8081, 8082)]
    for server in servers:
        threading.Thread(target=server.serve_forever, daemon=True).start()
    try:
        with tempfile.TemporaryDirectory() as tmp:
            folder = pathlib.Path(tmp)
            school, internal = free_port(), free_port()
            assert school != internal
            for name in ("body", "proxy", "fastcgi", "uwsgi", "scgi"):
                (folder / name).mkdir()
            config = folder / "nginx.conf"
            config.write_text(f"""
daemon off;
master_process off;
pid {folder}/nginx.pid;
error_log stderr warn;
events {{}}
http {{
  access_log off;
  client_body_temp_path {folder}/body;
  proxy_temp_path {folder}/proxy;
  fastcgi_temp_path {folder}/fastcgi;
  uwsgi_temp_path {folder}/uwsgi;
  scgi_temp_path {folder}/scgi;
  server {{ listen 127.0.0.1:{school}; include {ROOT}/nginx/school.conf.example; }}
  server {{ listen 127.0.0.1:{internal}; include {ROOT}/nginx/internal-tool.conf.example; }}
}}
""")
            proc = subprocess.Popen(["nginx", "-p", tmp, "-c", str(config)])
            try:
                for _ in range(50):
                    if proc.poll() is not None:
                        raise RuntimeError("Nginx exited before probe")
                    try:
                        with socket.create_connection(("127.0.0.1", school), timeout=.2):
                            break
                    except OSError:
                        time.sleep(.1)
                def call(port, path, method="GET", proof="Recovery example"):
                    request = urllib.request.Request(
                        f"http://127.0.0.1:{port}{path}", method=method,
                        headers={"Authorization": proof, "X-GH-Identity-Proof": "untrusted"},
                    )
                    try:
                        with urllib.request.urlopen(request, timeout=3) as response:
                            return response.status, json.load(response)
                    except urllib.error.HTTPError as exc:
                        return exc.code, None
                for kind, scheme in (("recovery", "Recovery"), ("invitations", "Invitation")):
                    for operation, method in (("validate", "GET"), ("exchange", "POST")):
                        path = f"/identity/public/auth/{kind}/{operation}"
                        status, data = call(school, path, method, scheme + " test")
                        assert status == 200 and data["authorization"] == scheme + " test" and data["proof"] == "", (path, status, data)
                for path in ("/identity/public/auth/login", "/identity/protected/auth/me", "/integrations/protected/connections"):
                    status, data = call(school, path)
                    assert status == 200 and data["authorization"] == data["proof"] == "", (path, status, data)
                status, data = call(internal, "/internal/protected/requests")
                assert status == 200 and data["authorization"] == data["proof"] == ""
                for path in ("/identity/restricted/auth/login", "/school/restricted/registry/lookup", "/internal/protected/requests", "/health/ready"):
                    assert call(school, path)[0] == 404, path
                for path in ("/identity/public/auth/login", "/internal/restricted/example", "/health/ready"):
                    assert call(internal, path)[0] == 404, path
            finally:
                proc.terminate()
                proc.wait(timeout=5)
    finally:
        for server in servers:
            server.shutdown()
            server.server_close()

if __name__ == "__main__":
    main()
