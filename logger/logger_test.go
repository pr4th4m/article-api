package logger

import (
	"testing"

	"go.uber.org/zap"
)

func TestNewLogger(t *testing.T) {

	lg := New(zap.DebugLevel)
	lg.Info("Testing logger")

}
