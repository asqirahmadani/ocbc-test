package logger

import (
	"testing"
)

type mockLogger struct {
	lastMessage string
}

func (m *mockLogger) Debugf(t string, a ...any)   {}
func (m *mockLogger) Debugw(m_ string, k ...any) {}
func (m *mockLogger) Infof(t string, a ...any)    {}
func (m *mockLogger) Infow(m_ string, k ...any)  {}
func (m *mockLogger) Warnf(t string, a ...any)    {}
func (m *mockLogger) Warnw(m_ string, k ...any)  {}
func (m *mockLogger) Errorf(t string, a ...any)   {}
func (m *mockLogger) Errorw(m_ string, k ...any) {
	m.lastMessage = m_ 
}
func (m *mockLogger) Sync() error { return nil }

func TestSetLogger(t *testing.T) {
	myMock := &mockLogger{}

	setLogger(myMock)

	if Log == nil {
		t.Fatal("expected Log to be set, but it remains nil")
	}

	testMsg := "test error message"
	Log.Errorw(testMsg)

	if myMock.lastMessage != testMsg {
		t.Errorf("expected %q, got %q", testMsg, myMock.lastMessage)
	}
}