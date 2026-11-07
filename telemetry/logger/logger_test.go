package logger

import (
	"testing"
)

func TestUnit_Logger(t *testing.T) {

	err := SetupLogging(
		DebugLevel,
		DebugLevel,
		"test.log",
	)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	Debug("Debug Message arg1: %#v arg2: %#v", 1, "test")
	DebugWithArgs("Debug Message With args", "arg1", 1, "arg2", "test")

	Info("Info Message arg1: %#v arg2: %#v", 1, "test")
	InfoWithArgs("Info Message With args", "arg1", 1, "arg2", "test")

	Warn("Warn Message arg1: %#v arg2: %#v", 1, "test")
	WarnWithArgs("Warn Message With args", "arg1", 1, "arg2", "test")

	Error("Error Message arg1: %#v arg2: %#v", 1, "test")
	ErrorWithArgs("Error Message With args", "arg1", 1, "arg2", "test")
}
