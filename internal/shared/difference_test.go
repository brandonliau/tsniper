package shared

import (
	"reflect"
	"testing"
)

func TestDifference_EmptySlices(t *testing.T) {
	var data1 []int
	var data2 []int

	got := Difference(data1, data2)
	want := []int{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Difference(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestDifference_NoOverlap(t *testing.T) {
	data1 := []int{1, 2, 3, 4}
	data2 := []int{5, 6, 7}

	got := Difference(data1, data2)
	want := []int{1, 2, 3, 4} // none of these appear in data2
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Difference(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestDifference_FullOverlap(t *testing.T) {
	data1 := []string{"apple", "banana", "cherry"}
	data2 := []string{"banana", "cherry", "apple"}

	got := Difference(data1, data2)
	want := []string{} // all are found in data2
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Difference(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestDifference_PartialOverlap(t *testing.T) {
	data1 := []int{1, 2, 3, 4, 5}
	data2 := []int{2, 4, 6}

	got := Difference(data1, data2)
	want := []int{1, 3, 5} // Only these are not in data2
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Difference(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestDifference_EmptySecondSlice(t *testing.T) {
	data1 := []int{10, 20, 30}
	data2 := []int{}

	got := Difference(data1, data2)
	want := []int{10, 20, 30}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Difference(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestDifference_EmptyFirstSlice(t *testing.T) {
	data1 := []int{}
	data2 := []int{100, 200, 300}

	got := Difference(data1, data2)
	want := []int{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Difference(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestDifference_StringsPartial(t *testing.T) {
	data1 := []string{"hello", "world", "foo", "bar"}
	data2 := []string{"world", "xyz"}

	got := Difference(data1, data2)
	want := []string{"hello", "foo", "bar"} // Only "world" is excluded
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Difference(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestDifference_UniqueElements(t *testing.T) {
	data1 := []int{1, 1, 2, 2, 3, 4, 4}
	data2 := []int{2}

	// This should remove all occurrences of 2, but not other numbers.
	// Duplicate elements in data1 that are not in data2 should still appear.
	got := Difference(data1, data2)
	want := []int{1, 1, 3, 4, 4}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Difference(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}
