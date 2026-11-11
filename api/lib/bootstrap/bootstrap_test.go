package bootstrap

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lejianwen/rustdesk-api/v2/model"
	"github.com/lejianwen/rustdesk-api/v2/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestSecretPersistsPrivately(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".bootstrap", "jwt.key")
	first, err := EnsureSecret(path)
	if err != nil {
		t.Fatal(err)
	}
	second, err := EnsureSecret(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(first) != 64 || first != second {
		t.Fatal("generated secret must remain unchanged")
	}
	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm() != 0600 {
		t.Fatal("secret file permissions must be 0600")
	}
	info, err = os.Stat(filepath.Dir(path))
	if err != nil || info.Mode().Perm() != 0700 {
		t.Fatal("secret directory permissions must be 0700")
	}
}

func TestRejectUnsafeSecretFiles(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target")
	if err := os.WriteFile(target, []byte("fixture-only-not-a-production-key-1234567890"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := EnsureSecret(target); err == nil {
		t.Fatal("publicly readable secret must be rejected")
	}
	if err := os.Chmod(target, 0600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "link")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	if _, err := EnsureSecret(link); err == nil {
		t.Fatal("symlink secret must be rejected")
	}
	dirLink := filepath.Join(dir, "directory-link")
	if err := os.Symlink(dir, dirLink); err != nil {
		t.Fatal(err)
	}
	if _, err := EnsureSecret(filepath.Join(dirLink, "outside", "jwt.key")); err == nil {
		t.Fatal("symlink directory must be rejected before creating files")
	}
	if _, err := os.Stat(filepath.Join(dir, "outside")); !os.IsNotExist(err) {
		t.Fatal("symlink directory must not cause writes outside the secret path")
	}
}

func TestAdminBootstrapDoesNotResetExistingDatabase(t *testing.T) {
	dir := t.TempDir()
	db, err := gorm.Open(sqlite.Open(filepath.Join(dir, "api.db")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.Group{}); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, ".bootstrap", "admin.password")
	if err := InitializeAdmin(db, "Default", "Share", path); err != nil {
		t.Fatal(err)
	}
	var admin model.User
	if err := db.First(&admin).Error; err != nil {
		t.Fatal(err)
	}
	password, err := ReadSecret(path)
	if err != nil {
		t.Fatal(err)
	}
	ok, _, err := utils.VerifyPassword(admin.Password, password)
	if err != nil || !ok || len(password) != 32 || admin.IsAdmin == nil || !*admin.IsAdmin {
		t.Fatal("initial administrator password must authenticate")
	}
	originalHash := admin.Password
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := InitializeAdmin(db, "Default", "Share", path); err != nil {
		t.Fatal(err)
	}
	if err := db.First(&admin).Error; err != nil || admin.Password != originalHash {
		t.Fatal("existing administrator must remain unchanged")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("existing database must not generate a new bootstrap password")
	}
}
