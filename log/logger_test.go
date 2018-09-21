package log

import (
	"github.com/riposa/utils/errors"
	"testing"
)

func TestLogger_EnableDebug(t *testing.T) {
	logger := New()
	logger.EnableDebug()
	logger.Info("info")
	logger.Debug("debug")
	logger.DisableDebug()
	logger.Info("info")
	logger.Debug("debug")
}

func TestLogger_Exception(t *testing.T) {
	logger := New()
	err := errors.New(2)
	logger.Exception(err)
}

func TestLogger_Print(t *testing.T) {
	logger := New()
	logger.Print("abcd")
}
