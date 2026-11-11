package global

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitI18nIgnoresShortFileNames(t *testing.T) {
	oldPath, oldLocalizer := Config.Gin.ResourcesPath, Localizer
	t.Cleanup(func() {
		Config.Gin.ResourcesPath, Localizer = oldPath, oldLocalizer
	})
	root := t.TempDir()
	directory := filepath.Join(root, "i18n")
	if err := os.Mkdir(directory, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "x"), []byte("ignored"), 0600); err != nil {
		t.Fatal(err)
	}
	Config.Gin.ResourcesPath = root
	InitI18n()
}
