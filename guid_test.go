package goutils

import (
	"testing"
)

func TestGenGuid(t *testing.T) {
	s := NewSnowflake(1)
	for i := 0; i < 100; i++ {
		t.Log(s.NextID())
	}
}
