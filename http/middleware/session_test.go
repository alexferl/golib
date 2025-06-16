package middleware

import (
	"testing"
)

func TestSession_FlagSet(t *testing.T) {
	config := &Session{
		Enabled: true,
		Store:   SessionStoreCookie,
		Cookie: SessionCookieStore{
			Secret: "testsecret",
		},
	}

	fs := config.FlagSet()

	if fs == nil {
		t.Fatal("FlagSet() returned nil")
	}

	enabledFlag := fs.Lookup(SessionEnabled)
	if enabledFlag == nil {
		t.Errorf("Flag %s not found", SessionEnabled)
	} else {
		if enabledFlag.DefValue != "true" {
			t.Errorf("Flag %s default value = %v, want true", SessionEnabled, enabledFlag.DefValue)
		}
	}

	storeFlag := fs.Lookup(SessionStoreType)
	if storeFlag == nil {
		t.Errorf("Flag %s not found", SessionStoreType)
	} else {
		if storeFlag.DefValue != "cookie" {
			t.Errorf("Flag %s default value = %v, want cookie", SessionStoreType, storeFlag.DefValue)
		}
	}

	cookieSecretFlag := fs.Lookup(SessionCookieSecret)
	if cookieSecretFlag == nil {
		t.Errorf("Flag %s not found", SessionCookieSecret)
	} else {
		if cookieSecretFlag.DefValue != "testsecret" {
			t.Errorf("Flag %s default value = %v, want testsecret", SessionCookieSecret, cookieSecretFlag.DefValue)
		}
	}
}

func TestSession_FlagSet_Parse(t *testing.T) {
	config := &Session{
		Enabled: false,
		Store:   SessionStoreCookie,
		Cookie: SessionCookieStore{
			Secret: "changeme",
		},
	}

	fs := config.FlagSet()

	args := []string{
		"--session-enabled",
		"--session-store", "cookie",
		"--session-cookie-secret", "newsecret",
	}

	err := fs.Parse(args)
	if err != nil {
		t.Fatalf("Failed to parse flags: %v", err)
	}

	if !config.Enabled {
		t.Errorf("Enabled = %v, want true", config.Enabled)
	}
	if config.Store != SessionStoreCookie {
		t.Errorf("Store = %v, want %v", config.Store, SessionStoreCookie)
	}
	if config.Cookie.Secret != "newsecret" {
		t.Errorf("Cookie.Secret = %v, want newsecret", config.Cookie.Secret)
	}
}

func TestDefaultSession(t *testing.T) {
	if DefaultSession == nil {
		t.Fatal("DefaultSession is nil")
	}

	if DefaultSession.Enabled != false {
		t.Errorf("DefaultSession.Enabled = %v, want false", DefaultSession.Enabled)
	}
	if DefaultSession.Store != SessionStoreCookie {
		t.Errorf("DefaultSession.Store = %v, want %v", DefaultSession.Store, SessionStoreCookie)
	}
	if DefaultSession.Cookie.Secret != "changeme" {
		t.Errorf("DefaultSession.Cookie.Secret = %v, want changeme", DefaultSession.Cookie.Secret)
	}
}

func TestSession_FlagSet_DefaultValues(t *testing.T) {
	config := &Session{
		Enabled: true,
		Store:   SessionStoreCookie,
		Cookie: SessionCookieStore{
			Secret: "testsecret123",
		},
	}

	fs := config.FlagSet()

	enabledFlag := fs.Lookup(SessionEnabled)
	if enabledFlag == nil {
		t.Fatal("Enabled flag not found")
	}
	if enabledFlag.DefValue != "true" {
		t.Errorf("Enabled flag default = %v, want true", enabledFlag.DefValue)
	}

	storeFlag := fs.Lookup(SessionStoreType)
	if storeFlag == nil {
		t.Fatal("Store flag not found")
	}
	if storeFlag.DefValue != "cookie" {
		t.Errorf("Store flag default = %v, want cookie", storeFlag.DefValue)
	}

	cookieSecretFlag := fs.Lookup(SessionCookieSecret)
	if cookieSecretFlag == nil {
		t.Fatal("CookieSecret flag not found")
	}
	if cookieSecretFlag.DefValue != "testsecret123" {
		t.Errorf("CookieSecret flag default = %v, want testsecret123", cookieSecretFlag.DefValue)
	}
}

func TestSession_FlagSet_DisabledByDefault(t *testing.T) {
	config := DefaultSession

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

func TestSessionStore_String(t *testing.T) {
	tests := []struct {
		store SessionStore
		want  string
	}{
		{SessionStoreCookie, "cookie"},
		{SessionStore("invalid"), "unknown store: invalid"},
	}

	for _, tt := range tests {
		got := tt.store.String()
		if got != tt.want {
			t.Errorf("SessionStore(%s).String() = %v, want %v", tt.store, got, tt.want)
		}
	}
}

func TestSessionStore_Set(t *testing.T) {
	tests := []struct {
		value   string
		want    SessionStore
		wantErr bool
	}{
		{"cookie", SessionStoreCookie, false},
		{"COOKIE", SessionStoreCookie, false},
		{"invalid", SessionStore(""), true},
	}

	for _, tt := range tests {
		var store SessionStore
		err := store.Set(tt.value)
		if tt.wantErr {
			if err == nil {
				t.Errorf("SessionStore.Set(%q) expected error but got nil", tt.value)
			}
		} else {
			if err != nil {
				t.Errorf("SessionStore.Set(%q) unexpected error: %v", tt.value, err)
			}
			if store != tt.want {
				t.Errorf("SessionStore.Set(%q) = %v, want %v", tt.value, store, tt.want)
			}
		}
	}
}

func TestSessionStore_Type(t *testing.T) {
	var store SessionStore
	if got := store.Type(); got != "string" {
		t.Errorf("SessionStore.Type() = %v, want string", got)
	}
}

func TestNewSession(t *testing.T) {
	config := &Session{
		Enabled: true,
		Store:   SessionStoreCookie,
		Cookie: SessionCookieStore{
			Secret: "testsecret",
		},
	}

	middleware := NewSession(config)
	if middleware == nil {
		t.Fatal("NewSession() returned nil")
	}
}

func TestNewSession_DefaultConfig(t *testing.T) {
	middleware := NewSession(DefaultSession)
	if middleware == nil {
		t.Fatal("NewSession() with DefaultSession returned nil")
	}
}

func TestNewSession_InvalidStore(t *testing.T) {
	config := &Session{
		Enabled: true,
		Store:   SessionStore("invalid"),
	}

	middleware := NewSession(config)
	if middleware != nil {
		t.Error("NewSession() with invalid store should return nil")
	}
}
