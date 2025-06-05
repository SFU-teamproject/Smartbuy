package mocklogger

import (
	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Errorf(format string, args ...any) { m.Called(format, args) }
func (m *MockLogger) Errorln(args ...any)               { m.Called(args...) }
func (m *MockLogger) Infof(format string, args ...any)  { m.Called(format, args) }
func (m *MockLogger) Infoln(args ...any)                { m.Called(args...) }
