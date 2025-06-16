package middleware

import (
	"testing"
)

func TestStatic_FlagSet(t *testing.T) {
	config := &Static{
		Enabled:    true,
		Root:       "/var/www",
		Index:      "home.html",
		HTML5:      true,
		Browse:     true,
		IgnoreBase: true,
	}

	fs := config.FlagSet()

	if fs == nil {
		t.Fatal("FlagSet() returned nil")
	}

	enabledFlag := fs.Lookup(StaticEnabled)
	if enabledFlag == nil {
		t.Errorf("Flag %s not found", StaticEnabled)
	} else {
		if enabledFlag.DefValue != "true" {
			t.Errorf("Flag %s default value = %v, want true", StaticEnabled, enabledFlag.DefValue)
		}
	}

	rootFlag := fs.Lookup(StaticRoot)
	if rootFlag == nil {
		t.Errorf("Flag %s not found", StaticRoot)
	} else {
		if rootFlag.DefValue != "/var/www" {
			t.Errorf("Flag %s default value = %v, want /var/www", StaticRoot, rootFlag.DefValue)
		}
	}

	indexFlag := fs.Lookup(StaticIndex)
	if indexFlag == nil {
		t.Errorf("Flag %s not found", StaticIndex)
	} else {
		if indexFlag.DefValue != "home.html" {
			t.Errorf("Flag %s default value = %v, want home.html", StaticIndex, indexFlag.DefValue)
		}
	}
}

func TestStatic_FlagSet_Parse(t *testing.T) {
	config := &Static{
		Enabled:    false,
		Root:       "",
		Index:      "index.html",
		HTML5:      false,
		Browse:     false,
		IgnoreBase: false,
	}

	fs := config.FlagSet()

	args := []string{
		"--static-enabled",
		"--static-root", "/public",
		"--static-index", "main.html",
		"--static-html5",
		"--static-browse",
		"--static-ignore-base",
	}

	err := fs.Parse(args)
	if err != nil {
		t.Fatalf("Failed to parse flags: %v", err)
	}

	if !config.Enabled {
		t.Errorf("Enabled = %v, want true", config.Enabled)
	}
	if config.Root != "/public" {
		t.Errorf("Root = %v, want /public", config.Root)
	}
	if config.Index != "main.html" {
		t.Errorf("Index = %v, want main.html", config.Index)
	}
	if !config.HTML5 {
		t.Errorf("HTML5 = %v, want true", config.HTML5)
	}
	if !config.Browse {
		t.Errorf("Browse = %v, want true", config.Browse)
	}
	if !config.IgnoreBase {
		t.Errorf("IgnoreBase = %v, want true", config.IgnoreBase)
	}
}

func TestDefaultStatic(t *testing.T) {
	if DefaultStatic == nil {
		t.Fatal("DefaultStatic is nil")
	}

	if DefaultStatic.Enabled != false {
		t.Errorf("DefaultStatic.Enabled = %v, want false", DefaultStatic.Enabled)
	}
	if DefaultStatic.Root != "" {
		t.Errorf("DefaultStatic.Root = %v, want empty string", DefaultStatic.Root)
	}
	if DefaultStatic.Index != "index.html" {
		t.Errorf("DefaultStatic.Index = %v, want index.html", DefaultStatic.Index)
	}
	if DefaultStatic.HTML5 != false {
		t.Errorf("DefaultStatic.HTML5 = %v, want false", DefaultStatic.HTML5)
	}
	if DefaultStatic.Browse != false {
		t.Errorf("DefaultStatic.Browse = %v, want false", DefaultStatic.Browse)
	}
	if DefaultStatic.IgnoreBase != false {
		t.Errorf("DefaultStatic.IgnoreBase = %v, want false", DefaultStatic.IgnoreBase)
	}
}

func TestStatic_FlagSet_DefaultValues(t *testing.T) {
	config := &Static{
		Enabled:    true,
		Root:       "/assets",
		Index:      "start.html",
		HTML5:      false,
		Browse:     true,
		IgnoreBase: false,
	}

	fs := config.FlagSet()

	enabledFlag := fs.Lookup(StaticEnabled)
	if enabledFlag == nil {
		t.Fatal("Enabled flag not found")
	}
	if enabledFlag.DefValue != "true" {
		t.Errorf("Enabled flag default = %v, want true", enabledFlag.DefValue)
	}

	rootFlag := fs.Lookup(StaticRoot)
	if rootFlag == nil {
		t.Fatal("Root flag not found")
	}
	if rootFlag.DefValue != "/assets" {
		t.Errorf("Root flag default = %v, want /assets", rootFlag.DefValue)
	}

	indexFlag := fs.Lookup(StaticIndex)
	if indexFlag == nil {
		t.Fatal("Index flag not found")
	}
	if indexFlag.DefValue != "start.html" {
		t.Errorf("Index flag default = %v, want start.html", indexFlag.DefValue)
	}
}

func TestStatic_FlagSet_DisabledByDefault(t *testing.T) {
	config := DefaultStatic

	fs := config.FlagSet()

	var args []string

	err := fs.Parse(args)
	if err != nil {
		t.Fatalf("Failed to parse empty flags: %v", err)
	}

	if config.Enabled {
		t.Errorf("Enabled = %v, want false (default)", config.Enabled)
	}
}

func TestNewStatic(t *testing.T) {
	config := &Static{
		Enabled:    true,
		Root:       "/public",
		Index:      "index.html",
		HTML5:      true,
		Browse:     false,
		IgnoreBase: false,
	}

	middleware := NewStatic(config)
	if middleware == nil {
		t.Fatal("NewStatic() returned nil")
	}
}

func TestNewStatic_DefaultConfig(t *testing.T) {
	middleware := NewStatic(DefaultStatic)
	if middleware == nil {
		t.Fatal("NewStatic() with DefaultStatic returned nil")
	}
}
