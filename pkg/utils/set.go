package utils

// Returns the elements in `data1` that aren't in `data2`
func Difference[T comparable](data1 []T, data2 []T) []T {
	diff := make([]T, 0)
	hash := make(map[T]struct{}, len(data2))
	for _, i := range data2 {
		hash[i] = struct{}{}
	}
	for _, j := range data1 {
		if _, ok := hash[j]; !ok {
			diff = append(diff, j)
		}
	}
	return diff
}

// Returns the elements in `data1` that are in `data2` by the given key function
func DifferenceBy[T any, K comparable](data1 []T, data2 []T, key func(T) K) []T {
	diff := make([]T, 0)
	hash := make(map[K]struct{}, len(data2))
	for _, i := range data2 {
		hash[key(i)] = struct{}{}
	}
	for _, j := range data1 {
		if _, ok := hash[key(j)]; !ok {
			diff = append(diff, j)
		}
	}
	return diff
}

// Returns the elements in `data1` that are in `data2`
func Intersection[T comparable](data1 []T, data2 []T) []T {
	inter := make([]T, 0)
	hash := make(map[T]struct{}, len(data2))
	for _, i := range data2 {
		hash[i] = struct{}{}
	}
	for _, j := range data1 {
		if _, ok := hash[j]; ok {
			inter = append(inter, j)
		}
	}
	return inter
}

// Returns the elements in `data1` that are in `data2` by the given key function
func IntersectionBy[T any, K comparable](data1 []T, data2 []T, key func(T) K) []T {
	inter := make([]T, 0)
	hash := make(map[K]struct{}, len(data2))
	for _, i := range data2 {
		hash[key(i)] = struct{}{}
	}
	for _, j := range data1 {
		if _, ok := hash[key(j)]; ok {
			inter = append(inter, j)
		}
	}
	return inter
}
