package main

import (
	"reflect"
	"testing"
)

func TestGenerator(t *testing.T) {
	g := NewGenerator()

	for _, tc := range []struct {
		in  uint64
		out []byte
	}{
		{0, []byte("0")},
		{61, []byte("Z")},
		{62, []byte("01")},
	} {
		got := g.EncodeID(tc.in)

		if !reflect.DeepEqual(tc.out, got) {
			t.Fatalf("expect %q, but got %q", tc.out, got)
		}

		id := g.DecodeID(got)

		if id != tc.in {
			t.Fatalf("expect %d, but got %d", tc.in, id)
		}
	}
}
