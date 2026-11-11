import base64
import hashlib
import hmac
import json
import os
import secrets
import subprocess
import tempfile
import time
import urllib.error
import urllib.request


IMAGE = os.environ["IMAGE"] + ":" + os.environ["IMAGE_TAG"]
NAME = "reviewed-fresh-install"
API = "http://127.0.0.1:21114"
DATA = tempfile.mkdtemp(prefix="rustdesk-fresh-", dir=os.environ["RUNNER_TEMP"])
JWT = "/app/data/.bootstrap/jwt.key"
PASSWORD = "/app/data/.bootstrap/admin.password"


def docker(*args, input_text=None):
    result = subprocess.run(
        ["docker", *args], input=input_text, text=True, capture_output=True
    )
    if result.returncode:
        raise RuntimeError("Docker operation failed: " + args[0])
    return result.stdout.strip()


def read_private(path):
    return docker("exec", NAME, "cat", path)


def start():
    docker(
        "run", "-d", "--name", NAME,
        "-p", "127.0.0.1:21114:21114",
        "-v", DATA + "/server:/data",
        "-v", DATA + "/api:/app/data",
        "-e", "MUST_LOGIN=Y",
        "-e", "RELAY=127.0.0.1:21117",
        "-e", "RUSTDESK_API_RUSTDESK_KEY_FILE=/data/id_ed25519.pub",
        IMAGE,
    )
    wait_healthy()


def wait_healthy():
    for _ in range(90):
        if docker("inspect", "-f", "{{.State.Health.Status}}", NAME) == "healthy":
            return
        time.sleep(2)
    raise RuntimeError("Fresh installation did not become healthy")


def login(password):
    request = urllib.request.Request(
        API + "/api/login",
        data=json.dumps({"username": "admin", "password": password}).encode(),
        headers={"Content-Type": "application/json"},
    )
    try:
        with urllib.request.urlopen(request, timeout=10) as response:
            return json.load(response)
    except urllib.error.HTTPError:
        return {}


def token_from(response):
    token = response.get("access_token", "")
    if not token and isinstance(response.get("data"), dict):
        token = response["data"].get("access_token", "")
    if not token:
        raise RuntimeError("Initial administrator could not authenticate")
    return token


def check_token(token, key):
    parts = token.split(".")
    if len(parts) != 3:
        raise RuntimeError("Mandatory-login API did not issue a JWT")
    signature = base64.urlsafe_b64decode(parts[2] + "=" * (-len(parts[2]) % 4))
    expected = hmac.new(key.encode(), ".".join(parts[:2]).encode(), hashlib.sha256).digest()
    if not hmac.compare_digest(signature, expected):
        raise RuntimeError("API did not use the persistent shared JWT key")


def check_no_secret_logs(values):
    logs = docker("logs", NAME)
    api_logs = docker("exec", NAME, "sh", "-c", "cat /app/runtime/log.txt 2>/dev/null || true")
    if any(value and (value in logs or value in api_logs) for value in values):
        raise RuntimeError("A credential was written to application logs")


try:
    start()
    key = read_private(JWT)
    password = read_private(PASSWORD)
    public_key = read_private("/data/id_ed25519.pub")
    private_key = read_private("/data/id_ed25519")
    if len(key) != 64 or len(password) != 32:
        raise RuntimeError("Bootstrap credentials have invalid lengths")
    docker(
        "exec", NAME, "sh", "-c",
        "test \"$(stat -c %a /app/data/.bootstrap)\" = 700 && "
        "test \"$(stat -c %a /app/data/.bootstrap/jwt.key)\" = 600 && "
        "test \"$(stat -c %a /app/data/.bootstrap/admin.password)\" = 600",
    )
    check_token(token_from(login(password)), key)
    check_no_secret_logs([key, password, private_key])
    docker("restart", NAME)
    wait_healthy()
    if read_private(JWT) != key or read_private("/data/id_ed25519.pub") != public_key:
        raise RuntimeError("Restart changed persistent server credentials")
    check_token(token_from(login(password)), key)
    docker("rm", "-f", NAME)
    start()
    if read_private(JWT) != key or read_private(PASSWORD) != password:
        raise RuntimeError("Container recreation changed bootstrap credentials")
    if read_private("/data/id_ed25519.pub") != public_key:
        raise RuntimeError("Container recreation changed the server public key")
    check_token(token_from(login(password)), key)
    replacement = secrets.token_hex(16)
    docker(
        "exec", "-i", NAME, "./apimain", "reset-admin-pwd", "--stdin",
        input_text=replacement + "\n",
    )
    check_token(token_from(login(replacement)), key)
    if login(password).get("access_token"):
        raise RuntimeError("Previous administrator password still authenticates")
    docker("restart", NAME)
    wait_healthy()
    check_token(token_from(login(replacement)), key)
    check_no_secret_logs([key, password, replacement, private_key])
    print("PASS: first install, admin login, JWT, restart, recreation, stdin reset, private files and logs")
finally:
    subprocess.run(["docker", "rm", "-f", NAME], capture_output=True, check=False)
