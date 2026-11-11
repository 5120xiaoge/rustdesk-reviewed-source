use sodiumoxide::crypto::sign;
use std::{fs, path::PathBuf, process::Command};

struct KeyFiles {
    directory: PathBuf,
    public: PathBuf,
    secret: PathBuf,
}

impl KeyFiles {
    fn new() -> Self {
        let directory = std::env::temp_dir().join(format!("rustdesk-keys-{}", uuid::Uuid::new_v4()));
        fs::create_dir(&directory).unwrap();
        Self {
            public: directory.join("id_ed25519.pub"),
            secret: directory.join("id_ed25519"),
            directory,
        }
    }

    fn validate(&self) -> std::process::Output {
        Command::new(env!("CARGO_BIN_EXE_rustdesk-utils"))
            .arg("validatekeypair-files")
            .arg(&self.public)
            .arg(&self.secret)
            .output()
            .unwrap()
    }
}

impl Drop for KeyFiles {
    fn drop(&mut self) {
        let _ = fs::remove_file(&self.public);
        let _ = fs::remove_file(&self.secret);
        let _ = fs::remove_dir(&self.directory);
    }
}

#[test]
fn validates_key_files_with_trailing_newlines() {
    sodiumoxide::init().unwrap();
    let files = KeyFiles::new();
    let (public, secret) = sign::gen_keypair();
    let secret_text = base64::encode(secret);
    fs::write(&files.public, format!("{}\n", base64::encode(public))).unwrap();
    fs::write(&files.secret, format!("{secret_text}\n")).unwrap();
    let output = files.validate();
    assert!(output.status.success());
    assert_eq!(
        String::from_utf8_lossy(&output.stdout).trim(),
        "Key pair is VALID"
    );
    assert!(!String::from_utf8_lossy(&output.stdout).contains(&secret_text));
    assert!(!String::from_utf8_lossy(&output.stderr).contains(&secret_text));
}

#[test]
fn rejects_mismatched_and_malformed_key_files() {
    sodiumoxide::init().unwrap();
    let files = KeyFiles::new();
    let (public, _) = sign::gen_keypair();
    let (_, other_secret) = sign::gen_keypair();
    fs::write(&files.public, base64::encode(public)).unwrap();
    fs::write(&files.secret, base64::encode(other_secret)).unwrap();
    assert!(!files.validate().status.success());
    fs::write(&files.secret, "not-a-key").unwrap();
    assert!(!files.validate().status.success());
}

#[test]
fn rejects_missing_key_files() {
    let files = KeyFiles::new();
    assert!(!files.validate().status.success());
}
