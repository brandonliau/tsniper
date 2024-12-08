package shared

import (
	"reflect"
	"testing"
)

func TestIntersection_EmptySlices(t *testing.T) {
	var data1 []int
	var data2 []int

	got := Intersection(data1, data2)
	want := []int{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Intersection(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestIntersection_NoOverlap(t *testing.T) {
	data1 := []int{1, 2, 3}
	data2 := []int{4, 5, 6}

	got := Intersection(data1, data2)
	want := []int{} // No common elements
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Intersection(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestIntersection_FullOverlap(t *testing.T) {
	data1 := []string{"a", "b", "c"}
	data2 := []string{"a", "b", "c"}

	got := Intersection(data1, data2)
	want := []string{"a", "b", "c"} // All elements are common
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Intersection(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestIntersection_PartialOverlap(t *testing.T) {
	data1 := []int{1, 2, 3, 4}
	data2 := []int{2, 4, 6}

	got := Intersection(data1, data2)
	want := []int{2, 4}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Intersection(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestIntersection_EmptyFirstSlice(t *testing.T) {
	data1 := []int{}
	data2 := []int{1, 2, 3}

	got := Intersection(data1, data2)
	want := []int{} // No intersection possible
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Intersection(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestIntersection_EmptySecondSlice(t *testing.T) {
	data1 := []int{1, 2, 3}
	data2 := []int{}

	got := Intersection(data1, data2)
	want := []int{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Intersection(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestIntersection_DuplicatesInBoth(t *testing.T) {
	data1 := []int{1, 1, 2, 2, 3}
	data2 := []int{1, 2, 2, 4}

	// Intersection should include each element of data2 that appears in data1, in data2 order.
	// Since data2 = {1,2,2,4}, and data1 has 1, 2, and 2 (multiple times),
	// each corresponding element should appear in the result.
	got := Intersection(data1, data2)
	want := []int{1, 2, 2}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Intersection(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestIntersection_OrderPreservation(t *testing.T) {
	data1 := []int{4, 3, 2, 1}
	data2 := []int{1, 2, 3, 4, 4}

	// Intersection should follow the order of data2.
	// data2 is {1, 2, 3, 4, 4}, and all appear in data1, so we should get {1, 2, 3, 4, 4} in that order.
	got := Intersection(data1, data2)
	want := []int{1, 2, 3, 4, 4}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Intersection(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestIntersection_StringsPartial(t *testing.T) {
	data1 := []string{"cat", "dog", "mouse"}
	data2 := []string{"dog", "elephant", "cat"}

	got := Intersection(data1, data2)
	want := []string{"dog", "cat"} // Both appear in data1, in the order they appear in data2
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Intersection(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}
