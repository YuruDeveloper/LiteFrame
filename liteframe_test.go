package liteframe

import (
	"testing"
)
func TestHello(t *testing.T) {
	want := "Hello, from LiteFrame!"
	if got := Hello(); got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}
