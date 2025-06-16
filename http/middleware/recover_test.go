package middleware

import (
	"testing"
)

func TestRecover_FlagSet(t *testing.T) {
	config := &Recover{
		Enabled:             false,
		StackSize:           8192,
		DisableStackAll:     true,
		DisablePrintStack:   true,
		DisableErrorHandler: true,
	}

	fs := config.FlagSet()

	if fs == nil {
		t.Fatal("FlagSet() returned nil")
	}

	enabledFlag := fs.Lookup(RecoverEnabled)
	if enabledFlag == nil {
		t.Errorf("Flag %s not found", RecoverEnabled)
	} else {
		if enabledFlag.DefValue != "false" {
			t.Errorf("Flag %s default value = %v, want false", RecoverEnabled, enabledFlag.DefValue)
		}
	}

	stackSizeFlag := fs.Lookup(RecoverStackSize)
	if stackSizeFlag == nil {
		t.Errorf("Flag %s not found", RecoverStackSize)
	} else {
		if stackSizeFlag.DefValue != "8192" {
			t.Errorf("Flag %s default value = %v, want 8192", RecoverStackSize, stackSizeFlag.DefValue)
		}
	}

	disableStackAllFlag := fs.Lookup(RecoverDisableStackAll)
	if disableStackAllFlag == nil {
		t.Errorf("Flag %s not found", RecoverDisableStackAll)
	} else {
		if disableStackAllFlag.DefValue != "true" {
			t.Errorf("Flag %s default value = %v, want true", RecoverDisableStackAll, disableStackAllFlag.DefValue)
		}
	}
}

func TestRecover_FlagSet_Parse(t *testing.T) {
	config := &Recover{
		Enabled:             true,
		StackSize:           4096,
		DisableStackAll:     false,
		DisablePrintStack:   false,
		DisableErrorHandler: false,
	}

	fs := config.FlagSet()

	args := []string{
		"--recover-enabled=false",
		"--recover-stack-size", "16384",
		"--recover-disable-stack-all",
		"--recover-disable-print-stack",
		"--recover-disable-error-handler",
	}

	err := fs.Parse(args)
	if err != nil {
		t.Fatalf("Failed to parse flags: %v", err)
	}

	if config.Enabled {
		t.Errorf("Enabled = %v, want false", config.Enabled)
	}
	if config.StackSize != 16384 {
		t.Errorf("StackSize = %v, want 16384", config.StackSize)
	}
	if !config.DisableStackAll {
		t.Errorf("DisableStackAll = %v, want true", config.DisableStackAll)
	}
	if !config.DisablePrintStack {
		t.Errorf("DisablePrintStack = %v, want true", config.DisablePrintStack)
	}
	if !config.DisableErrorHandler {
		t.Errorf("DisableErrorHandler = %v, want true", config.DisableErrorHandler)
	}
}

func TestDefaultRecover(t *testing.T) {
	if DefaultRecover == nil {
		t.Fatal("DefaultRecover is nil")
	}

	if DefaultRecover.Enabled != true {
		t.Errorf("DefaultRecover.Enabled = %v, want true", DefaultRecover.Enabled)
	}
	if DefaultRecover.StackSize != 4096 {
		t.Errorf("DefaultRecover.StackSize = %v, want 4096", DefaultRecover.StackSize)
	}
	if DefaultRecover.DisableStackAll != false {
		t.Errorf("DefaultRecover.DisableStackAll = %v, want false", DefaultRecover.DisableStackAll)
	}
	if DefaultRecover.DisablePrintStack != false {
		t.Errorf("DefaultRecover.DisablePrintStack = %v, want false", DefaultRecover.DisablePrintStack)
	}
	if DefaultRecover.DisableErrorHandler != false {
		t.Errorf("DefaultRecover.DisableErrorHandler = %v, want false", DefaultRecover.DisableErrorHandler)
	}
}

func TestRecover_FlagSet_DefaultValues(t *testing.T) {
	config := &Recover{
		Enabled:             false,
		StackSize:           2048,
		DisableStackAll:     true,
		DisablePrintStack:   false,
		DisableErrorHandler: true,
	}

	fs := config.FlagSet()

	enabledFlag := fs.Lookup(RecoverEnabled)
	if enabledFlag == nil {
		t.Fatal("Enabled flag not found")
	}
	if enabledFlag.DefValue != "false" {
		t.Errorf("Enabled flag default = %v, want false", enabledFlag.DefValue)
	}

	stackSizeFlag := fs.Lookup(RecoverStackSize)
	if stackSizeFlag == nil {
		t.Fatal("StackSize flag not found")
	}
	if stackSizeFlag.DefValue != "2048" {
		t.Errorf("StackSize flag default = %v, want 2048", stackSizeFlag.DefValue)
	}

	disableStackAllFlag := fs.Lookup(RecoverDisableStackAll)
	if disableStackAllFlag == nil {
		t.Fatal("DisableStackAll flag not found")
	}
	if disableStackAllFlag.DefValue != "true" {
		t.Errorf("DisableStackAll flag default = %v, want true", disableStackAllFlag.DefValue)
	}
}

func TestRecover_FlagSet_EnabledByDefault(t *testing.T) {
	config := DefaultRecover

	fs := config.FlagSet()

	var args []string

	err := fs.Parse(args)
	if err != nil {
		t.Fatalf("Failed to parse empty flags: %v", err)
	}

	if !config.Enabled {
		t.Errorf("Enabled = %v, want true (default)", config.Enabled)
	}
}

func TestNewRecover(t *testing.T) {
	config := &Recover{
		Enabled:             true,
		StackSize:           8192,
		DisableStackAll:     true,
		DisablePrintStack:   false,
		DisableErrorHandler: false,
	}

	middleware := NewRecover(config)
	if middleware == nil {
		t.Fatal("NewRecover() returned nil")
	}
}

func TestNewRecover_DefaultConfig(t *testing.T) {
	middleware := NewRecover(DefaultRecover)
	if middleware == nil {
		t.Fatal("NewRecover() with DefaultRecover returned nil")
	}
}
