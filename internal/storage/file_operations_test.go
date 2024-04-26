package storage

import (
	"testing"
)

func TestGetPathPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("GetPath did not panic")
		}
	}()
	GetPath("", "test")
}
