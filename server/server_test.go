package server

import (
	"reflect"
	"testing"
)

func TestCheckParamSlice(t *testing.T) {
	t.Run("nil pointer", func(t *testing.T) {
		var p *[]string
		got := checkParamSlice(p)
		want := make([]string, 0)
		assertSlice(t, got, want)
	})

	t.Run("nil pointer", func(t *testing.T) {
		p := []string{"test1", "test2"}
		got := checkParamSlice(&p)
		want := p
		assertSlice(t, got, want)
	})
}

func assertSlice(t *testing.T, got, want []string) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("mismatch: got %v, want %v", got, want)
	}
}
