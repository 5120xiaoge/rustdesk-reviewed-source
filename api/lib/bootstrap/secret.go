package bootstrap

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	JWTFile           = "./data/.bootstrap/jwt.key"
	AdminPasswordFile = "./data/.bootstrap/admin.password"
)

func ReadSecret(path string) (string, error) {
	if err := checkParentDirectories(path); err != nil {
		return "", err
	}
	before, err := os.Lstat(path)
	if err != nil {
		return "", err
	}
	if !before.Mode().IsRegular() || before.Mode().Perm()&0077 != 0 {
		return "", fmt.Errorf("secret file must be a private regular file: %s", path)
	}
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	after, err := file.Stat()
	if err != nil || !os.SameFile(before, after) {
		return "", fmt.Errorf("secret file changed while opening: %s", path)
	}
	data, err := io.ReadAll(io.LimitReader(file, 4097))
	if err != nil {
		return "", err
	}
	value := strings.TrimSpace(string(data))
	if len(data) > 4096 || len(value) < 32 {
		return "", fmt.Errorf("secret file has an invalid length: %s", path)
	}
	return value, nil
}

func checkParentDirectories(path string) error {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	for dir := filepath.Dir(absolute); ; dir = filepath.Dir(dir) {
		info, err := os.Lstat(dir)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
		if err == nil && !info.IsDir() {
			return fmt.Errorf("secret path contains a non-directory or symlink: %s", dir)
		}
		if filepath.Dir(dir) == dir {
			return nil
		}
	}
}

func EnsureSecret(path string) (string, error) {
	return ensureSecret(path, 32)
}

func ensureSecret(path string, randomBytes int) (string, error) {
	value, err := ReadSecret(path)
	if err == nil || !errors.Is(err, os.ErrNotExist) {
		return value, err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	info, err := os.Lstat(dir)
	if err != nil {
		return "", err
	}
	if !info.IsDir() || info.Mode().Perm()&0077 != 0 {
		return "", fmt.Errorf("secret directory must be private and not a symlink: %s", dir)
	}
	random := make([]byte, randomBytes)
	if _, err := rand.Read(random); err != nil {
		return "", err
	}
	value = hex.EncodeToString(random)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if errors.Is(err, os.ErrExist) {
		return ReadSecret(path)
	}
	if err != nil {
		return "", err
	}
	if _, err := io.WriteString(file, value+"\n"); err != nil {
		file.Close()
		return "", err
	}
	if err := file.Sync(); err != nil {
		file.Close()
		return "", err
	}
	if err := file.Close(); err != nil {
		return "", err
	}
	return value, nil
}
