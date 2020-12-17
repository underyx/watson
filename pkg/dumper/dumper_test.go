package dumper

import (
	"fmt"
	"math"
	"testing"

	"github.com/genkami/watson/pkg/vm"
	"github.com/google/go-cmp/cmp"
)

func TestSliceWritersInitialOpsIsEmpty(t *testing.T) {
	w := NewSliceWriter()
	ops := w.Ops()
	if len(ops) != 0 {
		t.Errorf("expected empty slice but got %#v", ops)
	}
}

func TestSliceWriterReturnsAllOpsThatAreWritten(t *testing.T) {
	w := NewSliceWriter()
	expected := []vm.Op{vm.Inew, vm.Iinc, vm.Ineg, vm.Fneg, vm.Snew, vm.Sadd}
	for _, op := range expected {
		err := w.Write(op)
		if err != nil {
			t.Fatal(err)
		}
	}
	actual := w.Ops()
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestDumpInt(t *testing.T) {
	test := func(n int64) {
		orig := vm.NewIntValue(n)
		w := NewSliceWriter()
		d := NewDumper(w)
		err := d.Dump(orig)
		if err != nil {
			t.Fatal(err)
		}
		dumped := w.Ops()
		v := vm.NewVM()
		for _, op := range dumped {
			err = v.Feed(op)
			if err != nil {
				t.Fatal(err)
			}
		}
		converted, err := v.Top()
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(orig, converted); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
	test(0)
	test(1)
	test(2)
	test(0x1234abcd)
	test(0x12345678abcdef0)
	test(-1)
}

func TestDumpFloat(t *testing.T) {
	test := func(n float64) {
		orig := vm.NewFloatValue(n)
		w := NewSliceWriter()
		d := NewDumper(w)
		err := d.Dump(orig)
		if err != nil {
			t.Fatal(err)
		}
		dumped := w.Ops()
		v := vm.NewVM()
		for _, op := range dumped {
			err = v.Feed(op)
			if err != nil {
				t.Fatal(err)
			}
		}
		converted, err := v.Top()
		if err != nil {
			t.Fatal(err)
		}
		if orig.IsNaN() {
			fmt.Printf("yey")
			if !converted.IsNaN() {
				t.Errorf("expected NaN but got %#v", converted)
			}
			return
		}
		if diff := cmp.Diff(orig, converted); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
	test(0)
	test(1)
	test(2)
	test(1.2345e67)
	test(1.2345e-67)
	test(-1.2345e67)
	test(-1.2345e-67)
	test(math.NaN())
	test(math.Inf(1))
	test(math.Inf(-1))
}
