package logger

import (
	"testing"

	"go.uber.org/zap/zaptest"
)

func TestInitZapSugaredLogger(t *testing.T) {
	Log = nil

	err := InitZapSugaredLogger()
	if err != nil {
		t.Fatalf("failed to initialize zap logger: %v", err)
	}

	if Log == nil {
		t.Fatal("expected Log to be initialized, but it is nil")
	}

	_, ok := Log.(*zapSugaredLogger)
	if !ok {
		t.Errorf("expected Log to be of type *zapSugaredLogger, got %T", Log)
	}

	Log.Infof("Testing zap initialization: %s", "success")
	Log.Infow("Testing zap w/ fields", "key", "value")
}

func TestInitZapSugaredLogger_Idempotency(t *testing.T) {
	Log = nil

	_ = InitZapSugaredLogger()
	firstInstance := Log

	err := InitZapSugaredLogger()
	if err != nil {
		t.Errorf("second call to InitZapSugaredLogger failed: %v", err)
	}

	if Log != firstInstance {
		t.Error("InitZapSugaredLogger should not overwrite Log if it is already set")
	}
}

func TestZapSugaredLogger_Methods(t *testing.T) {
	testLog := zaptest.NewLogger(t).Sugar()
	
	wrapper := &zapSugaredLogger{
		log: testLog,
	}

	t.Run("Logging Methods", func(t *testing.T) {
		wrapper.Debugf("test %s", "debugf")
		wrapper.Debugw("test debugw", "key", "val")
		
		wrapper.Infof("test %s", "infof")
		wrapper.Infow("test infow", "key", "val")
		
		wrapper.Warnf("test %s", "warnf")
		wrapper.Warnw("test warnw", "key", "val")
		
		wrapper.Errorf("test %s", "errorf")
		wrapper.Errorw("test errorw", "key", "val")
	})

	t.Run("Sync", func(t *testing.T) {
		_ = wrapper.Sync()
	})
}

func TestInitZapSugaredLogger_Coverage(t *testing.T) {
	Log = nil

	err := InitZapSugaredLogger()
	if err != nil {
		t.Fatalf("InitZapSugaredLogger failed: %v", err)
	}

	Log.Debugf("checking debugf")
	Log.Warnw("checking warnw", "attempt", 1)
}