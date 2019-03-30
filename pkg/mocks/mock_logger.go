package mocks

import "testing"

type TestLogger struct {
	T testing.TB
}

func (tl TestLogger) Log(args ...interface{}) error {
	tl.T.Log(args...)
	return nil
}
