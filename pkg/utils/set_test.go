package utils

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
	data1 := []int{1, 2, 3, 4}
	data2 := []int{5, 6, 7}

	got := Intersection(data1, data2)
	want := []int{} // no elements are common
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Intersection(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestIntersection_FullOverlap(t *testing.T) {
	data1 := []string{"apple", "banana", "cherry"}
	data2 := []string{"banana", "cherry", "apple"}

	got := Intersection(data1, data2)
	want := []string{"apple", "banana", "cherry"} // all elements are common
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Intersection(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestIntersection_PartialOverlap(t *testing.T) {
	data1 := []int{1, 2, 3, 4, 5}
	data2 := []int{2, 4, 6}

	got := Intersection(data1, data2)
	want := []int{2, 4} // only these are common
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Intersection(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestIntersection_EmptySecondSlice(t *testing.T) {
	data1 := []int{10, 20, 30}
	data2 := []int{}

	got := Intersection(data1, data2)
	want := []int{} // no elements in the second slice
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Intersection(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestIntersection_EmptyFirstSlice(t *testing.T) {
	data1 := []int{}
	data2 := []int{100, 200, 300}

	got := Intersection(data1, data2)
	want := []int{} // no elements in the first slice
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Intersection(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestIntersection_StringsPartial(t *testing.T) {
	data1 := []string{"hello", "world", "foo", "bar"}
	data2 := []string{"world", "xyz"}

	got := Intersection(data1, data2)
	want := []string{"world"} // only "world" is common
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Intersection(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestIntersection_UniqueElements(t *testing.T) {
	data1 := []int{1, 1, 2, 2, 3, 4, 4}
	data2 := []int{2, 4}

	// Duplicates should appear in the result only if they appear in both slices.
	got := Intersection(data1, data2)
	want := []int{2, 2, 4, 4}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Intersection(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

// --- DifferenceBy tests ---

type item struct {
	ID   int
	Name string
}

func byID(i item) int { return i.ID }

func TestDifferenceBy_EmptySlices(t *testing.T) {
	var data1 []item
	var data2 []item

	got := DifferenceBy(data1, data2, byID)
	want := []item{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("DifferenceBy(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestDifferenceBy_NoOverlap(t *testing.T) {
	data1 := []item{{1, "a"}, {2, "b"}}
	data2 := []item{{3, "c"}, {4, "d"}}

	got := DifferenceBy(data1, data2, byID)
	want := []item{{1, "a"}, {2, "b"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("DifferenceBy(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestDifferenceBy_FullOverlap(t *testing.T) {
	data1 := []item{{1, "a"}, {2, "b"}}
	data2 := []item{{2, "x"}, {1, "y"}}

	got := DifferenceBy(data1, data2, byID)
	want := []item{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("DifferenceBy(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestDifferenceBy_PartialOverlap(t *testing.T) {
	data1 := []item{{1, "a"}, {2, "b"}, {3, "c"}, {4, "d"}}
	data2 := []item{{2, "x"}, {4, "y"}}

	got := DifferenceBy(data1, data2, byID)
	want := []item{{1, "a"}, {3, "c"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("DifferenceBy(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestDifferenceBy_EmptySecondSlice(t *testing.T) {
	data1 := []item{{1, "a"}, {2, "b"}}
	data2 := []item{}

	got := DifferenceBy(data1, data2, byID)
	want := []item{{1, "a"}, {2, "b"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("DifferenceBy(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestDifferenceBy_EmptyFirstSlice(t *testing.T) {
	data1 := []item{}
	data2 := []item{{1, "a"}, {2, "b"}}

	got := DifferenceBy(data1, data2, byID)
	want := []item{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("DifferenceBy(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestDifferenceBy_DifferentValsSameKey(t *testing.T) {
	data1 := []item{{1, "original"}, {2, "original"}}
	data2 := []item{{1, "different"}}

	got := DifferenceBy(data1, data2, byID)
	want := []item{{2, "original"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("DifferenceBy(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestDifferenceBy_DuplicateKeys(t *testing.T) {
	data1 := []item{{1, "a"}, {1, "b"}, {2, "c"}, {2, "d"}}
	data2 := []item{{2, "x"}}

	got := DifferenceBy(data1, data2, byID)
	want := []item{{1, "a"}, {1, "b"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("DifferenceBy(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestDifferenceBy_StringKey(t *testing.T) {
	byName := func(i item) string { return i.Name }
	data1 := []item{{1, "alice"}, {2, "bob"}, {3, "charlie"}}
	data2 := []item{{99, "bob"}, {100, "charlie"}}

	got := DifferenceBy(data1, data2, byName)
	want := []item{{1, "alice"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("DifferenceBy(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

// --- IntersectionBy tests ---

func TestIntersectionBy_EmptySlices(t *testing.T) {
	var data1 []item
	var data2 []item

	got := IntersectionBy(data1, data2, byID)
	want := []item{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("IntersectionBy(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestIntersectionBy_NoOverlap(t *testing.T) {
	data1 := []item{{1, "a"}, {2, "b"}}
	data2 := []item{{3, "c"}, {4, "d"}}

	got := IntersectionBy(data1, data2, byID)
	want := []item{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("IntersectionBy(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestIntersectionBy_FullOverlap(t *testing.T) {
	data1 := []item{{1, "a"}, {2, "b"}}
	data2 := []item{{2, "x"}, {1, "y"}}

	got := IntersectionBy(data1, data2, byID)
	want := []item{{1, "a"}, {2, "b"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("IntersectionBy(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestIntersectionBy_PartialOverlap(t *testing.T) {
	data1 := []item{{1, "a"}, {2, "b"}, {3, "c"}, {4, "d"}}
	data2 := []item{{2, "x"}, {4, "y"}}

	got := IntersectionBy(data1, data2, byID)
	want := []item{{2, "b"}, {4, "d"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("IntersectionBy(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestIntersectionBy_EmptySecondSlice(t *testing.T) {
	data1 := []item{{1, "a"}, {2, "b"}}
	data2 := []item{}

	got := IntersectionBy(data1, data2, byID)
	want := []item{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("IntersectionBy(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestIntersectionBy_EmptyFirstSlice(t *testing.T) {
	data1 := []item{}
	data2 := []item{{1, "a"}, {2, "b"}}

	got := IntersectionBy(data1, data2, byID)
	want := []item{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("IntersectionBy(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestIntersectionBy_DifferentValsSameKey(t *testing.T) {
	data1 := []item{{1, "original"}, {2, "original"}}
	data2 := []item{{1, "different"}}

	got := IntersectionBy(data1, data2, byID)
	want := []item{{1, "original"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("IntersectionBy(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestIntersectionBy_DuplicateKeys(t *testing.T) {
	data1 := []item{{1, "a"}, {1, "b"}, {2, "c"}, {2, "d"}}
	data2 := []item{{2, "x"}}

	got := IntersectionBy(data1, data2, byID)
	want := []item{{2, "c"}, {2, "d"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("IntersectionBy(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}

func TestIntersectionBy_StringKey(t *testing.T) {
	byName := func(i item) string { return i.Name }
	data1 := []item{{1, "alice"}, {2, "bob"}, {3, "charlie"}}
	data2 := []item{{99, "bob"}, {100, "charlie"}}

	got := IntersectionBy(data1, data2, byName)
	want := []item{{2, "bob"}, {3, "charlie"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("IntersectionBy(%v, %v) = %v; want %v", data1, data2, got, want)
	}
}
